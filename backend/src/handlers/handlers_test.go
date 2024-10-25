package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var createdUserID primitive.ObjectID
var jwtToken string

type UserResponse struct {
	ID        string   `json:"id"`
	FirstName string   `json:"first_name"`
	Surname   string   `json:"surname"`
	Email     string   `json:"email"`
	Birthdate string   `json:"birthdate"`
	Password  string   `json:"password"`
	FileNames []string `json:"file_names"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func TestMain(m *testing.M) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB for testing: %v", err)
	}
	testDB := client.Database("docudefense_test")
	usersCollection = testDB.Collection("users")

	code := m.Run()

	err = testDB.Drop(context.TODO())
	if err != nil {
		log.Fatalf("Failed to drop test database: %v", err)
	}

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatalf("Failed to disconnect MongoDB client: %v", err)
	}

	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	newUser := UserResponse{
		FirstName: "Alice",
		Surname:   "Smith",
		Email:     "alice.smith@example.com",
		Birthdate: "1995-05-05",
		Password:  "alicePassword789",
		FileNames: []string{},
	}

	jsonData, err := json.Marshal(newUser)
	if err != nil {
		t.Fatalf("Failed to marshal new user: %v", err)
	}

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(CreateUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var createdUser UserResponse
	err = json.Unmarshal(rr.Body.Bytes(), &createdUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	createdUserID, err = primitive.ObjectIDFromHex(createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to convert ID to ObjectID: %v", err)
	}
}

func TestLoginUser(t *testing.T) {
	loginData := map[string]string{
		"email":    "alice.smith@example.com",
		"password": "alicePassword789",
	}
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		t.Fatalf("Failed to marshal login data: %v", err)
	}

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(LoginUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var loginResponse LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &loginResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	jwtToken = loginResponse.Token
}

func TestUpdateUser(t *testing.T) {
	updatedUser := UserResponse{
		FirstName: "Alicia",
		Surname:   "Smithson",
		Email:     "alice.smith@example.com",
		Birthdate: "1995-05-05",
		Password:  "alicePassword789",
		FileNames: []string{},
	}
	jsonData, err := json.Marshal(updatedUser)
	if err != nil {
		t.Fatalf("Failed to marshal updated user: %v", err)
	}

	req, err := http.NewRequest("PUT", "/users/"+createdUserID.Hex(), bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(UpdateUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var updatedResponse UserResponse
	err = json.Unmarshal(rr.Body.Bytes(), &updatedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if updatedResponse.FirstName != "Alicia" || updatedResponse.Surname != "Smithson" {
		t.Errorf("Handler returned unexpected body: got %+v want %+v", updatedResponse, updatedUser)
	}
}

func TestDeleteUser(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/users/"+createdUserID.Hex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(DeleteUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var deletedUser UserResponse
	err = usersCollection.FindOne(context.TODO(), bson.M{"_id": createdUserID}).Decode(&deletedUser)
	if err == nil {
		t.Fatal("Expected user to be deleted, but found one in database")
	}
}
