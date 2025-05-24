# Checkout API Documentation

This document outlines the Checkout API endpoints for the Commercify e-commerce system.

## Guest Checkout Endpoints

The following endpoints support guest checkout functionality, allowing users to create and manage checkout sessions without authentication.

### Get Current Checkout

```plaintext
GET /api/checkout
```

Retrieves the current checkout session for a user. If no checkout exists, a new one will be created.

**Response Body:**

```json
{
  "id": 123,
  "session_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "items": [
    {
      "id": 1,
      "product_id": 42,
      "variant_id": 7,
      "product_name": "Organic Cotton T-Shirt",
      "variant_name": "Medium / Blue",
      "sku": "TS-BL-M",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.3,
      "subtotal": 49.98,
      "created_at": "2025-05-24T10:30:00Z",
      "updated_at": "2025-05-24T10:30:00Z"
    }
  ],
  "status": "active",
  "shipping_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "billing_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "shipping_method_id": 0,
  "shipping_method": null,
  "payment_provider": "",
  "total_amount": 49.98,
  "shipping_cost": 0,
  "total_weight": 0.6,
  "customer_details": {
    "email": "",
    "phone": "",
    "full_name": ""
  },
  "currency": "USD",
  "discount_code": "",
  "discount_amount": 0,
  "final_amount": 49.98,
  "applied_discount": null,
  "created_at": "2025-05-24T10:30:00Z",
  "updated_at": "2025-05-24T10:30:00Z",
  "last_activity_at": "2025-05-24T10:30:00Z",
  "expires_at": "2025-05-24T11:30:00Z"
}
```

**Status Codes:**

- `200 OK`: Checkout retrieved or created successfully
- `500 Internal Server Error`: Server error

### Add Item to Checkout

```plaintext
POST /api/checkout/items
```

Adds a product item to the current checkout session.

**Request Body:**

```json
{
  "product_id": 42,
  "variant_id": 7,
  "quantity": 1
}
```

**Response Body:**

```json
{
  "id": 123,
  "session_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "items": [
    {
      "id": 1,
      "product_id": 42,
      "variant_id": 7,
      "product_name": "Organic Cotton T-Shirt",
      "variant_name": "Medium / Blue",
      "sku": "TS-BL-M",
      "price": 24.99,
      "quantity": 1,
      "weight": 0.3,
      "subtotal": 24.99,
      "created_at": "2025-05-24T10:30:00Z",
      "updated_at": "2025-05-24T10:30:00Z"
    }
  ],
  "status": "active",
  "shipping_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "billing_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "total_amount": 24.99,
  "shipping_cost": 0,
  "total_weight": 0.3,
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 24.99,
  "created_at": "2025-05-24T10:30:00Z",
  "updated_at": "2025-05-24T10:30:00Z",
  "last_activity_at": "2025-05-24T10:30:00Z",
  "expires_at": "2025-05-24T11:30:00Z"
}
```

**Status Codes:**

- `200 OK`: Item added successfully
- `400 Bad Request`: Invalid request body or product
- `404 Not Found`: Product not found
- `500 Internal Server Error`: Server error

### Update Checkout Item

```plaintext
PUT /api/checkout/items/{productId}
```

Updates the quantity or variant of an item in the current checkout.

**Path Parameters:**

- `productId`: ID of the product to update

**Request Body:**

```json
{
  "quantity": 2,
  "variant_id": 8
}
```

**Response Body:**

```json
{
  "id": 123,
  "session_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "items": [
    {
      "id": 1,
      "product_id": 42,
      "variant_id": 8,
      "product_name": "Organic Cotton T-Shirt",
      "variant_name": "Large / Blue",
      "sku": "TS-BL-L",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.35,
      "subtotal": 49.98,
      "created_at": "2025-05-24T10:30:00Z",
      "updated_at": "2025-05-24T10:35:00Z"
    }
  ],
  "status": "active",
  "total_amount": 49.98,
  "shipping_cost": 0,
  "total_weight": 0.7,
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 49.98,
  "updated_at": "2025-05-24T10:35:00Z",
  "last_activity_at": "2025-05-24T10:35:00Z",
  "expires_at": "2025-05-24T11:35:00Z"
}
```

**Status Codes:**

- `200 OK`: Item updated successfully
- `400 Bad Request`: Invalid request body
- `404 Not Found`: Product not found in checkout
- `500 Internal Server Error`: Server error

### Remove Item from Checkout

```plaintext
DELETE /api/checkout/items/{productId}
```

Removes an item from the current checkout session.

**Path Parameters:**

- `productId`: ID of the product to remove

**Response Body:**

```json
{
  "id": 123,
  "session_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "items": [],
  "status": "active",
  "shipping_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "billing_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "total_amount": 0,
  "shipping_cost": 0,
  "total_weight": 0,
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 0,
  "updated_at": "2025-05-24T10:40:00Z",
  "last_activity_at": "2025-05-24T10:40:00Z",
  "expires_at": "2025-05-24T11:40:00Z"
}
```

**Status Codes:**

- `200 OK`: Item removed successfully
- `404 Not Found`: Product not found in checkout
- `500 Internal Server Error`: Server error

### Clear Checkout

```plaintext
DELETE /api/checkout
```

Removes all items from the current checkout session.

**Response Body:**

```json
{
  "id": 123,
  "session_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "items": [],
  "status": "active",
  "shipping_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "billing_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "total_amount": 0,
  "shipping_cost": 0,
  "total_weight": 0,
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 0,
  "updated_at": "2025-05-24T10:45:00Z",
  "last_activity_at": "2025-05-24T10:45:00Z",
  "expires_at": "2025-05-24T11:45:00Z"
}
```

**Status Codes:**

- `200 OK`: Checkout cleared successfully
- `404 Not Found`: Checkout not found
- `500 Internal Server Error`: Server error

### Set Shipping Address

```plaintext
PUT /api/checkout/shipping-address
```

Sets the shipping address for the current checkout.

**Request Body:**

```json
{
  "address_line1": "123 Main Street",
  "address_line2": "Apt 4B",
  "city": "Springfield",
  "state": "IL",
  "postal_code": "62704",
  "country": "US"
}
```

**Response Body:**

```json
{
  "id": 123,
  "session_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "items": [],
  "status": "active",
  "shipping_address": {
    "address_line1": "123 Main Street",
    "address_line2": "Apt 4B",
    "city": "Springfield",
    "state": "IL",
    "postal_code": "62704",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "total_amount": 0,
  "shipping_cost": 0,
  "total_weight": 0,
  "currency": "USD",
  "updated_at": "2025-05-24T10:50:00Z",
  "last_activity_at": "2025-05-24T10:50:00Z",
  "expires_at": "2025-05-24T11:50:00Z"
}
```

**Status Codes:**

- `200 OK`: Shipping address set successfully
- `400 Bad Request`: Invalid address data
- `404 Not Found`: Checkout not found
- `500 Internal Server Error`: Server error

### Set Billing Address

```plaintext
PUT /api/checkout/billing-address
```

Sets the billing address for the current checkout.

**Request Body:**

```json
{
  "address_line1": "456 Commerce Ave",
  "address_line2": "Suite 300",
  "city": "Springfield",
  "state": "IL",
  "postal_code": "62704",
  "country": "US"
}
```

**Response Body:**

```json
{
  "id": 123,
  "session_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "items": [],
  "status": "active",
  "shipping_address": {
    "address_line1": "123 Main Street",
    "address_line2": "Apt 4B",
    "city": "Springfield",
    "state": "IL",
    "postal_code": "62704",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "456 Commerce Ave",
    "address_line2": "Suite 300",
    "city": "Springfield",
    "state": "IL",
    "postal_code": "62704",
    "country": "US"
  },
  "total_amount": 0,
  "shipping_cost": 0,
  "total_weight": 0,
  "currency": "USD",
  "updated_at": "2025-05-24T10:55:00Z",
  "last_activity_at": "2025-05-24T10:55:00Z",
  "expires_at": "2025-05-24T11:55:00Z"
}
```

**Status Codes:**

- `200 OK`: Billing address set successfully
- `400 Bad Request`: Invalid address data
- `404 Not Found`: Checkout not found
- `500 Internal Server Error`: Server error

### Set Customer Details

```plaintext
PUT /api/checkout/customer-details
```

Sets the customer contact information for the current checkout.

**Request Body:**

```json
{
  "email": "customer@example.com",
  "phone": "+1234567890",
  "full_name": "John Doe"
}
```

**Response Body:**

```json
{
  "id": 123,
  "session_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "items": [],
  "status": "active",
  "shipping_address": {
    "address_line1": "123 Main Street",
    "address_line2": "Apt 4B",
    "city": "Springfield",
    "state": "IL", 
    "postal_code": "62704",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "456 Commerce Ave",
    "address_line2": "Suite 300",
    "city": "Springfield",
    "state": "IL",
    "postal_code": "62704",
    "country": "US"
  },
  "customer_details": {
    "email": "customer@example.com",
    "phone": "+1234567890",
    "full_name": "John Doe"
  },
  "total_amount": 0,
  "shipping_cost": 0,
  "total_weight": 0,
  "currency": "USD",
  "updated_at": "2025-05-24T11:00:00Z",
  "last_activity_at": "2025-05-24T11:00:00Z",
  "expires_at": "2025-05-24T12:00:00Z"
}
```

**Status Codes:**

- `200 OK`: Customer details set successfully
- `400 Bad Request`: Invalid customer data
- `404 Not Found`: Checkout not found
- `500 Internal Server Error`: Server error

### Set Shipping Method

```plaintext
PUT /api/checkout/shipping-method
```

Sets the shipping method for the current checkout.

**Request Body:**

```json
{
  "shipping_method_id": 1
}
```

**Response Body:**

```json
{
  "id": 123,
  "session_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "items": [],
  "status": "active",
  "shipping_address": {
    "address_line1": "123 Main Street",
    "address_line2": "Apt 4B",
    "city": "Springfield",
    "state": "IL",
    "postal_code": "62704",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "456 Commerce Ave",
    "address_line2": "Suite 300",
    "city": "Springfield",
    "state": "IL",
    "postal_code": "62704",
    "country": "US"
  },
  "shipping_method_id": 1,
  "shipping_method": {
    "id": 1,
    "name": "Standard Shipping",
    "description": "Delivery in 5-7 business days",
    "cost": 5.99
  },
  "total_amount": 0,
  "shipping_cost": 5.99,
  "total_weight": 0,
  "currency": "USD",
  "final_amount": 5.99,
  "updated_at": "2025-05-24T11:05:00Z",
  "last_activity_at": "2025-05-24T11:05:00Z",
  "expires_at": "2025-05-24T12:05:00Z"
}
```

**Status Codes:**

- `200 OK`: Shipping method set successfully
- `400 Bad Request`: Invalid shipping method
- `404 Not Found`: Checkout or shipping method not found
- `500 Internal Server Error`: Server error

### Apply Discount

```plaintext
POST /api/checkout/discount
```

Applies a discount code to the current checkout.

**Request Body:**

```json
{
  "discount_code": "SUMMER25"
}
```

**Response Body:**

```json
{
  "id": 123,
  "session_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "items": [
    {
      "id": 1,
      "product_id": 42,
      "variant_id": 7,
      "product_name": "Organic Cotton T-Shirt",
      "variant_name": "Medium / Blue",
      "sku": "TS-BL-M",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.3,
      "subtotal": 49.98,
      "created_at": "2025-05-24T10:30:00Z",
      "updated_at": "2025-05-24T10:30:00Z"
    }
  ],
  "status": "active",
  "shipping_method_id": 1,
  "shipping_method": {
    "id": 1,
    "name": "Standard Shipping",
    "description": "Delivery in 5-7 business days",
    "cost": 5.99
  },
  "total_amount": 49.98,
  "shipping_cost": 5.99,
  "total_weight": 0.6,
  "currency": "USD",
  "discount_code": "SUMMER25",
  "discount_amount": 12.50,
  "final_amount": 43.47,
  "applied_discount": {
    "id": 5,
    "code": "SUMMER25",
    "type": "basket",
    "method": "percentage",
    "value": 25,
    "amount": 12.50
  },
  "updated_at": "2025-05-24T11:10:00Z",
  "last_activity_at": "2025-05-24T11:10:00Z",
  "expires_at": "2025-05-24T12:10:00Z"
}
```

**Status Codes:**

- `200 OK`: Discount applied successfully
- `400 Bad Request`: Invalid or expired discount code
- `404 Not Found`: Checkout not found
- `500 Internal Server Error`: Server error

### Remove Discount

```plaintext
DELETE /api/checkout/discount
```

Removes any applied discount code from the current checkout.

**Response Body:**

```json
{
  "id": 123,
  "session_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "items": [
    {
      "id": 1,
      "product_id": 42,
      "variant_id": 7,
      "product_name": "Organic Cotton T-Shirt",
      "variant_name": "Medium / Blue",
      "sku": "TS-BL-M",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.3,
      "subtotal": 49.98,
      "created_at": "2025-05-24T10:30:00Z",
      "updated_at": "2025-05-24T10:30:00Z"
    }
  ],
  "status": "active",
  "shipping_method_id": 1,
  "shipping_method": {
    "id": 1,
    "name": "Standard Shipping",
    "description": "Delivery in 5-7 business days",
    "cost": 5.99
  },
  "total_amount": 49.98,
  "shipping_cost": 5.99,
  "total_weight": 0.6,
  "currency": "USD",
  "discount_code": "",
  "discount_amount": 0,
  "final_amount": 55.97,
  "applied_discount": null,
  "updated_at": "2025-05-24T11:15:00Z",
  "last_activity_at": "2025-05-24T11:15:00Z",
  "expires_at": "2025-05-24T12:15:00Z"
}
```

**Status Codes:**

- `200 OK`: Discount removed successfully
- `404 Not Found`: Checkout not found
- `500 Internal Server Error`: Server error

### Complete Checkout

```plaintext
### Complete Checkout

```plaintext
POST /api/checkout/complete
```

**Request Body:**

```json
{
  "payment_provider": "stripe",
  "payment_data": {
    "card_details": {
      "card_number": "4111111111111111",
      "expiry_month": 12,
      "expiry_year": 2027,
      "cvv": "123",
      "cardholder_name": "John Doe",
      "token": "tok_visa_2024"
    }
  },
  "redirect_url": "https://example.com/order-confirmation"
}
```

Alternatively, for mobile payment methods:

```json
{
  "payment_provider": "mobilepay",
  "payment_data": {
    "phone_number": "+4512345678"
  },
  "redirect_url": "https://example.com/order-confirmation"
}
```

**Response Body:**

```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": 456,
    "user_id": 123,
    "order_number": "ORD-2025-0001",
    "status": "pending",
    "total_amount": 49.95,
    "final_amount": 49.95,
    "currency": "USD",
    "items": [
      {
        "id": 1,
        "product_id": 42,
        "variant_id": 7,
        "product_name": "Organic Cotton T-Shirt",
        "variant_name": "Medium / Blue",
        "sku": "TS-BL-M",
        "price": 24.99,
        "quantity": 2,
        "subtotal": 49.98
      }
    ],
    "shipping_address": {
      "address_line1": "123 Main Street",
      "address_line2": "Apt 4B",
      "city": "Springfield",
      "state": "IL",
      "postal_code": "62704",
      "country": "US"
    },
    "billing_address": {
      "address_line1": "456 Commerce Ave",
      "address_line2": "Suite 300",
      "city": "Springfield",
      "state": "IL",
      "postal_code": "62704",
      "country": "US"
    },
    "customer_details": {
      "email": "customer@example.com",
      "phone": "+1234567890",
      "full_name": "John Doe"
    },
    "shipping_method": "Standard Shipping",
    "shipping_cost": 5.99,
    "subtotal": 49.98,
    "total": 55.97,
    "discount_code": "",
    "discount_amount": 0,
    "final_amount": 55.97,
    "currency": "USD",
    "payment_provider": "stripe",
    "payment_status": "pending",
    "created_at": "2025-05-24T11:20:00Z"
  },
  "action_required": false,
  "redirect_url": ""
}
```

**Status Codes:**

- `201 Created`: Order created and payment processed successfully
- `400 Bad Request`: Invalid request body or payment processing failed
- `401 Unauthorized`: Not authenticated
- `404 Not Found`: Checkout session not found or expired
- `409 Conflict`: Order already exists for this checkout session
- `500 Internal Server Error`: Server error
- `402 Payment Required`: Payment failed
- `404 Not Found`: Checkout not found
- `409 Conflict`: Checkout is already completed
- `422 Unprocessable Entity`: Invalid payment data
- `500 Internal Server Error`: Server error

## Admin Checkout Endpoints

The following endpoints are available for administrative functions and require admin authentication.

### List All Checkouts (Admin)

```plaintext
GET /api/admin/checkouts
```

Returns a paginated list of all checkout sessions in the system.

**Query Parameters:**

- `offset`: Starting index for pagination (default: 0)
- `limit`: Maximum number of results to return (default: 20)
- `status`: Filter by checkout status (optional)

**Response Body:**

```json
{
  "items": [
    {
      "id": 123,
      "session_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
      "status": "active",
      "total_amount": 49.98,
      "discount_amount": 0,
      "final_amount": 55.97,
      "currency": "USD",
      "created_at": "2025-05-24T10:30:00Z",
      "updated_at": "2025-05-24T11:20:00Z",
      "expires_at": "2025-05-24T12:20:00Z"
    },
    {
      "id": 122,
      "session_id": "2cb94f54-6261-4522-a3fc-1b832f55ddd1",
      "status": "completed",
      "total_amount": 129.95,
      "discount_amount": 25.99,
      "final_amount": 110.94,
      "currency": "USD",
      "created_at": "2025-05-24T09:15:00Z",
      "updated_at": "2025-05-24T09:45:00Z",
      "expires_at": "2025-05-24T10:45:00Z",
      "completed_at": "2025-05-24T09:45:00Z",
      "converted_order_id": 455
    }
  ],
  "total": 52,
  "offset": 0,
  "limit": 20
}
```

**Status Codes:**

- `200 OK`: Checkouts retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `500 Internal Server Error`: Server error

### Get Checkout by ID (Admin)

```plaintext
GET /api/admin/checkouts/{checkoutId}
```

Returns detailed information about a specific checkout session.

**Path Parameters:**

- `checkoutId`: The unique identifier of the checkout

**Response Body:**

```json
{
  "id": 123,
  "session_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "items": [
    {
      "id": 1,
      "product_id": 42,
      "variant_id": 7,
      "product_name": "Organic Cotton T-Shirt",
      "variant_name": "Medium / Blue",
      "sku": "TS-BL-M",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.3,
      "subtotal": 49.98,
      "created_at": "2025-05-24T10:30:00Z",
      "updated_at": "2025-05-24T10:30:00Z"
    }
  ],
  "status": "active",
  "shipping_address": {
    "address_line1": "123 Main Street",
    "address_line2": "Apt 4B",
    "city": "Springfield",
    "state": "IL",
    "postal_code": "62704",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "456 Commerce Ave",
    "address_line2": "Suite 300",
    "city": "Springfield",
    "state": "IL",
    "postal_code": "62704",
    "country": "US"
  },
  "shipping_method_id": 1,
  "shipping_method": {
    "id": 1,
    "name": "Standard Shipping",
    "description": "Delivery in 5-7 business days",
    "cost": 5.99
  },
  "total_amount": 49.98,
  "shipping_cost": 5.99,
  "total_weight": 0.6,
  "customer_details": {
    "email": "customer@example.com",
    "phone": "+1234567890",
    "full_name": "John Doe"
  },
  "currency": "USD",
  "discount_code": "",
  "discount_amount": 0,
  "final_amount": 55.97,
  "created_at": "2025-05-24T10:30:00Z",
  "updated_at": "2025-05-24T11:20:00Z",
  "last_activity_at": "2025-05-24T11:20:00Z",
  "expires_at": "2025-05-24T12:20:00Z"
}
```

**Status Codes:**

- `200 OK`: Checkout retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Checkout not found
- `500 Internal Server Error`: Server error

### Delete Checkout (Admin)

```plaintext
DELETE /api/admin/checkouts/{checkoutId}
```

Deletes a specific checkout session from the system.

**Path Parameters:**

- `checkoutId`: The unique identifier of the checkout

**Status Codes:**

- `204 No Content`: Checkout deleted successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Checkout not found
- `500 Internal Server Error`: Server error
