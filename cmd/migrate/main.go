package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Napat/mcpserver-demo/internal/migrations"
	"github.com/Napat/mcpserver-demo/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	// Define program commands
	rollbackFlag := flag.Bool("rollback", false, "Rollback the last migration")
	rollbackAllFlag := flag.Bool("rollback-all", false, "Rollback all migrations")
	helpFlag := flag.Bool("help", false, "Show usage information")
	flag.Parse()

	// Show usage instructions
	if *helpFlag {
		fmt.Println("Usage: migrate [options]")
		fmt.Println("Options:")
		fmt.Println("  -rollback       Rollback the last migration")
		fmt.Println("  -rollback-all   Rollback all migrations")
		fmt.Println("  -help           Show this help message")
		return
	}

	// Load .env file
	err := godotenv.Load("configs/temp/.env")
	if err != nil {
		log.Printf("Warning: .env file not found or invalid: %v", err)
	}

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Execute according to command
	if *rollbackAllFlag {
		fmt.Println("Rolling back all migrations...")
		if err := migrations.RollbackAllMigrations(db); err != nil {
			log.Fatalf("Failed to rollback all migrations: %v", err)
		}
		fmt.Println("All migrations have been rolled back successfully")
	} else if *rollbackFlag {
		fmt.Println("Rolling back the last migration...")
		if err := migrations.RollbackMigration(db); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		fmt.Println("Last migration has been rolled back successfully")
	} else {
		fmt.Println("Running migrations...")
		if err := migrations.RunMigrations(db); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")
	}

	os.Exit(0)
}
