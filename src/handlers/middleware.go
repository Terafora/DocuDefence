package handlers

import (
	"DocuDefense/src/models"
	"log"
	"net/http"
)

func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve Basic Auth credentials (email and password)
		email, password, ok := r.BasicAuth()
		if !ok {
			log.Println("Authorization header missing or invalid")
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		// Find the user by email
		var foundUser *models.User
		for _, user := range users {
			if user.Email == email {
				foundUser = user
				break
			}
		}

		// If no user is found or the password check fails
		if foundUser == nil {
			log.Printf("User not found for email: %s", email)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Log the stored password hash and the provided password for comparison
		log.Printf("Stored password hash for user %s: %s", foundUser.Email, foundUser.Password)
		log.Printf("Provided password for user %s: %s", foundUser.Email, password)

		if err := foundUser.CheckPassword(password); err != nil {
			log.Printf("Invalid password for user %s: %v", foundUser.Email, err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Pass the request to the next handler if authentication succeeds
		log.Printf("User %s authenticated successfully", foundUser.Email)
		next.ServeHTTP(w, r)
	})
}
