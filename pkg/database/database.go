package database

import (
	"fmt"
	"log"
	"product-management/config"
	"product-management/internal/models"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database instance
var DB *gorm.DB

const maxRetries = 5
const retryDelay = 3 * time.Second

// Connect establishes a connection to the database with retry
func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		strconv.Itoa(cfg.DBPort),
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName)

	var err error

	for i := 1; i <= maxRetries; i++ {
		// Configure connection pooling
		dbConfig := &gorm.Config{
			PrepareStmt: true, // Enable prepared statement cache
		}

		// Open database connection with pooling
		DB, err = gorm.Open(postgres.Open(dsn), dbConfig)
		if err == nil {
			// Get underlying sql.DB
			sqlDB, err := DB.DB()
			if err != nil {
				return fmt.Errorf("failed to get database instance: %v", err)
			}

			// Set connection pool settings
			sqlDB.SetMaxIdleConns(10)                  // Maximum number of idle connections
			sqlDB.SetMaxOpenConns(100)                 // Maximum number of open connections
			sqlDB.SetConnMaxLifetime(time.Hour)        // Maximum lifetime of a connection
			sqlDB.SetConnMaxIdleTime(30 * time.Minute) // Maximum idle time of a connection

			// Check if DB is actually alive
			if err := sqlDB.Ping(); err == nil {
				log.Printf("✅ Connected to DB on attempt %d", i)
				break
			}
		}

		log.Printf("⚠️ Failed to connect to DB (attempt %d/%d): %v", i, maxRetries, err)
		time.Sleep(retryDelay)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	// Auto migrate models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Category{},
		&models.Review{},
		&models.Wishlist{},
		&models.ProductCategory{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %v", err)
	}

	log.Println("✅ Database connection established and migrations completed")
	return nil
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}
	return sqlDB.Close()
}
