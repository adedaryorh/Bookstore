package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/adedaryorh/bookstore-app/pkg/models"
	"github.com/julienschmidt/httprouter"
)

// Setup test database
func TestMain(m *testing.M) {
	// Set test database environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "bookstore_user")
	os.Setenv("DB_PASSWORD", "bookstore_pass")
	os.Setenv("DB_NAME", "bookstore_test")

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestGetAllBooks(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	GetAllBooks(rr, req, nil)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if response is valid JSON
	var books []models.Book
	if err := json.Unmarshal(rr.Body.Bytes(), &books); err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}

	// Check Content-Type
	expected := "application/json"
	if ctype := rr.Header().Get("Content-Type"); ctype != expected {
		t.Errorf("handler returned wrong content type: got %v want %v",
			ctype, expected)
	}
}

func TestGetBooks(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/book", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	GetBooks(rr, req, nil)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if response is valid JSON
	var books []models.Book
	if err := json.Unmarshal(rr.Body.Bytes(), &books); err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}
}

func TestCreateBook(t *testing.T) {
	// Create test book data
	testBook := models.Book{
		Title:           "Test Book",
		Author:          "Test Author",
		ISBN:            "1234567890123",
		PublicationYear: "2023",
		Genre:           stringPtr("Fiction"),
		Price:           float64Ptr(29.99),
	}

	// Convert to JSON
	jsonData, err := json.Marshal(testBook)
	if err != nil {
		t.Fatal(err)
	}

	// Create request
	req, err := http.NewRequest("POST", "/book", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler
	CreateBook(rr, req, nil)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check if response contains the created book
	var createdBook models.Book
	if err := json.Unmarshal(rr.Body.Bytes(), &createdBook); err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}

	// Verify book data
	if createdBook.Title != testBook.Title {
		t.Errorf("Created book title mismatch: got %v want %v",
			createdBook.Title, testBook.Title)
	}

	if createdBook.Author != testBook.Author {
		t.Errorf("Created book author mismatch: got %v want %v",
			createdBook.Author, testBook.Author)
	}
}

func TestCreateBookInvalidJSON(t *testing.T) {
	// Create request with invalid JSON
	req, err := http.NewRequest("POST", "/book", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler
	CreateBook(rr, req, nil)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check error response
	var errorResponse map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
		t.Errorf("Error response is not valid JSON: %v", err)
	}

	if errorResponse["error"] != "Invalid JSON format" {
		t.Errorf("Wrong error message: got %v want %v",
			errorResponse["error"], "Invalid JSON format")
	}
}

func TestGetBookByID(t *testing.T) {
	// Setup router with params
	router := httprouter.New()
	router.GET("/book/:bookId", GetBookByID)

	// Create request
	req, err := http.NewRequest("GET", "/book/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler through router
	router.ServeHTTP(rr, req)

	// Check the status code (could be 200 if book exists or 404 if not)
	if status := rr.Code; status != http.StatusOK && status != http.StatusNotFound {
		t.Errorf("handler returned unexpected status code: got %v want %v or %v",
			status, http.StatusOK, http.StatusNotFound)
	}

	// If book found, verify JSON structure
	if rr.Code == http.StatusOK {
		var book models.Book
		if err := json.Unmarshal(rr.Body.Bytes(), &book); err != nil {
			t.Errorf("Response is not valid JSON: %v", err)
		}
	}
}

func TestGetBookByIDInvalidID(t *testing.T) {
	// Setup router with params
	router := httprouter.New()
	router.GET("/book/:bookId", GetBookByID)

	// Create request with invalid ID
	req, err := http.NewRequest("GET", "/book/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler through router
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check error response
	var errorResponse map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
		t.Errorf("Error response is not valid JSON: %v", err)
	}

	if errorResponse["error"] != "Invalid book ID" {
		t.Errorf("Wrong error message: got %v want %v",
			errorResponse["error"], "Invalid book ID")
	}
}

func TestUpdateBook(t *testing.T) {
	// First, create a book to update
	testBook := models.Book{
		Title:           "Original Title",
		Author:          "Original Author",
		ISBN:            "1234567890124",
		PublicationYear: "2023",
	}

	// Create the book (this assumes CreateBook works)
	jsonData, _ := json.Marshal(testBook)
	createReq, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	createRR := httptest.NewRecorder()
	CreateBook(createRR, createReq, nil)

	if createRR.Code != http.StatusCreated {
		t.Skip("Skipping update test because create failed")
	}

	// Get the created book ID
	var createdBook models.Book
	json.Unmarshal(createRR.Body.Bytes(), &createdBook)

	// Prepare update data
	updateBook := models.Book{
		Title:           "Updated Title",
		Author:          "Updated Author",
		ISBN:            "1234567890124",
		PublicationYear: "2024",
	}

	// Convert to JSON
	updateData, err := json.Marshal(updateBook)
	if err != nil {
		t.Fatal(err)
	}

	// Setup router with params
	router := httprouter.New()
	router.PUT("/book/:bookId", UpdateBook)

	// Create request
	req, err := http.NewRequest("PUT", "/book/"+strconv.Itoa(int(createdBook.ID)), bytes.NewBuffer(updateData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler through router
	router.ServeHTTP(rr, req)

	// Check the status code (could be 200 if successful or 404 if book not found)
	if status := rr.Code; status != http.StatusOK && status != http.StatusNotFound {
		t.Errorf("handler returned unexpected status code: got %v want %v or %v",
			status, http.StatusOK, http.StatusNotFound)
	}
}

func TestUpdateBookInvalidID(t *testing.T) {
	// Prepare update data
	updateBook := models.Book{
		Title:  "Updated Title",
		Author: "Updated Author",
	}

	updateData, _ := json.Marshal(updateBook)

	// Setup router with params
	router := httprouter.New()
	router.PUT("/book/:bookId", UpdateBook)

	// Create request with invalid ID
	req, err := http.NewRequest("PUT", "/book/invalid", bytes.NewBuffer(updateData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler through router
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestDeleteBook(t *testing.T) {
	// Setup router with params
	router := httprouter.New()
	router.DELETE("/book/:bookId", DeleteBook)

	// Create request (using ID 999 which likely doesn't exist)
	req, err := http.NewRequest("DELETE", "/book/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler through router
	router.ServeHTTP(rr, req)

	// Check the status code (should be 404 for non-existent book)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestDeleteBookInvalidID(t *testing.T) {
	// Setup router with params
	router := httprouter.New()
	router.DELETE("/book/:bookId", DeleteBook)

	// Create request with invalid ID
	req, err := http.NewRequest("DELETE", "/book/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler through router
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check error response
	var errorResponse map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
		t.Errorf("Error response is not valid JSON: %v", err)
	}

	if errorResponse["error"] != "Invalid book ID" {
		t.Errorf("Wrong error message: got %v want %v",
			errorResponse["error"], "Invalid book ID")
	}
}

// Helper functions for pointer types
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
