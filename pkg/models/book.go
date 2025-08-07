package models

import (
	"time"

	"github.com/adedaryorh/bookstore-app/pkg/config"
	"github.com/jinzhu/gorm"
)

type Book struct {
	ID              uint      `json:"id" db:"id" gor:"primary_key"`
	Title           string    `json:"title" db:"title"`
	Author          string    `json:"author" db:"author"`
	ISBN            string    `json:"isbn" db:"isbn"`
	PublicationYear string    `json:"publication_year" db:"publication_year"`
	Genre           *string   `json:"genre" db:"genre"`
	Price           *float64  `json:"price" db:"price"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

var Db *gorm.DB

func init() {
	config.Connect()
	Db = config.GetDb()
	Db.AutoMigrate(&Book{})
}

func (b *Book) CreateBook() *Book {
	Db.Create(b)
	return b
}

func GetAllBooks() []Book {
	var books []Book
	Db.Find(&books)
	return books
}

func GetBookByID(id uint) (*Book, error) {
	var book Book
	if err := Db.First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (b *Book) UpdateBook() error {
	if err := Db.Save(b).Error; err != nil {
		return err
	}
	return nil
}

func DeleteBook(id uint) *Book {
	var book Book
	Db.Where("id = ?", id).First(&book)
	Db.Delete(&book)
	return &book
}
