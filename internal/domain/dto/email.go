package dto

// EmailTestDetails represents additional details in the email test response
type EmailTestDetails struct {
	TargetEmail string `json:"target_email"`
	OrderID     string `json:"order_id"`
}
