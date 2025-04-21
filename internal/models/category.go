package models

// Category represents a product category
type Category struct {
	BaseModel
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Products    []Product `gorm:"many2many:product_categories;" json:"products"`
}

// TableName specifies the table name for the Category model
func (Category) TableName() string {
	return "categories"
}
