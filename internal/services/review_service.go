package services

import (
	"product-management/internal/models"
	"product-management/internal/repositories"
)

// ReviewService handles business logic for reviews
type ReviewService struct {
	reviewRepo *repositories.ReviewRepository
}

// NewReviewService creates a new review service
func NewReviewService(reviewRepo *repositories.ReviewRepository) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
	}
}

// CreateReview creates a new review
func (s *ReviewService) CreateReview(review *models.Review) error {
	return s.reviewRepo.Create(review)
}

// GetReviewByID retrieves a review by its ID
func (s *ReviewService) GetReviewByID(id uint) (*models.Review, error) {
	return s.reviewRepo.GetByID(id)
}

// GetReviewsByProductID retrieves all reviews for a product
func (s *ReviewService) GetReviewsByProductID(productID uint) ([]models.Review, error) {
	return s.reviewRepo.GetByProductID(productID)
}

// GetReviewsByUserID retrieves all reviews by a user
func (s *ReviewService) GetReviewsByUserID(userID uint) ([]models.Review, error) {
	return s.reviewRepo.GetByUserID(userID)
}

// GetReviewByUserAndProduct retrieves a review by user ID and product ID
func (s *ReviewService) GetReviewByUserAndProduct(userID, productID uint) (*models.Review, error) {
	return s.reviewRepo.GetByUserAndProduct(userID, productID)
}

// UpdateReview updates a review
func (s *ReviewService) UpdateReview(review *models.Review) error {
	return s.reviewRepo.Update(review)
}

// DeleteReview deletes a review
func (s *ReviewService) DeleteReview(id uint) error {
	return s.reviewRepo.Delete(id)
}

// GetAverageRating calculates the average rating for a product
func (s *ReviewService) GetAverageRating(productID uint) (float64, error) {
	return s.reviewRepo.GetAverageRating(productID)
}

// GetReviewCount returns the number of reviews for a product
func (s *ReviewService) GetReviewCount(productID uint) (int64, error) {
	return s.reviewRepo.GetReviewCount(productID)
}

// SearchReviews retrieves reviews with pagination, filtering, and sorting
func (s *ReviewService) SearchReviews(page, pageSize int, productName, sortBy, order string) ([]models.Review, int64, error) {
	return s.reviewRepo.Search(page, pageSize, productName, sortBy, order)
}

// CountTotalReviews counts the total number of reviews
func (s *ReviewService) CountTotalReviews() (int64, error) {
	return s.reviewRepo.CountTotalReviews()
}

// CountReviewsWithUserID counts the number of reviews for a user
func (s *ReviewService) CountReviewsWithUserID(userID uint) (int64, error) {
	return s.reviewRepo.CountReviewsWithUserID(userID)
}
