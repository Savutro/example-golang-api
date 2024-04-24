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

const (
	dbUser     = "root"
	dbPassword = "w+dopkAx05452" // TODO Secure password
	dbHost     = "127.0.0.1"
	dbPort     = "3306"
	dbName     = "bookstore"
)

func Connect() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

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
