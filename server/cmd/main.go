package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ethanhosier/reel-farm/internal/handler"
	"github.com/ethanhosier/reel-farm/internal/repository"
	"github.com/ethanhosier/reel-farm/internal/service"
	"github.com/jackc/pgx/v5"
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
		Message: "Reel Farm API is healthy!",
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

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Register the health endpoint
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("GET /users/{id}", createUserHandler(conn).GetUserAccount)

	// Start the server
	fmt.Printf("🚀 Reel Farm server starting on port %s\n", port)
	fmt.Printf("📡 Health check available at: http://localhost:%s/health\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func createUserHandler(conn *pgx.Conn) *handler.UserHandler {
	return handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(conn)))
}
