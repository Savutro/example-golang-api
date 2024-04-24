package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/controllers"
	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/middleware"
	"github.com/gorilla/mux"
)

const (
	port = ":8080"
)

func main() {
	// Define API routes
	r := mux.NewRouter()

	// Authentication endpoints
	r.HandleFunc("/register", controllers.RegisterUserHandler).Methods("POST")
	r.HandleFunc("/login", controllers.LoginUserHandler).Methods("POST")
	r.HandleFunc("/logout", middleware.AuthRequired(controllers.LogoutUserHandler)).Methods("POST")

	// User creates a book associated to them
	r.HandleFunc("/book", controllers.CreateBookHandler).Methods("POST")

	// User gets all books associated to them
	r.HandleFunc("/book", controllers.GetBookHandler).Methods("GET")

	// User gets info about a specifc book
	r.HandleFunc("/book/{bookId}", controllers.GetBookByIdHandler).Methods("GET")

	// User overrides existing book
	r.HandleFunc("/book/{bookId}", controllers.UpdateBookHandler).Methods("PUT")

	// User deletes a specific book that's associated to them
	r.HandleFunc("/book/{bookId}", controllers.DeleteBookHandler).Methods("DELETE")

	// Setup TLS
	cert, err := tls.LoadX509KeyPair("servert.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to load key pair: %v", err)
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	server := &http.Server{
		Addr:      port,
		Handler:   r,
		TLSConfig: tlsCfg,
	}

	log.Printf("Listening on %s...", port)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
