package main

import (
	"fmt"
	"log"
	"product-management/config"
	"product-management/internal/models"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, strconv.Itoa(cfg.DBPort), cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Category{},
		&models.Review{},
		&models.Wishlist{},
		&models.ProductCategory{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	log.Println("Auto migration completed successfully")
}
