package contracts

type CapturePaymentRequest struct {
	Amount float64 `json:"amount,omitempty"` // Optional when is_full is true
	IsFull bool    `json:"is_full"`          // Whether to capture the full amount
}

type RefundPaymentRequest struct {
	Amount float64 `json:"amount,omitempty"` // Optional when is_full is true
	IsFull bool    `json:"is_full"`          // Whether to refund the full captured amount
}
