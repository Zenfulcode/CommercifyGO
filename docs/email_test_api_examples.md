# Email Test API Examples

## Test Email Endpoint

### POST /api/admin/test/email

Test endpoint to send order confirmation and notification emails to a specified email address.

#### Example Request - JSON Body

```bash
curl -X POST http://localhost:8080/api/admin/test/email \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "email": "test@example.com"
  }'
```

#### Example Request - Query Parameter

```bash
curl -X POST "http://localhost:8080/api/admin/test/email?email=test@example.com" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

#### Example Response - Success

```json
{
  "success": true,
  "message": "Both order confirmation and notification emails sent successfully",
  "details": {
    "target_email": "test@example.com",
    "order_id": "12345"
  }
}
```

#### Example Response - Validation Error

```json
{
  "success": false,
  "errors": ["Invalid email format"]
}
```

#### Example Response - Email Send Error

```json
{
  "success": false,
  "errors": [
    "Order confirmation: failed to send email: SMTP connection failed",
    "Order notification: failed to send email: SMTP connection failed"
  ]
}
```

## Email Content

The test emails use mock order data with the following characteristics:

- Order ID: 12345
- Order Number: ORD-12345
- Customer: John Doe
- Total Amount: $83.00 (after $15.00 discount)
- Shipping Cost: $8.50
- Items:
  - Test Product 1 (SKU: TEST-001) - Quantity: 2 - $25.00 each
  - Test Product 2 (SKU: TEST-002) - Quantity: 1 - $49.50 each
- Discount: TEST15 - $15.00 off
- Shipping Address: 123 Test Street, Apt 4B, Test City, Test State, 12345, Test Country
- Billing Address: 456 Billing Ave, Billing City, Billing State, 67890, Billing Country

## Notes

- This endpoint is only accessible to admin users
- If no email is provided in the request, it falls back to the admin email from configuration
- Both order confirmation and notification emails are sent to the specified address
- The endpoint validates email format before attempting to send emails
- Individual email failures are reported in the error array
