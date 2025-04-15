package migrations

import (
	"os"

	"github.com/Napat/mcpserver-demo/models"
	"gorm.io/gorm"
)

type SeedInitialUsers_20250413111743 struct{}

// Name returns the name of the migration
func (m *SeedInitialUsers_20250413111743) Name() string {
	return "20250413111743_seed_initial_users"
}

// Up is the function to upgrade database
func (m *SeedInitialUsers_20250413111743) Up(tx *gorm.DB) error {
	// Run migration in transaction
	return tx.Transaction(func(tx *gorm.DB) error {
		// Seed admin user
		adminEmail := os.Getenv("ADMIN_EMAIL")
		adminPassword := os.Getenv("ADMIN_PASSWORD")

		if adminEmail == "" {
			adminEmail = "admin@example.com"
		}
		if adminPassword == "" {
			adminPassword = "admin123"
		}

		var count int64
		tx.Model(&models.User{}).Where("email = ?", adminEmail).Count(&count)

		if count == 0 {
			adminUser := models.User{
				Email:     adminEmail,
				Password:  adminPassword,
				FirstName: "Admin",
				LastName:  "User",
				Role:      models.RoleAdmin,
				Active:    true,
			}

			if err := tx.Create(&adminUser).Error; err != nil {
				return err
			}
		}

		// Seed test users
		testUsers := []struct {
			Email     string
			Password  string
			FirstName string
			LastName  string
			Role      models.UserRole
		}{
			{
				Email:     "user@example.com",
				Password:  "user123",
				FirstName: "Test",
				LastName:  "User",
				Role:      models.RoleUser,
			},
			{
				Email:     "staff@example.com",
				Password:  "staff123",
				FirstName: "Test",
				LastName:  "Staff",
				Role:      models.RoleStaff,
			},
			{
				Email:     "manager@example.com",
				Password:  "manager123",
				FirstName: "Test",
				LastName:  "Manager",
				Role:      models.RoleManager,
			},
			{
				Email:     "superadmin@example.com",
				Password:  "superadmin123",
				FirstName: "Super",
				LastName:  "Admin",
				Role:      models.RoleSuperAdmin,
			},
		}

		for _, userData := range testUsers {
			var count int64
			tx.Model(&models.User{}).Where("email = ?", userData.Email).Count(&count)

			if count == 0 {
				user := models.User{
					Email:     userData.Email,
					Password:  userData.Password,
					FirstName: userData.FirstName,
					LastName:  userData.LastName,
					Role:      userData.Role,
					Active:    true,
				}

				if err := tx.Create(&user).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// Down is the function to downgrade database
func (m *SeedInitialUsers_20250413111743) Down(tx *gorm.DB) error {
	// Run migration in transaction
	return tx.Transaction(func(tx *gorm.DB) error {
		// Remove test users
		testEmails := []string{
			"user@example.com",
			"staff@example.com",
			"manager@example.com",
			"superadmin@example.com",
		}

		for _, email := range testEmails {
			if err := tx.Where("email = ?", email).Delete(&models.User{}).Error; err != nil {
				return err
			}
		}

		// Remove admin user
		adminEmail := os.Getenv("ADMIN_EMAIL")
		if adminEmail == "" {
			adminEmail = "admin@example.com"
		}

		if err := tx.Where("email = ?", adminEmail).Delete(&models.User{}).Error; err != nil {
			return err
		}

		return nil
	})
}
