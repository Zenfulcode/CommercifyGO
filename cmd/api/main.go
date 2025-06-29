package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/infrastructure/database"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		// Using standard log since logger isn't initialized yet
		log.Println("No .env file found, using environment variables")
	}

	// Initialize logger
	logger := logger.NewLogger()
	logger.Info("Starting Commercify backend service")

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

	// Run database migrations
	if err := database.RunMigrations(db, cfg.Database); err != nil {
		logger.Fatal("Failed to run database migrations: %v", err)
	}

	// Initialize API server
	server := api.NewServer(cfg, db, logger)

	// Start background checkout expiry process
	go startCheckoutExpiryProcess(server, logger)

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server on port %s", cfg.Server.Port)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited properly")
}

// startCheckoutExpiryProcess runs a background process to expire old checkouts
func startCheckoutExpiryProcess(server *api.Server, logger logger.Logger) {
	// Run every 15 minutes
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	// Run immediately on startup
	expireCheckouts(server, logger)

	for range ticker.C {
		expireCheckouts(server, logger)
	}
}

// expireCheckouts expires old checkouts
func expireCheckouts(server *api.Server, logger logger.Logger) {
	checkoutUseCase := server.GetContainer().UseCases().CheckoutUseCase()
	if checkoutUseCase == nil {
		logger.Error("CheckoutUseCase not available")
		return
	}

	amountExpired, err := checkoutUseCase.ExpireOldCheckouts()
	if err != nil {
		logger.Error("Failed to expire old checkouts: %v", err)
	} else {
		logger.Info("Expired %d old checkouts", amountExpired)
	}
}
