package types

import "product-management/internal/models"

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`           // Whether the request was successful
	Message string      `json:"message,omitempty"` // Optional message
	Error   string      `json:"error,omitempty"`   // Error message if success is false
	Data    interface{} `json:"data,omitempty"`    // Response data
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Items      interface{} `json:"items"`       // List of items
	Total      int64       `json:"total"`       // Total number of items
	Page       int         `json:"page"`        // Current page number
	PageSize   int         `json:"page_size"`   // Number of items per page
	TotalPages int         `json:"total_pages"` // Total number of pages
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(items interface{}, total int64, page, pageSize int) PaginatedResponse {
	totalPages := (int(total) + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	return PaginatedResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error       string `json:"error"`                 // Error message
	Code        string `json:"code,omitempty"`        // Error code for client handling
	Description string `json:"description,omitempty"` // Detailed error description
}

// SuccessResponse represents a success response with a message
type SuccessResponse struct {
	Message string `json:"message"` // Success message
}

// ProductListResponse represents a paginated list of products
type ProductListResponse struct {
	PaginatedResponse
	Items []models.Product `json:"items"` // Override Items with specific type
}

// WishlistResponse represents a paginated list of wishlist items
type WishlistResponse struct {
	PaginatedResponse
	Items []models.Wishlist `json:"items"` // Override Items with specific type
}

// NewProductListResponse creates a new product list response
func NewProductListResponse(products []models.Product, total int64, page, pageSize int) ProductListResponse {
	return ProductListResponse{
		PaginatedResponse: NewPaginatedResponse(products, total, page, pageSize),
		Items:             products,
	}
}

// NewWishlistResponse creates a new wishlist response
func NewWishlistResponse(wishlist []models.Wishlist, total int64, page, pageSize int) WishlistResponse {
	return WishlistResponse{
		PaginatedResponse: NewPaginatedResponse(wishlist, total, page, pageSize),
		Items:             wishlist,
	}
}

// CategoryDistributionResponse represents the response for category distribution
type CategoryDistributionResponse struct {
	Name         string `json:"name"`
	ProductCount int64  `json:"product_count"`
}

// LoginResponse represents the response for login
type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user"`
}
