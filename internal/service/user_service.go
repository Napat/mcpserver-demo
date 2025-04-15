package service

import (
	"mime/multipart"
	"time"

	"github.com/Napat/mcpserver-demo/internal/repository"
	"github.com/Napat/mcpserver-demo/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=./user_service.go -destination=./mocks/mock_user_service.go -package=mocks

// IUserService interface for managing user business logic
type IUserService interface {
	Register(user *models.User) error
	Login(email, password string) (*models.User, error)
	UpdateProfile(user *models.User) error
	UpdateProfileImage(userID uint, file *multipart.FileHeader) (string, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetLoginHistory(userID uint, limit int) ([]models.LoginHistory, error)
	RecordLogin(userID uint, ipAddress, userAgent string) error
}

// UserService struct for handling user business logic
type UserService struct {
	userRepo repository.IUserRepository
	logger   *zap.Logger
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.IUserRepository, logger *zap.Logger) IUserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// Register registers a new user
func (s *UserService) Register(user *models.User) error {
	return s.userRepo.Create(user)
}

// Login authenticates a user and returns user information
func (s *UserService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := user.VerifyPassword(password); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateProfile updates user profile information
func (s *UserService) UpdateProfile(user *models.User) error {
	return s.userRepo.Update(user)
}

// UpdateProfileImage updates user profile image
func (s *UserService) UpdateProfileImage(userID uint, file *multipart.FileHeader) (string, error) {
	// Use Repository for uploading and managing image files
	return s.userRepo.UpdateProfileImage(userID, file)
}

// GetUserByID retrieves user by ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// GetUserByEmail retrieves user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

// GetLoginHistory retrieves login history
func (s *UserService) GetLoginHistory(userID uint, limit int) ([]models.LoginHistory, error) {
	return s.userRepo.GetLoginHistory(userID, limit)
}

// RecordLogin records a login event
func (s *UserService) RecordLogin(userID uint, ipAddress, userAgent string) error {
	now := time.Now()
	history := &models.LoginHistory{
		UserID:    uint64(userID),
		LoginTime: now,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: now,
	}
	return s.userRepo.RecordLogin(history)
}
