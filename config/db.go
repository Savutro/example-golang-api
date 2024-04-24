package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB
)

const (
	dbUser     = "root"
	dbPassword = "password12345" // TODO Secure password
	dbHost     = "127.0.0.1"
	dbPort     = "3306"
	dbName     = "bookstore"
)

func Connect() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Could not connect to the database: ", err)
		os.Exit(1)
	}
}

func GetDB() *gorm.DB {
	return db
}
