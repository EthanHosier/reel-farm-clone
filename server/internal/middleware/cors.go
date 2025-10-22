package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// CORSMiddleware creates a CORS middleware for development
func CORSMiddleware() func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173", // Vite dev server
			"http://localhost:3000", // Backend (for self-requests)
			"https://reel-farm-clone-frontend.vercel.app/",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodPatch,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		AllowCredentials: true,
		Debug:            false, // Enable debug to see CORS logs
	})

	return c.Handler
}
