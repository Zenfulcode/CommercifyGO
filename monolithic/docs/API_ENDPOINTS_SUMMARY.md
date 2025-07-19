# API Endpoints Summary

This document provides a quick overview of all available API endpoints in the Commercify system.

## Base URL

```
/api
```

## Public Endpoints

### Health Check

- `GET /health` - Health check endpoint

### Authentication

- `POST /api/auth/register` - Register new user
- `POST /api/auth/signin` - User login

### Products

- `GET /api/products/{productId}` - Get product by ID
- `GET /api/products/search` - Search products

### Categories

- `GET /api/categories` - List all categories
- `GET /api/categories/{id}` - Get category by ID
- `GET /api/categories/{id}/children` - Get child categories

### Payment Providers

- `GET /api/payment/providers` - Get available payment providers

### Discounts

- `POST /api/discounts/validate` - Validate discount code

### Currencies

- `GET /api/currencies` - List enabled currencies
- `GET /api/currencies/default` - Get default currency
- `POST /api/currencies/convert` - Convert amount between currencies

### Shipping

- `POST /api/shipping/options` - Calculate shipping options

### Checkout (Guest)

- `GET /api/checkout` - Get current checkout
- `POST /api/checkout/items` - Add item to checkout
- `PUT /api/checkout/items/{sku}` - Update checkout item
- `DELETE /api/checkout/items/{sku}` - Remove item from checkout
- `DELETE /api/checkout` - Clear checkout
- `PUT /api/checkout/shipping-address` - Set shipping address
- `PUT /api/checkout/billing-address` - Set billing address
- `PUT /api/checkout/customer-details` - Set customer details
- `PUT /api/checkout/shipping-method` - Set shipping method
- `PUT /api/checkout/currency` - Set checkout currency
- `POST /api/checkout/discount` - Apply discount
- `DELETE /api/checkout/discount` - Remove discount
- `POST /api/checkout/complete` - Complete checkout

## Authenticated User Endpoints

### User Profile

- `GET /api/users/me` - Get user profile
- `PUT /api/users/me` - Update user profile
- `PUT /api/users/me/password` - Change password

### Orders

- `GET /api/orders` - List user orders
- `GET /api/orders/{orderId}` - Get order by ID (also accessible via checkout session)

## Admin Endpoints

All admin endpoints require authentication and admin role.

### User Management

- `GET /api/admin/users` - List all users

### Order Management

- `GET /api/admin/orders` - List all orders
- `PUT /api/admin/orders/{orderId}/status` - Update order status

### Checkout Management

- `GET /api/admin/checkouts` - List all checkouts
- `GET /api/admin/checkouts/{checkoutId}` - Get checkout by ID
- `DELETE /api/admin/checkouts/{checkoutId}` - Delete checkout

### Currency Management

- `GET /api/admin/currencies/all` - List all currencies
- `POST /api/admin/currencies` - Create currency
- `PUT /api/admin/currencies` - Update currency
- `DELETE /api/admin/currencies` - Delete currency
- `PUT /api/admin/currencies/default` - Set default currency

### Category Management

- `POST /api/admin/categories` - Create category
- `PUT /api/admin/categories/{id}` - Update category
- `DELETE /api/admin/categories/{id}` - Delete category

### Product Management

- `GET /api/admin/products` - List all products
- `POST /api/admin/products` - Create product
- `PUT /api/admin/products/{productId}` - Update product
- `DELETE /api/admin/products/{productId}` - Delete product

### Product Variant Management

- `POST /api/admin/products/{productId}/variants` - Add product variant
- `PUT /api/admin/products/{productId}/variants/{variantId}` - Update variant
- `DELETE /api/admin/products/{productId}/variants/{variantId}` - Delete variant

### Shipping Management

- `POST /api/admin/shipping/methods` - Create shipping method
- `POST /api/admin/shipping/zones` - Create shipping zone
- `POST /api/admin/shipping/rates` - Create shipping rate
- `POST /api/admin/shipping/rates/weight` - Create weight-based rate
- `POST /api/admin/shipping/rates/value` - Create value-based rate

### Discount Management

- `POST /api/admin/discounts` - Create discount
- `GET /api/admin/discounts/{discountId}` - Get discount
- `PUT /api/admin/discounts/{discountId}` - Update discount
- `DELETE /api/admin/discounts/{discountId}` - Delete discount
- `GET /api/admin/discounts` - List all discounts
- `GET /api/admin/discounts/active` - List active discounts
- `POST /api/admin/discounts/apply/{orderId}` - Apply discount to order
- `DELETE /api/admin/discounts/remove/{orderId}` - Remove discount from order

### Payment Management

- `POST /api/admin/payments/{paymentId}/capture` - Capture payment
- `POST /api/admin/payments/{paymentId}/cancel` - Cancel payment
- `POST /api/admin/payments/{paymentId}/refund` - Refund payment
- `POST /api/admin/payments/{paymentId}/force-approve` - Force approve MobilePay payment

### Payment Provider Management

- `GET /api/admin/payment-providers` - Get all payment providers
- `GET /api/admin/payment-providers/enabled` - Get enabled providers
- `PUT /api/admin/payment-providers/{providerType}/enable` - Enable/disable provider
- `PUT /api/admin/payment-providers/{providerType}/configuration` - Update configuration
- `POST /api/admin/payment-providers/{providerType}/webhook` - Register webhook
- `DELETE /api/admin/payment-providers/{providerType}/webhook` - Delete webhook
- `GET /api/admin/payment-providers/{providerType}/webhook` - Get webhook info

### Email Testing

- `POST /api/admin/test/email` - Send test email

## Webhook Endpoints

Server-to-server communication endpoints (no authentication required):

- `POST /api/webhooks/stripe` - Stripe webhook
- `POST /api/webhooks/mobilepay` - MobilePay webhook

## Authentication

Most endpoints require authentication via JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Permission Levels

1. **Public** - No authentication required
2. **Authenticated** - Valid JWT token required
3. **Admin** - JWT token with admin role required
4. **Webhook** - Server-to-server, signature verification

## Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists or conflict
- `500 Internal Server Error` - Server error
