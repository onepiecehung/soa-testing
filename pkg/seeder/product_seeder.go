package seeder

import (
	"log"
	"product-management/internal/models"

	"gorm.io/gorm"
)

// SeedProducts seeds initial product data if the products table is empty
func SeedProducts(db *gorm.DB) error {
	// Check if products table is empty
	var count int64
	if err := db.Model(&models.Product{}).Count(&count).Error; err != nil {
		return err
	}

	// If products exist, skip seeding
	if count > 0 {
		log.Println("Products table already has data, skipping seeding")
		return nil
	}

	// Sample products data
	products := []models.Product{
		{
			Name:          "SmartWatch Pro",
			Description:   "Advanced smartwatch with fitness tracking.",
			Price:         290.00,
			StockQuantity: 60,
			Status:        models.StatusActive,
		},
		{
			Name:          "Wireless Mouse X",
			Description:   "Ergonomic wireless mouse with silent clicks.",
			Price:         25.50,
			StockQuantity: 0,
			Status:        models.StatusInactive,
		},
		{
			Name:          "UltraBook Air",
			Description:   "Lightweight laptop with long battery life.",
			Price:         1199.00,
			StockQuantity: 42,
			Status:        models.StatusActive,
		},
		{
			Name:          "Vision 24 Monitor",
			Description:   "24-inch Full HD monitor with slim bezel.",
			Price:         179.99,
			StockQuantity: 5,
			Status:        models.StatusActive,
		},
		{
			Name:          "NoiseAway Earbuds",
			Description:   "Wireless earbuds with active noise cancellation.",
			Price:         79.95,
			StockQuantity: 89,
			Status:        models.StatusActive,
		},
		{
			Name:          "Keyboard Master",
			Description:   "Mechanical keyboard with customizable RGB lighting.",
			Price:         79.95,
			StockQuantity: 45,
			Status:        models.StatusActive,
		},
		{
			Name:          "PowerLap 15",
			Description:   "15-inch gaming laptop with powerful specs.",
			Price:         1349.00,
			StockQuantity: 18,
			Status:        models.StatusActive,
		},
		{
			Name:          "CurveView 34",
			Description:   "34-inch ultrawide curved monitor for immersive experience.",
			Price:         599.00,
			StockQuantity: 30,
			Status:        models.StatusActive,
		},
		{
			Name:          "Portable SSD 1TB",
			Description:   "1TB external solid-state drive.",
			Price:         129.00,
			StockQuantity: 95,
			Status:        models.StatusActive,
		},
		{
			Name:          "SoundWave Speaker",
			Description:   "Bluetooth speaker with 360-degree sound.",
			Price:         69.99,
			StockQuantity: 70,
			Status:        models.StatusActive,
		},
	}

	// Create sample categories
	categories := []models.Category{
		{
			Name:        "Electronics",
			Description: "Electronic devices and gadgets",
		},
		{
			Name:        "Accessories",
			Description: "Computer and device accessories",
		},
		{
			Name:        "Laptops",
			Description: "Notebook computers and laptops",
		},
		{
			Name:        "Monitors",
			Description: "Computer monitors and displays",
		},
	}

	// Insert categories
	if err := db.Create(&categories).Error; err != nil {
		return err
	}

	// Insert products
	if err := db.Create(&products).Error; err != nil {
		return err
	}

	// Associate products with categories
	productCategories := []map[string]interface{}{
		{"product_name": "SmartWatch Pro", "category_name": "Electronics"},
		{"product_name": "Wireless Mouse X", "category_name": "Accessories"},
		{"product_name": "UltraBook Air", "category_name": "Laptops"},
		{"product_name": "Vision 24 Monitor", "category_name": "Monitors"},
		{"product_name": "NoiseAway Earbuds", "category_name": "Electronics"},
		{"product_name": "Keyboard Master", "category_name": "Accessories"},
		{"product_name": "PowerLap 15", "category_name": "Laptops"},
		{"product_name": "CurveView 34", "category_name": "Monitors"},
		{"product_name": "Portable SSD 1TB", "category_name": "Accessories"},
		{"product_name": "SoundWave Speaker", "category_name": "Electronics"},
	}

	for _, pc := range productCategories {
		var product models.Product
		var category models.Category

		if err := db.Where("name = ?", pc["product_name"]).First(&product).Error; err != nil {
			return err
		}
		if err := db.Where("name = ?", pc["category_name"]).First(&category).Error; err != nil {
			return err
		}

		if err := db.Model(&product).Association("Categories").Append(&category); err != nil {
			return err
		}
	}

	log.Println("Successfully seeded products and categories")
	return nil
}
