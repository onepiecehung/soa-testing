package services

import (
	"errors"
	"log"
	"time"

	"product-management/internal/dto"
	"product-management/internal/models"
	"product-management/internal/repositories"
	"product-management/pkg/database"
	"product-management/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repositories.NewUserRepository(database.DB),
	}
}

// Login authenticates a user and returns JWT tokens
func (s *AuthService) Login(req dto.LoginRequest) (*models.User, string, string, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, "", "", err
	}
	// update last login
	if err = s.userRepo.UpdateLastLogin(user); err != nil {
		// Log the error but continue with login
		log.Printf("Failed to update last login time for user %d: %v", user.ID, err)
	}

	return user, accessToken, refreshToken, nil
}

// generateAccessToken creates a new JWT access token
func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(utils.GetEnv("JWT_SECRET", "your-secret-key")))
}

// generateRefreshToken creates a new JWT refresh token
func (s *AuthService) generateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(utils.GetEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key")))
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(utils.GetEnv("JWT_SECRET", "your-secret-key")), nil
	})
}

// ValidateRefreshToken validates a refresh token
func (s *AuthService) ValidateRefreshToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(utils.GetEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key")), nil
	})
}

// GetCurrentUser returns the current user from the token
func (s *AuthService) GetCurrentUser(userID uint) (*models.User, error) {
	return s.userRepo.GetByID(uint(userID))
}

// UpdatePassword updates a user's password
func (s *AuthService) UpdatePassword(userID uint, req dto.UpdatePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password - we have BeforeSave hook in User model to hash the password
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	// if err != nil {
	// 	return err
	// }

	user.Password = string(req.NewPassword)
	return s.userRepo.Update(user)
}

// UpdateUser updates a user's information
func (s *AuthService) UpdateUser(userID uint, req dto.UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	updateFields := make(map[string]interface{})
	// Update fields if provided
	if req.Username != "" {
		updateFields["username"] = req.Username
	}
	if req.Email != "" {
		updateFields["email"] = req.Email
	}
	if req.FullName != "" {
		updateFields["full_name"] = req.FullName
	}

	if len(updateFields) == 0 {
		return nil
	}
	return s.userRepo.UpdateFields(user.ID, updateFields)
}

// CheckUserNameExists checks if a username exists
func (s *AuthService) CheckUserNameExists(username string) (bool, error) {
	user, err := s.userRepo.GetByUsername2(username)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

// CheckEmailExists checks if an email exists
func (s *AuthService) CheckEmailExists(email string) (bool, error) {
	user, err := s.userRepo.GetByEmail2(email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

// UpdateUserRole updates a user's role
func (s *AuthService) UpdateUserRole(userID uint, role models.Role) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Update only the role field
	return s.userRepo.UpdateFields(user.ID, map[string]interface{}{
		"role": role,
	})
}

// DeleteUser performs a soft delete on a user
func (s *AuthService) DeleteUser(userID uint) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Don't allow deleting admin users
	if user.Role == models.RoleAdmin {
		return errors.New("cannot delete admin user")
	}

	return s.userRepo.Delete(userID)
}
