package main

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
)

type MixpanelEvent struct {
    Event      string                 `json:"event"`
    Properties map[string]interface{} `json:"properties"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    // 🔓 Allow all origins and methods
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // 🛑 Handle preflight request
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    body, _ := io.ReadAll(r.Body)
    defer r.Body.Close()

    fmt.Println("🔔 New Request:")
    fmt.Printf("Method: %s\n", r.Method)
    fmt.Printf("Path: %s\n", r.URL.Path)
    fmt.Printf("Headers: %v\n", r.Header)
    fmt.Printf("Raw Body:\n%s\n", string(body))

    form, err := url.ParseQuery(string(body))
    if err != nil {
        fmt.Println("❌ Failed to parse form:", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    encodedData := form.Get("data")
    if encodedData == "" {
        fmt.Println("⚠️ No 'data' field found")
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    decodedURL, err := url.QueryUnescape(encodedData)
    if err != nil {
        fmt.Println("❌ URL decode error:", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    decodedBase64, err := base64.StdEncoding.DecodeString(decodedURL)
    if err != nil {
        fmt.Println("❌ Base64 decode error:", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    var events []MixpanelEvent
    if err := json.Unmarshal(decodedBase64, &events); err != nil {
        fmt.Println("❌ JSON unmarshal error:", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    for _, e := range events {
        fmt.Printf("🎯 Event: %s\n", e.Event)
        fmt.Println("📦 Properties:")
        for k, v := range e.Properties {
            fmt.Printf("  - %s: %v\n", k, v)
        }
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Received and decoded"))
}

func main() {
    http.HandleFunc("/track/", handler)
    fmt.Println("🚀 Listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}