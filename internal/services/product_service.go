package services

import (
	"errors"
	"product-management/internal/models"
	"product-management/internal/repositories"
	"product-management/pkg/database"
)

// ProductService handles business logic for products
type ProductService struct {
	productRepo *repositories.ProductRepository
}

// NewProductService creates a new ProductService instance
func NewProductService() *ProductService {
	return &ProductService{
		productRepo: repositories.NewProductRepository(database.DB),
	}
}

// CreateProduct creates a new product with validation
func (s *ProductService) CreateProduct(product *models.Product, categories []models.Category) error {
	// Validate required fields
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	if product.StockQuantity < 0 {
		return errors.New("stock quantity cannot be negative")
	}
	if product.Status == "" {
		product.Status = models.StatusActive
	}

	return s.productRepo.Create(product, categories)
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(id uint) (*models.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// UpdateProduct updates an existing product with validation
func (s *ProductService) UpdateProduct(product *models.Product, categoryIDs []uint) error {
	// Validate required fields
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	if product.StockQuantity < 0 {
		return errors.New("stock quantity cannot be negative")
	}

	return s.productRepo.Update(product, categoryIDs)
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(id uint) error {
	return s.productRepo.Delete(id)
}

// ListProducts retrieves a paginated list of products with filters
func (s *ProductService) ListProducts(page, limit int, categoryID uint, search string, sort string, statuses []string) ([]models.Product, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.productRepo.List(page, limit, categoryID, search, sort, statuses)
}

// AddToWishlist adds a product to a user's wishlist
func (s *ProductService) AddToWishlist(userID, productID uint) error {
	// Check if product exists
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	return s.productRepo.AddToWishlist(userID, productID)
}

// RemoveFromWishlist removes a product from a user's wishlist
func (s *ProductService) RemoveFromWishlist(userID, productID uint) error {
	return s.productRepo.RemoveFromWishlist(userID, productID)
}

// GetWishlist retrieves a user's wishlist
func (s *ProductService) GetWishlist(userID uint, page, limit int) ([]models.Wishlist, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.productRepo.GetWishlist(userID, page, limit)
}

// IsProductInWishlist checks if a product is already in the user's wishlist
func (s *ProductService) IsProductInWishlist(userID, productID uint) (bool, error) {
	var count int64
	err := s.productRepo.DB().Model(&models.Wishlist{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
