package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"DocuDefense/src/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

	// Get the MongoDB URI from the environment variable
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MongoDB URI not found in environment variables")
	}

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	// Close the client connection when the function is done
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		log.Fatal("Ping failed:", err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	// Pass the MongoDB client to handlers
	handlers.SetMongoClient(client)

	// Set up routes
	r := mux.NewRouter()

	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	// Protect routes with JWTAuthMiddleware
	r.Handle("/users/{id}", handlers.JWTAuthMiddleware(http.HandlerFunc(handlers.UpdateUser))).Methods("PUT")
	r.Handle("/users/{id}", handlers.JWTAuthMiddleware(http.HandlerFunc(handlers.DeleteUser))).Methods("DELETE")
	r.Handle("/users/{id}/upload", handlers.JWTAuthMiddleware(http.HandlerFunc(handlers.UploadFile))).Methods("POST")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", r))
}
