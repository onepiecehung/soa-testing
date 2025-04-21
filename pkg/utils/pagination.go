package utils

import (
	"strconv"
)

// PaginationParams represents the pagination parameters from request
type PaginationParams struct {
	Page  int
	Limit int
}

// ParsePaginationParams parses pagination parameters from request query
// Returns default values if parameters are invalid or not provided
func ParsePaginationParams(pageStr, limitStr string) PaginationParams {
	// Parse page number
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Parse limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// PaginationResponse represents the pagination response structure
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int64       `json:"total_pages"`
}

// NewPaginationResponse creates a new pagination response
func NewPaginationResponse(data interface{}, total int64, page, limit int) PaginationResponse {
	// Calculate total pages
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// CalculateOffset calculates the offset for pagination
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}
