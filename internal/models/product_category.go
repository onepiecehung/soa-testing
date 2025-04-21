package models

import (
	"time"
)

// ProductCategory represents the many-to-many relationship between products and categories
type ProductCategory struct {
	ProductID  uint      `gorm:"primaryKey;onDelete:CASCADE"`
	CategoryID uint      `gorm:"primaryKey;onDelete:CASCADE"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Product    Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Category   Category  `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
}
