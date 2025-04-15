package migrations

import (
	"github.com/Napat/mcpserver-demo/models"
	"gorm.io/gorm"
)

type CreateInitialTables_20250413111742 struct{}

// Name returns the name of the migration
func (m *CreateInitialTables_20250413111742) Name() string {
	return "20250413111742_create_initial_tables"
}

// Up is the function to upgrade database
func (m *CreateInitialTables_20250413111742) Up(tx *gorm.DB) error {
	// Run migration in transaction
	return tx.Transaction(func(tx *gorm.DB) error {
		// Create users table
		if err := tx.AutoMigrate(&models.User{}); err != nil {
			return err
		}

		// Create login_histories table
		if err := tx.AutoMigrate(&models.LoginHistory{}); err != nil {
			return err
		}

		// Create notes table
		if err := tx.AutoMigrate(&models.Note{}); err != nil {
			return err
		}

		return nil
	})
}

// Down is the function to downgrade database
func (m *CreateInitialTables_20250413111742) Down(tx *gorm.DB) error {
	// Run migration in transaction
	return tx.Transaction(func(tx *gorm.DB) error {
		// Drop tables in reverse order
		if err := tx.Migrator().DropTable("notes"); err != nil {
			return err
		}

		if err := tx.Migrator().DropTable("login_histories"); err != nil {
			return err
		}

		if err := tx.Migrator().DropTable("users"); err != nil {
			return err
		}

		return nil
	})
}
