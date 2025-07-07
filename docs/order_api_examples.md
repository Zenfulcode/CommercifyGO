# Order API Examples

This document provides example request bodies for the order system API endpoints.

# Order API Examples

This document provides example request bodies for the order system API endpoints.

## Public Order Endpoints

### Get Order

```plaintext
GET /api/orders/{orderId}
```

Retrieve a specific order. This endpoint supports optional authentication - users can access their own orders, while non-authenticated users can access orders via checkout session cookie.

**Path Parameters:**

- `orderId` (required): Order ID

**Authorization:**

- Authenticated users can access their own orders
- Admin users can access any order
- Non-authenticated users can access orders if they have a valid checkout session cookie

**Response Body:**

```json
{
  "success": true,
  "data": {
    "id": 123,
    "order_number": "ORD-20250707-123",
    "user_id": 1,
    "checkout_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "items": [
      {
        "id": 1,
        "order_id": 123,
        "product_id": 42,
        "variant_id": 7,
        "sku": "PROD-001-L",
        "product_name": "Organic Cotton T-Shirt",
        "variant_name": "Size L, Color Blue",
        "quantity": 2,
        "unit_price": 24.99,
        "total_price": 49.98,
        "image_url": "https://example.com/image.jpg",
        "created_at": "2025-07-07T10:30:45Z",
        "updated_at": "2025-07-07T10:30:45Z"
      }
    ],
    "status": "paid",
    "payment_status": "captured",
    "total_amount": 49.98,
    "shipping_cost": 9.99,
    "discount_amount": 5.0,
    "final_amount": 54.97,
    "currency": "USD",
    "shipping_address": {
      "address_line1": "123 Main St",
      "address_line2": "Apt 4B",
      "city": "New York",
      "state": "NY",
      "postal_code": "10001",
      "country": "US"
    },
    "billing_address": {
      "address_line1": "123 Main St",
      "address_line2": "Apt 4B",
      "city": "New York",
      "state": "NY",
      "postal_code": "10001",
      "country": "US"
    },
    "shipping_details": {
      "shipping_rate_id": 1,
      "shipping_method_id": 1,
      "name": "Standard Shipping",
      "description": "5-7 business days",
      "estimated_delivery_days": 7,
      "cost": 9.99,
      "free_shipping": false
    },
    "discount_details": {
      "id": 1,
      "code": "SUMMER2025",
      "type": "basket",
      "method": "percentage",
      "value": 10.0,
      "amount": 5.0
    },
    "payment_transactions": [
      {
        "id": 1,
        "transaction_id": "txn_123456789",
        "external_id": "pi_1234567890",
        "type": "authorize",
        "status": "successful",
        "amount": 54.97,
        "currency": "USD",
        "provider": "stripe",
        "created_at": "2025-07-07T10:30:45Z",
        "updated_at": "2025-07-07T10:30:45Z"
      }
    ],
    "customer": {
      "email": "customer@example.com",
      "phone": "+1234567890",
      "full_name": "John Smith"
    },
    "action_required": false,
    "action_url": null,
    "created_at": "2025-07-07T10:30:45Z",
    "updated_at": "2025-07-07T10:30:45Z"
  }
}
```

**Status Codes:**

- `200 OK`: Order retrieved successfully
- `401 Unauthorized`: Not authenticated and no valid checkout session
- `403 Forbidden`: Not authorized to access this order
- `404 Not Found`: Order not found

## Authenticated User Endpoints

### List User Orders

```plaintext
GET /api/orders
```

List orders for the authenticated user.

**Query Parameters:**

- `page` (number, optional): Page number (default: 1)
- `pageSize` (number, optional): Items per page (default: 10)

**Response Body:**

```json
{
  "success": true,
  "data": [
    {
      "id": 123,
      "order_number": "ORD-20250707-123",
      "checkout_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
      "user_id": 1,
      "customer": {
        "email": "customer@example.com",
        "phone": "+1234567890",
        "full_name": "John Smith"
      },
      "status": "paid",
      "payment_status": "captured",
      "total_amount": 49.98,
      "shipping_cost": 9.99,
      "discount_amount": 5.0,
      "final_amount": 54.97,
      "order_lines_amount": 2,
      "currency": "USD",
      "created_at": "2025-07-07T10:30:45Z",
      "updated_at": "2025-07-07T10:30:45Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 1
  }
}
```

**Status Codes:**

- `200 OK`: Orders retrieved successfully
- `401 Unauthorized`: Not authenticated

## Admin Order Endpoints

All admin order endpoints require authentication and admin role.

### List All Orders

```plaintext
GET /api/admin/orders
```

List all orders in the system (admin only).

**Query Parameters:**

- `page` (number, optional): Page number (default: 1)
- `pageSize` (number, optional): Items per page (default: 10)
- `status` (string, optional): Filter by order status

**Status Codes:**

- `200 OK`: Orders retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

### Update Order Status

```plaintext
PUT /api/admin/orders/{orderId}/status
```

Update the status of an order (admin only).

**Path Parameters:**

- `orderId` (required): Order ID

**Request Body:**

```json
{
  "status": "shipped"
}
```

**Status Codes:**

- `200 OK`: Order status updated successfully
- `400 Bad Request`: Invalid request body or status
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Order not found
  "transaction_id": "TXN-AUTH-2025-001",
  "external_id": "pi_1234567890",
  "type": "authorize",
  "status": "successful",
  "amount": 2514.97,
  "currency": "USD",
  "provider": "stripe",
  "created_at": "2024-03-20T11:00:00Z",
  "updated_at": "2024-03-20T11:00:00Z"
  },
  {
  "id": 2,
  "transaction_id": "TXN-CAPTURE-2025-001",
  "external_id": "ch_1234567890",
  "type": "capture",
  "status": "successful",
  "amount": 2514.97,
  "currency": "USD",
  "provider": "stripe",
  "created_at": "2024-03-20T11:05:00Z",
  "updated_at": "2024-03-20T11:05:00Z"
  }
  ],
  "shipping_address": {...},
  "billing_address": {...},
  "payment_method": "credit_card",
  "shipping_method": "express",
  "shipping_cost": 14.99,
  "tax_amount": 0,
  "discount_amount": 0,
  "created_at": "2024-03-20T11:00:00Z",
  "updated_at": "2024-03-20T11:05:00Z"
  }
  }

````

**Status Codes:**

- `200 OK`: Order retrieved successfully
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized for this order
- `404 Not Found`: Order not found
- `500 Internal Server Error`: Failed to retrieve order

### List User Orders

```plaintext
GET /api/orders
````

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
