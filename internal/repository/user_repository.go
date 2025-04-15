package repository

import (
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/Napat/mcpserver-demo/models"
	"github.com/Napat/mcpserver-demo/pkg/storage"
	"gorm.io/gorm"
)

//go:generate mockgen -source=./user_repository.go -destination=./mocks/mock_user_repository.go -package=mocks

// IUserRepository is an interface for managing all user data (Facade Pattern)
type IUserRepository interface {
	// Database operations
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	GetLoginHistory(userID uint, limit int) ([]models.LoginHistory, error)
	RecordLogin(history *models.LoginHistory) error

	// File storage operations combined with database
	UpdateProfileImage(userID uint, file *multipart.FileHeader) (string, error)
	DeleteProfileImage(userID uint) error
}

// UserRepository implements IUserRepository following the Facade Pattern
// combining access to both database and file storage
type UserRepository struct {
	db          *gorm.DB
	fileStorage storage.IFileStorage
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB, fileStorage storage.IFileStorage) IUserRepository {
	return &UserRepository{
		db:          db,
		fileStorage: fileStorage,
	}
}

// Create adds a new user to the database
func (r *UserRepository) Create(user *models.User) error {
	// Check if email already exists
	var count int64
	r.db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		return errors.New("email already exists")
	}

	// GORM will automatically call BeforeSave hook which will hash the password
	return r.db.Create(user).Error
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// Update updates user information
func (r *UserRepository) Update(user *models.User) error {
	// Don't update critical fields like password, email, role
	return r.db.Model(user).
		Select("first_name", "last_name", "gender", "birth_date", "active").
		Updates(user).Error
}

// UpdateProfileImage updates profile image, handling both file storage and database
func (r *UserRepository) UpdateProfileImage(userID uint, file *multipart.FileHeader) (string, error) {
	// Find user first to get the existing profile image URL (if any)
	user, err := r.FindByID(userID)
	if err != nil {
		return "", err
	}

	// สร้างชื่อไฟล์ที่ไม่ซ้ำด้วยเวลาปัจจุบันและชื่อไฟล์
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

	// ไม่ต้องมีคำนำหน้า "profiles/"
	objectName := fileName

	// Upload file to storage
	imageURL, err := r.fileStorage.UploadFile("profiles", objectName, file)
	if err != nil {
		return "", err
	}

	// Update URL in database
	err = r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("profile_image_url", imageURL).Error

	if err != nil {
		// If database update fails, delete the uploaded file
		_ = r.fileStorage.DeleteFile("profiles", objectName)
		return "", err
	}

	// If user already had a profile image, delete the old one
	if user.ProfileImageURL != "" && user.ProfileImageURL != imageURL {
		// Extract only the filename from URL
		oldObjectName := extractObjectNameFromURL(user.ProfileImageURL)
		if oldObjectName != "" {
			_ = r.fileStorage.DeleteFile("profiles", oldObjectName)
		}
	}

	return imageURL, nil
}

// DeleteProfileImage deletes a user's profile image
func (r *UserRepository) DeleteProfileImage(userID uint) error {
	// Find user first
	user, err := r.FindByID(userID)
	if err != nil {
		return err
	}

	// If no profile image, do nothing
	if user.ProfileImageURL == "" {
		return nil
	}

	// Extract only the filename from URL
	objectName := extractObjectNameFromURL(user.ProfileImageURL)
	if objectName == "" {
		return errors.New("invalid profile image URL format")
	}

	// Delete file from storage
	err = r.fileStorage.DeleteFile("profiles", objectName)
	if err != nil {
		return err
	}

	// Update database
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("profile_image_url", "").Error
}

// Delete removes a user
func (r *UserRepository) Delete(id uint) error {
	// Find user first
	user, err := r.FindByID(id)
	if err != nil {
		return err
	}

	// If user has a profile image, delete it
	if user.ProfileImageURL != "" {
		// Extract only the filename from URL
		objectName := extractObjectNameFromURL(user.ProfileImageURL)
		if objectName != "" {
			_ = r.fileStorage.DeleteFile("profiles", objectName)
		}
	}

	// Delete user data
	return r.db.Delete(&models.User{}, id).Error
}

// GetLoginHistory retrieves login history
func (r *UserRepository) GetLoginHistory(userID uint, limit int) ([]models.LoginHistory, error) {
	var histories []models.LoginHistory
	result := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories)
	return histories, result.Error
}

// RecordLogin records login history
func (r *UserRepository) RecordLogin(history *models.LoginHistory) error {
	return r.db.Create(history).Error
}

// extractObjectNameFromURL extracts object name from URL
// Example: http://localhost:9000/profiles/avatar.jpg -> avatar.jpg
func extractObjectNameFromURL(url string) string {
	// This is a simple example, you may need to adjust for your URL format
	// In this example, assume URL format is: http://host:port/bucket/file
	if url == "" {
		return ""
	}

	// Find position of last "/"
	lastSlashIndex := -1
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '/' {
			lastSlashIndex = i
			break
		}
	}

	if lastSlashIndex == -1 || lastSlashIndex >= len(url)-1 {
		return ""
	}

	// Extract filename
	fileName := url[lastSlashIndex+1:]

	// Return only the filename, ไม่ต้องเพิ่ม "profiles/" นำหน้า
	return fileName
}
