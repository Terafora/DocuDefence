package handlers

import (
	"DocuDefense/src/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var users = make(map[string]models.User)

// Get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var userList []models.User
	for _, user := range users {
		userList = append(userList, user)
	}
	json.NewEncoder(w).Encode(userList)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	// Initialize the FileNames slice
	user.FileNames = []string{}
	users[user.ID] = user

	json.NewEncoder(w).Encode(user)
}

// Update an existing user by ID
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedUser models.User
	_ = json.NewDecoder(r.Body).Decode(&updatedUser)

	if _, ok := users[params["id"]]; ok {
		users[params["id"]] = updatedUser
		json.NewEncoder(w).Encode(updatedUser)
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}

// Delete a user by ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := users[params["id"]]; ok {
		delete(users, params["id"])
		json.NewEncoder(w).Encode(users)
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse the form to get the file data
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		log.Printf("Error parsing form: %v", err) // Added logging for form parsing error
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file and handler from the form
	file, handler, err := r.FormFile("contract")
	if err != nil {
		log.Printf("Error retrieving file from form: %v", err) // Added logging for file retrieval error
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure the file is a PDF
	if handler.Header.Get("Content-Type") != "application/pdf" {
		log.Printf("Invalid file type: %v, expected application/pdf", handler.Header.Get("Content-Type")) // Logging invalid file type
		http.Error(w, "Only PDF files are allowed", http.StatusBadRequest)
		return
	}

	// Get the user ID from the URL
	params := mux.Vars(r)
	userID := params["id"]

	// Check if the user exists
	user, ok := users[userID]
	if !ok {
		log.Printf("User not found: %s", userID) // Logging if user not found
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Create a file path to save the file on disk
	filePath := "./uploads/" + handler.Filename

	// Ensure the uploads directory exists
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		log.Printf("Error creating uploads directory: %v", err) // Logging directory creation error
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return
	}

	// Create the file on disk
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file on disk: %v", err) // Logging file creation error
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file data to the file on disk
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Error saving file to disk: %v", err) // Logging file saving error
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Add the new file name to the user's list of files
	user.FileNames = append(user.FileNames, handler.Filename)
	users[userID] = user // Save the updated user

	// Log the successful upload
	log.Printf("Successfully uploaded file for user %s: %s", userID, handler.Filename)

	// Return the updated user information
	json.NewEncoder(w).Encode(user)
}
