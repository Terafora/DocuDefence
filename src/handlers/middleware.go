package handlers

import (
	"DocuDefense/src/models"
	"net/http"
)

func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		// Find the user by email
		var foundUser *models.User
		for _, user := range users {
			if user.Email == email {
				foundUser = &user
				break
			}
		}

		if foundUser == nil || foundUser.CheckPassword(password) != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Pass the request to the next handler if authentication succeeds
		next.ServeHTTP(w, r)
	})
}
