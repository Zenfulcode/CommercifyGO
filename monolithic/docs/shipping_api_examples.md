# Shipping API Examples

This document provides example request bodies for the shipping system API endpoints.

## Public Shipping Endpoints

### Calculate Shipping Options

```plaintext
POST /api/shipping/options
```

Calculate available shipping options for an address and order details.

**Request Body:**

```json
{
  "address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "order_value": 150.0,
  "order_weight": 2.5
}
```

**Response Body:**

```json
{
  "success": true,
  "data": {
    "options": [
      {
        "shipping_rate_id": 1,
        "shipping_method_id": 1,
        "name": "Standard Shipping",
        "description": "Delivery in 3-5 business days",
        "estimated_delivery_days": 5,
        "cost": 7.99,
        "free_shipping": false
      },
      {
        "shipping_rate_id": 2,
        "shipping_method_id": 2,
        "name": "Express Shipping",
        "description": "Delivery in 1-2 business days",
        "estimated_delivery_days": 2,
        "cost": 14.99,
        "free_shipping": false
      },
      {
        "shipping_rate_id": 3,
        "shipping_method_id": 3,
        "name": "Free Ground Shipping",
        "description": "Free shipping for orders over $100",
        "estimated_delivery_days": 7,
        "cost": 0.0,
        "free_shipping": true
      }
    ]
  }
}
```

**Status Codes:**

- `200 OK`: Shipping options calculated successfully
- `400 Bad Request`: Invalid request body or address

## Admin Shipping Endpoints

All admin shipping endpoints require authentication and admin role.

### Create Shipping Method

```plaintext
POST /api/admin/shipping/methods
```

Create a new shipping method (admin only).

**Request Body:**

```json
{
  "name": "Premium Overnight",
  "description": "Next day delivery guaranteed",
  "estimated_delivery_days": 1
}
```

**Status Codes:**

- `201 Created`: Shipping method created successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

### Create Shipping Zone

```plaintext
POST /api/admin/shipping/zones
```

Create a new shipping zone (admin only).

**Request Body:**

```json
{
  "name": "US West Coast",
  "description": "West coast shipping zone",
  "countries": ["US"],
  "states": ["CA", "OR", "WA"],
  "zip_codes": ["9****", "8****"]
}
```

**Status Codes:**

- `201 Created`: Shipping zone created successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

### Create Shipping Rate

```plaintext
POST /api/admin/shipping/rates
```

Create a new shipping rate (admin only).

**Request Body:**

```json
{
  "shipping_method_id": 1,
  "shipping_zone_id": 1,
  "base_rate": 9.99,
  "min_order_value": 0.0,
  "free_shipping_threshold": 100.0,
  "active": true
}
```

**Status Codes:**

- `201 Created`: Shipping rate created successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

### Create Weight-Based Rate

```plaintext
POST /api/admin/shipping/rates/weight
```

Create a weight-based shipping rate (admin only).

**Request Body:**

```json
{
  "shipping_rate_id": 1,
  "min_weight": 0.0,
  "max_weight": 5.0,
  "rate": 5.99
}
```

**Status Codes:**

- `201 Created`: Weight-based rate created successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

### Create Value-Based Rate

```plaintext
POST /api/admin/shipping/rates/value
```

Create a value-based shipping rate (admin only).

**Request Body:**

```json
{
  "shipping_rate_id": 1,
  "min_order_value": 0.0,
  "max_order_value": 50.0,
  "rate": 9.99
}
```

**Status Codes:**

- `201 Created`: Value-based rate created successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

### Update Shipping Method

`PUT /api/admin/shipping/methods/{id}`

Update an existing shipping method.

```json
{
  "name": "Premium Overnight",
  "description": "Next day delivery by 10:00 AM guaranteed",
  "estimated_delivery_days": 1,
  "active": true
}
```

### Create Shipping Zone

`POST /api/admin/shipping/zones`

Create a new shipping zone.

```json
{
  "name": "US West Coast",
  "description": "CA, OR, WA, NV, AZ",
  "countries": ["US"],
  "states": ["CA", "OR", "WA", "NV", "AZ"],
  "zip_codes": []
}
```

### Update Shipping Zone

`PUT /api/admin/shipping/zones/{id}`

Update an existing shipping zone.

```json
{
  "name": "US West Coast",
  "description": "CA, OR, WA, NV, AZ, HI",
  "countries": ["US"],
  "states": ["CA", "OR", "WA", "NV", "AZ", "HI"],
  "zip_codes": [],
  "active": true
}
```

### Create Shipping Rate

`POST /api/admin/shipping/rates`

Create a new shipping rate connecting a method and zone.

```json
{
  "shipping_method_id": 1,
  "shipping_zone_id": 1,
  "base_rate": 8.99,
  "min_order_value": 0.0,
  "free_shipping_threshold": 100.0,
  "active": true
}
```

### Update Shipping Rate

`PUT /api/admin/shipping/rates/{id}`

Update an existing shipping rate.

```json
{
  "base_rate": 7.99,
  "min_order_value": 0.0,
  "free_shipping_threshold": 75.0,
  "active": true
}
```

### Create Weight-Based Rate

`POST /api/admin/shipping/rates/weight`

Add a weight-based rate to an existing shipping rate.

```json
{
  "shipping_rate_id": 1,
  "min_weight": 5.0,
  "max_weight": 10.0,
  "rate": 3.99
}
```

### Create Value-Based Rate

`POST /api/admin/shipping/rates/value`

Add a value-based rate to an existing shipping rate.

```json
{
  "shipping_rate_id": 1,
  "min_order_value": 50.0,
  "max_order_value": 100.0,
  "rate": -1.5
}
```

## Example Workflow

### Shipping Configuration Flow (Admin)

1. Admin creates shipping methods (Standard, Express, etc.)
2. Admin creates shipping zones (US Domestic, International, etc.)
3. Admin creates shipping rates connecting methods to zones
4. Admin adds weight-based or value-based rules to rates as needed

### Customer Shipping Selection Flow

1. When a customer enters their shipping address and has items in cart, the system calls the shipping options endpoint
2. Available shipping options are displayed to the customer based on their location and order details
3. Customer selects a shipping option during checkout
4. The selected shipping method is included in the order
5. Shipping cost is calculated and added to the order total
