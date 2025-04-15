package models

import (
	"github.com/usama1031/book-management/pkg/config"
	"gorm.io/gorm"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `gorm:"" json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
	UserID      string `json:"uid"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{})
}

func (b *Book) CreateBook() *Book {
	db.Create(&b)
	return b
}

func GetAllBooks() []Book {
	var books []Book

	db.Find(&books)
	return books
}

func GetAllBooksByUserID(Id string) ([]Book, *gorm.DB) {
	var books []Book

	res := db.Where("user_id = ?", Id).Find(&books)
	return books, res
}

func GetBookByID(Id int64) (*Book, *gorm.DB) {
	var getBook Book
	res := db.First(&getBook, Id)
	return &getBook, res
}

func DeleteBook(Id int64) Book {
	var book Book
	db.Where("ID=?", Id).Delete(&book)
	return book
}
