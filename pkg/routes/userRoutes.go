package routes

import (
	"github.com/gorilla/mux"
	"github.com/usama1031/book-management/pkg/controllers"
)

var UserRoutes = func(router *mux.Router) {
	router.HandleFunc("/users/signup/", controllers.SignUp).Methods("POST")
	router.HandleFunc("/users/login/", controllers.Login).Methods("POST")
}
