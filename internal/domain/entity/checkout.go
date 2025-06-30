package entity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/dto"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"gorm.io/gorm"
)

// CheckoutStatus represents the current status of a checkout
type CheckoutStatus string

const (
	// CheckoutStatusActive represents an active checkout that is being modified
	CheckoutStatusActive CheckoutStatus = "active"
	// CheckoutStatusCompleted represents a checkout that has been converted to an order
	CheckoutStatusCompleted CheckoutStatus = "completed"
	// CheckoutStatusAbandoned represents a checkout that was abandoned by the user
	CheckoutStatusAbandoned CheckoutStatus = "abandoned"
	// CheckoutStatusExpired represents a checkout that has expired due to inactivity
	CheckoutStatusExpired CheckoutStatus = "expired"
)

// Checkout represents a user's checkout session
type Checkout struct {
	gorm.Model
	UserID              *uint           `gorm:"index"`
	User                *User           `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	SessionID           string          `gorm:"uniqueIndex;not null;size:255"`
	Items               []CheckoutItem  `gorm:"foreignKey:CheckoutID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Status              CheckoutStatus  `gorm:"not null;size:50;default:'active'"`
	ShippingAddr        Address         `gorm:"embedded;embeddedPrefix:shipping_"`
	BillingAddr         Address         `gorm:"embedded;embeddedPrefix:billing_"`
	ShippingOptionJSON  *string         `gorm:"column:shipping_option_json;type:text"`
	PaymentProvider     string          `gorm:"size:100"`
	TotalAmount         int64           `gorm:"default:0"`
	ShippingCost        int64           `gorm:"default:0"`
	TotalWeight         float64         `gorm:"default:0"`
	CustomerDetails     CustomerDetails `gorm:"embedded;embeddedPrefix:customer_"`
	Currency            string          `gorm:"not null;size:3"`
	DiscountCode        string          `gorm:"size:100"`
	DiscountAmount      int64           `gorm:"default:0"`
	FinalAmount         int64           `gorm:"default:0"`
	AppliedDiscountJSON *string         `gorm:"column:applied_discount_json;type:text"`
	LastActivityAt      time.Time       `gorm:"index"`
	ExpiresAt           time.Time       `gorm:"index"`
	CompletedAt         *time.Time
	ConvertedOrderID    *uint  `gorm:"index"`
	ConvertedOrder      *Order `gorm:"foreignKey:ConvertedOrderID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
}

func (c *Checkout) CalculateTotals() {
	c.recalculateTotals()
}

// CheckoutItem represents an item in a checkout
type CheckoutItem struct {
	gorm.Model
	CheckoutID       uint           `gorm:"index;not null"`
	Checkout         Checkout       `gorm:"foreignKey:CheckoutID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	ProductID        uint           `gorm:"index;not null"`
	Product          Product        `gorm:"foreignKey:ProductID;constraint:OnDelete:RESTRICT,OnUpdate:CASCADE"`
	ProductVariantID uint           `gorm:"index;not null"`
	ProductVariant   ProductVariant `gorm:"foreignKey:ProductVariantID;constraint:OnDelete:RESTRICT,OnUpdate:CASCADE"`
	ImageURL         string         `gorm:"size:500"`
	Quantity         int            `gorm:"not null"`
	Price            int64          `gorm:"not null"` // Price at time of adding to cart
	Weight           float64        `gorm:"default:0"`
	ProductName      string         `gorm:"not null;size:255"`
	VariantName      string         `gorm:"size:255"`
	SKU              string         `gorm:"not null;size:100"`
}

// AppliedDiscount represents a discount applied to a checkout
type AppliedDiscount struct {
	DiscountID     uint      `gorm:"index"`
	Discount       *Discount `gorm:"foreignKey:DiscountID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	DiscountCode   string    `gorm:"size:100"`
	DiscountAmount int64     `gorm:"default:0"`
}

// NewCheckout creates a new checkout for a guest user
func NewCheckout(sessionID string, currency string) (*Checkout, error) {
	if sessionID == "" {
		return nil, errors.New("session ID cannot be empty")
	}

	if currency == "" {
		return nil, errors.New("currency cannot be empty")
	}

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) // Checkouts expire after 24 hours by default

	return &Checkout{
		SessionID:      sessionID,
		Items:          []CheckoutItem{},
		Status:         CheckoutStatusActive,
		Currency:       currency,
		TotalAmount:    0,
		ShippingCost:   0,
		DiscountAmount: 0,
		FinalAmount:    0,
		LastActivityAt: now,
		ExpiresAt:      expiresAt,
	}, nil
}

// AddItem adds a product to the checkout
func (c *Checkout) AddItem(productID uint, variantID uint, quantity int, price int64, weight float64, productName string, variantName string, sku string) error {
	if productID == 0 {
		return errors.New("product ID cannot be empty")
	}
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	if price < 0 {
		return errors.New("price cannot be negative")
	}

	// Check if the product is already in the checkout
	for i, item := range c.Items {
		// Match by both product ID and variant ID (if variant ID is provided)
		if item.ProductID == productID &&
			(variantID == 0 || item.ProductVariantID == variantID) {
			// Update quantity if product already exists
			c.Items[i].Quantity += quantity

			// Update checkout
			c.recalculateTotals()

			c.LastActivityAt = time.Now()

			return nil
		}
	}

	// Add new item if product doesn't exist in checkout
	now := time.Now()
	c.Items = append(c.Items, CheckoutItem{
		CheckoutID:       c.ID, // Set the checkout ID for the foreign key
		ProductID:        productID,
		ProductVariantID: variantID,
		Quantity:         quantity,
		Price:            price,
		Weight:           weight,
		ProductName:      productName,
		VariantName:      variantName,
		SKU:              sku,
	})

	// Update checkout
	c.recalculateTotals()
	c.LastActivityAt = now

	return nil
}

// UpdateItem updates the quantity of a product in the checkout
func (c *Checkout) UpdateItem(productID uint, variantID uint, quantity int) error {
	if productID == 0 {
		return errors.New("product ID cannot be empty")
	}
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	for i, item := range c.Items {
		// Match by both product ID and variant ID (if variant ID is provided)
		if item.ProductID == productID &&
			(variantID == 0 || item.ProductVariantID == variantID) {
			c.Items[i].Quantity = quantity

			// Update checkout
			c.recalculateTotals()

			c.LastActivityAt = time.Now()

			return nil
		}
	}

	return errors.New("product not found in checkout")
}

// RemoveItem removes a product from the checkout
func (c *Checkout) RemoveItem(productID uint, variantID uint) error {
	if productID == 0 {
		return errors.New("product ID cannot be empty")
	}

	for i, item := range c.Items {
		// Match by both product ID and variant ID (if variant ID is provided)
		if item.ProductID == productID &&
			(variantID == 0 || item.ProductVariantID == variantID) {
			// Remove item from slice
			c.Items = append(c.Items[:i], c.Items[i+1:]...)

			// Update checkout
			c.recalculateTotals()

			c.LastActivityAt = time.Now()

			return nil
		}
	}

	return errors.New("product not found in checkout")
}

// SetShippingAddress sets the shipping address for the checkout
func (c *Checkout) SetShippingAddress(address Address) {
	c.ShippingAddr = address

	c.LastActivityAt = time.Now()
}

// SetBillingAddress sets the billing address for the checkout
func (c *Checkout) SetBillingAddress(address Address) {
	c.BillingAddr = address

	c.LastActivityAt = time.Now()
}

// SetCustomerDetails sets the customer details for the checkout
func (c *Checkout) SetCustomerDetails(details CustomerDetails) {
	c.CustomerDetails = details

	c.LastActivityAt = time.Now()
}

// SetShippingMethod sets the shipping method for the checkout
func (c *Checkout) SetShippingMethod(option *ShippingOption) {
	if option != nil {
		c.ShippingCost = option.Cost

		// Store shipping option as JSON
		c.SetShippingOptionJSON(option)
	} else {
		c.ShippingCost = 0

		// Clear shipping option JSON
		c.SetShippingOptionJSON(nil)
	}

	c.recalculateTotals()

	c.LastActivityAt = time.Now()
}

// SetPaymentProvider sets the payment provider for the checkout
func (c *Checkout) SetPaymentProvider(provider string) {
	c.PaymentProvider = provider

	c.LastActivityAt = time.Now()
}

// SetCurrency changes the currency of the checkout and converts all prices
func (c *Checkout) SetCurrency(newCurrency string, fromCurrency *Currency, toCurrency *Currency) {
	if c.Currency == newCurrency {
		return
	}

	// Convert all item prices
	for i := range c.Items {
		c.Items[i].Price = fromCurrency.ConvertAmount(c.Items[i].Price, toCurrency)
	}

	// Convert shipping cost
	c.ShippingCost = fromCurrency.ConvertAmount(c.ShippingCost, toCurrency)

	// Convert discount amount
	c.DiscountAmount = fromCurrency.ConvertAmount(c.DiscountAmount, toCurrency)

	// Update currency
	c.Currency = newCurrency

	// Recalculate totals with new currency prices
	c.recalculateTotals()

	c.LastActivityAt = time.Now()
}

// ApplyDiscount applies a discount to the checkout
func (c *Checkout) ApplyDiscount(discount *Discount) {
	if discount == nil {
		// Remove any existing discount
		c.DiscountCode = ""
		c.DiscountAmount = 0
		c.SetAppliedDiscountJSON(nil)
	} else {
		// Calculate discount amount
		discountAmount := discount.CalculateDiscount(&Order{
			TotalAmount: c.TotalAmount,
			Items:       convertCheckoutItemsToOrderItems(c.Items),
		})

		// Apply the discount
		c.DiscountCode = discount.Code
		c.DiscountAmount = discountAmount

		// Store applied discount as JSON
		appliedDiscount := &AppliedDiscount{
			DiscountID:     discount.ID,
			DiscountCode:   discount.Code,
			DiscountAmount: discountAmount,
		}
		c.SetAppliedDiscountJSON(appliedDiscount)
	}

	c.recalculateTotals()

	c.LastActivityAt = time.Now()
}

// Clear empties the checkout
func (c *Checkout) Clear() {
	c.Items = []CheckoutItem{}
	c.TotalAmount = 0
	c.TotalWeight = 0
	c.DiscountAmount = 0
	c.FinalAmount = 0
	c.SetAppliedDiscountJSON(nil)

	c.LastActivityAt = time.Now()
}

// MarkAsCompleted marks the checkout as completed and sets the completed_at timestamp
func (c *Checkout) MarkAsCompleted(orderID uint) {
	c.Status = CheckoutStatusCompleted
	c.ConvertedOrderID = &orderID
	now := time.Now()
	c.CompletedAt = &now
	c.UpdatedAt = now
	c.LastActivityAt = now
}

// MarkAsAbandoned marks the checkout as abandoned
func (c *Checkout) MarkAsAbandoned() {
	c.Status = CheckoutStatusAbandoned

	c.LastActivityAt = time.Now()
}

// MarkAsExpired marks the checkout as expired
func (c *Checkout) MarkAsExpired() {
	c.Status = CheckoutStatusExpired

	c.LastActivityAt = time.Now()
}

// IsExpired checks if the checkout has expired
func (c *Checkout) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// ExtendExpiry extends the expiry time of the checkout
func (c *Checkout) ExtendExpiry(duration time.Duration) {
	c.ExpiresAt = time.Now().Add(duration)

	c.LastActivityAt = time.Now()
}

// TotalItems returns the total number of items in the checkout
func (c *Checkout) TotalItems() int {
	total := 0
	for _, item := range c.Items {
		total += item.Quantity
	}
	return total
}

// HasCustomerInfo returns true if the checkout has customer information
func (c *Checkout) HasCustomerInfo() bool {
	return c.CustomerDetails.Email != "" ||
		c.CustomerDetails.Phone != "" ||
		c.CustomerDetails.FullName != ""
}

// HasShippingInfo returns true if the checkout has shipping address information
func (c *Checkout) HasShippingInfo() bool {
	return c.ShippingAddr.Street1 != "" ||
		c.ShippingAddr.City != "" ||
		c.ShippingAddr.PostalCode != "" ||
		c.ShippingAddr.Country != ""
}

// HasCustomerOrShippingInfo returns true if the checkout has either customer or shipping information
func (c *Checkout) HasCustomerOrShippingInfo() bool {
	return c.HasCustomerInfo() || c.HasShippingInfo()
}

// IsEmpty returns true if the checkout has no items and no customer/shipping information
func (c *Checkout) IsEmpty() bool {
	return len(c.Items) == 0 && !c.HasCustomerOrShippingInfo()
}

// ShouldBeAbandoned returns true if the checkout should be marked as abandoned
// (has customer/shipping info and hasn't been active for 15 minutes)
func (c *Checkout) ShouldBeAbandoned() bool {
	if c.Status != CheckoutStatusActive {
		return false
	}

	if !c.HasCustomerOrShippingInfo() {
		return false
	}

	abandonThreshold := time.Now().Add(-15 * time.Minute)
	return c.LastActivityAt.Before(abandonThreshold)
}

// ShouldBeDeleted returns true if the checkout should be deleted
func (c *Checkout) ShouldBeDeleted() bool {
	now := time.Now()

	// Delete empty checkouts after 24 hours
	if !c.HasCustomerOrShippingInfo() {
		deleteThreshold := now.Add(-24 * time.Hour)
		return c.LastActivityAt.Before(deleteThreshold)
	}

	// Delete abandoned checkouts after 7 days in abandoned state
	if c.Status == CheckoutStatusAbandoned {
		deleteThreshold := now.Add(-7 * 24 * time.Hour)
		return c.UpdatedAt.Before(deleteThreshold)
	}

	// Delete all expired checkouts
	if c.Status == CheckoutStatusExpired {
		return true
	}

	return false
}

// recalculateTotals recalculates the total amount, weight, and final amount
func (c *Checkout) recalculateTotals() {
	// Calculate total amount and weight
	totalAmount := int64(0)
	totalWeight := float64(0)
	for _, item := range c.Items {
		itemTotal := item.Price * int64(item.Quantity)
		totalAmount += itemTotal
		totalWeight += item.Weight * float64(item.Quantity)
	}

	c.TotalAmount = totalAmount
	c.TotalWeight = totalWeight

	// Calculate final amount with explicit calculation to avoid floating point inconsistencies
	c.FinalAmount = max(totalAmount+c.ShippingCost-c.DiscountAmount, 0)
}

// ToOrder converts a checkout to an order
func (c *Checkout) ToOrder() *Order {
	// Create order items from checkout items
	items := make([]OrderItem, len(c.Items))
	for i, item := range c.Items {
		items[i] = OrderItem{
			ProductID:        item.ProductID,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			Price:            item.Price,
			Subtotal:         item.Price * int64(item.Quantity),
			Weight:           item.Weight,
			ProductName:      item.ProductName,
			SKU:              item.SKU,
		}
	}

	// Determine if this is a guest order
	isGuestOrder := c.UserID == nil

	var userID uint
	if c.UserID != nil {
		userID = *c.UserID
	}

	// Create the order
	var shippingMethodID uint
	if storedOption := c.GetShippingOption(); storedOption != nil {
		shippingMethodID = storedOption.ShippingMethodID
	}

	// Create ShippingOption from JSON if available
	var shippingOption *ShippingOption
	if storedOption := c.GetShippingOption(); storedOption != nil {
		shippingOption = storedOption
	}

	// Create AppliedDiscount from JSON if available
	var appliedDiscount *AppliedDiscount
	if storedDiscount := c.GetAppliedDiscount(); storedDiscount != nil {
		appliedDiscount = storedDiscount
	}

	order := &Order{
		UserID:            userID, // This will be 0 for guest orders
		Items:             items,
		Currency:          c.Currency,
		TotalAmount:       c.TotalAmount,
		TotalWeight:       c.TotalWeight,
		ShippingCost:      c.ShippingCost,
		DiscountAmount:    c.DiscountAmount,
		FinalAmount:       c.FinalAmount,
		Status:            OrderStatusPending,
		PaymentStatus:     PaymentStatusPending, // Initialize payment status
		ShippingAddr:      c.ShippingAddr,
		BillingAddr:       c.BillingAddr,
		CustomerDetails:   &c.CustomerDetails,
		ShippingMethodID:  shippingMethodID,
		ShippingOption:    shippingOption,
		PaymentProvider:   c.PaymentProvider,
		IsGuestOrder:      isGuestOrder,
		PaymentMethod:     "wallet", // Default payment method
		AppliedDiscount:   appliedDiscount,
		CheckoutSessionID: c.SessionID,
	}

	// Generate a friendly order number (will be replaced with actual ID after creation)
	order.OrderNumber = generateOrderNumber()

	return order
}

// convertCheckoutItemsToOrderItems is a helper function to convert checkout items to order items
func convertCheckoutItemsToOrderItems(checkoutItems []CheckoutItem) []OrderItem {
	orderItems := make([]OrderItem, len(checkoutItems))
	for i, item := range checkoutItems {
		orderItems[i] = OrderItem{
			ProductID:        item.ProductID,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			Price:            item.Price,
			Subtotal:         item.Price * int64(item.Quantity),
			Weight:           item.Weight,
			ProductName:      item.ProductName,
			SKU:              item.SKU,
		}
	}
	return orderItems
}

// generateOrderNumber generates a temporary order number
func generateOrderNumber() string {
	return "ORD-" + time.Now().Format("20060102") + "-TEMP"
}

func (c *Checkout) ToCheckoutDTO() *dto.CheckoutDTO {
	var userID uint
	if c.UserID != nil {
		userID = *c.UserID
	}

	var shippingMethodID uint
	var shippingOption *dto.ShippingOptionDTO
	if storedOption := c.GetShippingOption(); storedOption != nil {
		shippingMethodID = storedOption.ShippingMethodID
		shippingOption = &dto.ShippingOptionDTO{
			ShippingRateID:        storedOption.ShippingRateID,
			ShippingMethodID:      storedOption.ShippingMethodID,
			Name:                  storedOption.Name,
			Description:           storedOption.Description,
			EstimatedDeliveryDays: storedOption.EstimatedDeliveryDays,
			Cost:                  money.FromCents(storedOption.Cost),
			FreeShipping:          storedOption.FreeShipping,
		}
	}

	return &dto.CheckoutDTO{
		ID:               c.ID,
		SessionID:        c.SessionID,
		UserID:           userID,
		Status:           string(c.Status),
		ShippingAddress:  *c.ShippingAddr.ToAddressDTO(),
		BillingAddress:   *c.BillingAddr.ToAddressDTO(),
		ShippingMethodID: shippingMethodID,
		ShippingOption:   shippingOption,
		CustomerDetails:  *c.CustomerDetails.ToCustomerDetailsDTO(),
		PaymentProvider:  c.PaymentProvider,
		TotalAmount:      money.FromCents(c.TotalAmount),
		ShippingCost:     money.FromCents(c.ShippingCost),
		TotalWeight:      c.TotalWeight,
		Currency:         c.Currency,
		DiscountCode:     c.DiscountCode,
		DiscountAmount:   money.FromCents(c.DiscountAmount),
		FinalAmount:      money.FromCents(c.FinalAmount),
		LastActivityAt:   c.LastActivityAt,
		ExpiresAt:        c.ExpiresAt,
	}
}

// ToAppliedDiscountDTO converts AppliedDiscount to DTO
func (a *AppliedDiscount) ToAppliedDiscountDTO() *dto.AppliedDiscountDTO {
	if a == nil {
		return nil
	}

	var discountType, discountMethod string
	var discountValue float64

	if a.Discount != nil {
		discountType = string(a.Discount.Type)
		discountMethod = string(a.Discount.Method)
		discountValue = a.Discount.Value
	}

	return &dto.AppliedDiscountDTO{
		ID:     a.DiscountID,
		Code:   a.DiscountCode,
		Type:   discountType,
		Method: discountMethod,
		Value:  discountValue,
		Amount: money.FromCents(a.DiscountAmount),
	}
}

// GetShippingOption retrieves the shipping option from JSON
func (c *Checkout) GetShippingOption() *ShippingOption {
	if c.ShippingOptionJSON == nil || *c.ShippingOptionJSON == "" {
		return nil
	}

	var option ShippingOption
	if err := json.Unmarshal([]byte(*c.ShippingOptionJSON), &option); err != nil {
		return nil
	}

	return &option
}

// SetShippingOptionJSON stores the shipping option as JSON
func (c *Checkout) SetShippingOptionJSON(option *ShippingOption) error {
	if option == nil {
		c.ShippingOptionJSON = nil
		return nil
	}

	data, err := json.Marshal(option)
	if err != nil {
		return err
	}

	jsonStr := string(data)
	c.ShippingOptionJSON = &jsonStr
	return nil
}

// GetAppliedDiscount retrieves the applied discount from JSON
func (c *Checkout) GetAppliedDiscount() *AppliedDiscount {
	if c.AppliedDiscountJSON == nil || *c.AppliedDiscountJSON == "" {
		return nil
	}

	var discount AppliedDiscount
	if err := json.Unmarshal([]byte(*c.AppliedDiscountJSON), &discount); err != nil {
		return nil
	}

	return &discount
}

// SetAppliedDiscountJSON stores the applied discount as JSON
func (c *Checkout) SetAppliedDiscountJSON(discount *AppliedDiscount) error {
	if discount == nil {
		c.AppliedDiscountJSON = nil
		return nil
	}

	data, err := json.Marshal(discount)
	if err != nil {
		return err
	}

	jsonStr := string(data)
	c.AppliedDiscountJSON = &jsonStr
	return nil
}
