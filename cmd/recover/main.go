package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/infrastructure/repository/postgres"
	"github.com/zenfulcode/commercify/internal/infrastructure/service"
)

func main() {
	// Parse command line flags
	flag.Parse()

	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	checkoutRepo := postgres.NewCheckoutRepository(db)

	// Initialize email service
	emailService := service.NewEmailServiceFromEnv()

	// Determine template path
	templatePath := "templates"
	if val := os.Getenv("TEMPLATE_PATH"); val != "" {
		templatePath = val
	}

	// Store configuration
	storeName := os.Getenv("STORE_NAME")
	if storeName == "" {
		storeName = "Commercify"
	}

	storeLogoURL := os.Getenv("STORE_LOGO_URL")
	if storeLogoURL == "" {
		storeLogoURL = "https://example.com/logo.png"
	}

	storeURL := os.Getenv("STORE_URL")
	if storeURL == "" {
		storeURL = "https://example.com"
	}

	privacyPolicyURL := os.Getenv("PRIVACY_POLICY_URL")
	if privacyPolicyURL == "" {
		privacyPolicyURL = "https://example.com/privacy"
	}

	// Initialize checkout recovery use case
	recoveryUseCase := usecase.NewCheckoutRecoveryUseCase(
		checkoutRepo,
		emailService,
		filepath.Join(templatePath, "emails"),
		storeName,
		storeLogoURL,
		storeURL,
		privacyPolicyURL,
	)

	// Process abandoned checkouts
	count, err := recoveryUseCase.ProcessAbandonedCheckouts()
	if err != nil {
		log.Fatalf("Failed to process abandoned checkouts: %v", err)
	}

	fmt.Printf("Successfully processed %d abandoned checkouts\n", count)
}
