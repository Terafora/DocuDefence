package handlers

import (
	"DocuDefense/backend/src/models"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB collections
var usersCollection *mongo.Collection

// SetMongoClient initializes the MongoDB client and sets the users collection
func SetMongoClient(client *mongo.Client) {
	usersCollection = client.Database("docudefense").Collection("users")
}

// JWT secret key
var jwtKey = []byte("your_secret_key")

// Get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := usersCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error retrieving users: %v", err)
		http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var userList []models.User
	if err = cursor.All(ctx, &userList); err != nil {
		log.Printf("Error decoding users: %v", err)
		http.Error(w, "Error decoding users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(userList)
}

// Updated CreateUser function
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding user data: %v", err)
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	// Ensure user has a unique ID and initialize file_names as an empty array
	user.ID = primitive.NewObjectID()
	user.FileNames = []string{} // Initialize file_names as an empty array

	if err := user.HashPassword(user.Password); err != nil {
		log.Printf("Error hashing password for user: %v", err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = usersCollection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	user.Password = "" // Remove password before returning response
	json.NewEncoder(w).Encode(user)
}

// GetUserByEmail retrieves a user document by email and returns their ID
func GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	var user models.User
	err := usersCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		log.Printf("User not found for email %s: %v", email, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := map[string]string{"id": user.ID.Hex()}
	json.NewEncoder(w).Encode(response)
}

// GetUserFiles retrieves files for a specific user by ID
func GetUserFiles(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var user models.User
	err = usersCollection.FindOne(context.Background(), bson.M{"_id": userIDObj}).Decode(&user)
	if err != nil {
		log.Printf("User not found for ID %s: %v", userID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string][]string{"file_names": user.FileNames})
}

// UpdateUser updates the user information if the requester is the account owner
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	claims, ok := r.Context().Value("userClaims").(*Claims)
	if !ok || claims == nil {
		log.Println("Unauthorized access: Unable to retrieve user claims")
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
	requesterEmail := claims.Email

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var targetUser models.User
	err = usersCollection.FindOne(context.Background(), bson.M{"_id": userIDObj}).Decode(&targetUser)
	if err != nil {
		log.Printf("User not found for ID %s: %v", userID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if requesterEmail != targetUser.Email {
		log.Printf("Unauthorized update attempt by %s on account %s", requesterEmail, targetUser.Email)
		http.Error(w, "You are not authorized to update this account", http.StatusForbidden)
		return
	}

	var updatedUser models.User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		log.Printf("Error decoding user update data: %v", err)
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	if updatedUser.Password != "" {
		if err := updatedUser.HashPassword(updatedUser.Password); err != nil {
			log.Printf("Error hashing updated password: %v", err)
			http.Error(w, "Error updating password", http.StatusInternalServerError)
			return
		}
	}

	update := bson.M{
		"$set": bson.M{
			"first_name": updatedUser.FirstName,
			"surname":    updatedUser.Surname,
			"email":      updatedUser.Email,
			"password":   updatedUser.Password,
		},
	}

	filter := bson.M{"_id": userIDObj}
	result, err := usersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil || result.MatchedCount == 0 {
		log.Printf("Error updating user %s: %v", userID, err)
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	log.Printf("User updated: %v", updatedUser)
	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser deletes the user if they are the account owner
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	claims := r.Context().Value("userClaims").(*Claims)
	requesterEmail := claims.Email

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var targetUser models.User
	err = usersCollection.FindOne(context.Background(), bson.M{"_id": userIDObj}).Decode(&targetUser)
	if err != nil {
		log.Printf("User not found for ID %s: %v", userID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if requesterEmail != targetUser.Email {
		log.Printf("Unauthorized delete attempt by %s on account %s", requesterEmail, targetUser.Email)
		http.Error(w, "You are not authorized to delete this account", http.StatusForbidden)
		return
	}

	result, err := usersCollection.DeleteOne(context.Background(), bson.M{"_id": userIDObj})
	if err != nil || result.DeletedCount == 0 {
		log.Printf("Error deleting user %s: %v", userID, err)
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	log.Printf("User deleted with ID: %s", userID)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted",
	})
}

// UploadFile allows a user to upload a PDF file
func UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB size limit
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("contract")
	if err != nil {
		log.Printf("Error retrieving file from form: %v", err)
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if the file is a PDF
	if handler.Header.Get("Content-Type") != "application/pdf" {
		log.Printf("Invalid file type: %v, expected application/pdf", handler.Header.Get("Content-Type"))
		http.Error(w, "Only PDF files are allowed", http.StatusBadRequest)
		return
	}

	// Retrieve user ID from URL
	params := mux.Vars(r)
	userID := params["id"]
	log.Printf("Received upload request for user ID: %s", userID)

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Save file to disk
	filePath := "./uploads/" + handler.Filename
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		log.Printf("Error creating uploads directory: %v", err)
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file on disk: %v", err)
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Error saving file to disk: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Retry logic for updating the file list in MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$push": bson.M{"file_names": handler.Filename}}
	attempts := 3
	for attempts > 0 {
		result, err := usersCollection.UpdateOne(ctx, bson.M{"_id": userIDObj}, update)
		if err == nil && result.MatchedCount > 0 {
			// Success
			log.Printf("Successfully uploaded file for user %s: %s", userID, handler.Filename)
			json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded"})
			return
		} else if err != nil {
			log.Printf("Error updating user's file list, retrying... Attempts left: %d, Error: %v", attempts-1, err)
			time.Sleep(500 * time.Millisecond) // Short delay before retrying
			attempts--
		} else {
			http.Error(w, "Error updating file list", http.StatusInternalServerError)
			return
		}
	}

	// Final failure response if retries are exhausted
	log.Printf("Failed to update user's file list after retries for user %s", userID)
	http.Error(w, "Failed to update file list after retries", http.StatusInternalServerError)
}

// Search for users by first name or surname
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	filter := bson.M{}

	// Check if there is a first_name or surname in the query parameters
	if firstName, ok := queryParams["first_name"]; ok {
		filter["first_name"] = bson.M{"$regex": firstName[0], "$options": "i"}
	}
	if surname, ok := queryParams["surname"]; ok {
		filter["surname"] = bson.M{"$regex": surname[0], "$options": "i"}
	}

	// Search users in MongoDB with the filter
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := usersCollection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Decode the result into a slice of User objects
	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		http.Error(w, "Error decoding users", http.StatusInternalServerError)
		return
	}

	// Return the matched users as JSON
	json.NewEncoder(w).Encode(users)
}

// GenerateJWT generates a JWT token for authenticated users
func GenerateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// LoginUser logs in a user and generates a JWT
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var foundUser models.User
	err = usersCollection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&foundUser)
	if err != nil {
		log.Printf("Login failed: user with email %s not found", loginData.Email)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = foundUser.CheckPassword(loginData.Password)
	if err != nil {
		log.Printf("Login failed: invalid password for user %s. Error: %v", foundUser.Email, err)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(&foundUser)
	if err != nil {
		log.Printf("Error generating token for user %s: %v", foundUser.Email, err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	log.Printf("Generated token for user %s: %s", foundUser.Email, token)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   token,
	})
}
