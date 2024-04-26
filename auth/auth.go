package auth

import (
	"fmt"

	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/config"
	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/models"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var (
	cfg         config.Config
	Store       *sessions.CookieStore
	SessionName string
	db          *gorm.DB
)

func Init(database *gorm.DB, config config.Config) {
	db = database
	cfg = config
	SessionName = cfg.SessionName
	Store = sessions.NewCookieStore([]byte(cfg.CookiePassword))
}

func RegisterNewUser(username, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		UUID:     uuid.New().String(),
		Username: username,
		Password: string(passwordHash),
		Role:     "standard",
	}

	db.Create(&user)

	return nil
}

func LoginUser(username, password string) error {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return fmt.Errorf("username couldn't be found: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return fmt.Errorf("password doesn't match the hash: %v", err)
	}

	return nil
}

func LogoutUser() error {
	// TODO implement

	return nil
}
