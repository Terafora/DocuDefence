package handlers

import (
	"DocuDefense/src/models"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test GetUsers
func TestGetUsers(t *testing.T) {
	// Create a new request
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUsers)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that the response body is a valid JSON array
	var users []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &users)
	if err != nil {
		t.Errorf("handler returned invalid JSON array: %v", err)
	}

	// Check if the array is empty
	if len(users) != 0 {
		t.Errorf("handler returned non-empty array: got %v users, want 0", len(users))
	}
}

// Test CreateUser
func TestCreateUser(t *testing.T) {
	// Create a new POST request with user data
	user := map[string]string{
		"id":         "1",
		"first_name": "Charlotte",
		"surname":    "Stone",
		"email":      "charlotte@example.com",
		"birthdate":  "1992-01-08",
	}
	userJson, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(userJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUser)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the response contains an empty FileNames slice
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	if response["file_names"] == nil {
		t.Errorf("Expected 'file_names' field, got nil")
	}
}

// / Test UploadFile using an actual PDF file
func TestUploadFile(t *testing.T) {
	// First, create a user for the test
	users["1"] = models.User{
		ID:        "1",
		FirstName: "Charlotte",
		Surname:   "Stone",
		Email:     "charlotte@example.com",
		Birthdate: "1992-01-08",
		FileNames: []string{},
	}

	// Prepare a new POST request for file upload
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Instead of loading a real file, simulate a PDF file with valid PDF header
	part, err := writer.CreateFormFile("contract", "testfile.pdf")
	if err != nil {
		t.Fatal(err)
	}

	// Simulate a PDF by starting the file with the %PDF header
	_, err = part.Write([]byte("%PDF-1.4\nThis is a test PDF file\n"))
	if err != nil {
		t.Fatal(err)
	}

	// Close the writer to finalize the form
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Create the request with the correct headers for multipart/form-data
	req, err := http.NewRequest("POST", "/users/1/upload", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UploadFile)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check if the status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the file was added to the user's file list
	user := users["1"]
	if len(user.FileNames) != 1 {
		t.Errorf("Expected 1 file name, got %v", len(user.FileNames))
	}
	if user.FileNames[0] != "testfile.pdf" {
		t.Errorf("Expected file name 'testfile.pdf', got %v", user.FileNames[0])
	}
}

// Test UploadMultipleFiles using actual PDF files
func TestUploadMultipleFiles(t *testing.T) {
	// First, create a user for the test
	users["1"] = models.User{
		ID:        "1",
		FirstName: "Charlotte",
		Surname:   "Stone",
		Email:     "charlotte@example.com",
		Birthdate: "1992-01-08",
		FileNames: []string{},
	}

	// Function to upload a file
	uploadFile := func(filename string, content string) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		// Simulate the PDF content with valid header
		part, err := writer.CreateFormFile("contract", filename)
		if err != nil {
			t.Fatal(err)
		}
		_, err = part.Write([]byte(content))
		if err != nil {
			t.Fatal(err)
		}

		err = writer.Close()
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/users/1/upload", body)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(UploadFile)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	}

	// Upload two simulated PDF files
	uploadFile("file1.pdf", "%PDF-1.4\nThis is a test PDF file\n")
	uploadFile("file2.pdf", "%PDF-1.4\nAnother test PDF file\n")

	// Check if both files were added to the user's file list
	user := users["1"]
	if len(user.FileNames) != 2 {
		t.Errorf("Expected 2 file names, got %v", len(user.FileNames))
	}
	if user.FileNames[0] != "file1.pdf" || user.FileNames[1] != "file2.pdf" {
		t.Errorf("Expected 'file1.pdf' and 'file2.pdf', got %v", user.FileNames)
	}
}
