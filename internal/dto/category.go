package dto

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CategoryResponse represents the response for category operations
type CategoryResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ProductCount int    `json:"product_count"`
}

// CategoryDistributionResponse represents the distribution of products across categories
type CategoryDistributionResponse struct {
	Name         string `json:"name"`
	ProductCount int    `json:"product_count"`
}
