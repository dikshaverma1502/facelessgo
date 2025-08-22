package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Event struct {
	Event      string                 `json:"event"`
	Properties map[string]interface{} `json:"properties"`
}

func trackHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("📩 Received POST request for /track")

	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Printf("❌ Failed to decode body: %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Printf("✅ Event: %s\n", event.Event)
	log.Printf("🔹 Properties: %+v\n", event.Properties)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	http.HandleFunc("/track", trackHandler)
	log.Println("🚀 Server starting on port 10000")
	log.Fatal(http.ListenAndServe(":10000", nil))
}
