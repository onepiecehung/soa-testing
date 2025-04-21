package handlers

import (
	"net/http"
	"strconv"
	"time"

	"product-management/internal/dto"
	"product-management/internal/models"
	"product-management/internal/services"
	"product-management/internal/types"
	"product-management/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ReviewHandler handles review-related HTTP requests
type ReviewHandler struct {
	reviewService *services.ReviewService
}

// NewReviewHandler creates a new review handler
func NewReviewHandler(reviewService *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// CreateReview handles the creation of a new review
// @Summary Create a new review
// @Description Create a new review for a product
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body dto.CreateReviewRequest true "Review data"
// @Success 201 {object} dto.ReviewResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 409 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security Bearer
// @Router /reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req dto.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	// Check if user has already reviewed this product
	existingReview, err := h.reviewService.GetReviewByUserAndProduct(userID.(uint), req.ProductID)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to check existing review",
		})
		return
	}

	if existingReview != nil {
		c.JSON(http.StatusConflict, types.ErrorResponse{
			Error: "You have already reviewed this product",
		})
		return
	}

	review := &models.Review{
		UserID:    userID.(uint),
		ProductID: req.ProductID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := h.reviewService.CreateReview(review); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to create review",
		})
		return
	}

	response := dto.ReviewResponse{
		ID:        review.ID,
		UserID:    review.UserID,
		ProductID: review.ProductID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt.Format(time.RFC3339),
		UpdatedAt: review.UpdatedAt.Format(time.RFC3339),
	}

	logger.WithFields(logrus.Fields{
		"review_id":  review.ID,
		"product_id": review.ProductID,
		"user_id":    review.UserID,
	}).Info("Review created successfully")

	c.JSON(http.StatusCreated, response)
}

// GetReviewByID godoc
// @Summary      Get a review
// @Description  Get a review by its ID
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      int  true  "Review ID"
// @Success      200  {object}  models.Review
// @Failure      400  {object}  types.ErrorResponse
// @Failure      404  {object}  types.ErrorResponse
// @Router       /reviews/{id} [get]
func (h *ReviewHandler) GetReviewByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.Param("id"),
		}).Error("Invalid review ID")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid review ID"})
		return
	}

	review, err := h.reviewService.GetReviewByID(uint(id))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Review not found")
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Review not found"})
		return
	}

	logger.WithFields(logrus.Fields{
		"review_id": review.ID,
	}).Info("Review retrieved successfully")

	c.JSON(http.StatusOK, review)
}

// DeleteReview godoc
// @Summary      Delete a review
// @Description  Delete a review by its ID
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      int  true  "Review ID"
// @Success      204  {object}  types.SuccessResponse
// @Failure      400  {object}  types.ErrorResponse
// @Failure      404  {object}  types.ErrorResponse
// @Failure      500  {object}  types.ErrorResponse
// @Router       /reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := h.reviewService.DeleteReview(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchReviews handles searching for reviews with pagination and filtering
// @Summary Search reviews
// @Description Search reviews with pagination, product name filter, and sorting
// @Tags reviews
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param product_name query string false "Product name to filter by"
// @Param sort_by query string false "Field to sort by (created_at, rating)" default(created_at)
// @Param order query string false "Sort order (asc, desc)" default(desc)
// @Success 200 {object} dto.ReviewListResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security Bearer
// @Router /reviews/ [get]
func (h *ReviewHandler) SearchReviews(c *gin.Context) {
	var req dto.ReviewSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error: "Invalid query parameters",
		})
		return
	}

	reviews, total, err := h.reviewService.SearchReviews(
		req.Page,
		req.PageSize,
		req.ProductName,
		req.SortBy,
		req.Order,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to search reviews",
		})
		return
	}

	// Convert reviews to response format
	items := make([]dto.ReviewResponse, len(reviews))
	for i, review := range reviews {
		items[i] = dto.ReviewResponse{
			ID:        review.ID,
			UserID:    review.UserID,
			ProductID: review.ProductID,
			Rating:    review.Rating,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt.Format(time.RFC3339),
			UpdatedAt: review.UpdatedAt.Format(time.RFC3339),
			User: &dto.UserOutput{
				ID:       review.User.ID,
				Username: review.User.Username,
				Email:    review.User.Email,
				FullName: review.User.FullName,
			},
			Product: &dto.ProductResponse{
				ID:          review.Product.ID,
				Name:        review.Product.Name,
				Description: review.Product.Description,
				Price:       review.Product.Price,
				Quantity:    review.Product.StockQuantity,
				Status:      string(review.Product.Status),
			},
		}
	}

	// Calculate total pages
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	response := dto.ReviewListResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// GetTotalReviews godoc
// @Summary      Get total review count
// @Description  Get the total number of reviews for all products
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Success      200  {object}  types.APIResponse
// @Failure      500  {object}  types.ErrorResponse
// @Security     Bearer
// @Router       /reviews/count [get]
func (h *ReviewHandler) GetTotalReviews(c *gin.Context) {
	// Đếm tổng số review
	count, err := h.reviewService.CountTotalReviews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to count total reviews",
		})
		return
	}

	// Đếm số review của user hiện tại
	userID := c.GetUint("userID")
	myReviewCount, err := h.reviewService.CountReviewsWithUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "Failed to count user's reviews",
		})
		return
	}

	// Trả kết quả dạng JSON
	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "Total reviews retrieved successfully",
		Data: gin.H{
			"total_reviews":   count,
			"my_review_count": myReviewCount,
		},
	})
}
