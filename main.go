package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var mixpanelToken string

// Mixpanel event format
type MixpanelEvent struct {
	Event      string                 `json:"event"`
	Properties map[string]interface{} `json:"properties"`
}

// API handler
func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s\n", r.Method, r.URL.Path)

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Request body: %s\n", string(body))

	// Unmarshal into MixpanelEvent
	var event MixpanelEvent
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Ensure Mixpanel token is set
	if mixpanelToken == "" {
		http.Error(w, "Mixpanel token not configured", http.StatusInternalServerError)
		return
	}

	// Add token to event properties
	event.Properties["token"] = mixpanelToken

	// Encode event as JSON
	eventData, _ := json.Marshal(event)
	reqBody := []byte(fmt.Sprintf(`{"event": "%s", "properties": %s}`, event.Event, string(eventData)))

	// Send to Mixpanel track API
	resp, err := http.Post("https://api.mixpanel.com/track?ip=1", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Error sending to Mixpanel: %v", err)
		http.Error(w, "Failed to send to Mixpanel", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("Mixpanel response: %s", resp.Status)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Event logged and sent to Mixpanel")
}

func main() {
	// Load Mixpanel token from environment variable
	mixpanelToken = os.Getenv("MIXPANEL_TOKEN")
	if mixpanelToken == "" {
		log.Fatal("MIXPANEL_TOKEN environment variable not set")
	}

	http.HandleFunc("/track", handler)

	port := ":10000"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = ":" + fromEnv
	}

	log.Printf("Starting server on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
