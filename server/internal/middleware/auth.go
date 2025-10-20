package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/ethanhosier/reel-farm/internal/api"
	"github.com/ethanhosier/reel-farm/internal/context_keys"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens and extracts user ID
func AuthMiddleware(devMode bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health endpoint
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			// Development mode: use hardcoded user ID
			if devMode {
				userID := "65a950f6-a3b0-4be2-824a-b99051d5a62f"
				ctx := context_keys.SetUserID(r.Context(), userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				http.Error(w, "JWT_SECRET is required", http.StatusInternalServerError)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(api.ErrorResponse{
					Error:   "unauthorized",
					Message: "Missing or invalid authorization header",
				})
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(api.ErrorResponse{
					Error:   "unauthorized",
					Message: "Invalid or expired token",
				})
				return
			}

			// Extract user ID from token claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(api.ErrorResponse{
					Error:   "unauthorized",
					Message: "Invalid token claims",
				})
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(api.ErrorResponse{
					Error:   "unauthorized",
					Message: "User ID not found in token",
				})
				return
			}

			// Add user ID to context
			ctx := context_keys.SetUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
