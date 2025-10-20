package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ethanhosier/reel-farm/internal/api"
	"github.com/ethanhosier/reel-farm/internal/handler"
	"github.com/ethanhosier/reel-farm/internal/repository"
	"github.com/ethanhosier/reel-farm/internal/service"
	"github.com/jackc/pgx/v5"
)

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

	// Create services and handlers
	userRepo := repository.NewUserRepository(conn)
	userService := service.NewUserService(userRepo)
	apiServer := handler.NewAPIServer(userService)

	// Create HTTP handler using generated code with custom error handler
	handler := api.HandlerWithOptions(apiServer, api.StdHTTPServerOptions{
		BaseRouter: http.NewServeMux(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{
				Error:   "invalid_parameter",
				Message: err.Error(),
			})
		},
	})

	// Start the server
	fmt.Printf("🚀 Reel Farm server starting on port %s\n", port)
	fmt.Printf("📡 Health check available at: http://localhost:%s/health\n", port)
	fmt.Printf("👤 User endpoint available at: http://localhost:%s/users/{id}\n", port)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
