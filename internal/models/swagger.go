package models

// SwaggerProduct represents a product for Swagger documentation
type SwaggerProduct struct {
	ID            uint              `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Price         float64           `json:"price"`
	StockQuantity int               `json:"stock_quantity"`
	Status        ProductStatus     `json:"status"`
	Categories    []SwaggerCategory `json:"categories"`
	Reviews       []SwaggerReview   `json:"reviews"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
}

// SwaggerCategory represents a category for Swagger documentation
type SwaggerCategory struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// SwaggerReview represents a review for Swagger documentation
type SwaggerReview struct {
	ID        uint   `json:"id"`
	ProductID uint   `json:"product_id"`
	UserID    uint   `json:"user_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// SwaggerWishlist represents a wishlist item for Swagger documentation
type SwaggerWishlist struct {
	UserID    uint           `json:"user_id"`
	ProductID uint           `json:"product_id"`
	AddedAt   string         `json:"added_at"`
	Product   SwaggerProduct `json:"product"`
}

// SwaggerResponse represents a standard API response
type SwaggerResponse struct {
	Data  interface{} `json:"data"`
	Meta  interface{} `json:"meta,omitempty"`
	Error string      `json:"error,omitempty"`
}
