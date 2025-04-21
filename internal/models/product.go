package models

// ProductStatus represents the possible statuses of a product
type ProductStatus string

const (
	StatusActive   ProductStatus = "active"
	StatusInactive ProductStatus = "inactive"
	StatusDraft    ProductStatus = "draft"
)

// Product represents a product in the store
type Product struct {
	BaseModel
	Name          string        `gorm:"not null" json:"name"`
	Description   string        `json:"description"`
	Price         float64       `gorm:"not null" json:"price"`
	StockQuantity int           `gorm:"not null;default:0" json:"stock_quantity"`
	Status        ProductStatus `gorm:"default:active" json:"status"`
	Reviews       []Review      `json:"reviews"`
	Categories    []Category    `gorm:"many2many:product_categories;" json:"categories"`
	Wishlists     []Wishlist    `json:"wishlists"`
}

// TableName specifies the table name for the Product model
func (Product) TableName() string {
	return "products"
}
