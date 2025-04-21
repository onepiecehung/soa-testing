package dto

// CreateProductRequest represents the request body for creating a new product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required" example:"SmartWatch Pro"`    // Product name
	Description string  `json:"description" example:"Advanced smartwatch"`           // Product description
	Price       float64 `json:"price" binding:"required,gt=0" example:"299.99"`      // Product price
	Quantity    int     `json:"quantity" binding:"required,gte=0" example:"100"`     // Stock quantity
	Categories  []uint  `json:"categories" binding:"required,min=1" example:"1,2,3"` // Category IDs
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"required" example:"SmartWatch Pro 2"`                     // Product name
	Description string  `json:"description" example:"Updated smartwatch features"`                      // Product description
	Price       float64 `json:"price" binding:"required,gt=0" example:"349.99"`                         // Product price
	Quantity    int     `json:"quantity" binding:"required,gte=0" example:"150"`                        // Stock quantity
	Categories  []uint  `json:"categories" binding:"required,min=1" example:"1,2,3"`                    // Category IDs
	Status      string  `json:"status" binding:"required,oneof=active inactive draft" example:"active"` // Product status
}

// ProductResponse represents the response for product operations
type ProductResponse struct {
	ID          uint             `json:"id" example:"1"`                            // Product ID
	Name        string           `json:"name" example:"SmartWatch Pro"`             // Product name
	Description string           `json:"description" example:"Advanced smartwatch"` // Product description
	Price       float64          `json:"price" example:"299.99"`                    // Product price
	Quantity    int              `json:"quantity" example:"100"`                    // Stock quantity
	Status      string           `json:"status" example:"active"`                   // Product status
	Categories  []CategoryOutput `json:"categories"`                                // Associated categories
}

// CategoryOutput represents the category data in product responses
type CategoryOutput struct {
	ID   uint   `json:"id" example:"1"`             // Category ID
	Name string `json:"name" example:"Electronics"` // Category name
}

// ProductListResponse represents the response for listing products
type ProductListResponse struct {
	Products []ProductResponse `json:"products"`               // List of products
	Total    int64             `json:"total" example:"100"`    // Total number of products
	Page     int               `json:"page" example:"1"`       // Current page number
	PageSize int               `json:"page_size" example:"10"` // Number of items per page
}

// ProductSearchRequest represents the request for searching products
type ProductSearchRequest struct {
	Search     string   `form:"search"`               // Search query
	CategoryID uint     `form:"category"`             // Filter by category ID
	Statuses   []string `form:"status"`               // Filter by statuses
	Sort       string   `form:"sort"`                 // Sort field
	Page       int      `form:"page,default=1"`       // Page number
	PageSize   int      `form:"page_size,default=10"` // Items per page
}
