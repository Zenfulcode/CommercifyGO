# Payment API Examples

This document provides example request bodies for the payment system API endpoints.

## Public Payment Endpoints

### Get Available Payment Providers

```plaintext
GET /api/payment/providers
GET /api/payment/providers?currency=<currency_code>
```

Retrieves the list of available payment providers for the store. Optionally filter by currency to get only providers that support the specified currency.

**Query Parameters:**

- `currency` (optional): Three-letter ISO currency code (e.g., "USD", "EUR", "NOK") to filter providers by supported currency

Example response:

```json
[
  {
    "type": "stripe",
    "name": "Stripe",
    "description": "Pay with credit or debit card",
    "enabled": true,
    "methods": ["credit_card"],
    "supported_currencies": [
      "USD",
      "EUR",
      "GBP",
      "JPY",
      "CAD",
      "AUD",
      "CHF",
      "SEK",
      "NOK",
      "DKK",
      "PLN",
      "CZK",
      "HUF",
      "BGN",
      "RON",
      "HRK",
      "ISK",
      "MXN",
      "BRL",
      "SGD",
      "HKD",
      "INR",
      "MYR",
      "PHP",
      "THB",
      "TWD",
      "KRW",
      "NZD",
      "ILS",
      "ZAR"
    ]
  },
  {
    "type": "mobilepay",
    "name": "MobilePay",
    "description": "Pay with MobilePay app",
    "enabled": true,
    "methods": ["wallet"],
    "supported_currencies": ["NOK", "DKK", "EUR"]
  }
]
```

Example request filtering by currency:

```plaintext
GET /api/payment/providers?currency=NOK
```

This would return only payment providers that support Norwegian Krone (NOK), which would include Stripe and MobilePay but exclude providers that don't support NOK.

**Status Codes:**

- `200 OK`: Providers retrieved successfully

## Admin Payment Management Endpoints

### Capture Payment

```plaintext
POST /api/admin/payments/{paymentId}/capture
```

Capture a previously authorized payment (admin only).

**Request Body (Partial Capture):**

```json
{
  "amount": 1500.00,
  "is_full": false
}
```

**Request Body (Full Capture):**

```json
{
  "is_full": true
}
```

**Note:** 
- When `is_full` is `true`, the `amount` field is ignored and the full authorized amount is captured
- When `is_full` is `false` (or omitted), the `amount` field is required
- If both `amount` and `is_full: true` are provided, `is_full` takes precedence

Example response:

```json
{
  "status": "success",
  "message": "Payment captured successfully"
}
```

**Status Codes:**

- `200 OK`: Payment captured successfully
- `400 Bad Request`: Invalid request or capture not allowed
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized (not an admin)
- `404 Not Found`: Payment not found
- `500 Internal Server Error`: Failed to capture payment

### Cancel Payment

```plaintext
POST /api/admin/payments/{paymentId}/cancel
```

Cancel a payment that requires action but hasn't been completed (admin only).

Example response:

```json
{
  "status": "success",
  "message": "Payment cancelled successfully"
}
```

**Status Codes:**

- `200 OK`: Payment cancelled successfully
- `400 Bad Request`: Payment cancellation not allowed
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized (not an admin)
- `404 Not Found`: Payment not found
- `500 Internal Server Error`: Failed to cancel payment

### Refund Payment

```plaintext
POST /api/admin/payments/{paymentId}/refund
```

Refund a captured payment (admin only).

**Request Body (Partial Refund):**

```json
{
  "amount": 1500.00,
  "is_full": false
}
```

**Request Body (Full Refund):**

```json
{
  "is_full": true
}
```

**Note:** 
- When `is_full` is `true`, the `amount` field is ignored and the full captured amount is refunded
- When `is_full` is `false` (or omitted), the `amount` field is required
- If both `amount` and `is_full: true` are provided, `is_full` takes precedence

Example response:

```json
{
  "status": "success",
  "message": "Payment refunded successfully"
}
```

**Status Codes:**

- `200 OK`: Payment refunded successfully
- `400 Bad Request`: Invalid request or refund not allowed
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized (not an admin)
- `404 Not Found`: Payment not found
- `500 Internal Server Error`: Failed to refund payment

## Payment Webhook Endpoints

### Stripe Webhook

```plaintext
POST /api/webhooks/stripe
```

Endpoint for receiving Stripe payment event webhooks.

**Note:** This endpoint is for Stripe's server-to-server communication and should not be called directly by clients.

### MobilePay Webhook

```plaintext
POST /api/webhooks/mobilepay
```

Endpoint for receiving MobilePay payment event webhooks.

**Note:** This endpoint is for MobilePay's server-to-server communication and should not be called directly by clients.

## Payment Workflow Examples

### Credit Card Payment Flow (with 3D Secure)

1. Customer enters payment information and submits order
2. System sends payment request to Stripe
3. If 3D Secure is required:
   - Order status is set to "pending_action"
   - Customer is redirected to 3D Secure authentication page via action_url
   - After authentication, customer is redirected back to the store
   - Stripe sends webhook notification to confirm payment status
   - System updates order status to "paid"
4. If 3D Secure is not required:
   - Payment is processed immediately
   - Order status is set to "paid"

### MobilePay Payment Flow

1. Customer selects MobilePay as payment method and provides phone number
2. System creates payment request with MobilePay
3. Customer is redirected to MobilePay app or web interface via action_url
4. Customer approves payment in MobilePay app
5. MobilePay sends webhook notification confirming payment authorization
6. System updates order status to "paid"
7. Admin can later capture the payment to complete the transaction

### PayPal Payment Flow

1. Customer selects PayPal as payment method
2. System creates payment request with PayPal
3. Customer is redirected to PayPal login page via action_url
4. Customer logs in to PayPal and approves payment
5. PayPal redirects customer back to store's return URL
6. System verifies payment status with PayPal API
7. Order status is updated to "paid"

### Admin Payment Management Flow

1. Customer places order and authorizes payment
2. Admin reviews order and decides to capture the payment
3. Admin uses the capture endpoint to process the payment
4. If needed, admin can issue partial or full refunds using the refund endpoint
5. For problematic payments, admin can cancel pending payments using the cancel endpoint
