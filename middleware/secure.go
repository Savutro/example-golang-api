package middleware

import (
	"net/http"

	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/auth"
	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/config"
	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/models"
)

// Checks if user is logged in otherwise returns error
func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.Store.Get(r, auth.SessionName)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if username, ok := session.Values["username"].(string); !ok || username == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

// Checks if user is logged in and has admin role
func AuthAndRoleRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		db := config.GetDB()
		session, err := auth.Store.Get(r, auth.SessionName)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		if err := db.Where("username = ?", session.Values["username"]).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		handler.ServeHTTP(w, r)
	}
}
