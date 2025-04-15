package migrations

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// IMigration is an interface for all migrations
type IMigration interface {
	Name() string
	Up(tx *gorm.DB) error
	Down(tx *gorm.DB) error
}

// MigrationRecord is a model for storing migration history
type MigrationRecord struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null;uniqueIndex"`
	AppliedAt time.Time `gorm:"not null"`
}

// Registry stores all migrations
type Registry struct {
	migrations []IMigration
}

// NewRegistry creates a new registry and registers all migrations
func NewRegistry() *Registry {
	registry := &Registry{
		migrations: []IMigration{},
	}

	// Register migrations in chronological order
	registry.Register(
		&CreateInitialTables_20250413111742{},
		&SeedInitialUsers_20250413111743{},
	)

	return registry
}

// Register registers migrations
func (r *Registry) Register(migrations ...IMigration) {
	r.migrations = append(r.migrations, migrations...)
}

// GetMigrations returns all migrations
func (r *Registry) GetMigrations() []IMigration {
	return r.migrations
}

// CreateMigrationTable creates a table for storing migration history
func CreateMigrationTable(db *gorm.DB) error {
	return db.AutoMigrate(&MigrationRecord{})
}

// RunMigrations migrates database up
func RunMigrations(db *gorm.DB) error {
	// Create migration table if it doesn't exist
	if err := CreateMigrationTable(db); err != nil {
		return err
	}

	registry := NewRegistry()
	migrations := registry.GetMigrations()

	for _, migration := range migrations {
		var record MigrationRecord
		result := db.Where("name = ?", migration.Name()).First(&record)

		// If not yet migrated, perform migration
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				fmt.Printf("Running migration: %s\n", migration.Name())

				// Perform migration
				if err := migration.Up(db); err != nil {
					return fmt.Errorf("error running migration %s: %w", migration.Name(), err)
				}

				// Record migration history
				db.Create(&MigrationRecord{
					Name:      migration.Name(),
					AppliedAt: time.Now(),
				})

				fmt.Printf("Migration applied: %s\n", migration.Name())
			} else {
				return result.Error
			}
		}
	}

	return nil
}

// RollbackMigration rolls back the latest migration
func RollbackMigration(db *gorm.DB) error {
	// Create migration table if it doesn't exist
	if err := CreateMigrationTable(db); err != nil {
		return err
	}

	registry := NewRegistry()
	migrations := registry.GetMigrations()

	// Find the latest applied migration
	var lastRecord MigrationRecord
	result := db.Order("applied_at DESC").First(&lastRecord)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("No migrations to rollback")
			return nil
		}
		return result.Error
	}

	// Find migration that matches the name
	for _, migration := range migrations {
		if migration.Name() == lastRecord.Name {
			fmt.Printf("Rolling back migration: %s\n", migration.Name())

			// Perform rollback
			if err := migration.Down(db); err != nil {
				return fmt.Errorf("error rolling back migration %s: %w", migration.Name(), err)
			}

			// Delete migration history
			db.Delete(&lastRecord)

			fmt.Printf("Migration rolled back: %s\n", migration.Name())
			return nil
		}
	}

	return fmt.Errorf("migration not found: %s", lastRecord.Name)
}

// RollbackAllMigrations rolls back all migrations
func RollbackAllMigrations(db *gorm.DB) error {
	// Create migration table if it doesn't exist
	if err := CreateMigrationTable(db); err != nil {
		return err
	}

	registry := NewRegistry()
	migrations := registry.GetMigrations()

	// Map migrations by name for faster lookup
	migrationMap := make(map[string]IMigration)
	for _, migration := range migrations {
		migrationMap[migration.Name()] = migration
	}

	// Find all applied migrations ordered by most recent
	var records []MigrationRecord
	result := db.Order("applied_at DESC").Find(&records)

	if result.Error != nil {
		return result.Error
	}

	if len(records) == 0 {
		fmt.Println("No migrations to rollback")
		return nil
	}

	// Rollback all migrations in reverse order
	for _, record := range records {
		migration, exists := migrationMap[record.Name]
		if !exists {
			return fmt.Errorf("migration not found: %s", record.Name)
		}

		fmt.Printf("Rolling back migration: %s\n", migration.Name())

		// Perform rollback
		if err := migration.Down(db); err != nil {
			return fmt.Errorf("error rolling back migration %s: %w", migration.Name(), err)
		}

		// Delete migration history
		db.Delete(&record)

		fmt.Printf("Migration rolled back: %s\n", migration.Name())
	}

	return nil
}
