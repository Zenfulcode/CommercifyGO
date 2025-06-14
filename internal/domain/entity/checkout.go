package entity

import (
	"errors"
	"time"
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
	ID               uint             `json:"id"`
	UserID           uint             `json:"user_id,omitempty"`
	SessionID        string           `json:"session_id,omitempty"`
	Items            []CheckoutItem   `json:"items"`
	Status           CheckoutStatus   `json:"status"`
	ShippingAddr     Address          `json:"shipping_address"`
	BillingAddr      Address          `json:"billing_address"`
	ShippingMethodID uint             `json:"shipping_method_id,omitempty"`
	ShippingOption   *ShippingOption  `json:"shipping_option,omitempty"`
	PaymentProvider  string           `json:"payment_provider,omitempty"`
	TotalAmount      int64            `json:"total_amount"`  // stored in cents
	ShippingCost     int64            `json:"shipping_cost"` // stored in cents
	TotalWeight      float64          `json:"total_weight"`
	CustomerDetails  CustomerDetails  `json:"customer_details"`
	Currency         string           `json:"currency"`
	DiscountCode     string           `json:"discount_code,omitempty"`
	DiscountAmount   int64            `json:"discount_amount"` // stored in cents
	FinalAmount      int64            `json:"final_amount"`    // stored in cents
	AppliedDiscount  *AppliedDiscount `json:"applied_discount,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	LastActivityAt   time.Time        `json:"last_activity_at"`
	ExpiresAt        time.Time        `json:"expires_at"`
	CompletedAt      *time.Time       `json:"completed_at,omitempty"`
	ConvertedOrderID uint             `json:"converted_order_id,omitempty"`
}

func (c *Checkout) CalculateTotals() {
	c.recalculateTotals()
}

// CheckoutItem represents an item in a checkout
type CheckoutItem struct {
	ID               uint      `json:"id"`
	CheckoutID       uint      `json:"checkout_id"`
	ProductID        uint      `json:"product_id"`
	ProductVariantID uint      `json:"product_variant_id,omitempty"`
	ImageURL         string    `json:"image_url,omitempty"`
	Quantity         int       `json:"quantity"`
	Price            int64     `json:"price"` // stored in cents
	Weight           float64   `json:"weight"`
	ProductName      string    `json:"product_name"`
	VariantName      string    `json:"variant_name,omitempty"`
	SKU              string    `json:"sku,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// AppliedDiscount represents a discount applied to a checkout
type AppliedDiscount struct {
	DiscountID     uint   `json:"discount_id"`
	DiscountCode   string `json:"discount_code"`
	DiscountAmount int64  `json:"discount_amount"` // stored in cents
}

// NewCheckout creates a new checkout for a guest user
func NewCheckout(sessionID string) (*Checkout, error) {
	if sessionID == "" {
		return nil, errors.New("session ID cannot be empty")
	}

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) // Checkouts expire after 24 hours by default

	return &Checkout{
		SessionID:      sessionID,
		Items:          []CheckoutItem{},
		Status:         CheckoutStatusActive,
		Currency:       "USD", // Default currency
		TotalAmount:    0,
		ShippingCost:   0,
		DiscountAmount: 0,
		FinalAmount:    0,
		CreatedAt:      now,
		UpdatedAt:      now,
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
			c.Items[i].UpdatedAt = time.Now()

			// Update checkout
			c.recalculateTotals()
			c.UpdatedAt = time.Now()
			c.LastActivityAt = time.Now()

			return nil
		}
	}

	// Add new item if product doesn't exist in checkout
	now := time.Now()
	c.Items = append(c.Items, CheckoutItem{
		ProductID:        productID,
		ProductVariantID: variantID,
		Quantity:         quantity,
		Price:            price,
		Weight:           weight,
		ProductName:      productName,
		VariantName:      variantName,
		SKU:              sku,
		CreatedAt:        now,
		UpdatedAt:        now,
	})

	// Update checkout
	c.recalculateTotals()
	c.UpdatedAt = now
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
			c.Items[i].UpdatedAt = time.Now()

			// Update checkout
			c.recalculateTotals()
			c.UpdatedAt = time.Now()
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
			c.UpdatedAt = time.Now()
			c.LastActivityAt = time.Now()

			return nil
		}
	}

	return errors.New("product not found in checkout")
}

// SetShippingAddress sets the shipping address for the checkout
func (c *Checkout) SetShippingAddress(address Address) {
	c.ShippingAddr = address
	c.UpdatedAt = time.Now()
	c.LastActivityAt = time.Now()
}

// SetBillingAddress sets the billing address for the checkout
func (c *Checkout) SetBillingAddress(address Address) {
	c.BillingAddr = address
	c.UpdatedAt = time.Now()
	c.LastActivityAt = time.Now()
}

// SetCustomerDetails sets the customer details for the checkout
func (c *Checkout) SetCustomerDetails(details CustomerDetails) {
	c.CustomerDetails = details
	c.UpdatedAt = time.Now()
	c.LastActivityAt = time.Now()
}

// SetShippingMethod sets the shipping method for the checkout
func (c *Checkout) SetShippingMethod(option *ShippingOption) {
	c.ShippingMethodID = option.ShippingMethodID
	c.ShippingCost = option.Cost
	c.ShippingOption = option

	c.recalculateTotals()

	c.UpdatedAt = time.Now()
	c.LastActivityAt = time.Now()
}

// SetPaymentProvider sets the payment provider for the checkout
func (c *Checkout) SetPaymentProvider(provider string) {
	c.PaymentProvider = provider
	c.UpdatedAt = time.Now()
	c.LastActivityAt = time.Now()
}

// SetCurrency changes the currency of the checkout and converts all prices
func (c *Checkout) SetCurrency(newCurrency string, fromCurrency *Currency, toCurrency *Currency) {
	if c.Currency == newCurrency {
		// No change needed
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
	c.UpdatedAt = time.Now()
	c.LastActivityAt = time.Now()
}

// ApplyDiscount applies a discount to the checkout
func (c *Checkout) ApplyDiscount(discount *Discount) {
	if discount == nil {
		// Remove any existing discount
		c.DiscountCode = ""
		c.DiscountAmount = 0
		c.AppliedDiscount = nil
	} else {
		// Calculate discount amount
		discountAmount := discount.CalculateDiscount(&Order{
			TotalAmount: c.TotalAmount,
			Items:       convertCheckoutItemsToOrderItems(c.Items),
		})

		// Apply the discount
		c.DiscountCode = discount.Code
		c.DiscountAmount = discountAmount
		c.AppliedDiscount = &AppliedDiscount{
			DiscountID:     discount.ID,
			DiscountCode:   discount.Code,
			DiscountAmount: discountAmount,
		}
	}

	c.recalculateTotals()
	c.UpdatedAt = time.Now()
	c.LastActivityAt = time.Now()
}

// Clear empties the checkout
func (c *Checkout) Clear() {
	c.Items = []CheckoutItem{}
	c.TotalAmount = 0
	c.TotalWeight = 0
	c.DiscountAmount = 0
	c.FinalAmount = 0
	c.AppliedDiscount = nil
	c.UpdatedAt = time.Now()
	c.LastActivityAt = time.Now()
}

// MarkAsCompleted marks the checkout as completed and sets the completed_at timestamp
func (c *Checkout) MarkAsCompleted(orderID uint) {
	c.Status = CheckoutStatusCompleted
	c.ConvertedOrderID = orderID
	now := time.Now()
	c.CompletedAt = &now
	c.UpdatedAt = now
	c.LastActivityAt = now
}

// MarkAsAbandoned marks the checkout as abandoned
func (c *Checkout) MarkAsAbandoned() {
	c.Status = CheckoutStatusAbandoned
	c.UpdatedAt = time.Now()
	c.LastActivityAt = time.Now()
}

// MarkAsExpired marks the checkout as expired
func (c *Checkout) MarkAsExpired() {
	c.Status = CheckoutStatusExpired
	c.UpdatedAt = time.Now()
	c.LastActivityAt = time.Now()
}

// IsExpired checks if the checkout has expired
func (c *Checkout) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// ExtendExpiry extends the expiry time of the checkout
func (c *Checkout) ExtendExpiry(duration time.Duration) {
	c.ExpiresAt = time.Now().Add(duration)
	c.UpdatedAt = time.Now()
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
			ProductID:   item.ProductID,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Price * int64(item.Quantity),
			Weight:      item.Weight,
			ProductName: item.ProductName,
			SKU:         item.SKU,
		}
	}

	// Determine if this is a guest order
	isGuestOrder := c.UserID == 0

	// Create the order
	order := &Order{
		UserID:            c.UserID, // This will be 0 for guest orders
		Items:             items,
		TotalAmount:       c.TotalAmount,
		TotalWeight:       c.TotalWeight,
		ShippingCost:      c.ShippingCost,
		DiscountAmount:    c.DiscountAmount,
		FinalAmount:       c.FinalAmount,
		Status:            OrderStatusPending,
		ShippingAddr:      c.ShippingAddr,
		BillingAddr:       c.BillingAddr,
		CustomerDetails:   &c.CustomerDetails,
		ShippingMethodID:  c.ShippingMethodID,
		ShippingOption:    c.ShippingOption,
		PaymentProvider:   c.PaymentProvider,
		IsGuestOrder:      isGuestOrder,
		PaymentMethod:     "wallet", // Default payment method
		AppliedDiscount:   c.AppliedDiscount,
		CheckoutSessionID: c.SessionID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
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
			ProductID:   item.ProductID,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Price * int64(item.Quantity),
			Weight:      item.Weight,
			ProductName: item.ProductName,
			SKU:         item.SKU,
		}
	}
	return orderItems
}

// generateOrderNumber generates a temporary order number
func generateOrderNumber() string {
	return "ORD-" + time.Now().Format("20060102") + "-TEMP"
}
