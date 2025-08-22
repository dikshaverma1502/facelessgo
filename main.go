package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Event struct {
	Event      string                 `json:"event"`
	Properties map[string]interface{} `json:"properties"`
}

func trackHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s\n", r.Method, r.URL.Path)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Body: %s\n", string(body))

	var event Event
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Event received: %s", event.Event)
	log.Printf("Properties: %+v", event.Properties)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func main() {
	http.HandleFunc("/track", trackHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	log.Printf("Starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
