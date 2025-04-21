package repositories

import (
	"product-management/internal/models"

	"gorm.io/gorm"
)

// ReviewRepository handles database operations for reviews
type ReviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository creates a new review repository
func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// Create creates a new review
func (r *ReviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

// GetByID retrieves a review by its ID
func (r *ReviewRepository) GetByID(id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("User").First(&review, id).Error
	return &review, err
}

// GetByProductID retrieves all reviews for a product
func (r *ReviewRepository) GetByProductID(productID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Preload("User").
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

// GetByUserID retrieves all reviews by a user
func (r *ReviewRepository) GetByUserID(userID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Preload("Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

// GetByUserAndProduct retrieves a review by user ID and product ID
func (r *ReviewRepository) GetByUserAndProduct(userID, productID uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// Update updates a review
func (r *ReviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

// Delete deletes a review
func (r *ReviewRepository) Delete(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}

// GetAverageRating calculates the average rating for a product
func (r *ReviewRepository) GetAverageRating(productID uint) (float64, error) {
	var avg float64
	err := r.db.Model(&models.Review{}).
		Where("product_id = ?", productID).
		Select("AVG(rating)").
		Row().
		Scan(&avg)
	return avg, err
}

// GetReviewCount returns the number of reviews for a product
func (r *ReviewRepository) GetReviewCount(productID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Review{}).
		Where("product_id = ?", productID).
		Count(&count).Error
	return count, err
}

// Search retrieves reviews with pagination, filtering, and sorting
func (r *ReviewRepository) Search(page, pageSize int, productName, sortBy, order string) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := r.db.Model(&models.Review{}).
		Preload("User").
		Preload("Product")

	// Apply product name filter if provided
	if productName != "" {
		query = query.Joins("JOIN products ON products.id = reviews.product_id").
			Where("products.name LIKE ?", "%"+productName+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sortBy != "" {
		query = query.Order(sortBy + " " + order)
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&reviews).Error

	return reviews, total, err
}

// CountTotalReviews counts the total number of reviews for all products
func (r *ReviewRepository) CountTotalReviews() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Review{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountReviewsWithUserID counts the number of reviews for a user
func (r *ReviewRepository) CountReviewsWithUserID(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Review{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
