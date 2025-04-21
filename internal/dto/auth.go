package dto

import "time"

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Email           string `json:"email" binding:"required,email" example:"john@example.com"`
	FullName        string `json:"full_name" binding:"required" example:"John Doe"`
	Password        string `json:"password" binding:"required,min=6" example:"password123"`
	ConfirmPassword string `json:"confirm_password" binding:"required" example:"password123"`
	Role            string `json:"role,omitempty" example:"user" enums:"user,admin"`
}

// RegisterResponse represents the response for successful registration
type RegisterResponse struct {
	Message string     `json:"message" example:"user registered successfully"`
	User    UserOutput `json:"user"`
}

// UserOutput represents the user data to be returned in responses
type UserOutput struct {
	ID        uint      `json:"id" example:"1"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john@example.com"`
	FullName  string    `json:"full_name" example:"John Doe"`
	Role      string    `json:"role" example:"user" enums:"user,admin"`
	LastLogin time.Time `json:"last_login" example:"2021-01-01T00:00:00Z"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}
