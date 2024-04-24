package models

import (
	"github.com/jinzhu/gorm"
)

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func (b *Book) CreateBook(db *gorm.DB) *Book {
	db.NewRecord(b)
	db.Create(&b)
	return b
}

func GetAllBooks(db *gorm.DB) []Book {
	var books []Book
	db.Find(&books)
	return books
}

func GetBookById(Id int64, db *gorm.DB) (*Book, *gorm.DB) {
	var getBook Book
	dbred := db.Where("ID=?", Id).Find(&getBook)
	return &getBook, dbred
}

func DeleteBook(Id int64, db *gorm.DB) Book {
	var book Book
	db.Unscoped().Where("ID=?", Id).Delete(book)
	return book
}
