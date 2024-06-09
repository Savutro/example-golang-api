package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
)

func CorsHandler() func(http.Handler) http.Handler {
	// Apply CORS middleware to the router
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}), // Allow requests from your frontend URL
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	return corsHandler
}
