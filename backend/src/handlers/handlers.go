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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define an interface for MongoDB collection methods used in handlers (testing purposes)
type UserCollection interface {
	Find(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)
}

// MongoDB collections
var usersCollection *mongo.Collection

// SetMongoClient initializes the MongoDB client and sets the users collection
func SetMongoClient(client *mongo.Client) {
	usersCollection = client.Database("docudefense").Collection("users")
}

// Allows setting a mock collection for testing
func SetUsersCollection(collection *mongo.Collection) {
	usersCollection = collection
}

// Secret key for signing the JWT
var jwtKey = []byte("your_secret_key") // Ensure this is the same as used in JWTAuthMiddleware

// Get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all users in the collection
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

	log.Printf("Returning list of users: %v", userList)
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

	// Ensure a unique ObjectID is assigned if ID is zero
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	// **Hash the password before saving it**
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

	// Clear the password before responding
	user.Password = ""
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	// Decode the updated user details
	var updatedUser models.User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		log.Printf("Error decoding user update data: %v", err)
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	// Hash the password if provided
	if updatedUser.Password != "" {
		if err := updatedUser.HashPassword(updatedUser.Password); err != nil {
			log.Printf("Error hashing updated password: %v", err)
			http.Error(w, "Error updating password", http.StatusInternalServerError)
			return
		}
	}

	// Update the user in the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIDObj, err := primitive.ObjectIDFromHex(userID) // Ensure this is a valid ObjectID
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": userIDObj}
	update := bson.M{"$set": updatedUser}

	result, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating user %s: %v", userID, err)
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		log.Printf("No user found with ID: %s", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	log.Printf("User updated: %v", updatedUser)
	json.NewEncoder(w).Encode(updatedUser)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	// Delete the user from the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	result, err := usersCollection.DeleteOne(ctx, bson.M{"_id": userIDObj})
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

	// Update the user's file list in the database
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var foundUser models.User
	err = usersCollection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&foundUser)
	if err != nil {
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
	token, err := GenerateJWT(&foundUser)
	if err != nil {
		log.Printf("Error generating token for user %s: %v", foundUser.Email, err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	log.Printf("Login successful for user: %v", foundUser)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   token,
	})
}
