package handlers

import (
	"DocuDefense/src/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

var users = make(map[string]*models.User) // Store pointers to users

// Secret key for signing the JWT
var jwtKey = []byte("your_secret_key")

// Get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var userList []*models.User
	for _, user := range users {
		userList = append(userList, user)
	}
	log.Printf("Returning list of users: %v", userList) // Log the list of users
	json.NewEncoder(w).Encode(userList)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// Decode the request body into the user object
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding user data: %v", err)
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	// Check if a password was provided
	if user.Password == "" {
		log.Printf("Password is missing for user %v", user)
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Hash the password before saving it
	if err := user.HashPassword(user.Password); err != nil {
		log.Printf("Error hashing password for user %v: %v", user, err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Initialize the FileNames slice
	user.FileNames = []string{}

	// Save the user with the hashed password to the users map
	users[user.ID] = &user
	log.Printf("User created with hashed password: %v", users[user.ID].Password)

	// Create a response object (copy) to prevent exposing the password hash
	responseUser := user
	responseUser.Password = "" // Clear the password only for the response

	// Send the created user as a response (without the password)
	json.NewEncoder(w).Encode(responseUser)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	existingUser, ok := users[params["id"]]
	if !ok {
		log.Printf("User not found: %s", params["id"]) // Log user not found
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Decode the updated user details
	var updatedUser models.User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		log.Printf("Error decoding user update data: %v", err) // Log decoding error
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	// Update fields of the existing user while keeping its pointer intact
	existingUser.FirstName = updatedUser.FirstName
	existingUser.Surname = updatedUser.Surname
	existingUser.Email = updatedUser.Email
	existingUser.Birthdate = updatedUser.Birthdate

	if updatedUser.Password != "" {
		// Rehash password if it's being updated
		err := existingUser.HashPassword(updatedUser.Password)
		if err != nil {
			log.Printf("Error hashing updated password for user %s: %v", existingUser.Email, err) // Log hashing error
			http.Error(w, "Error updating password", http.StatusInternalServerError)
			return
		}
	}

	log.Printf("User updated: %v", existingUser) // Log user update
	// Respond with the updated user
	json.NewEncoder(w).Encode(existingUser)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := users[params["id"]]; ok {
		log.Printf("Deleting user with ID: %s", params["id"]) // Log user deletion
		delete(users, params["id"])
		json.NewEncoder(w).Encode(users)
	} else {
		log.Printf("Attempt to delete non-existing user with ID: %s", params["id"]) // Log non-existing user
		http.Error(w, "User not found", http.StatusNotFound)
	}
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file and handler from the form
	file, handler, err := r.FormFile("contract")
	if err != nil {
		log.Printf("Error retrieving file from form: %v", err)
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure the file is a PDF
	if handler.Header.Get("Content-Type") != "application/pdf" {
		log.Printf("Invalid file type: %v, expected application/pdf", handler.Header.Get("Content-Type"))
		http.Error(w, "Only PDF files are allowed", http.StatusBadRequest)
		return
	}

	// Get the user ID from the URL
	params := mux.Vars(r)
	userID := params["id"]

	// Check if the user exists
	user, ok := users[userID]
	if !ok {
		log.Printf("User not found: %s", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Create a file path to save the file on disk
	filePath := "./uploads/" + handler.Filename

	// Ensure the uploads directory exists
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		log.Printf("Error creating uploads directory: %v", err)
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return
	}

	// Create the file on disk
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file on disk: %v", err)
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file data to the file on disk
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Error saving file to disk: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Add the new file name to the user's list of files
	user.FileNames = append(user.FileNames, handler.Filename)
	users[userID] = user

	// Log the successful upload
	log.Printf("Successfully uploaded file for user %s: %s", userID, handler.Filename)

	// Return the updated user information
	json.NewEncoder(w).Encode(user)
}

// GenerateJWT generates a new JWT token for an authenticated user
func GenerateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// LoginUser authenticates a user by email and password and returns a JWT
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		log.Printf("Error decoding login data: %v", err)
		http.Error(w, "Invalid login data", http.StatusBadRequest)
		return
	}

	// Find the user by email
	var foundUser *models.User
	for _, user := range users {
		if user.Email == loginData.Email {
			foundUser = user
			break
		}
	}

	if foundUser == nil {
		log.Printf("Login failed: user with email %s not found", loginData.Email)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Check if the password is correct
	log.Printf("Stored password hash for user: %s", foundUser.Password)    // Log the stored hash
	log.Printf("Provided password for comparison: %s", loginData.Password) // Log the provided password

	err = foundUser.CheckPassword(loginData.Password)
	if err != nil {
		log.Printf("Login failed: invalid password for user %s. Error: %v", foundUser.Email, err)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// If login is successful, generate JWT
	token, err := GenerateJWT(foundUser)
	if err != nil {
		log.Printf("Error generating token for user %s: %v", foundUser.Email, err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// If login is successful
	log.Printf("Login successful for user: %v", foundUser)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   token,
	})
}
