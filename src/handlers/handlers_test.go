package handlers

import (
	"DocuDefense/src/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test GetUsers
func TestGetUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUsers)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &users)
	if err != nil {
		t.Errorf("handler returned invalid JSON array: %v", err)
	}
}

// Test CreateUser
func TestCreateUser(t *testing.T) {
	user := map[string]string{
		"id":         "1",
		"first_name": "Charlotte",
		"surname":    "Stone",
		"email":      "charlotte@example.com",
		"birthdate":  "1992-01-08",
		"password":   "YourSuperSecretPassword",
	}

	userJson, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(userJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response["password"] != "" {
		t.Errorf("Expected password to be empty in response, got %v", response["password"])
	}
}

// Test LoginUser
func TestLoginUser(t *testing.T) {
	users["1"] = &models.User{
		ID:        "1",
		FirstName: "Charlotte",
		Surname:   "Stone",
		Email:     "charlotte@example.com",
		Birthdate: "1992-01-08",
	}
	users["1"].HashPassword("YourSuperSecretPassword")

	loginData := map[string]string{
		"email":    "charlotte@example.com",
		"password": "YourSuperSecretPassword",
	}
	loginJson, _ := json.Marshal(loginData)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response["message"] != "Login successful" {
		t.Errorf("Expected login success message, got %v", response["message"])
	}
}

// Test UpdateUser
func TestUpdateUser(t *testing.T) {
	// Create a user before updating
	users["1"] = &models.User{
		ID:        "1",
		FirstName: "Charlotte",
		Surname:   "Stone",
		Email:     "charlotte@example.com",
		Birthdate: "1992-01-08",
	}
	users["1"].HashPassword("YourSuperSecretPassword")

	updatedUser := map[string]string{
		"first_name": "UpdatedName",
		"surname":    "Stone",
		"email":      "updated@example.com",
		"password":   "NewSecretPassword",
	}
	updatedUserJson, _ := json.Marshal(updatedUser)

	// Correct the request URL to include the ID
	req, err := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(updatedUserJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that the user's details have been updated correctly
	user := users["1"]
	if user.FirstName != "UpdatedName" || user.Email != "updated@example.com" {
		t.Errorf("User was not updated correctly")
	}
}

// Test DeleteUser
func TestDeleteUser(t *testing.T) {
	// Ensure the user exists before deleting
	users["1"] = &models.User{
		ID:        "1",
		FirstName: "Charlotte",
		Surname:   "Stone",
		Email:     "charlotte@example.com",
		Birthdate: "1992-01-08",
	}

	// Correct the request URL to include the ID
	req, err := http.NewRequest("DELETE", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that the user has been deleted
	if _, ok := users["1"]; ok {
		t.Errorf("User was not deleted")
	}
}

// Test Basic Authentication Middleware
func TestBasicAuthMiddleware(t *testing.T) {
	users["1"] = &models.User{
		ID:        "1",
		FirstName: "Charlotte",
		Surname:   "Stone",
		Email:     "charlotte@example.com",
	}
	users["1"].HashPassword("YourSuperSecretPassword")

	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("charlotte@example.com", "YourSuperSecretPassword")

	rr := httptest.NewRecorder()
	handler := BasicAuthMiddleware(http.HandlerFunc(GetUsers))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
