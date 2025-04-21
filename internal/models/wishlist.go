package models

import (
	"time"
)

// Wishlist represents a user's wishlist item
type Wishlist struct {
	BaseModel
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
	AddedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"added_at"`
}

// TableName specifies the table name for the Wishlist model
func (Wishlist) TableName() string {
	return "wishlists"
}
