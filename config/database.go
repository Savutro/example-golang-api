package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

func Connect(cfg Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Could not connect to the database: ", err)
		os.Exit(1)
	}
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
