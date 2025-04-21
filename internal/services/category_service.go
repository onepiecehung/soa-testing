package services

import (
	"errors"
	"product-management/internal/dto"
	"product-management/internal/models"
	"product-management/internal/repositories"
	"product-management/pkg/database"

	"gorm.io/gorm"
)

// CategoryService handles business logic for categories
type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

// NewCategoryService creates a new CategoryService instance
func NewCategoryService() *CategoryService {
	return &CategoryService{
		categoryRepo: repositories.NewCategoryRepository(database.DB),
	}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(req dto.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoryByID retrieves a category by ID
func (s *CategoryService) GetCategoryByID(id uint) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return category, nil
}

// GetAllCategories retrieves all categories
func (s *CategoryService) GetAllCategories() ([]dto.CategoryResponse, error) {
	return s.categoryRepo.GetAllWithProductCount()
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(id uint, req dto.UpdateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		BaseModel:   models.BaseModel{ID: id},
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryRepo.Update(category); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(id uint) error {
	// Check if category has any products
	var count int64
	if err := s.categoryRepo.DB().Model(&models.ProductCategory{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("cannot delete category with associated products")
	}

	return s.categoryRepo.Delete(id)
}

// GetProductsByCategoryID retrieves all products in a category
func (s *CategoryService) GetProductsByCategoryID(categoryID uint) ([]models.Product, error) {
	return s.categoryRepo.GetProductsByCategoryID(categoryID)
}

// AddProductToCategory adds a product to a category
func (s *CategoryService) AddProductToCategory(categoryID, productID uint) error {
	return s.categoryRepo.AddProductToCategory(categoryID, productID)
}

// RemoveProductFromCategory removes a product from a category
func (s *CategoryService) RemoveProductFromCategory(categoryID, productID uint) error {
	return s.categoryRepo.RemoveProductFromCategory(categoryID, productID)
}

// GetCategoryDistribution gets the distribution of products across categories
func (s *CategoryService) GetCategoryDistribution() ([]dto.CategoryDistributionResponse, error) {
	return s.categoryRepo.GetCategoryDistribution()
}
