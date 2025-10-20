package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Response represents a simple JSON response
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Port    string `json:"port"`
}

// healthHandler handles the /health endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Create response
	response := Response{
		Message: "Reel Farm API is healthy!!",
		Status:  "ok",
		Port:    port,
	}

	// Convert to JSON and send
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Register the health endpoint
	http.HandleFunc("/health", healthHandler)

	// Start the server
	fmt.Printf("🚀 Reel Farm server starting on port %s\n", port)
	fmt.Printf("📡 Health check available at: http://localhost:%s/health\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
