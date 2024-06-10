package models

import (
	//"strconv"

	"github.com/jinzhu/gorm"
)

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func GetAllBooks(db *gorm.DB) []Book {
	var books []Book
	db.Find(&books)
	return books
}

func GetBookById(Id int64, db *gorm.DB) (*Book, *gorm.DB) {
	var getBook Book
	// NOTE: gorm query
	dbred := db.Where("ID=?", Id).Find(&getBook)

	// deliberate security issue
	// dbred := db.Exec("SELECT * FROM books WHERE ID=" + strconv.FormatInt(Id, 10))
	return &getBook, dbred
}

func (b *Book) CreateBook(db *gorm.DB) *Book {
	db.NewRecord(b)
	db.Create(&b)
	return b
}

func DeleteBook(Id int64, db *gorm.DB) Book {
	var book Book
	db.Unscoped().Where("ID=?", Id).Delete(book)
	return book
}
