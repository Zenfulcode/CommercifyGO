# Webhook API Documentation

This document outlines the webhook endpoints for receiving payment provider notifications.

## Webhook Endpoints

Webhook endpoints are designed for server-to-server communication and do not require authentication or CORS handling. They are used by payment providers to notify the system about payment events.

### Stripe Webhook

```plaintext
POST /api/webhooks/stripe
```

Stripe webhook endpoint for receiving payment events from Stripe.

**Description:**
This endpoint handles incoming webhook events from Stripe, including:

- Payment authorization events
- Payment capture events
- Payment failure events
- Payment cancellation events

**Headers:**

- `Stripe-Signature`: Webhook signature for verification

**Request Body:**
The request body contains the Stripe event data in JSON format. The exact structure depends on the event type.

**Response:**

- `200 OK`: Event processed successfully
- `400 Bad Request`: Invalid signature or malformed payload
- `500 Internal Server Error`: Error processing the event

### MobilePay Webhook

```plaintext
POST /api/webhooks/mobilepay
```

MobilePay webhook endpoint for receiving payment events from MobilePay.

**Description:**
This endpoint handles incoming webhook events from MobilePay, including:

- Payment authorization events
- Payment capture events
- Payment cancellation events
- Payment expiration events
- Payment refund events

**Headers:**

- Content verification headers as required by MobilePay

**Request Body:**
The request body contains the MobilePay event data in JSON format. The exact structure depends on the event type and follows the MobilePay webhook specification.

**Response:**

- `200 OK`: Event processed successfully
- `400 Bad Request`: Invalid signature or malformed payload
- `500 Internal Server Error`: Error processing the event

## Security

### Signature Verification

Both webhook endpoints implement signature verification to ensure the authenticity of incoming requests:

- **Stripe**: Uses the `Stripe-Signature` header with HMAC-SHA256 verification
- **MobilePay**: Uses MobilePay's signature verification mechanism

### Event Processing

Webhook events are processed asynchronously and include:

1. **Signature Verification**: Verify the request comes from the legitimate payment provider
2. **Event Parsing**: Parse the event data and extract relevant information
3. **Order Updates**: Update order status and payment transactions
4. **Database Recording**: Record payment transactions in the database
5. **Response**: Return appropriate HTTP status code

### Error Handling

If webhook processing fails:

- The endpoint returns an appropriate error status code
- The error is logged for debugging purposes
- Payment providers typically retry failed webhook deliveries

### Event Types

#### Stripe Events

- `payment_intent.succeeded`: Payment successfully captured
- `payment_intent.payment_failed`: Payment failed
- `payment_intent.canceled`: Payment was canceled

#### MobilePay Events

- `payment.authorized`: Payment was authorized
- `payment.captured`: Payment was captured
- `payment.cancelled`: Payment was cancelled
- `payment.expired`: Payment authorization expired
- `payment.refunded`: Payment was refunded

## Testing Webhooks

### Development

During development, you can use tools like ngrok to expose your local webhook endpoints to payment providers for testing.

### Production

In production, ensure webhook endpoints are:

- Accessible over HTTPS
- Properly configured in the payment provider dashboard
- Monitoring webhook delivery success rates
