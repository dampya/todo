package database

import (
	"log"
	"os"
	"fmt"

	"go/todo/models"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)


func ConnectDB() *gorm.DB {
	dsn:= fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	if err:= db.AutoMigrate(&models.User{}, &models.Todo{}); err != nil {
		log.Fatal("migration failed:", err)
	}
	
	return db
}
