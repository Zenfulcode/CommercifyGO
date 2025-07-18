package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	Server          ServerConfig
	Database        DatabaseConfig
	Auth            AuthConfig
	Payment         PaymentConfig
	Email           EmailConfig
	Stripe          StripeConfig
	MobilePay       MobilePayConfig
	CORS            CORSConfig
	DefaultCurrency string // Default currency for the store
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Driver   string // Database driver: "sqlite" or "postgres"
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Debug    string // Silent, Info, Warn, Error
}

// AuthConfig holds authentication-specific configuration
type AuthConfig struct {
	JWTSecret     string
	TokenDuration int
}

// PaymentConfig holds payment-specific configuration
type PaymentConfig struct {
	EnabledProviders []string // List of enabled payment providers
}

// EmailConfig holds email-specific configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	AdminEmail   string
	ContactEmail string // Contact email for customer support
	StoreName    string // Store name used in emails
	Enabled      bool
}

// StripeConfig holds Stripe-specific configuration
type StripeConfig struct {
	SecretKey          string
	PublicKey          string
	WebhookSecret      string
	PaymentDescription string
	ReturnURL          string
	Enabled            bool
}

// MobilePayConfig holds MobilePay-specific configuration
type MobilePayConfig struct {
	MerchantSerialNumber string
	SubscriptionKey      string
	ClientID             string
	ClientSecret         string
	ReturnURL            string
	WebhookURL           string
	PaymentDescription   string
	Enabled              bool
	IsTestMode           bool
}

// CORSConfig holds CORS-specific configuration
type CORSConfig struct {
	AllowedOrigins  []string
	AllowAllOrigins bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	readTimeout, err := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "15"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_READ_TIMEOUT: %w", err)
	}

	writeTimeout, err := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "15"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_WRITE_TIMEOUT: %w", err)
	}

	tokenDuration, err := strconv.Atoi(getEnv("AUTH_TOKEN_DURATION", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid AUTH_TOKEN_DURATION: %w", err)
	}

	smtpPort, err := strconv.Atoi(getEnv("EMAIL_SMTP_PORT", "587"))
	if err != nil {
		return nil, fmt.Errorf("invalid EMAIL_SMTP_PORT: %w", err)
	}

	emailEnabled, err := strconv.ParseBool(getEnv("EMAIL_ENABLED", "false"))
	if err != nil {
		return nil, fmt.Errorf("invalid EMAIL_ENABLED: %w", err)
	}

	stripeEnabled, err := strconv.ParseBool(getEnv("STRIPE_ENABLED", "false"))
	if err != nil {
		return nil, fmt.Errorf("invalid STRIPE_ENABLED: %w", err)
	}

	mobilePayEnabled, err := strconv.ParseBool(getEnv("MOBILEPAY_ENABLED", "false"))
	if err != nil {
		return nil, fmt.Errorf("invalid MOBILEPAY_ENABLED: %w", err)
	}

	mobilePayTestMode, err := strconv.ParseBool(getEnv("MOBILEPAY_TEST_MODE", "true"))
	if err != nil {
		return nil, fmt.Errorf("invalid MOBILEPAY_TEST_MODE: %w", err)
	}

	// Parse enabled payment providers
	enabledProviders := []string{"mock"} // Always enable mock provider for testing
	if stripeEnabled {
		enabledProviders = append(enabledProviders, "stripe")
	}
	if mobilePayEnabled {
		enabledProviders = append(enabledProviders, "mobilepay")
	}

	config := Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "6091"),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "sqlite"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "commercify.db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			Debug:    getEnv("DB_DEBUG", "false"),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("AUTH_JWT_SECRET", "your-secret-key"),
			TokenDuration: tokenDuration,
		},
		Payment: PaymentConfig{
			EnabledProviders: enabledProviders,
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("EMAIL_SMTP_HOST", "smtp.example.com"),
			SMTPPort:     smtpPort,
			SMTPUsername: getEnv("EMAIL_SMTP_USERNAME", ""),
			SMTPPassword: getEnv("EMAIL_SMTP_PASSWORD", ""),
			FromEmail:    getEnv("EMAIL_FROM_ADDRESS", "noreply@example.com"),
			FromName:     getEnv("EMAIL_FROM_NAME", "Commercify Store"),
			AdminEmail:   getEnv("EMAIL_ADMIN_ADDRESS", "admin@example.com"),
			StoreName:    getEnv("STORE_NAME", "Commercify Store"),
			Enabled:      emailEnabled,
		},
		Stripe: StripeConfig{
			SecretKey:          getEnv("STRIPE_SECRET_KEY", ""),
			PublicKey:          getEnv("STRIPE_PUBLIC_KEY", ""),
			WebhookSecret:      getEnv("STRIPE_WEBHOOK_SECRET", ""),
			PaymentDescription: getEnv("STRIPE_PAYMENT_DESCRIPTION", "Commercify Store Purchase"),
			ReturnURL:          getEnv("RETURN_URL", ""),
			Enabled:            stripeEnabled,
		},
		MobilePay: MobilePayConfig{
			MerchantSerialNumber: getEnv("MOBILEPAY_MERCHANT_SERIAL_NUMBER", ""),
			SubscriptionKey:      getEnv("MOBILEPAY_SUBSCRIPTION_KEY", ""),
			ClientID:             getEnv("MOBILEPAY_CLIENT_ID", ""),
			ClientSecret:         getEnv("MOBILEPAY_CLIENT_SECRET", ""),
			ReturnURL:            getEnv("RETURN_URL", ""),
			WebhookURL:           getEnv("MOBILEPAY_WEBHOOK_URL", ""),
			PaymentDescription:   getEnv("MOBILEPAY_PAYMENT_DESCRIPTION", "Commercify Store Purchase"),
			Enabled:              mobilePayEnabled,
			IsTestMode:           mobilePayTestMode,
		},
		CORS: CORSConfig{
			AllowedOrigins:  strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "*"), ","),
			AllowAllOrigins: true,
		},
		DefaultCurrency: getEnv("DEFAULT_CURRENCY", "USD"),
	}

	config.Email.ContactEmail = getEnv("EMAIL_CONTACT_ADDRESS", config.Email.FromEmail)

	return &config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
