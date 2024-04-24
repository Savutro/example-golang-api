package controllers

import (
	"net/http"

	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/auth"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	err = auth.RegisterNewUser(username, password)
	if err != nil {
		http.Error(w, "Could not register new user", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Write([]byte("Registration successful"))
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	err = auth.LoginUser(username, password)
	if err != nil {
		// Respond with error invalid credentials
		return
	}

	session, _ := auth.Store.Get(r, auth.SessionName)
	session.Values["username"] = username
	session.Save(r, w)

	// Respond with a success message
}

func LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.Store.Get(r, auth.SessionName)
	session.Values["username"] = nil
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		// respond with an error
		return
	}

	// Respond with a success message
}
