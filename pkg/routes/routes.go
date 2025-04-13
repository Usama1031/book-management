package routes

import (
	"github.com/gorilla/mux"
	"github.com/usama1031/book-management/pkg/controllers"
)

var RegisterBookStoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/", controllers.GetBook).Methods("GET")
	router.HandleFunc("/{bookId}", controllers.GetBookByID).Methods("GET")
	router.HandleFunc("/{bookId}", controllers.UpdateBook).Methods("PUT")
	router.HandleFunc("/{bookId}", controllers.DeleteBook).Methods("DELETE")
}
