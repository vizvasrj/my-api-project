package middelware

import (
	"fmt"
	"net/http"
	"src/internal/helpers"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// Define a secret key for JWT signing
var secretKey = []byte("some-key")

// Authentication Middleware
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the JWT token from the Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is null", http.StatusUnauthorized)
			return
		}

		// Parse and validate the JWT token
		claims, err := helpers.ValidateToken(tokenString)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		context.Set(r, "uid", claims.Uid)
		context.Set(r, "username", claims.Username)

		// Continue to the next middleware or handler
		next.ServeHTTP(w, r)
	}
}

// Authorization Middleware
func AuthorizeMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user roles or permissions from the JWT claims
		claims, ok := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Implement authorization logic based on user roles or permissions
		// Example: Check if the user has "admin" role to access a resource
		if role, ok := claims["role"].(string); ok && role == "admin" {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
