package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// TraceIDKey is the context key for trace ID
type TraceIDKey struct{}

// wrappedWriter wraps http.ResponseWriter to capture status code
type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Logging creates a middleware that logs HTTP requests with trace ID
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate or extract trace ID
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()[:8] // Use first 8 characters for brevity
		}

		// Add trace ID to context
		ctx := context.WithValue(r.Context(), TraceIDKey{}, traceID)
		r = r.WithContext(ctx)

		// Add trace ID to response headers
		w.Header().Set("X-Trace-ID", traceID)

		wrapped := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}

		log.Printf("[%s] Received %s %s", traceID, r.Method, r.URL.Path)
		next.ServeHTTP(wrapped, r)
		log.Printf("[%s] Completed operation %v %s %s in %v", traceID, wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}
