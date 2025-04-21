package dto

// UpdatePasswordRequest represents the request body for updating password
type UpdatePasswordRequest struct {
	CurrentPassword    string `json:"current_password" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required,min=6"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required,eqfield=NewPassword"`
}

// UpdateUserRequest represents the request body for updating user information
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3"`
	Email    string `json:"email" binding:"omitempty,email"`
	FullName string `json:"full_name" binding:"omitempty"`
}

// UserResponse represents the response for user information
type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	LastLogin string `json:"last_login"`
}

// ListUsersRequest represents the request parameters for listing users
type ListUsersRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search" binding:"omitempty"`
	Role     string `form:"role" binding:"omitempty,oneof=user admin"`
}

// UpdateUserRoleRequest represents the request body for updating user role
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user admin"`
}
