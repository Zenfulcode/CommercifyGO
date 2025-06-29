package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/infrastructure/container"
	"github.com/zenfulcode/commercify/internal/infrastructure/database"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize logger
	logger := logger.NewLogger()
	logger.Info("Starting checkout expiry cleanup tool")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize dependency container
	diContainer := container.NewContainer(cfg, db, logger)

	// Get checkout use case
	checkoutUseCase := diContainer.UseCases().CheckoutUseCase()

	// Expire old checkouts
	result, err := checkoutUseCase.ExpireOldCheckouts()
	if err != nil {
		logger.Fatal("Failed to expire old checkouts: %v", err)
	}

	logger.Info("Checkout cleanup completed:")
	logger.Info("- Abandoned checkouts: %d", result.AbandonedCount)
	logger.Info("- Deleted checkouts: %d", result.DeletedCount)
	logger.Info("- Expired checkouts: %d", result.ExpiredCount)
	logger.Info("Total processed: %d", result.AbandonedCount+result.DeletedCount+result.ExpiredCount)
}
