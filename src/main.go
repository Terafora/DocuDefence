package main

import (
	"DocuDefense/src/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.Handle("/users/{id}", handlers.BasicAuthMiddleware(http.HandlerFunc(handlers.UpdateUser))).Methods("PUT")
	r.Handle("/users/{id}", handlers.BasicAuthMiddleware(http.HandlerFunc(handlers.DeleteUser))).Methods("DELETE")
	r.Handle("/users/{id}/upload", handlers.BasicAuthMiddleware(http.HandlerFunc(handlers.UploadFile))).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
