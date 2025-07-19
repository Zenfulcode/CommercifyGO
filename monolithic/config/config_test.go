package config

import (
	"os"
	"testing"
)

func TestEmailConfigDefaults(t *testing.T) {
	// Clear any existing environment variables
	envVars := []string{
		"EMAIL_CONTACT_ADDRESS",
		"STORE_NAME",
		"EMAIL_FROM_ADDRESS",
		"EMAIL_FROM_NAME",
	}

	originalValues := make(map[string]string)
	for _, envVar := range envVars {
		originalValues[envVar] = os.Getenv(envVar)
		os.Unsetenv(envVar)
	}

	// Restore environment variables after test
	defer func() {
		for envVar, value := range originalValues {
			if value != "" {
				os.Setenv(envVar, value)
			}
		}
	}()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test that ContactEmail falls back to FromEmail when not set
	expectedContactEmail := "noreply@example.com" // This is the default FromEmail
	if config.Email.ContactEmail != expectedContactEmail {
		t.Errorf("Expected ContactEmail to be %s, got %s", expectedContactEmail, config.Email.ContactEmail)
	}

	// Test that StoreName falls back to FromName when not set
	expectedStoreName := "Commercify Store" // This is the default FromName
	if config.Email.StoreName != expectedStoreName {
		t.Errorf("Expected StoreName to be %s, got %s", expectedStoreName, config.Email.StoreName)
	}
}

func TestEmailConfigCustomValues(t *testing.T) {
	// Set custom environment variables
	os.Setenv("EMAIL_CONTACT_ADDRESS", "support@custom.com")
	os.Setenv("STORE_NAME", "Custom Store")
	os.Setenv("EMAIL_FROM_ADDRESS", "from@custom.com")
	os.Setenv("EMAIL_FROM_NAME", "Custom From Name")

	// Clean up after test
	defer func() {
		os.Unsetenv("EMAIL_CONTACT_ADDRESS")
		os.Unsetenv("STORE_NAME")
		os.Unsetenv("EMAIL_FROM_ADDRESS")
		os.Unsetenv("EMAIL_FROM_NAME")
	}()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test that custom values are used
	if config.Email.ContactEmail != "support@custom.com" {
		t.Errorf("Expected ContactEmail to be support@custom.com, got %s", config.Email.ContactEmail)
	}

	if config.Email.StoreName != "Custom Store" {
		t.Errorf("Expected StoreName to be Custom Store, got %s", config.Email.StoreName)
	}

	if config.Email.FromEmail != "from@custom.com" {
		t.Errorf("Expected FromEmail to be from@custom.com, got %s", config.Email.FromEmail)
	}

	if config.Email.FromName != "Custom From Name" {
		t.Errorf("Expected FromName to be Custom From Name, got %s", config.Email.FromName)
	}
}
