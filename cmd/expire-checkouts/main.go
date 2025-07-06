package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/infrastructure/container"
	"github.com/zenfulcode/commercify/internal/infrastructure/database"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

func main() {
	// Parse command line flags
	forceDelete := flag.Bool("force", false, "Force delete all expired, abandoned, and old completed checkouts")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize logger
	logger := logger.NewLogger()
	if *forceDelete {
		logger.Info("Starting checkout force deletion tool")
	} else {
		logger.Info("Starting checkout expiry cleanup tool")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Initialize dependency container
	diContainer := container.NewContainer(cfg, db, logger)

	// Get checkout use case
	checkoutUseCase := diContainer.UseCases().CheckoutUseCase()

	if *forceDelete {
		// Force delete all expired checkouts
		deleteResult, err := checkoutUseCase.ForceDeleteAllExpiredCheckouts()
		if err != nil {
			logger.Fatal("Failed to force delete expired checkouts: %v", err)
		}

		logger.Info("Force deletion completed:")
		logger.Info("- Force deleted checkouts: %d", deleteResult.DeletedCount)
		logger.Info("Total processed: %d", deleteResult.DeletedCount)
	} else {
		// Regular expire old checkouts
		expireResult, err := checkoutUseCase.ExpireOldCheckouts()
		if err != nil {
			logger.Fatal("Failed to expire old checkouts: %v", err)
		}

		logger.Info("Checkout cleanup completed:")
		logger.Info("- Abandoned checkouts: %d", expireResult.AbandonedCount)
		logger.Info("- Deleted checkouts: %d", expireResult.DeletedCount)
		logger.Info("- Expired checkouts: %d", expireResult.ExpiredCount)
		logger.Info("Total processed: %d", expireResult.AbandonedCount+expireResult.DeletedCount+expireResult.ExpiredCount)
	}
}
