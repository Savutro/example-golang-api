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
		session, _ := auth.Store.Get(r, auth.SessionName)
		if session.Values != nil { // TODO Potential error
			if session.Values["username"] == nil {
				// User is not logged in
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}

		handler.ServeHTTP(w, r)
	}
}

// Checks if user is logged in and has admin role
func AuthAndRoleRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := auth.Store.Get(r, auth.SessionName)
		if session.Values["username"] == nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
		// Check role
		var user models.User
		db := config.GetDB()
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
