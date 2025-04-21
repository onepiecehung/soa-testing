package handlers

import (
	"errors"
	"net/http"
	"product-management/internal/dto"
	"product-management/internal/models"
	"product-management/internal/repositories"
	"product-management/internal/services"
	"product-management/internal/types"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	userRepo    *repositories.UserRepository
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo *repositories.UserRepository, authService *services.AuthService) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, authService: authService}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration details"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate passwords match
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	// Basic validation
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.FullName = strings.TrimSpace(req.FullName)

	if req.Username == "" || req.Email == "" || req.FullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, email and full name cannot be empty"})
		return
	}

	// Set default role as user
	userRole := models.RoleUser

	// If role is provided, validate it
	if req.Role != "" {
		switch models.Role(req.Role) {
		case models.RoleUser:
			userRole = models.RoleUser
		case models.RoleAdmin:
			// Here you might want to add additional checks to ensure only authorized users can create admin accounts
			// For example, check if the request comes from an existing admin user
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to create admin account"})
			return
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
			return
		}
	}

	// Create user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Password: req.Password,
		Role:     userRole,
	}

	if err := h.userRepo.Create(user); err != nil {
		if strings.Contains(err.Error(), "username already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		if strings.Contains(err.Error(), "email already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// Create response
	response := dto.RegisterResponse{
		Message: "user registered successfully",
		User: dto.UserOutput{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     string(user.Role),
		},
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      dto.LoginRequest  true  "Login credentials"
// @Success      200         {object}   types.APIResponse
// @Failure      400         {object}   types.ErrorResponse
// @Failure      401         {object}   types.ErrorResponse
// @Failure      500         {object}   types.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	user, accessToken, refreshToken, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: err.Error()})
		return
	}

	// Create user output without sensitive data
	userOutput := dto.UserOutput{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      string(user.Role),
		LastLogin: user.LastLogin,
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data: types.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         userOutput,
		},
	})
}

// GetCurrentUser godoc
// @Summary      Get current user information
// @Description  Get information of the currently logged-in user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  types.APIResponse
// @Failure      401  {object}  types.ErrorResponse
// @Failure      500  {object}  types.ErrorResponse
// @Router       /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      string(user.Role),
		LastLogin: user.LastLogin.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    response,
	})
}

// GetUserByID godoc
// @Summary      Get user information by ID
// @Description  Get information of a user by their ID
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  types.APIResponse
// @Failure      401  {object}  types.ErrorResponse
// @Failure      404  {object}  types.ErrorResponse
// @Failure      500  {object}  types.ErrorResponse
// @Router       /auth/users/{id} [get]
func (h *AuthHandler) GetUserByID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "invalid user ID"})
		return
	}

	user, err := h.userRepo.GetByID(uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      string(user.Role),
		LastLogin: user.LastLogin.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    response,
	})
}

// UpdatePassword godoc
// @Summary      Update user password
// @Description  Update the password of the currently logged-in user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request  body      dto.UpdatePasswordRequest  true  "Password update details"
// @Success      200     {object}   types.SuccessResponse
// @Failure      400     {object}   types.ErrorResponse
// @Failure      401     {object}   types.ErrorResponse
// @Failure      500     {object}   types.ErrorResponse
// @Router       /auth/password [put]
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	userID := c.GetUint("userID")

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.authService.UpdatePassword(user.ID, req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{Message: "password updated successfully"})
}

// UpdateUser godoc
// @Summary      Update user information
// @Description  Update information of the currently logged-in user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request  body      dto.UpdateUserRequest  true  "User update details"
// @Success      200     {object}   types.SuccessResponse
// @Failure      400     {object}   types.ErrorResponse
// @Failure      401     {object}   types.ErrorResponse
// @Failure      500     {object}   types.ErrorResponse
// @Router       /auth/me [put]
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	userID := c.GetUint("userID")

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	// check if username exists
	if req.Username != "" {
		exists, err := h.authService.CheckUserNameExists(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
			return
		}
		if exists {
			c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "username already exists"})
			return
		}
	}
	// check if email exists
	if req.Email != "" {
		exists, err := h.authService.CheckEmailExists(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
			return
		}
		if exists {
			c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "email already exists"})
			return
		}
	}
	if err := h.authService.UpdateUser(user.ID, req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{Message: "user updated successfully"})
}

// ListUsers godoc
// @Summary      List users
// @Description  Get a paginated list of users with search and filter options
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page      query     int     false  "Page number (default: 1)"
// @Param        page_size query     int     false  "Number of items per page (default: 10, max: 100)"
// @Param        search    query     string  false  "Search by username or email"
// @Param        role      query     string  false  "Filter by role (user/admin)"
// @Success      200      {object}   types.APIResponse
// @Failure      400      {object}   types.ErrorResponse
// @Failure      401      {object}   types.ErrorResponse
// @Failure      500      {object}   types.ErrorResponse
// @Router       /auth/users [get]
func (h *AuthHandler) ListUsers(c *gin.Context) {
	var req dto.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	// Set default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// Convert role string to models.Role
	var role models.Role
	if req.Role != "" {
		role = models.Role(req.Role)
	}

	users, total, err := h.userRepo.ListUsers(req.Page, req.PageSize, req.Search, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	// Convert users to response format
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			Role:      string(user.Role),
			LastLogin: user.LastLogin.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data: types.PaginatedResponse{
			Items:      userResponses,
			Total:      total,
			Page:       req.Page,
			PageSize:   req.PageSize,
			TotalPages: int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
		},
	})
}

// UpdateUserRole godoc
// @Summary      Update user role
// @Description  Update the role of a user (only admin can do this)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id      path      int                     true  "User ID"
// @Param        request body      dto.UpdateUserRoleRequest  true  "Role update details"
// @Success      200    {object}   types.SuccessResponse
// @Failure      400    {object}   types.ErrorResponse
// @Failure      401    {object}   types.ErrorResponse
// @Failure      403    {object}   types.ErrorResponse
// @Failure      404    {object}   types.ErrorResponse
// @Failure      500    {object}   types.ErrorResponse
// @Router       /auth/users/{id}/role [put]
func (h *AuthHandler) UpdateUserRole(c *gin.Context) {
	// Get user ID from path
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "invalid user ID"})
		return
	}

	// Get current user
	currentUserID := c.GetUint("userID")
	currentUser, err := h.authService.GetCurrentUser(currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	/*
		Even though a middleware is used to check the user's role before allowing access to the route, we still perform an additional role check within the handler itself.
		This redundancy is intentional and important for security reasons.
		During the execution of a request, there is a possibility that the user's role may change — for example, the user might lose their admin privileges and be downgraded to a regular user. If we rely solely on the role check performed by the middleware at the start of the request, we might miss such changes that occur mid-request.
		By verifying the user's role again in the handler using the most up-to-date information from the database, we ensure that access control remains accurate and consistent, even if the user's role changes during the request lifecycle.
		To avoid this issue, all of the user's active sessions should be revoked immediately after their role is updated.
	*/
	// Check if current user is admin
	if currentUser.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "only admin can update user roles"})
		return
	}

	// Parse request body
	var req dto.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	// Update user role
	if err := h.authService.UpdateUserRole(uint(userID), models.Role(req.Role)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{Message: "user role updated successfully"})
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Soft delete a user (only admin can do this)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  types.SuccessResponse
// @Failure      400  {object}  types.ErrorResponse
// @Failure      401  {object}  types.ErrorResponse
// @Failure      403  {object}  types.ErrorResponse
// @Failure      404  {object}  types.ErrorResponse
// @Failure      500  {object}  types.ErrorResponse
// @Router       /auth/users/{id} [delete]
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	// Get user ID from path
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "invalid user ID"})
		return
	}

	// Get current user
	currentUserID := c.GetUint("userID")
	currentUser, err := h.authService.GetCurrentUser(currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	// Check if current user is admin
	if currentUser.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "only admin can delete users"})
		return
	}

	// Delete user
	if err := h.authService.DeleteUser(uint(userID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "user not found"})
			return
		}
		if err.Error() == "cannot delete admin user" {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{Message: "user deleted successfully"})
}
