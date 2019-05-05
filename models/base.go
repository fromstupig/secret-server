package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}

	conn, err := gorm.Open("sqlite3", "./secret.db")
	if err != nil {
		log.Fatal(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Secret{})
}

func DB() *gorm.DB {
	return db
}
