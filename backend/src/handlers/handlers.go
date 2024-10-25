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

// Create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding user data: %v", err)
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

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

	user.Password = ""
	json.NewEncoder(w).Encode(user)
}

// Update user information (only if the requester is the account owner)
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	// Retrieve the user claims from the JWT token
	claims := r.Context().Value("userClaims").(*Claims)
	requesterEmail := claims.Email

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Find the user by ID
	var targetUser models.User
	err = usersCollection.FindOne(context.Background(), bson.M{"_id": userIDObj}).Decode(&targetUser)
	if err != nil {
		log.Printf("User not found for ID %s: %v", userID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Ensure the user can only update their own account
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

	// Update the user
	filter := bson.M{"_id": userIDObj}
	update := bson.M{"$set": updatedUser}
	result, err := usersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil || result.MatchedCount == 0 {
		log.Printf("Error updating user %s: %v", userID, err)
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	log.Printf("User updated: %v", updatedUser)
	json.NewEncoder(w).Encode(updatedUser)
}

// Delete user (only if the requester is the account owner)
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	// Retrieve the user claims from the JWT token
	claims := r.Context().Value("userClaims").(*Claims)
	requesterEmail := claims.Email

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Find the user by ID
	var targetUser models.User
	err = usersCollection.FindOne(context.Background(), bson.M{"_id": userIDObj}).Decode(&targetUser)
	if err != nil {
		log.Printf("User not found for ID %s: %v", userID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Ensure the user can only delete their own account
	if requesterEmail != targetUser.Email {
		log.Printf("Unauthorized delete attempt by %s on account %s", requesterEmail, targetUser.Email)
		http.Error(w, "You are not authorized to delete this account", http.StatusForbidden)
		return
	}

	// Delete the user
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

// Search for users by first name or surname
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	filter := bson.M{}

	if firstName, ok := queryParams["first_name"]; ok {
		filter["first_name"] = bson.M{"$regex": firstName[0], "$options": "i"}
	}
	if surname, ok := queryParams["surname"]; ok {
		filter["surname"] = bson.M{"$regex": surname[0], "$options": "i"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := usersCollection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		http.Error(w, "Error decoding users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// Upload PDF file for a user
func UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
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

	if handler.Header.Get("Content-Type") != "application/pdf" {
		log.Printf("Invalid file type: %v, expected application/pdf", handler.Header.Get("Content-Type"))
		http.Error(w, "Only PDF files are allowed", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	userID := params["id"]
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$push": bson.M{"filenames": handler.Filename}}
	result, err := usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil || result.MatchedCount == 0 {
		log.Printf("Error updating user's file list for user %s: %v", userID, err)
		http.Error(w, "Error updating file list", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully uploaded file for user %s: %s", userID, handler.Filename)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File uploaded",
	})
}

// Generate JWT token for authenticated user
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

// Login user and generate JWT
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

	// Generate JWT
	token, err := GenerateJWT(&foundUser)
	if err != nil {
		log.Printf("Error generating token for user %s: %v", foundUser.Email, err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Log the token and response for debugging
	log.Printf("Generated token for user %s: %s", foundUser.Email, token)
	log.Printf("Login successful for user: %v", foundUser)

	// Send response with the token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   token,
	})
}
