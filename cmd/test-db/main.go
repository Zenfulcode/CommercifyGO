package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/infrastructure/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Testing database connection with driver: %s\n", cfg.Database.Driver)

	// Initialize database
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	fmt.Println("Database connection successful!")

	// Test basic database operations
	var result string
	switch cfg.Database.Driver {
	case "sqlite":
		if err := db.Raw("SELECT sqlite_version()").Scan(&result).Error; err != nil {
			log.Fatalf("Failed to query SQLite version: %v", err)
		}
		fmt.Printf("SQLite version: %s\n", result)
	case "postgres":
		if err := db.Raw("SELECT version()").Scan(&result).Error; err != nil {
			log.Fatalf("Failed to query PostgreSQL version: %v", err)
		}
		fmt.Printf("PostgreSQL version: %s\n", result)
	}

	// Close database connection
	if err := database.Close(db); err != nil {
		log.Printf("Warning: Failed to close database: %v", err)
	}

	fmt.Println("Database test completed successfully!")
}
