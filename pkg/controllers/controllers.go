package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/adedaryorh/bookstore-app/pkg/models"
	"github.com/julienschmidt/httprouter"
)

func GetAllBooks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	books := models.GetAllBooks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func GetBooks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	GetAllBooks(w, r, nil)
}

func GetBookByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bookIdStr := ps.ByName("bookId")
	bookId, err := strconv.ParseUint(bookIdStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid book ID"})
		return
	}

	book, err := models.GetBookByID(uint(bookId))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Book not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func CreateBook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON format"})
		return
	}

	createdBook := book.CreateBook()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdBook)
}

func UpdateBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bookIdStr := ps.ByName("bookId")
	bookId, err := strconv.ParseUint(bookIdStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid book ID"})
		return
	}

	existingBook, err := models.GetBookByID(uint(bookId))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Book not found"})
		return
	}

	var updatedBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON format"})
		return
	}

	updatedBook.ID = existingBook.ID
	updatedBook.CreatedAt = existingBook.CreatedAt

	if err := updatedBook.UpdateBook(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update book"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBook)
}

func DeleteBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bookIdStr := ps.ByName("bookId")
	bookId, err := strconv.ParseUint(bookIdStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid book ID"})
		return
	}
	_, err = models.GetBookByID(uint(bookId))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Book not found"})
		return
	}

	deletedBook := models.DeleteBook(uint(bookId))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Book deleted successfully",
		"book":    deletedBook,
	})
}
