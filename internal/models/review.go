package models

// Review represents a product review
type Review struct {
	BaseModel
	ProductID uint    `gorm:"not null" json:"product_id"`
	UserID    uint    `gorm:"not null" json:"user_id"`
	Rating    int     `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment   string  `json:"comment"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	User      User    `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for the Review model
func (Review) TableName() string {
	return "reviews"
}
