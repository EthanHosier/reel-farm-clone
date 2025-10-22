package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ethanhosier/reel-farm/internal/api"
	"github.com/ethanhosier/reel-farm/internal/handler"
	"github.com/ethanhosier/reel-farm/internal/middleware"
	"github.com/ethanhosier/reel-farm/internal/repository"
	"github.com/ethanhosier/reel-farm/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Parse command line flags
	noAuth := flag.Bool("noAuth", false, "Disable authentication (for development/testing)")
	flag.Parse()

	// Load .env file if it exists
	envPath := filepath.Join(".", ".env")
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Warning: Failed to load .env file: %v", err)
		} else {
			log.Println("✅ Loaded environment variables from .env file")
		}
	} else {
		log.Println("ℹ️  No .env file found, using system environment variables")
	}

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
	subscriptionService := service.NewSubscriptionService(userRepo)
	apiServer := handler.NewAPIServer(userService, subscriptionService)

	// Create HTTP handler using generated code with auth middleware
	apiHandler := api.HandlerWithOptions(apiServer, api.StdHTTPServerOptions{
		BaseRouter: http.NewServeMux(),
		Middlewares: []api.MiddlewareFunc{
			api.MiddlewareFunc(middleware.AuthMiddleware(*noAuth)),
		},
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{
				Error:   "invalid_parameter",
				Message: err.Error(),
			})
		},
	})

	// Create webhook handler
	webhookHandler := handler.NewWebhookHandler(subscriptionService)

	// Create main router
	mux := http.NewServeMux()

	// Add API routes (with auth middleware)
	mux.Handle("/", middleware.CORSMiddleware()(apiHandler))

	// Add webhook routes (no auth middleware, but with CORS)
	mux.Handle("/webhooks/stripe", middleware.CORSMiddleware()(webhookHandler))

	// Start the server
	fmt.Printf("🚀 Reel Farm server starting on port %s\n", port)
	if *noAuth {
		fmt.Printf("🔧 No-auth mode enabled - authentication disabled\n")
	}
	fmt.Printf("📡 Health check available at: http://localhost:%s/health\n", port)
	fmt.Printf("👤 User endpoint available at: http://localhost:%s/user\n", port)
	fmt.Printf("🔗 Stripe webhook available at: http://localhost:%s/webhooks/stripe\n", port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
