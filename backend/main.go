package main

import (
	"DocuDefense/backend/src/handlers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get MongoDB URI from the environment variable
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MongoDB URI not found in environment variables")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		log.Fatal("Ping failed:", err)
	}
	fmt.Println("Pinged MongoDB successfully!")

	// Pass the MongoDB client to handlers
	handlers.SetMongoClient(client)

	// Set up routes
	r := mux.NewRouter()
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	// Endpoint for fetching user data by email (e.g., for user ID lookup)
	r.HandleFunc("/users/email", handlers.GetUserByEmail).Methods("GET")

	// User-specific routes that require JWT authentication
	r.Handle("/users/{id}", handlers.JWTAuthMiddleware(http.HandlerFunc(handlers.UpdateUser))).Methods("PUT")
	r.Handle("/users/{id}", handlers.JWTAuthMiddleware(http.HandlerFunc(handlers.DeleteUser))).Methods("DELETE")
	r.Handle("/users/{id}/upload", handlers.JWTAuthMiddleware(http.HandlerFunc(handlers.UploadFile))).Methods("POST")
	r.Handle("/users/{id}/files", handlers.JWTAuthMiddleware(http.HandlerFunc(handlers.GetUserFiles))).Methods("GET")
	r.Handle("/users/{id}/files/{filename}/download", handlers.JWTAuthMiddleware(http.HandlerFunc(handlers.DownloadFile))).Methods("GET")
	r.Handle("/users/{id}/files/{filename}/delete", handlers.JWTAuthMiddleware(http.HandlerFunc(handlers.DeleteFile))).Methods("DELETE")
	r.HandleFunc("/api/users/search", handlers.SearchUsers).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Fatal(http.ListenAndServe(":8000", handler))
}
