package dto

// CreateReviewRequest represents the request body for creating a review
type CreateReviewRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment" binding:"required,min=1,max=500"`
}

// ReviewResponse represents the response for review operations
type ReviewResponse struct {
	ID        uint             `json:"id"`
	ProductID uint             `json:"product_id"`
	UserID    uint             `json:"user_id"`
	Rating    int              `json:"rating"`
	Comment   string           `json:"comment"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
	User      *UserOutput      `json:"user,omitempty"`
	Product   *ProductResponse `json:"product,omitempty"`
}

// ReviewSearchRequest represents the request parameters for searching reviews
type ReviewSearchRequest struct {
	Page        int    `form:"page" binding:"min=1" default:"1"`
	PageSize    int    `form:"page_size" binding:"min=1,max=100" default:"10"`
	ProductName string `form:"product_name"`
	SortBy      string `form:"sort_by" binding:"oneof=created_at rating" default:"created_at"`
	Order       string `form:"order" binding:"oneof=asc desc" default:"desc"`
}

// ReviewListResponse represents the response for a list of reviews
type ReviewListResponse struct {
	Items      []ReviewResponse `json:"items"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}
