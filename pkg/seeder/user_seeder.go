package seeder

import (
	"log"
	"product-management/internal/models"

	"gorm.io/gorm"
)

// SeedUsers creates initial admin and user accounts
func SeedUsers(db *gorm.DB) error {
	// Check and create admin user
	// adminPassword, err := utils.HashPassword("password123")
	// if err != nil {
	// 	return err
	// }

	admin := &models.User{
		Username: "admin_test",
		Email:    "admin@soa.com",
		Password: "password123",
		FullName: "Admin Test",
		Role:     models.RoleAdmin,
	}

	// Check if admin email exists
	var existingAdmin models.User
	if err := db.Where("email = ?", admin.Email).First(&existingAdmin).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		// Create admin if not exists
		if err := db.Create(admin).Error; err != nil {
			return err
		}
		log.Printf("Created admin user: %s", admin.Email)
	} else {
		log.Printf("Admin user already exists: %s", admin.Email)
	}

	// Check and create regular user
	// userPassword, err := utils.HashPassword("password123")
	// if err != nil {
	// 	return err
	// }

	user := &models.User{
		Username: "user_test",
		Email:    "user@soa.com",
		Password: "password123",
		FullName: "User Test",
		Role:     models.RoleUser,
	}

	// Check if user email exists
	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		// Create user if not exists
		if err := db.Create(user).Error; err != nil {
			return err
		}
		log.Printf("Created regular user: %s", user.Email)
	} else {
		log.Printf("Regular user already exists: %s", user.Email)
	}

	return nil
}
