package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims defines the structure for JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		// Extract the token from the Authorization header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // No "Bearer" prefix
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Referencing jwtKey from handlers.go, no redeclaration needed here
			return jwtKey, nil
		})

		if err != nil || !token.Valid || time.Now().After(claims.ExpiresAt.Time) {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Token is valid, pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
