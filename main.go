package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/auth"
	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/config"
	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/controllers"
	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/middleware"
	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/models"
	"github.com/gorilla/mux"
)

const (
	port = ":8080"
)

func main() {
	log.Print("Starting App...")

	// Get configuration values from a yaml
	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		log.Fatalf("Can't load config: %v", err)
	}

	db := config.Connect(cfg)

	if err := db.AutoMigrate(&models.User{}, &models.Book{}).Error; err != nil {
		log.Fatal("Could not migrate database: ", err)
	}

	auth.Init(db, cfg)

	// Define API routes
	r := mux.NewRouter()

	// Authentication endpoints
	r.HandleFunc("/register", controllers.RegisterUserHandler).Methods("POST")
	r.HandleFunc("/login", controllers.LoginUserHandler).Methods("POST")     // returns invalid session cookie
	r.HandleFunc("/twofactor", controllers.TwoFactorHandler).Methods("POST") // makes session cookie valid after 2FA
	r.HandleFunc("/logout", middleware.AuthRequired(controllers.LogoutUserHandler)).Methods("POST")

	// Example endpoint for role based authorization
	r.HandleFunc("/admin", middleware.AuthAndRoleRequired(controllers.AdminHandler)).Methods("GET")

	// User creates a book in the database
	r.HandleFunc("/book", middleware.AuthRequired(controllers.CreateBookHandler)).Methods("POST")

	// User gets all books from database
	r.HandleFunc("/book", middleware.AuthRequired(controllers.GetBookHandler)).Methods("GET")

	// User gets info about a specifc book from database
	r.HandleFunc("/book/{bookId}", middleware.AuthRequired(controllers.GetBookByIdHandler)).Methods("GET")

	// User overrides existing book with provided info
	r.HandleFunc("/book/{bookId}", middleware.AuthRequired(controllers.UpdateBookHandler)).Methods("PUT")

	// User deletes a specific book
	r.HandleFunc("/book/{bookId}", middleware.AuthRequired(controllers.DeleteBookHandler)).Methods("DELETE")

	log.Print("Registered Routes.")

	// Setup TLS
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to load key pair: %v", err)
	}
	log.Print("Loaded key pair.")

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
