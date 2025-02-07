package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/config"
	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/models"
	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/utils"
	"github.com/gorilla/mux"
)

var NewBook models.Book

func GetBookHandler(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()
	newBooks := models.GetAllBooks(db)
	res, _ := json.Marshal(newBooks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBookByIdHandler(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
	}
	bookDetails, _ := models.GetBookById(ID, db)
	res, _ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()
	CreateBook := &models.Book{}
	utils.ParseBody(r, CreateBook)
	b := CreateBook.CreateBook(db)
	res, _ := json.Marshal(b)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
	}
	book := models.DeleteBook(ID, db)
	res, _ := json.Marshal(book)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()
	var updateBook = &models.Book{}
	utils.ParseBody(r, updateBook)
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
	}
	bookDetails, db := models.GetBookById(ID, db)
	if updateBook.Name != "" {
		bookDetails.Name = updateBook.Name
	}
	if updateBook.Author != "" {
		bookDetails.Author = updateBook.Author
	}
	if updateBook.Publication != "" {
		bookDetails.Publication = updateBook.Publication
	}
	db.Save(&bookDetails)
	res, _ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
