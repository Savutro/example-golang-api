package controllers

import (
	"encoding/json"
	"net/http"

	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/xlzd/gotp"
)

// Helper function to send JSON response
func sendJSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Could not parse form"})
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")
	randomSecret := gotp.RandomSecret(16)

	uri := auth.GenerateTOTPWithSecret(randomSecret, username)

	err = auth.RegisterNewUser(username, password, randomSecret)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not register new user"})
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]string{"uri": uri})
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Could not parse form"})
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	err = auth.LoginUser(username, password)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Credentials are invalid"})
		return
	}

	// Create intermediate jwt token for 2FA
	tokenString, err := auth.CreateJWT(username)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not create JWT"})
		return
	}

	// Save token and username to session
	session, _ := auth.Store.Get(r, auth.SessionName)
	session.Values["token"] = tokenString
	session.Values["username"] = username
	session.Save(r, w)

	sendJSONResponse(w, http.StatusOK, map[string]string{"message": "Login successful"})
}

func TwoFactorHandler(w http.ResponseWriter, r *http.Request) {
	session, err := auth.Store.Get(r, auth.SessionName)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Couldn't get session"})
		return
	}
	tokenString, ok := session.Values["token"].(string)
	if !ok {
		sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid session token"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("super-secret-code2"), nil // TODO hide secret
	})

	if err != nil {
		sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Authentication failed"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid session token"})
			return
		}

		err := r.ParseForm()
		if err != nil {
			sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Could not parse form"})
			return
		}

		// Get code from user and secret from db
		code := r.Form.Get("2FACode")
		secret, err := auth.GetSecretFromDB(username)
		if err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not find user in token"})
			return
		}

		// Verify code with secret
		ok = auth.VerifyOTP(secret, code)
		if !ok {
			sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Couldn't verify code"})
			return
		}

		// Add authenticated flag to session and remove jwt
		session.Values["authenticated"] = true
		session.Values["token"] = nil
		session.Save(r, w)

		sendJSONResponse(w, http.StatusOK, map[string]string{"message": "2FA successful"})
	}
}

func LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.Store.Get(r, auth.SessionName)
	session.Values["username"] = nil
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not invalidate session"})
		return
	}
	sendJSONResponse(w, http.StatusOK, map[string]string{"message": "Logout successful"})
}
