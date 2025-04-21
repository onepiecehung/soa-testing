package repositories

import (
	"errors"
	"product-management/internal/models"
	"time"

	"gorm.io/gorm"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	// Check if username already exists
	var count int64
	if err := r.db.Model(&models.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("username already exists")
	}

	// Check if email already exists
	if err := r.db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("email already exists")
	}

	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username, returns nil if not found
func (r *UserRepository) GetByUsername2(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email, returns nil if not found
func (r *UserRepository) GetByEmail2(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll retrieves all users
func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Update fields
func (r *UserRepository) UpdateFields(userID uint, fields map[string]interface{}) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(fields).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// UpdateLastLogin updates the last login time for a user
func (r *UserRepository) UpdateLastLogin(user *models.User) error {
	// Set the current time
	user.LastLogin = time.Now()
	// Update only the LastLogin field
	result := r.db.Model(user).Update("last_login", user.LastLogin)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ListUsers retrieves a paginated list of users with search and filter options
func (r *UserRepository) ListUsers(page, pageSize int, search string, role models.Role) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Build query
	query := r.db.Model(&models.User{})

	// Apply search filter
	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply role filter
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
