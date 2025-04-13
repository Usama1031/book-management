package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/usama1031/book-management/pkg/models"
	"github.com/usama1031/book-management/pkg/routes"
)

func main() {

	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)

	log.Println("Server started at http://localhost:9011")
	log.Fatal(http.ListenAndServe(":9011", r))

}
