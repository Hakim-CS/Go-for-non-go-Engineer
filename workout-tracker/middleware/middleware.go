package middleware

import (
	"context"
	"net/http"
	"strings"
	"workout-tracker/auth"
)

// what is a middleware?
// A middleware is a function that wraps an HTTP handler to perform some processing
// on the request before passing it to the handler, or on the response before sending it back to the client.
// Middlewares are commonly used for tasks like authentication, logging, CORS handling, etc.

// contextKey is a custom type for context keys to avoid conflicts
type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware checks if the request has a valid JWT token
// If valid, it adds the user info to the request context
// If invalid, it returns a 401 Unauthorized response
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		// Expected format: "Bearer <token>"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Split the header to get just the token part
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Validate the token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add the user info to the request context
		// This allows handlers to access the current user's information
		ctx := context.WithValue(r.Context(), UserContextKey, claims)

		// Call the next handler with the updated context
		next(w, r.WithContext(ctx))
	}
}

// GetUserFromContext extracts the user claims from the request context
// This should be called in handlers that are protected by AuthMiddleware
func GetUserFromContext(r *http.Request) *auth.Claims {
	claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}

// CORSMiddleware adds CORS headers to allow frontend apps to call the API
// This is useful when your frontend is on a different domain/port
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin (in production, specify your frontend URL)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
