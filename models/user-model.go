package models

import (
	"errors"
	"net/mail"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UUID     string `gorm:"type:varchar(36);unique"`
	Username string `gorm:"type:varchar(100);unique"`
	Password string `gorm:"not null;size:72"`
	Role     string `gorm:"type:varchar(20)"`
}

func (u *User) BeforeSave() error {
	if _, err := mail.ParseAddress(u.Username); err != nil {
		return errors.New("invalid email")
	}

	return nil
}
