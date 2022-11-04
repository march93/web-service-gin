package database

import (
	"fmt"
	"web-service-gin/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DB_USERNAME = "root"
const DB_PASSWORD = "password"
const DB_NAME = "go_db"
const DB_HOST = "127.0.0.1"
const DB_PORT = "3306"

var Db *gorm.DB

func connectDB() *gorm.DB {
	dbString := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
	fmt.Println("dbString: ", dbString)

	db, err := gorm.Open(mysql.Open(dbString), &gorm.Config{})

	if err != nil {
		fmt.Printf("Error connecting to database : error=%v", err)
		return nil
	}

	db.AutoMigrate(&models.Document{})
	db.AutoMigrate(&models.Repository{})
	return db
}

func InitDB() {
	Db = connectDB()
}

func GetDB() *gorm.DB {
	return Db
}

func ClearTable() {
	Db.Exec("DELETE FROM repositories;")
	Db.Exec("DELETE FROM documents;")
}
