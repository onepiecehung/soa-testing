package repositories

import (
	"product-management/internal/models"
	"strings"

	"gorm.io/gorm"
)

// ProductRepository handles database operations for products
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository instance
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create creates a new product with categories
func (r *ProductRepository) Create(product *models.Product, categories []models.Category) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return err
		}
		if len(categories) > 0 {
			return tx.Model(product).Association("Categories").Append(categories)
		}
		return nil
	})
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Categories").Preload("Reviews").First(&product, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// GetAll retrieves all products
func (r *ProductRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Categories").Preload("Reviews").Find(&products).Error
	return products, err
}

// Update updates a product and its categories
func (r *ProductRepository) Update(product *models.Product, categoryIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(product).Select("name", "description", "price", "stock_quantity", "status").Updates(product).Error; err != nil {
			return err
		}

		if err := tx.Model(product).Association("Categories").Clear(); err != nil {
			return err
		}

		if len(categoryIDs) > 0 {
			var categories []models.Category
			if err := tx.Find(&categories, categoryIDs).Error; err != nil {
				return err
			}
			return tx.Model(product).Association("Categories").Append(categories)
		}
		return nil
	})
}

// Delete deletes a product
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// List retrieves a paginated list of products with filters
func (r *ProductRepository) List(page, limit int, categoryID uint, search string, sort string, statuses []string) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})

	// Apply status filter if provided
	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}

	// Apply category filter if provided
	if categoryID > 0 {
		query = query.Joins("JOIN product_categories ON products.id = product_categories.product_id").
			Where("product_categories.category_id = ?", categoryID)
	}

	// Apply search filter if provided
	if search != "" {
		search = "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", search, search)
	}

	// Apply sorting
	switch sort {
	case "name":
		query = query.Order("name")
	case "price":
		query = query.Order("price")
	case "created_at":
		query = query.Order("created_at desc")
	default:
		query = query.Order("created_at desc")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Preload("Categories").Preload("Reviews").
		Offset(offset).Limit(limit).Find(&products).Error

	return products, total, err
}

// AddToWishlist adds a product to a user's wishlist
func (r *ProductRepository) AddToWishlist(userID, productID uint) error {
	wishlist := &models.Wishlist{
		UserID:    userID,
		ProductID: productID,
	}
	return r.db.Create(wishlist).Error
}

// RemoveFromWishlist removes a product from a user's wishlist
func (r *ProductRepository) RemoveFromWishlist(userID, productID uint) error {
	return r.db.Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.Wishlist{}).Error
}

// GetWishlist retrieves a user's wishlist
func (r *ProductRepository) GetWishlist(userID uint, page, limit int) ([]models.Wishlist, int64, error) {
	var wishlist []models.Wishlist
	var total int64

	// Count total records
	if err := r.db.Model(&models.Wishlist{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and preload product with its categories
	offset := (page - 1) * limit
	err := r.db.Preload("Product.Categories").
		Where("user_id = ?", userID).
		Offset(offset).Limit(limit).
		Find(&wishlist).Error

	return wishlist, total, err
}

// CountTotalWishlistItems counts the total number of wishlist items
func (r *ProductRepository) CountTotalWishlistItems() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Wishlist{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountUserWishlistItems counts the number of wishlist items for a specific user
func (r *ProductRepository) CountUserWishlistItems(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Wishlist{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// DB returns the database instance
func (r *ProductRepository) DB() *gorm.DB {
	return r.db
}
