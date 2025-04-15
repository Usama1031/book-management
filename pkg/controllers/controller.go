package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/usama1031/book-management/pkg/models"
	"github.com/usama1031/book-management/pkg/utils"
)

func GetBook(w http.ResponseWriter, r *http.Request) {

	userType := r.Context().Value("user_type").(string)
	userID := r.Context().Value("uid").(string)

	var books []models.Book

	if userType == "ADMIN" {

		books = models.GetAllBooks()

	} else {
		var res []models.Book

		res, result := models.GetAllBooksByUserID(userID)

		if result.Error != nil {
			http.Error(w, "Could not retrieve books", http.StatusInternalServerError)
			return
		}

		books = res
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)

}

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	Id, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
		return
	}

	bookDetails, res := models.GetBookByID(Id)

	if res.Error != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	userType := r.Context().Value("user_type").(string)
	userID := r.Context().Value("uid").(string)

	if userType == "ADMIN" || bookDetails.UserID == userID {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bookDetails)
	} else {
		http.Error(w, "Do not have the permission", http.StatusUnauthorized)
	}

}

func CreateBook(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("uid").(string)

	log.Println(userID)

	CreateBook := &models.Book{}
	utils.ParseBody(r, CreateBook)

	CreateBook.UserID = userID

	b := CreateBook.CreateBook()
	res, _ := json.Marshal(b)

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]

	Id, err := strconv.ParseInt(bookId, 0, 0)

	if err != nil {
		fmt.Println("error while parsing")
	}

	userType := r.Context().Value("user_type").(string)
	userID := r.Context().Value("uid").(string)

	bookDetails, result := models.GetBookByID(Id)

	if result.Error != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	if userType != "ADMIN" && userID != bookDetails.UserID {
		http.Error(w, "Do not have the permission to perform this action!", http.StatusUnauthorized)
		return
	}

	book := models.DeleteBook(Id)

	res, _ := json.Marshal(book)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {

	var updateBook = &models.Book{}
	utils.ParseBody(r, updateBook)

	vars := mux.Vars(r)
	bookId := vars["bookId"]

	Id, err := strconv.ParseInt(bookId, 0, 0)

	if err != nil {
		fmt.Println("error while parsing")
	}

	bookDetails, db := models.GetBookByID(Id)

	userType := r.Context().Value("user_type").(string)
	userID := r.Context().Value("uid").(string)

	if userType != "ADMIN" && userID != bookDetails.UserID {
		http.Error(w, "Do not have the permission to perform this action!", http.StatusUnauthorized)
		return
	}

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
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
