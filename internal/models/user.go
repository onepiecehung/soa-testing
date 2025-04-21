package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Role represents user roles in the system
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// User represents a user in the system
type User struct {
	BaseModel
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	FullName  string    `json:"full_name"`
	Password  string    `json:"-" gorm:"not null"` // "-" means this field won't be included in JSON
	Role      Role      `json:"role" gorm:"type:varchar(10);default:'user'"`
	LastLogin time.Time `json:"last_login"`
	Reviews   []Review  `json:"reviews"` // One-to-many relationship with Review
}

// BeforeSave is a GORM hook that hashes the password before saving
func (u *User) BeforeSave(tx *gorm.DB) error {
	// Only hash the password if it has been changed
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// ValidatePassword checks if the provided password matches the stored hash
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
