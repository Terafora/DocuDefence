package handlers

import (
	"DocuDefense/backend/src/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB collections
var usersCollection *mongo.Collection
var documentsCollection *mongo.Collection

// SetMongoClient initializes the MongoDB client and sets the users collection
func SetMongoClient(client *mongo.Client) {
	usersCollection = client.Database("docudefense").Collection("users")
	documentsCollection = client.Database("docudefense").Collection("documents")
}

// JWT secret key
var jwtKey = []byte("your_secret_key")

// Get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Default limit per page
	}

	// Calculate skip and limit
	skip := (page - 1) * limit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find users with pagination
	cursor, err := usersCollection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
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

	// Ensure user has a unique ID
	user.ID = primitive.NewObjectID()

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

	var documents []models.Document
	cursor, err := documentsCollection.Find(context.Background(), bson.M{"user_id": userIDObj})
	if err != nil {
		log.Printf("Error retrieving documents for user %s: %v", userID, err)
		http.Error(w, "Error retrieving documents", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &documents); err != nil {
		log.Printf("Error decoding documents: %v", err)
		http.Error(w, "Error decoding documents", http.StatusInternalServerError)
		return
	}

	// Log each document to confirm structure
	for _, doc := range documents {
		log.Printf("Document fetched: Filename: %s, Version: %d, UploadDate: %v", doc.Filename, doc.Version, doc.UploadDate)
	}

	json.NewEncoder(w).Encode(documents)
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

// UploadFile allows a user to upload a PDF file with version control
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

	// Get user ID from URL parameters
	params := mux.Vars(r)
	userID := params["id"]
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Standardize filename by replacing spaces with underscores
	filename := strings.ReplaceAll(handler.Filename, " ", "_")
	filePath := "./uploads/" + filename

	// Save file to disk
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

	// Determine the new version number by finding the latest version of this file
	var lastDoc models.Document
	cursor, err := documentsCollection.Find(context.Background(), bson.M{
		"user_id":  userIDObj,
		"filename": filename,
	}, options.Find().SetSort(bson.D{{"version", -1}}).SetLimit(1))
	if err != nil {
		log.Printf("Error finding latest version for file %s: %v", filename, err)
		http.Error(w, "Error finding latest file version", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	newVersion := 1 // Default to version 1 if no previous version is found
	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&lastDoc); err != nil {
			log.Printf("Error decoding document: %v", err)
			http.Error(w, "Error decoding document", http.StatusInternalServerError)
			return
		}
		newVersion = lastDoc.Version + 1 // Increment version based on the latest document found
	}

	// Store the standardized filename and version in MongoDB
	newDoc := models.Document{
		ID:         primitive.NewObjectID(),
		UserID:     userIDObj,
		Filename:   filename,
		Version:    newVersion,
		UploadDate: time.Now(),
	}
	_, err = documentsCollection.InsertOne(context.Background(), newDoc)
	if err != nil {
		log.Printf("Error creating document entry: %v", err)
		http.Error(w, "Error creating document entry", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully uploaded file: %s, version: %d", filename, newVersion)
	json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded", "filename": filename, "version": fmt.Sprint(newVersion)})
}

// DownloadFile allows a user to download a file by filename
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	encodedFilename := params["filename"]
	filename, err := url.QueryUnescape(encodedFilename)
	if err != nil {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := "./uploads/" + filename
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, file)
}

// DeleteFile allows a user to delete a file by filename
func DeleteFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	encodedFilename := params["filename"]
	filename, err := url.QueryUnescape(encodedFilename)
	if err != nil {
		log.Printf("Error decoding filename: %v", err)
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	userID := params["id"]
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Retrieve all versions of the file by filename and user ID
	cursor, err := documentsCollection.Find(context.Background(), bson.M{
		"user_id":  userIDObj,
		"filename": filename,
	})
	if err != nil {
		log.Printf("Error retrieving file versions for deletion: %v", err)
		http.Error(w, "Error retrieving file versions", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Track errors if file versions can't be deleted from the file system
	var deletionErrors []string
	for cursor.Next(context.Background()) {
		var doc models.Document
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("Error decoding document during deletion: %v", err)
			continue
		}

		// Remove the file from disk
		filePath := "./uploads/" + doc.Filename
		if err := os.Remove(filePath); err != nil {
			log.Printf("Error deleting file from disk: %v", err)
			deletionErrors = append(deletionErrors, fmt.Sprintf("Error deleting version %d", doc.Version))
		}

		// Remove the document from the database
		_, err := documentsCollection.DeleteOne(context.Background(), bson.M{"_id": doc.ID})
		if err != nil {
			log.Printf("Error deleting file from database: %v", err)
			deletionErrors = append(deletionErrors, fmt.Sprintf("Error removing version %d from database", doc.Version))
		}
	}

	if len(deletionErrors) > 0 {
		http.Error(w, strings.Join(deletionErrors, "; "), http.StatusInternalServerError)
	} else {
		log.Printf("Successfully deleted all versions of file: %s", filename)
		json.NewEncoder(w).Encode(map[string]string{"message": "File deleted successfully"})
	}
}

// Search for users by first name or surname
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	// Filter setup
	queryParams := r.URL.Query()
	filter := bson.M{}
	if term, ok := queryParams["term"]; ok && term[0] != "" {
		filter["$or"] = []bson.M{
			{"first_name": bson.M{"$regex": term[0], "$options": "i"}},
			{"surname": bson.M{"$regex": term[0], "$options": "i"}},
		}
	} else {
		http.Error(w, `{"error": "Search term is required"}`, http.StatusBadRequest)
		return
	}

	// Search with pagination
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := usersCollection.Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		http.Error(w, `{"error": "Error fetching users"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	users := []models.User{}
	if err := cursor.All(ctx, &users); err != nil {
		http.Error(w, `{"error": "Error decoding users"}`, http.StatusInternalServerError)
		return
	}

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
