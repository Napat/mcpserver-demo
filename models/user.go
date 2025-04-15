package models

import (
	"errors"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRole is a union type for user roles
type UserRole uint8

const (
	// RoleUser is a regular user role
	RoleUser UserRole = 1 << iota
	// RoleStaff is a staff role
	RoleStaff
	// RoleManager is a manager role
	RoleManager
	// RoleAdmin is an administrator role
	RoleAdmin
	// RoleSuperAdmin is a super administrator role
	RoleSuperAdmin
)

// RoleNames maps UserRole values to role names
var RoleNames = map[UserRole]string{
	RoleUser:       "User",
	RoleStaff:      "Staff",
	RoleManager:    "Manager",
	RoleAdmin:      "Admin",
	RoleSuperAdmin: "Super Admin",
}

// HasRole checks if a role has the specified role
func (role UserRole) HasRole(r UserRole) bool {
	return (role & r) != 0
}

// AddRole adds a role
func (role *UserRole) AddRole(r UserRole) {
	*role |= r
}

// RemoveRole removes a role
func (role *UserRole) RemoveRole(r UserRole) {
	*role &= ^r
}

// GetRoleNames returns the user's role names as a slice
func (role UserRole) GetRoleNames() []string {
	var roles []string
	for r, name := range RoleNames {
		if role.HasRole(r) {
			roles = append(roles, name)
		}
	}
	return roles
}

// LoginHistory stores user login history
type LoginHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`
	LoginTime time.Time `gorm:"not null" json:"login_time"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName defines the table name
func (LoginHistory) TableName() string {
	return "login_histories"
}

// User represents a user in the system
type User struct {
	ID              uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Email           string     `gorm:"column:email;uniqueIndex;type:varchar(255);not null" json:"email"`
	Password        string     `gorm:"column:password;type:varchar(255);not null" json:"-"`
	FirstName       string     `gorm:"column:first_name;type:varchar(100);not null" json:"first_name"`
	LastName        string     `gorm:"column:last_name;type:varchar(100);not null" json:"last_name"`
	Role            UserRole   `gorm:"column:role;index;type:smallint;not null" json:"role"`
	Active          bool       `gorm:"column:active;type:boolean;not null;default:true" json:"active"`
	Gender          string     `gorm:"column:gender;type:varchar(10)" json:"gender"`
	ProfileImageURL string     `gorm:"column:profile_image_url;type:varchar(255)" json:"profile_image_url"`
	LastLoginTime   *time.Time `gorm:"column:last_login_time;index;type:timestamp" json:"last_login_time"`
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName defines the table name
func (User) TableName() string {
	return "users"
}

// VerifyPassword checks password validity
func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// BeforeSave runs before saving the data
func (u *User) BeforeSave(tx *gorm.DB) error {
	// If the password is not yet hashed, hash it
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role.HasRole(RoleAdmin) || u.Role.HasRole(RoleSuperAdmin)
}

// IsSuperAdmin checks if the user is a super admin
func (u *User) IsSuperAdmin() bool {
	return u.Role.HasRole(RoleSuperAdmin)
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Active
}

// AfterFind runs after retrieving the data
func (u *User) AfterFind(tx *gorm.DB) error {
	if !u.Active {
		return errors.New("user is inactive")
	}
	return nil
}

// GetLoginHistoryLimit retrieves the limit for login history storage from .env
func GetLoginHistoryLimit() int {
	limitStr := os.Getenv("LOGIN_HISTORY_LIMIT")
	if limitStr == "" {
		return 10 // default value
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return 10
	}

	return limit
}

// RecordLogin records login history
func (u *User) RecordLogin(db *gorm.DB, ipAddress, userAgent string) error {
	now := time.Now()
	u.LastLoginTime = &now

	// Create new login history
	loginHistory := LoginHistory{
		UserID:    u.ID,
		LoginTime: now,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: now,
	}

	// Save the latest login history
	if err := db.Create(&loginHistory).Error; err != nil {
		return err
	}

	// Update user's LastLoginTime
	if err := db.Model(u).Update("last_login_time", now).Error; err != nil {
		return err
	}

	// Delete old history, keeping only as many as specified in .env
	limit := GetLoginHistoryLimit()
	var count int64
	db.Model(&LoginHistory{}).Where("user_id = ?", u.ID).Count(&count)

	if count > int64(limit) {
		var oldestLoginIDs []uint
		db.Model(&LoginHistory{}).
			Where("user_id = ?", u.ID).
			Order("login_time ASC").
			Limit(int(count)-limit).
			Pluck("id", &oldestLoginIDs)

		if len(oldestLoginIDs) > 0 {
			db.Delete(&LoginHistory{}, oldestLoginIDs)
		}
	}

	return nil
}

// GetRecentLoginHistory retrieves recent login history
func (u *User) GetRecentLoginHistory(db *gorm.DB, limit int) ([]LoginHistory, error) {
	if limit <= 0 {
		limit = GetLoginHistoryLimit()
	}

	var loginHistories []LoginHistory
	err := db.Where("user_id = ?", u.ID).
		Order("login_time DESC").
		Limit(limit).
		Find(&loginHistories).Error

	return loginHistories, err
}
