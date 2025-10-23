package middleware

import (
	"net/http"
	"strings"
)

// PathStrip creates a middleware that strips a prefix from the request path
func PathStrip(prefix string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Strip the prefix from the path
			if strings.HasPrefix(r.URL.Path, prefix) {
				r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
				// Ensure path starts with /
				if r.URL.Path == "" {
					r.URL.Path = "/"
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
