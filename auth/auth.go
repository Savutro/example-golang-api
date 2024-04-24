package auth

import (
	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/models"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var (
	key   = []byte("super-secret-key") // TODO solve with env
	Store = sessions.NewCookieStore([]byte(key))
	db    *gorm.DB
)

const SessionName = "user-session" // TODO should be user specific

func Init(database *gorm.DB) {
	db = database
}

func RegisterNewUser(username, password string) error {
	// Hash the password using bcrypt
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create a new user instance
	user := &models.User{
		Username: username,
		Password: string(passwordHash),
	}

	db.Create(&user)

	return nil
}

func LoginUser(username, password string) error {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return err
	}

	return nil
}

func LogoutUser() error {
	// TODO implement

	return nil
}
