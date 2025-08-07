package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/adedaryorh/bookstore-app/pkg/models"
	"github.com/adedaryorh/bookstore-app/pkg/routes"
	"github.com/julienschmidt/httprouter"
)

func TestMain(m *testing.M) {

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "bookstore_user")
	os.Setenv("DB_PASSWORD", "bookstore_pass")
	os.Setenv("DB_NAME", "bookstore_test")

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestBookCRUDFlow(t *testing.T) {

	router := httprouter.New()
	routes.RegisterRoutes(router)

	testBook := models.Book{
		Title:           "Integration Test Book",
		Author:          "Lincoln Author",
		ISBN:            "9999999999999",
		PublicationYear: "2023",
		Genre:           stringPtr("Testing"),
		Price:           float64Ptr(99.99),
	}

	var createdBookID uint

	t.Run("Create Book", func(t *testing.T) {
		jsonData, _ := json.Marshal(testBook)
		req, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Fatalf("Create failed with status %d", status)
		}

		var createdBook models.Book
		json.Unmarshal(rr.Body.Bytes(), &createdBook)
		createdBookID = createdBook.ID

		if createdBook.Title != testBook.Title {
			t.Errorf("Title mismatch: got %v want %v", createdBook.Title, testBook.Title)
		}
	})

	t.Run("Get Book by ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/book/"+strconv.Itoa(int(createdBookID)), nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("Get by ID failed with status %d", status)
		}

		var retrievedBook models.Book
		json.Unmarshal(rr.Body.Bytes(), &retrievedBook)

		if retrievedBook.ID != createdBookID {
			t.Errorf("ID mismatch: got %v want %v", retrievedBook.ID, createdBookID)
		}
	})

	t.Run("Get All Books", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/books", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("Get all books failed with status %d", status)
		}

		var books []models.Book
		json.Unmarshal(rr.Body.Bytes(), &books)

		found := false
		for _, book := range books {
			if book.ID == createdBookID {
				found = true
				break
			}
		}

		if !found {
			t.Error("Created book not found in the list of all books")
		}
	})

	t.Run("Update Book", func(t *testing.T) {
		updatedBook := testBook
		updatedBook.Title = "Updated Integration Test Book"
		updatedBook.Price = float64Ptr(199.99)

		jsonData, _ := json.Marshal(updatedBook)
		req, _ := http.NewRequest("PUT", "/book/"+strconv.Itoa(int(createdBookID)), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("Update failed with status %d", status)
		}

		var resultBook models.Book
		json.Unmarshal(rr.Body.Bytes(), &resultBook)

		if resultBook.Title != "Updated Integration Test Book" {
			t.Errorf("Title not updated: got %v want %v", resultBook.Title, "Updated Integration Test Book")
		}
	})

	t.Run("Delete Book", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/book/"+strconv.Itoa(int(createdBookID)), nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("Delete failed with status %d", status)
		}

		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)

		if response["message"] != "Book deleted successfully" {
			t.Errorf("Unexpected delete message: %v", response["message"])
		}
	})

	t.Run("Verify Book Deleted", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/book/"+strconv.Itoa(int(createdBookID)), nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected 404 for deleted book, got %d", status)
		}
	})
}

func TestHealthCheck(t *testing.T) {
	router := httprouter.New()
	routes.RegisterRoutes(router)

	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Health check failed with status %d", status)
	}

	if body := rr.Body.String(); body != "OK" {
		t.Errorf("Health check returned wrong body: got %v want %v", body, "OK")
	}
}

func TestErrorScenarios(t *testing.T) {
	router := httprouter.New()
	routes.RegisterRoutes(router)

	t.Run("Get Non-existent Book", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/book/99999", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected 404 for non-existent book, got %d", status)
		}
	})

	t.Run("Create Book with Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/book", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected 400 for invalid JSON, got %d", status)
		}
	})

	t.Run("Update Non-existent Book", func(t *testing.T) {
		testBook := models.Book{Title: "Test", Author: "Test"}
		jsonData, _ := json.Marshal(testBook)
		req, _ := http.NewRequest("PUT", "/book/99999", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected 404 for updating non-existent book, got %d", status)
		}
	})

	t.Run("Delete Non-existent Book", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/book/99999", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected 404 for deleting non-existent book, got %d", status)
		}
	})

	t.Run("Invalid Book ID Format", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/book/invalid", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected 400 for invalid book ID, got %d", status)
		}
	})
}

func TestMultipleGetEndpoints(t *testing.T) {
	router := httprouter.New()
	routes.RegisterRoutes(router)

	t.Run("GET /book", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/book", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("GET /book failed with status %d", status)
		}

		var books []models.Book
		if err := json.Unmarshal(rr.Body.Bytes(), &books); err != nil {
			t.Errorf("Response is not valid JSON array: %v", err)
		}
	})

	t.Run("GET /books", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/books", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("GET /books failed with status %d", status)
		}

		var books []models.Book
		if err := json.Unmarshal(rr.Body.Bytes(), &books); err != nil {
			t.Errorf("Response is not valid JSON array: %v", err)
		}
	})
}

func TestConcurrentOperations(t *testing.T) {
	router := httprouter.New()
	routes.RegisterRoutes(router)
	bookCount := 5
	results := make(chan error, bookCount)

	for i := 0; i < bookCount; i++ {
		go func(index int) {
			testBook := models.Book{
				Title:           fmt.Sprintf("Concurrent Book %d", index),
				Author:          fmt.Sprintf("Author %d", index),
				ISBN:            fmt.Sprintf("123456789012%d", index),
				PublicationYear: "2023",
			}

			jsonData, _ := json.Marshal(testBook)
			req, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusCreated {
				results <- fmt.Errorf("concurrent create failed with status %d", rr.Code)
				return
			}
			results <- nil
		}(i)
	}

	for i := 0; i < bookCount; i++ {
		if err := <-results; err != nil {
			t.Errorf("Concurrent operation failed: %v", err)
		}
	}
}

func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
