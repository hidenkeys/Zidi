package config

import (
	"github.com/hidenkeys/zidibackend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("📦 Connected to database")
}

func MigrateDatabase() {
	err := DB.AutoMigrate(&models.Organization{}, &models.User{}, &models.Campaign{}, &models.Customer{}, &models.Question{}, &models.Response{}, &models.Coupon{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("✅ Database migration successful")
}
