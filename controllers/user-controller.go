package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/xlzd/gotp"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")
	randomSecret := gotp.RandomSecret(16)

	uri := auth.GenerateTOTPWithSecret(randomSecret) //TODO Return URI instead of printing at the end
	fmt.Print(uri)

	err = auth.RegisterNewUser(username, password, randomSecret)
	if err != nil {
		http.Error(w, "Could not register new user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"uri": uri})
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
		http.Error(w, "Credentials are invalid", http.StatusBadRequest)
		return
	}

	// Create intermediate jwt token for 2FA
	tokenString, err := auth.CreateJWT(username)
	if err != nil {
		http.Error(w, "Could not create JWT", http.StatusInternalServerError)
		return
	}

	// Save token to session
	session, _ := auth.Store.Get(r, auth.SessionName)
	session.Values["token"] = tokenString
	session.Save(r, w)

	w.Write([]byte("Login successful"))
}

func TwoFactorHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.Store.Get(r, auth.SessionName)
	tokenString, ok := session.Values["token"].(string)
	if !ok {
		http.Error(w, "Invalid session token.", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your JWT secret"), nil
	})

	if err != nil {
		http.Error(w, "Authentication failed.", http.StatusUnauthorized)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			http.Error(w, "Invalid session token.", http.StatusUnauthorized)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Could not parse form", http.StatusBadRequest)
			return
		}

		// Get code from user and secret from db
		code := r.Form.Get("2FACode")
		secret, err := auth.GetSecretFromDB(username)
		if err != nil {
			http.Error(w, "Could not find user in token", http.StatusInternalServerError)
			return
		}

		// Verify code with secret
		ok = auth.VerifyOTP(secret, code)
		if !ok {
			http.Error(w, "Couldn't verify code", http.StatusBadRequest)
			return
		}
	}

}

func LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.Store.Get(r, auth.SessionName)
	session.Values["username"] = nil
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Could not invalidate session.", http.StatusInternalServerError)
	}
	w.Write([]byte("Logout successful"))
}
