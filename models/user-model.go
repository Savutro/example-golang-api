package models

import (
	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/config"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
}

func init() {
	db = config.GetDB()
	db.AutoMigrate(&User{})
}
