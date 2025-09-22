package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=127.0.0.1 user=postgres password=admin123 dbname=booking_db port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // <--- Silence GORM info logs
	})
	if err != nil {
		log.Fatal("Failed to connect to database! ", err)
	}

	DB = db
	fmt.Println("Connected to Postgres Database")
}
