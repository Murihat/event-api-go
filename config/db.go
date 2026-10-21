package config

import (
	"log"
	"os"

	"example.com/event-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URI")
	if dsn == "" {
		log.Fatal("Environment Variable is required")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connection database", err)
	}

	err = database.AutoMigrate(&models.Event{}, models.User{})
	if err != nil {
		log.Fatal("Failed to migration database", err)
	}

	DB = database
	log.Println("Successfully connection to database")
}
