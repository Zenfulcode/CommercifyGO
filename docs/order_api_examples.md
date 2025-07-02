# Order API Examples

This document provides example request bodies for the order system API endpoints.

### Get Order

```plaintext
GET /api/orders/{id}
```

Retrieve a specific order for the authenticated user.

Example response:

```json
{
  "success": true,
  "message": "Order retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "user_id": "550e8400-e29b-41d4-a716-446655440004",
    "status": "paid",
    "total_amount": 2514.97,
    "currency": "USD",
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440005",
        "order_id": "550e8400-e29b-41d4-a716-446655440003",
        "product_id": "550e8400-e29b-41d4-a716-446655440006",
        "name": "Premium Product",
        "sku": "PROD-002",
        "quantity": 1,
        "unit_price": 2499.99,
        "total_price": 2499.99
      }
    ],
    "shipping_address": {
      "first_name": "Sarah",
      "last_name": "Johnson",
      "address_line1": "456 Oak Avenue",
      "address_line2": "Suite 100",
      "city": "Seattle",
      "state": "WA",
      "postal_code": "98101",
      "country": "US",
      "phone_number": "+1987654321"
    },
    "billing_address": {
      "first_name": "Sarah",
      "last_name": "Johnson",
      "address_line1": "456 Oak Avenue",
      "address_line2": "Suite 100",
      "city": "Seattle",
      "state": "WA",
      "postal_code": "98101",
      "country": "US",
      "phone_number": "+1987654321"
    },
    "payment_method": "wallet",
    "payment_status": "captured",
    "shipping_method": "express",
    "shipping_cost": 14.99,
    "tax_amount": 0,
    "discount_amount": 0,
    "created_at": "2024-03-20T11:00:00Z",
    "updated_at": "2024-03-20T11:05:00Z"
  }
}
```

**Status Codes:**

- `200 OK`: Order retrieved successfully
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized for this order
- `404 Not Found`: Order not found
- `500 Internal Server Error`: Failed to retrieve order

### List User Orders

```plaintext
GET /api/orders
```

List all orders for the authenticated user.

**Query Parameters:**

- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)

Example response:

```json
{
  "success": true,
  "message": "Orders retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "user_id": "550e8400-e29b-41d4-a716-446655440004",
      "status": "paid",
      "total_amount": 2514.97,
      "currency": "USD",
      "payment_method": "wallet",
      "payment_status": "captured",
      "shipping_method": "express",
      "shipping_cost": 14.99,
      "tax_amount": 0,
      "discount_amount": 0,
      "created_at": "2024-03-20T11:00:00Z",
      "updated_at": "2024-03-20T11:05:00Z"
    }
  ],
  "pagination": {
    "total": 1,
    "offset": 0,
    "limit": 10
  }
}
```

**Status Codes:**

- `200 OK`: Orders retrieved successfully
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to retrieve orders

## Admin Order Endpoints

### List All Orders

```plaintext
GET /api/admin/orders
```

List all orders in the system (admin only).

**Query Parameters:**

- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)
- `status` (optional): Filter by order status

Example response:

```json
{
  "success": true,
  "message": "Orders retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "user_id": "550e8400-e29b-41d4-a716-446655440004",
      "status": "paid",
      "total_amount": 2514.97,
      "currency": "USD",
      "payment_method": "wallet",
      "payment_status": "captured",
      "shipping_method": "express",
      "shipping_cost": 14.99,
      "tax_amount": 0,
      "discount_amount": 0,
      "created_at": "2024-03-20T11:00:00Z",
      "updated_at": "2024-03-20T11:05:00Z"
    }
  ],
  "pagination": {
    "total": 1,
    "offset": 0,
    "limit": 10
  }
}
```

**Status Codes:**

- `200 OK`: Orders retrieved successfully
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized (not an admin)
- `500 Internal Server Error`: Failed to retrieve orders

### Update Order Status

```plaintext
PUT /api/admin/orders/{id}/status
```

Update an order's status (admin only).

**Request Body:**

```json
{
  "status": "shipped"
}
```

Example response:

```json
{
  "success": true,
  "message": "Order status updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "user_id": "550e8400-e29b-41d4-a716-446655440004",
    "status": "shipped",
    "total_amount": 2514.97,
    "currency": "USD",
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440005",
        "order_id": "550e8400-e29b-41d4-a716-446655440003",
        "product_id": "550e8400-e29b-41d4-a716-446655440006",
        "name": "Premium Product",
        "sku": "PROD-002",
        "quantity": 1,
        "unit_price": 2499.99,
        "total_price": 2499.99
      }
    ],
    "shipping_address": {
      "first_name": "Sarah",
      "last_name": "Johnson",
      "address_line1": "456 Oak Avenue",
      "address_line2": "Suite 100",
      "city": "Seattle",
      "state": "WA",
      "postal_code": "98101",
      "country": "US",
      "phone_number": "+1987654321"
    },
    "billing_address": {
      "first_name": "Sarah",
      "last_name": "Johnson",
      "address_line1": "456 Oak Avenue",
      "address_line2": "Suite 100",
      "city": "Seattle",
      "state": "WA",
      "postal_code": "98101",
      "country": "US",
      "phone_number": "+1987654321"
    },
    "payment_method": "wallet",
    "payment_status": "captured",
    "shipping_method": "express",
    "shipping_cost": 14.99,
    "tax_amount": 0,
    "discount_amount": 0,
    "created_at": "2024-03-20T11:00:00Z",
    "updated_at": "2024-03-20T14:30:00Z"
  }
}
```

**Status Codes:**

- `200 OK`: Order status updated successfully
- `400 Bad Request`: Invalid order status
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized (not an admin)
- `404 Not Found`: Order not found
- `500 Internal Server Error`: Failed to update order status
