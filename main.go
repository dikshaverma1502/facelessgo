package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// handler function
func handler(w http.ResponseWriter, r *http.Request) {
	// Log request method and path
	log.Printf("Received %s request for %s\n", r.Method, r.URL.Path)

	// Log headers
	log.Println("Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("%s: %s\n", name, value)
		}
	}

	// Read and log body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	} else {
		log.Printf("Body: %s\n", string(body))
	}
	defer r.Body.Close()

	// Respond back
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "✅ Request received and logged")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func main() {
	// Routes
	http.HandleFunc("/", handler)
	http.HandleFunc("/healthz", healthCheck)

	// Render provides PORT via environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000" // fallback for local run
	}

	log.Printf("🚀 Starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ Server failed: %v\n", err)
	}
}
