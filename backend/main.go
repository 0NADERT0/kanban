package main

import (
	"fmt"
	"log"
	"net/http"

	gHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"backend/db"
	"backend/handlers"
	"backend/middleware"
)

func main() {
	middleware.LoadEnv()
	db.InitDB()

	r := mux.NewRouter()
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	s := r.PathPrefix("/").Subrouter()
	s.Use(middleware.AuthMiddleware)
	s.HandleFunc("/board", handlers.BoardHandler).Methods("GET")
	s.HandleFunc("/cards", handlers.CreateCardHandler).Methods("POST")
	s.HandleFunc("/cards/{id}", handlers.DeleteCardHandler).Methods("DELETE")
	s.HandleFunc("/cards/{id}", handlers.UpdateCardHandler).Methods("PUT")

	corsHandler := gHandlers.CORS(
		gHandlers.AllowedOrigins([]string{"*"}),
		gHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		gHandlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS", "PUT"}),
	)(r)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
