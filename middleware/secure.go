package middleware

import (
	"net/http"

	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/auth"
)

// Checks if user is logged in otherwise returns error
func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := auth.Store.Get(r, auth.SessionName)
		if session.Values["username"] == nil {
			// User is not logged in
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		handler.ServeHTTP(w, r)
	}
}
