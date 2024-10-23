package handlers

import (
	"DocuDefense/src/models"
	"encoding/json"
	"io"
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

// Create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
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

// Upload a file and associate it with a user
func UploadFile(w http.ResponseWriter, r *http.Request) {

	// Parse the form to get the file data
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file and handler from the form
	file, handler, err := r.FormFile("contract")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get the user ID from the URL
	params := mux.Vars(r)
	userID := params["id"]

	// Check if the user exists
	user, ok := users[userID]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Create a file path to save the file on disk
	filePath := "./uploads/" + handler.Filename

	// Ensure the uploads directory exists
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return
	}

	// Create the file on disk
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file data to the file on disk
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	user.FileName = handler.Filename
	users[userID] = user

	json.NewEncoder(w).Encode(user)
}
