# Commercify REST API Documentation

This document provides comprehensive documentation for the Commercify e-commerce backend API.

## Base URL

```
https://api.commercify.com/api
```

All API endpoints are prefixed with `/api` unless otherwise specified.

## Authentication

The API uses JWT (JSON Web Token) authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Response Format

All API responses follow a consistent format:

```json
{
  "success": boolean,
  "message": "string (optional)",
  "data": object | array | null,
  "error": "string (optional)",
  "pagination": {
    "page": number,
    "page_size": number,
    "total": number
  }
}
```

## Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists or conflict
- `500 Internal Server Error`: Server error

## Public Endpoints

### Health Check

```plaintext
GET /health
```

Health check endpoint for load balancers and monitoring.

**Response:**

```json
{
  "status": "healthy",
  "timestamp": "2025-07-07T10:30:45Z",
  "services": {
    "database": "healthy"
  }
}
```

### Authentication

#### Register User

```plaintext
POST /api/auth/register
```

Register a new user account.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "Password123!",
  "first_name": "John",
  "last_name": "Smith"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Smith",
      "role": "user",
      "created_at": "2025-07-07T10:30:45Z",
      "updated_at": "2025-07-07T10:30:45Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "",
    "expires_in": 3600
  }
}
```

#### Login

```plaintext
POST /api/auth/signin
```

Authenticate a user and retrieve a JWT token.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "Password123!"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Smith",
      "role": "user",
      "created_at": "2025-07-07T10:30:45Z",
      "updated_at": "2025-07-07T10:30:45Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "",
    "expires_in": 3600
  }
}
```

### Products

#### Get Product

```plaintext
GET /api/products/{productId}
```

Get a product by ID.

**Response:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Product Name",
    "description": "Product description",
    "currency": "USD",
    "price": 99.99,
    "sku": "PROD-001",
    "total_stock": 100,
    "category": "Electronics",
    "category_id": 1,
    "images": ["https://example.com/image1.jpg"],
    "has_variants": true,
    "active": true,
    "variants": [
      {
        "id": 1,
        "product_id": 1,
        "variant_name": "Size L",
        "sku": "PROD-001-L",
        "stock": 50,
        "attributes": {
          "size": "L",
          "color": "blue"
        },
        "images": ["https://example.com/variant1.jpg"],
        "is_default": true,
        "weight": 0.5,
        "price": 99.99,
        "currency": "USD",
        "created_at": "2025-07-07T10:30:45Z",
        "updated_at": "2025-07-07T10:30:45Z"
      }
    ],
    "created_at": "2025-07-07T10:30:45Z",
    "updated_at": "2025-07-07T10:30:45Z"
  }
}
```

#### Search Products

```plaintext
GET /api/products/search
```

Search for products with optional filters.

**Query Parameters:**

- `query` (string, optional): Search term
- `category_id` (number, optional): Filter by category ID
- `min_price` (number, optional): Minimum price filter
- `max_price` (number, optional): Maximum price filter
- `currency` (string, optional): Currency code
- `active_only` (boolean, optional): Show only active products
- `page` (number, optional): Page number (default: 1)
- `page_size` (number, optional): Items per page (default: 10)

### Categories

#### List Categories

```plaintext
GET /api/categories
```

Get all categories.

#### Get Category

```plaintext
GET /api/categories/{id}
```

Get a category by ID.

#### Get Child Categories

```plaintext
GET /api/categories/{id}/children
```

Get child categories of a parent category.

### Payment Providers

#### Get Available Payment Providers

```plaintext
GET /api/payment/providers
```

Get available payment providers.

### Discounts

#### Validate Discount Code

```plaintext
POST /api/discounts/validate
```

Validate a discount code.

**Request Body:**

```json
{
  "discount_code": "SUMMER2025"
}
```

### Currencies

#### List Enabled Currencies

```plaintext
GET /api/currencies
```

Get all enabled currencies.

#### Get Default Currency

```plaintext
GET /api/currencies/default
```

Get the default currency.

#### Convert Amount

```plaintext
POST /api/currencies/convert
```

Convert amount between currencies.

**Request Body:**

```json
{
  "amount": 100.0,
  "from_currency": "USD",
  "to_currency": "EUR"
}
```

### Shipping

#### Calculate Shipping Options

```plaintext
POST /api/shipping/options
```

Calculate available shipping options.

**Request Body:**

```json
{
  "address": {
    "address_line1": "123 Main St",
    "address_line2": "",
    "city": "New York",
    "state": "NY",
    "postal_code": "10001",
    "country": "US"
  },
  "order_value": 150.0,
  "order_weight": 2.5
}
```

### Checkout (Guest)

#### Get Checkout

```plaintext
GET /api/checkout
```

Get current checkout session.

#### Add to Checkout

```plaintext
POST /api/checkout/items
```

Add item to checkout.

**Request Body:**

```json
{
  "sku": "PROD-001-L",
  "quantity": 2,
  "currency": "USD"
}
```

#### Update Checkout Item

```plaintext
PUT /api/checkout/items/{sku}
```

Update checkout item quantity.

**Request Body:**

```json
{
  "quantity": 3
}
```

#### Remove from Checkout

```plaintext
DELETE /api/checkout/items/{sku}
```

Remove item from checkout.

#### Clear Checkout

```plaintext
DELETE /api/checkout
```

Clear entire checkout.

#### Set Shipping Address

```plaintext
PUT /api/checkout/shipping-address
```

Set checkout shipping address.

**Request Body:**

```json
{
  "address_line1": "123 Main St",
  "address_line2": "Apt 4B",
  "city": "New York",
  "state": "NY",
  "postal_code": "10001",
  "country": "US"
}
```

#### Set Billing Address

```plaintext
PUT /api/checkout/billing-address
```

Set checkout billing address.

#### Set Customer Details

```plaintext
PUT /api/checkout/customer-details
```

Set customer information.

**Request Body:**

```json
{
  "email": "customer@example.com",
  "phone": "+1234567890",
  "full_name": "John Smith"
}
```

#### Set Shipping Method

```plaintext
PUT /api/checkout/shipping-method
```

Set shipping method.

**Request Body:**

```json
{
  "shipping_method_id": 1
}
```

#### Set Currency

```plaintext
PUT /api/checkout/currency
```

Set checkout currency.

**Request Body:**

```json
{
  "currency": "EUR"
}
```

#### Apply Discount

```plaintext
POST /api/checkout/discount
```

Apply discount code to checkout.

**Request Body:**

```json
{
  "discount_code": "SUMMER2025"
}
```

#### Remove Discount

```plaintext
DELETE /api/checkout/discount
```

Remove applied discount from checkout.

#### Complete Checkout

```plaintext
POST /api/checkout/complete
```

Complete checkout and create order.

**Request Body:**

```json
{
  "payment_provider": "stripe",
  "payment_data": {
    "card_details": {
      "card_number": "4242424242424242",
      "expiry_month": 12,
      "expiry_year": 2025,
      "cvv": "123",
      "cardholder_name": "John Smith"
    }
  }
}
```

## Authenticated Endpoints

### User Profile

#### Get User Profile

```plaintext
GET /api/users/me
```

Get authenticated user's profile.

#### Update User Profile

```plaintext
PUT /api/users/me
```

Update authenticated user's profile.

**Request Body:**

```json
{
  "first_name": "John",
  "last_name": "Smith"
}
```

#### Change Password

```plaintext
PUT /api/users/me/password
```

Change user password.

**Request Body:**

```json
{
  "current_password": "oldPassword123!",
  "new_password": "newPassword123!"
}
```

### Orders

#### List User Orders

```plaintext
GET /api/orders
```

List orders for authenticated user.

**Query Parameters:**

- `page` (number, optional): Page number
- `pageSize` (number, optional): Items per page

#### Get Order

```plaintext
GET /api/orders/{orderId}
```

Get order by ID (accessible by order owner or via checkout session).

## Admin Endpoints

All admin endpoints require authentication and admin role.

### User Management

#### List All Users

```plaintext
GET /api/admin/users
```

List all users (admin only).

**Query Parameters:**

- `page` (number, optional): Page number
- `page_size` (number, optional): Items per page

### Order Management

#### List All Orders

```plaintext
GET /api/admin/orders
```

List all orders (admin only).

**Query Parameters:**

- `page` (number, optional): Page number
- `pageSize` (number, optional): Items per page
- `status` (string, optional): Filter by order status

#### Update Order Status

```plaintext
PUT /api/admin/orders/{orderId}/status
```

Update order status (admin only).

### Checkout Management

#### List Checkouts

```plaintext
GET /api/admin/checkouts
```

List all checkouts (admin only).

#### Get Checkout

```plaintext
GET /api/admin/checkouts/{checkoutId}
```

Get checkout by ID (admin only).

#### Delete Checkout

```plaintext
DELETE /api/admin/checkouts/{checkoutId}
```

Delete checkout (admin only).

### Currency Management

#### List All Currencies

```plaintext
GET /api/admin/currencies/all
```

List all currencies including disabled ones (admin only).

#### Create Currency

```plaintext
POST /api/admin/currencies
```

Create new currency (admin only).

**Request Body:**

```json
{
  "code": "EUR",
  "name": "Euro",
  "symbol": "€",
  "exchange_rate": 0.85,
  "is_enabled": true,
  "is_default": false
}
```

#### Update Currency

```plaintext
PUT /api/admin/currencies
```

Update currency (admin only).

#### Delete Currency

```plaintext
DELETE /api/admin/currencies
```

Delete currency (admin only).

#### Set Default Currency

```plaintext
PUT /api/admin/currencies/default
```

Set default currency (admin only).

### Category Management

#### Create Category

```plaintext
POST /api/admin/categories
```

Create new category (admin only).

**Request Body:**

```json
{
  "name": "Electronics",
  "description": "Electronic devices and accessories",
  "parent_id": 1
}
```

#### Update Category

```plaintext
PUT /api/admin/categories/{id}
```

Update category (admin only).

#### Delete Category

```plaintext
DELETE /api/admin/categories/{id}
```

Delete category (admin only).

### Product Management

#### List Products

```plaintext
GET /api/admin/products
```

List all products (admin only).

**Query Parameters:**

- `page` (number, optional): Page number
- `page_size` (number, optional): Items per page
- `query` (string, optional): Search term
- `category_id` (number, optional): Filter by category
- `min_price` (number, optional): Minimum price
- `max_price` (number, optional): Maximum price
- `currency` (string, optional): Currency code
- `active_only` (boolean, optional): Show only active products

#### Create Product

```plaintext
POST /api/admin/products
```

Create new product (admin only).

**Request Body:**

```json
{
  "name": "Product Name",
  "description": "Product description",
  "currency": "USD",
  "category_id": 1,
  "images": ["https://example.com/image1.jpg"],
  "active": true,
  "variants": [
    {
      "sku": "PROD-001-L",
      "stock": 100,
      "attributes": [
        { "name": "size", "value": "L" },
        { "name": "color", "value": "blue" }
      ],
      "images": ["https://example.com/variant1.jpg"],
      "is_default": true,
      "weight": 0.5,
      "price": 99.99
    }
  ]
}
```

#### Update Product

```plaintext
PUT /api/admin/products/{productId}
```

Update product (admin only).

#### Delete Product

```plaintext
DELETE /api/admin/products/{productId}
```

Delete product (admin only).

#### Add Product Variant

```plaintext
POST /api/admin/products/{productId}/variants
```

Add variant to product (admin only).

#### Update Product Variant

```plaintext
PUT /api/admin/products/{productId}/variants/{variantId}
```

Update product variant (admin only).

#### Delete Product Variant

```plaintext
DELETE /api/admin/products/{productId}/variants/{variantId}
```

Delete product variant (admin only).

### Shipping Management

#### Create Shipping Method

```plaintext
POST /api/admin/shipping/methods
```

Create new shipping method (admin only).

**Request Body:**

```json
{
  "name": "Standard Shipping",
  "description": "5-7 business days",
  "estimated_delivery_days": 7
}
```

#### Create Shipping Zone

```plaintext
POST /api/admin/shipping/zones
```

Create new shipping zone (admin only).

**Request Body:**

```json
{
  "name": "US Zone",
  "description": "United States shipping zone",
  "countries": ["US"],
  "states": ["NY", "CA", "TX"],
  "zip_codes": ["10001", "90210"]
}
```

#### Create Shipping Rate

```plaintext
POST /api/admin/shipping/rates
```

Create new shipping rate (admin only).

**Request Body:**

```json
{
  "shipping_method_id": 1,
  "shipping_zone_id": 1,
  "base_rate": 9.99,
  "min_order_value": 0,
  "free_shipping_threshold": 100.0,
  "active": true
}
```

#### Create Weight-Based Rate

```plaintext
POST /api/admin/shipping/rates/weight
```

Create weight-based shipping rate (admin only).

**Request Body:**

```json
{
  "shipping_rate_id": 1,
  "min_weight": 0,
  "max_weight": 5.0,
  "rate": 5.99
}
```

#### Create Value-Based Rate

```plaintext
POST /api/admin/shipping/rates/value
```

Create value-based shipping rate (admin only).

**Request Body:**

```json
{
  "shipping_rate_id": 1,
  "min_order_value": 0,
  "max_order_value": 50.0,
  "rate": 9.99
}
```

### Discount Management

#### Create Discount

```plaintext
POST /api/admin/discounts
```

Create new discount (admin only).

**Request Body:**

```json
{
  "code": "SUMMER2025",
  "type": "basket",
  "method": "percentage",
  "value": 15.0,
  "min_order_value": 50.0,
  "max_discount_value": 30.0,
  "product_ids": [],
  "category_ids": [],
  "start_date": "2025-05-01T00:00:00Z",
  "end_date": "2025-08-31T23:59:59Z",
  "usage_limit": 500
}
```

#### Get Discount

```plaintext
GET /api/admin/discounts/{discountId}
```

Get discount by ID (admin only).

#### Update Discount

```plaintext
PUT /api/admin/discounts/{discountId}
```

Update discount (admin only).

#### Delete Discount

```plaintext
DELETE /api/admin/discounts/{discountId}
```

Delete discount (admin only).

#### List Discounts

```plaintext
GET /api/admin/discounts
```

List all discounts (admin only).

#### List Active Discounts

```plaintext
GET /api/admin/discounts/active
```

List active discounts (admin only).

#### Apply Discount to Order

```plaintext
POST /api/admin/discounts/apply/{orderId}
```

Apply discount to order (admin only).

#### Remove Discount from Order

```plaintext
DELETE /api/admin/discounts/remove/{orderId}
```

Remove discount from order (admin only).

### Payment Management

#### Capture Payment

```plaintext
POST /api/admin/payments/{paymentId}/capture
```

Capture authorized payment (admin only).

#### Cancel Payment

```plaintext
POST /api/admin/payments/{paymentId}/cancel
```

Cancel payment (admin only).

#### Refund Payment

```plaintext
POST /api/admin/payments/{paymentId}/refund
```

Refund payment (admin only).

#### Force Approve MobilePay Payment

```plaintext
POST /api/admin/payments/{paymentId}/force-approve
```

Force approve MobilePay payment (admin only).

### Payment Provider Management

#### Get Payment Providers

```plaintext
GET /api/admin/payment-providers
```

Get all payment providers (admin only).

#### Get Enabled Payment Providers

```plaintext
GET /api/admin/payment-providers/enabled
```

Get enabled payment providers (admin only).

#### Enable Payment Provider

```plaintext
PUT /api/admin/payment-providers/{providerType}/enable
```

Enable/disable payment provider (admin only).

**Request Body:**

```json
{
  "enabled": true
}
```

#### Update Provider Configuration

```plaintext
PUT /api/admin/payment-providers/{providerType}/configuration
```

Update payment provider configuration (admin only).

#### Register Webhook

```plaintext
POST /api/admin/payment-providers/{providerType}/webhook
```

Register webhook for payment provider (admin only).

#### Delete Webhook

```plaintext
DELETE /api/admin/payment-providers/{providerType}/webhook
```

Delete webhook for payment provider (admin only).

#### Get Webhook Info

```plaintext
GET /api/admin/payment-providers/{providerType}/webhook
```

Get webhook information for payment provider (admin only).

### Email Testing

#### Test Email

```plaintext
POST /api/admin/test/email
```

Send test email (admin only).

## Webhook Endpoints

### Stripe Webhook

```plaintext
POST /api/webhooks/stripe
```

Stripe webhook endpoint for payment events.

### MobilePay Webhook

```plaintext
POST /api/webhooks/mobilepay
```

MobilePay webhook endpoint for payment events.

## Error Responses

All error responses follow this format:

```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

Common error scenarios:

- Authentication required (401)
- Insufficient permissions (403)
- Resource not found (404)
- Validation errors (400)
- Server errors (500)
