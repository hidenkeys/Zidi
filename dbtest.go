package main

import (
	"github.com/hidenkeys/zidibackend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDatabase() {
	//dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
	//	os.Getenv("DB_HOST"),
	//	os.Getenv("DB_USER"),
	//	os.Getenv("DB_PASSWORD"),
	//	os.Getenv("DB_NAME"),
	//	os.Getenv("DB_PORT"),
	//)

	db, err := gorm.Open(postgres.Open("postgres://avnadmin:AVNS_T-ehBEdXEPR6M3dQWeX@pg-18b4c785-letimapro23-87d3.h.aivencloud.com:20123/defaultdb?sslmode=require"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("📦 Connected to database")
}

func MigrateDatabase() {
	err := DB.AutoMigrate(&models.Organization{}, &models.User{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("✅ Database migration successful")
}
