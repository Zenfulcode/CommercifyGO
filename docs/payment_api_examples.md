# Payment API Examples

This document provides example request bodies for the payment system API endpoints.

## Public Payment Endpoints

# Payment API Examples

This document provides example request bodies for the payment system API endpoints.

## Public Payment Endpoints

### Get Available Payment Providers

```plaintext
GET /api/payment/providers
```

Retrieves the list of available payment providers for the store.

**Query Parameters:**

- `currency` (string, optional): Three-letter ISO currency code to filter providers by supported currency

**Response Body:**

```json
{
  "success": true,
  "data": [
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
}
```

**Status Codes:**

- `200 OK`: Providers retrieved successfully
- `500 Internal Server Error`: Failed to retrieve providers

## Admin Payment Management Endpoints

All admin payment endpoints require authentication and admin role.

### Capture Payment

```plaintext
POST /api/admin/payments/{paymentId}/capture
```

Capture a previously authorized payment (admin only).

**Path Parameters:**

- `paymentId` (required): Payment ID

**Request Body (Optional for partial capture):**

```json
{
  "amount": 150.0
}
```

**Status Codes:**

- `200 OK`: Payment captured successfully
- `400 Bad Request`: Invalid request or payment cannot be captured
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Payment not found

### Cancel Payment

```plaintext
POST /api/admin/payments/{paymentId}/cancel
```

Cancel an authorized payment (admin only).

**Path Parameters:**

- `paymentId` (required): Payment ID

**Status Codes:**

- `200 OK`: Payment cancelled successfully
- `400 Bad Request`: Payment cannot be cancelled
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Payment not found

### Refund Payment

```plaintext
POST /api/admin/payments/{paymentId}/refund
```

Refund a captured payment (admin only).

**Path Parameters:**

- `paymentId` (required): Payment ID

**Request Body (Optional for partial refund):**

```json
{
  "amount": 75.0,
  "reason": "Customer requested refund"
}
```

**Status Codes:**

- `200 OK`: Payment refunded successfully
- `400 Bad Request`: Invalid request or payment cannot be refunded
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Payment not found

### Force Approve MobilePay Payment

```plaintext
POST /api/admin/payments/{paymentId}/force-approve
```

Force approve a MobilePay payment (admin only). This is typically used for testing purposes.

**Path Parameters:**

- `paymentId` (required): Payment ID

**Status Codes:**

- `200 OK`: Payment force approved successfully
- `400 Bad Request`: Payment cannot be force approved or not a MobilePay payment
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Payment not found

## Admin Payment Provider Management Endpoints

### Get Payment Providers

```plaintext
GET /api/admin/payment-providers
```

Get all payment providers with their configuration (admin only).

**Status Codes:**

- `200 OK`: Payment providers retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

### Get Enabled Payment Providers

```plaintext
GET /api/admin/payment-providers/enabled
```

Get only enabled payment providers (admin only).

**Status Codes:**

- `200 OK`: Enabled payment providers retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

### Enable/Disable Payment Provider

```plaintext
PUT /api/admin/payment-providers/{providerType}/enable
```

Enable or disable a payment provider (admin only).

**Path Parameters:**

- `providerType` (required): Provider type (e.g., "stripe", "mobilepay")

**Request Body:**

```json
{
  "enabled": true
}
```

**Status Codes:**

- `200 OK`: Provider status updated successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Provider not found

### Update Provider Configuration

```plaintext
PUT /api/admin/payment-providers/{providerType}/configuration
```

Update payment provider configuration (admin only).

**Path Parameters:**

- `providerType` (required): Provider type (e.g., "stripe", "mobilepay")

**Request Body:**

```json
{
  "configuration": {
    "api_key": "sk_test_...",
    "webhook_secret": "whsec_...",
    "sandbox_mode": true
  }
}
```

**Status Codes:**

- `200 OK`: Configuration updated successfully
- `400 Bad Request`: Invalid configuration
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Provider not found

### Register Webhook

```plaintext
POST /api/admin/payment-providers/{providerType}/webhook
```

Register a webhook for a payment provider (admin only).

**Path Parameters:**

- `providerType` (required): Provider type (e.g., "stripe", "mobilepay")

**Request Body:**

```json
{
  "url": "https://api.example.com/webhooks/stripe",
  "events": ["payment_intent.succeeded", "payment_intent.payment_failed"]
}
```

**Status Codes:**

- `201 Created`: Webhook registered successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

### Delete Webhook

```plaintext
DELETE /api/admin/payment-providers/{providerType}/webhook
```

Delete webhook for a payment provider (admin only).

**Path Parameters:**

- `providerType` (required): Provider type (e.g., "stripe", "mobilepay")

**Status Codes:**

- `200 OK`: Webhook deleted successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Webhook not found

### Get Webhook Info

```plaintext
GET /api/admin/payment-providers/{providerType}/webhook
```

Get webhook information for a payment provider (admin only).

**Path Parameters:**

- `providerType` (required): Provider type (e.g., "stripe", "mobilepay")

**Status Codes:**

- `200 OK`: Webhook information retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Webhook not found
  "is_full": false
  }

````

**Request Body (Full Capture):**

```json
{
  "is_full": true
}
````

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
  "amount": 1500.0,
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
