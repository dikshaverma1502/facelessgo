package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// TrackEvent represents the JSON Mixpanel sends
type TrackEvent struct {
	Event      string                 `json:"event"`
	Properties map[string]interface{} `json:"properties"`
}

// handler for all other routes
func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("📥 Received %s request for %s\n", r.Method, r.URL.Path)
	log.Println("🔹 Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("%s: %s\n", name, value)
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	} else if len(body) > 0 {
		log.Printf("📦 Body: %s\n", string(body))
	} else {
		log.Println("📦 Body: <empty>")
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"success","message":"Request received and logged"}`)
}

// mixpanel track endpoint
func trackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var event TrackEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Printf("❌ Failed to decode JSON: %v", err)
		return
	}

	log.Printf("🎯 Mixpanel Event: %s", event.Event)
	log.Printf("   Properties: %+v", event.Properties)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"ok","message":"Event logged"}`)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/track", trackHandler)

	port := ":10000" // fallback
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = ":" + fromEnv
	}

	log.Printf("🚀 Starting faceless server on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
