package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Credentials — in production, use environment variables
	dsn := "host=127.0.0.1 user=postgres password=admin123 dbname=booking_db port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database! ", err)
	}

	DB = database
	fmt.Println("Connected to Postgres Database")
}
