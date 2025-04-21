package repositories

import (
	"product-management/internal/dto"
	"product-management/internal/models"

	"gorm.io/gorm"
)

// CategoryRepository handles database operations for categories
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// GetByID retrieves a category by its ID
func (r *CategoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	// err := r.db.Preload("Products").First(&category, id).Error
	err := r.db.First(&category, id).Error
	return &category, err
}

// GetAll retrieves all categories
func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	// err := r.db.Preload("Products").Find(&categories).Error
	return categories, err
}

// Update updates a category
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete deletes a category
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

// GetProductsByCategoryID retrieves all products in a category
func (r *CategoryRepository) GetProductsByCategoryID(categoryID uint) ([]models.Product, error) {
	var category models.Category
	err := r.db.Preload("Products").First(&category, categoryID).Error
	if err != nil {
		return nil, err
	}
	return category.Products, nil
}

// AddProductToCategory adds a product to a category
func (r *CategoryRepository) AddProductToCategory(categoryID, productID uint) error {
	var category models.Category
	var product models.Product

	if err := r.db.First(&category, categoryID).Error; err != nil {
		return err
	}
	if err := r.db.First(&product, productID).Error; err != nil {
		return err
	}

	return r.db.Model(&category).Association("Products").Append(&product)
}

// RemoveProductFromCategory removes a product from a category
func (r *CategoryRepository) RemoveProductFromCategory(categoryID, productID uint) error {
	var category models.Category
	var product models.Product

	if err := r.db.First(&category, categoryID).Error; err != nil {
		return err
	}
	if err := r.db.First(&product, productID).Error; err != nil {
		return err
	}

	return r.db.Model(&category).Association("Products").Delete(&product)
}

// DB returns the database instance
func (r *CategoryRepository) DB() *gorm.DB {
	return r.db
}

// GetCategoryDistribution gets the distribution of products across categories
func (r *CategoryRepository) GetCategoryDistribution() ([]dto.CategoryDistributionResponse, error) {
	var distributions []dto.CategoryDistributionResponse

	err := r.db.Table("categories").
		Select("categories.name, COUNT(DISTINCT product_categories.product_id) as product_count").
		Joins("LEFT JOIN product_categories ON categories.id = product_categories.category_id").
		Group("categories.id, categories.name").
		Find(&distributions).Error

	return distributions, err
}

// GetAllWithProductCount retrieves all categories with their product counts
func (r *CategoryRepository) GetAllWithProductCount() ([]dto.CategoryResponse, error) {
	var responses []dto.CategoryResponse

	err := r.db.Table("categories").
		Select("categories.id, categories.name, categories.description, COUNT(DISTINCT product_categories.product_id) as product_count").
		Joins("LEFT JOIN product_categories ON categories.id = product_categories.category_id").
		Group("categories.id, categories.name, categories.description").
		Find(&responses).Error

	return responses, err
}
