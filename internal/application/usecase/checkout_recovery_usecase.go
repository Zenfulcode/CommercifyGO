package usecase

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// CheckoutRecoveryUseCase represents the use case for recovering abandoned checkouts
type CheckoutRecoveryUseCase struct {
	checkoutRepo     repository.CheckoutRepository
	emailService     EmailService
	templatePath     string
	storeName        string
	storeLogoURL     string
	storeURL         string
	privacyPolicyURL string
}

// EmailService defines the interface for sending emails
type EmailService interface {
	SendEmail(to, subject string, body []byte) error
}

// CheckoutRecoveryData contains data for checkout recovery email
type CheckoutRecoveryData struct {
	StoreName           string
	StoreLogoURL        string
	CustomerName        string
	CustomerEmail       string
	Items               []CheckoutItemData
	FormattedTotal      string
	FormattedDiscount   string
	FormattedFinalTotal string
	CheckoutURL         string
	AppliedDiscount     *entity.AppliedDiscount
	DiscountOffer       *DiscountOfferData
	CurrentYear         int
	UnsubscribeURL      string
	PrivacyPolicyURL    string
}

// CheckoutItemData contains formatted data for a checkout item in email
type CheckoutItemData struct {
	ProductName    string
	VariantName    string
	Quantity       int
	FormattedPrice string
}

// DiscountOfferData contains data for special discount offer
type DiscountOfferData struct {
	Code        string
	Description string
	Value       string
	ExpiryDate  string
}

// NewCheckoutRecoveryUseCase creates a new instance of CheckoutRecoveryUseCase
func NewCheckoutRecoveryUseCase(
	checkoutRepo repository.CheckoutRepository,
	emailService EmailService,
	templatePath string,
	storeName string,
	storeLogoURL string,
	storeURL string,
	privacyPolicyURL string,
) *CheckoutRecoveryUseCase {
	return &CheckoutRecoveryUseCase{
		checkoutRepo:     checkoutRepo,
		emailService:     emailService,
		templatePath:     templatePath,
		storeName:        storeName,
		storeLogoURL:     storeLogoURL,
		storeURL:         storeURL,
		privacyPolicyURL: privacyPolicyURL,
	}
}

// SendRecoveryEmail sends a recovery email for an abandoned checkout
func (uc *CheckoutRecoveryUseCase) SendRecoveryEmail(checkout *entity.Checkout) error {
	// Check if checkout has valid customer information
	if checkout.CustomerDetails.Email == "" {
		return fmt.Errorf("cannot send recovery email: customer email is missing")
	}

	// Format currency values
	formattedTotal := formatCurrency(checkout.TotalAmount, checkout.Currency)
	formattedDiscount := formatCurrency(checkout.DiscountAmount, checkout.Currency)
	formattedFinalTotal := formatCurrency(checkout.FinalAmount, checkout.Currency)

	// Create checkout item data
	var itemsData []CheckoutItemData
	for _, item := range checkout.Items {
		itemData := CheckoutItemData{
			ProductName:    item.ProductName,
			VariantName:    item.VariantName,
			Quantity:       item.Quantity,
			FormattedPrice: formatCurrency(item.Price*int64(item.Quantity), checkout.Currency),
		}
		itemsData = append(itemsData, itemData)
	}

	// Create checkout URL
	checkoutURL := fmt.Sprintf("%s/checkout/%d", uc.storeURL, checkout.ID)
	if checkout.SessionID != "" {
		checkoutURL = fmt.Sprintf("%s/checkout?session=%s", uc.storeURL, checkout.SessionID)
	}

	// Create unsubscribe URL
	unsubscribeURL := fmt.Sprintf("%s/unsubscribe?email=%s", uc.storeURL, checkout.CustomerDetails.Email)

	// Create recovery data
	data := CheckoutRecoveryData{
		StoreName:           uc.storeName,
		StoreLogoURL:        uc.storeLogoURL,
		CustomerName:        checkout.CustomerDetails.FullName,
		CustomerEmail:       checkout.CustomerDetails.Email,
		Items:               itemsData,
		FormattedTotal:      formattedTotal,
		FormattedDiscount:   formattedDiscount,
		FormattedFinalTotal: formattedFinalTotal,
		CheckoutURL:         checkoutURL,
		AppliedDiscount:     checkout.AppliedDiscount,
		DiscountOffer:       createDiscountOffer(),
		CurrentYear:         time.Now().Year(),
		UnsubscribeURL:      unsubscribeURL,
		PrivacyPolicyURL:    uc.privacyPolicyURL,
	}

	// Parse template
	tmpl, err := template.ParseFiles(filepath.Join(uc.templatePath, "checkout_recovery.html"))
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Render template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// Send email
	subject := "Complete Your Purchase at " + uc.storeName
	if err := uc.emailService.SendEmail(checkout.CustomerDetails.Email, subject, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to send recovery email: %w", err)
	}

	return nil
}

// ProcessAbandonedCheckouts finds and processes abandoned checkouts
func (uc *CheckoutRecoveryUseCase) ProcessAbandonedCheckouts() (int, error) {
	// Get expired checkouts
	expiredCheckouts, err := uc.checkoutRepo.GetExpiredCheckouts()
	if err != nil {
		return 0, fmt.Errorf("failed to get expired checkouts: %w", err)
	}

	// Count of emails sent
	sentCount := 0

	// Process each expired checkout
	for _, checkout := range expiredCheckouts {
		// Only process active checkouts with customer email
		if checkout.Status == entity.CheckoutStatusActive &&
			checkout.CustomerDetails.Email != "" &&
			len(checkout.Items) > 0 {

			// Mark as abandoned
			checkout.Status = entity.CheckoutStatusAbandoned
			if err := uc.checkoutRepo.Update(checkout); err != nil {
				// Log error but continue with other checkouts
				fmt.Printf("Error updating checkout status: %v\n", err)
				continue
			}

			// Send recovery email
			if err := uc.SendRecoveryEmail(checkout); err != nil {
				// Log error but continue with other checkouts
				fmt.Printf("Error sending recovery email: %v\n", err)
				continue
			}

			sentCount++
		}
	}

	return sentCount, nil
}

// Helper function to create a discount offer
func createDiscountOffer() *DiscountOfferData {
	// Create a discount code that expires in 48 hours
	expiryDate := time.Now().Add(48 * time.Hour).Format("January 2, 2006")

	return &DiscountOfferData{
		Code:        "COMEBACK15",
		Description: "Complete your purchase now and get 15% off!",
		Value:       "15% off",
		ExpiryDate:  expiryDate,
	}
}

// Helper function to format currency
func formatCurrency(amount int64, currency string) string {
	// Format amount as decimal
	decimal := float64(amount) / 100

	// Format based on currency
	switch currency {
	case "USD":
		return fmt.Sprintf("$%.2f", decimal)
	case "EUR":
		return fmt.Sprintf("€%.2f", decimal)
	case "GBP":
		return fmt.Sprintf("£%.2f", decimal)
	default:
		return fmt.Sprintf("%.2f %s", decimal, currency)
	}
}
