package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

// Response structure for clean JSON replies
type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// handler function
func handler(w http.ResponseWriter, r *http.Request) {
	// Log request method and path
	log.Printf("📥 Received %s request for %s\n", r.Method, r.URL.Path)

	// Log headers
	log.Println("🔹 Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("%s: %s\n", name, value)
		}
	}

	// Read and log body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Error reading body: %v", err)
	} else if len(body) > 0 {
		log.Printf("📦 Body: %s\n", string(body))
	} else {
		log.Println("📦 Body: <empty>")
	}
	defer r.Body.Close()

	// Always respond in JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := ApiResponse{
		Status:  "success",
		Message: "Request received and logged",
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", handler)

	// Use env PORT if provided (for Render, Vercel, etc.)
	port := ":10000" // fallback
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = ":" + fromEnv
	}

	log.Printf("🚀 Starting faceless server on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("❌ Server failed: %v\n", err)
	}
}
