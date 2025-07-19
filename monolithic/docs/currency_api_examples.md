# Currency API Examples

This document provides example request bodies for the currency system API endpoints.

## Public Currency Endpoints

### List Enabled Currencies

```plaintext
GET /api/currencies
```

Retrieve all currencies that are enabled in the system.

Example response:

```json
[
  {
    "code": "USD",
    "name": "US Dollar",
    "symbol": "$",
    "exchange_rate": 1.0,
    "is_enabled": true,
    "is_default": true,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  {
    "code": "EUR",
    "name": "Euro",
    "symbol": "€",
    "exchange_rate": 0.85,
    "is_enabled": true,
    "is_default": false,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  {
    "code": "GBP",
    "name": "British Pound",
    "symbol": "£",
    "exchange_rate": 0.75,
    "is_enabled": true,
    "is_default": false,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
]
```

**Status Codes:**

- `200 OK`: Currencies retrieved successfully
- `500 Internal Server Error`: Failed to retrieve currencies

### Get Default Currency

```plaintext
GET /api/currencies/default
```

Retrieve the default currency used in the system.

Example response:

```json
{
  "code": "USD",
  "name": "US Dollar",
  "symbol": "$",
  "exchange_rate": 1.0,
  "is_enabled": true,
  "is_default": true,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

**Status Codes:**

- `200 OK`: Default currency retrieved successfully
- `404 Not Found`: Default currency not found
- `500 Internal Server Error`: Failed to retrieve default currency

### Convert Currency Amount

```plaintext
POST /api/currencies/convert
```

Convert an amount from one currency to another.

**Request Body:**

```json
{
  "amount": 100.0,
  "from_currency": "USD",
  "to_currency": "EUR"
}
```

Example response:

```json
{
  "from": {
    "currency": "USD",
    "amount": 100.0,
    "cents": 10000
  },
  "to": {
    "currency": "EUR",
    "amount": 85.0,
    "cents": 8500
  }
}
```

**Status Codes:**

- `200 OK`: Amount converted successfully
- `400 Bad Request`: Invalid request body or currency not found
- `500 Internal Server Error`: Failed to convert amount

## Admin Currency Endpoints

### List All Currencies

```plaintext
GET /api/admin/currencies/all
```

List all currencies in the system, including disabled ones (admin only).

Example response:

```json
[
  {
    "code": "USD",
    "name": "US Dollar",
    "symbol": "$",
    "exchange_rate": 1.0,
    "is_enabled": true,
    "is_default": true,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  {
    "code": "EUR",
    "name": "Euro",
    "symbol": "€",
    "exchange_rate": 0.85,
    "is_enabled": true,
    "is_default": false,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  {
    "code": "GBP",
    "name": "British Pound",
    "symbol": "£",
    "exchange_rate": 0.75,
    "is_enabled": true,
    "is_default": false,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  {
    "code": "JPY",
    "name": "Japanese Yen",
    "symbol": "¥",
    "exchange_rate": 110.0,
    "is_enabled": false,
    "is_default": false,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
]
```

**Status Codes:**

- `200 OK`: Currencies retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `500 Internal Server Error`: Failed to retrieve currencies

### Create Currency

```plaintext
POST /api/admin/currencies
```

Create a new currency (admin only).

**Request Body:**

```json
{
  "code": "CAD",
  "name": "Canadian Dollar",
  "symbol": "C$",
  "exchange_rate": 1.25,
  "is_enabled": true,
  "is_default": false
}
```

Example response:

```json
{
  "code": "CAD",
  "name": "Canadian Dollar",
  "symbol": "C$",
  "exchange_rate": 1.25,
  "is_enabled": true,
  "is_default": false,
  "created_at": "2025-05-08T15:30:45Z",
  "updated_at": "2025-05-08T15:30:45Z"
}
```

**Status Codes:**

- `201 Created`: Currency created successfully
- `400 Bad Request`: Invalid request body or currency code already exists
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `500 Internal Server Error`: Failed to create currency

### Update Currency

```plaintext
PUT /api/admin/currencies?code={code}
```

Update an existing currency (admin only).

**Request Body:**

```json
{
  "name": "Canadian Dollar",
  "symbol": "CA$",
  "exchange_rate": 1.27,
  "is_enabled": true,
  "is_default": false
}
```

Example response:

```json
{
  "code": "CAD",
  "name": "Canadian Dollar",
  "symbol": "CA$",
  "exchange_rate": 1.27,
  "is_enabled": true,
  "is_default": false,
  "created_at": "2025-05-08T15:30:45Z",
  "updated_at": "2025-05-08T15:45:22Z"
}
```

**Status Codes:**

- `200 OK`: Currency updated successfully
- `400 Bad Request`: Invalid request body or currency not found
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `500 Internal Server Error`: Failed to update currency

### Delete Currency

```plaintext
DELETE /api/admin/currencies?code={code}
```

Delete a currency (admin only). Cannot delete the default currency.

Example response:

```json
{
  "status": "success",
  "message": "Currency deleted successfully"
}
```

**Status Codes:**

- `200 OK`: Currency deleted successfully
- `400 Bad Request`: Cannot delete default currency
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Currency not found
- `500 Internal Server Error`: Failed to delete currency

### Set Default Currency

```plaintext
PUT /api/admin/currencies/default?code={code}
```

Set a currency as the default currency (admin only).

Example response:

```json
{
  "code": "EUR",
  "name": "Euro",
  "symbol": "€",
  "exchange_rate": 0.85,
  "is_enabled": true,
  "is_default": true,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2025-05-08T16:15:30Z"
}
```

**Status Codes:**

- `200 OK`: Default currency set successfully
- `400 Bad Request`: Invalid currency code
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Currency not found
- `500 Internal Server Error`: Failed to set default currency

## Multi-Currency Support

The system supports selling products in multiple currencies. When creating or updating products, you can specify prices in different currencies.

### Product with Multi-Currency Prices

When retrieving products, you can view prices in the store's default currency or a specific currency by using the appropriate endpoints.

Example of a product with multiple currency prices:

```json
{
  "id": 1,
  "product_number": "PROD-000001",
  "name": "Smartphone",
  "description": "Latest smartphone model",
  "price": 999.99,
  "prices": [
    {
      "currency_code": "USD",
      "price": 999.99,
      "compare_price": 1099.99
    },
    {
      "currency_code": "EUR",
      "price": 849.99,
      "compare_price": 934.99
    },
    {
      "currency_code": "GBP",
      "price": 749.99,
      "compare_price": 824.99
    }
  ],
  "stock": 50,
  "weight": 0.35,
  "category_id": 1,
  "seller_id": 2,
  "images": ["smartphone.jpg"],
  "has_variants": false,
  "created_at": "2023-04-15T10:00:00Z",
  "updated_at": "2023-04-15T10:00:00Z"
}
```

## Example Workflows

### Setting Up Multi-Currency Support

1. Admin creates different currencies with appropriate exchange rates
2. Admin sets one currency as the default currency
3. Sellers can specify prices in different currencies for their products
4. Customers can view prices in their preferred currency

### Currency Conversion Process

1. Customer selects a non-default currency
2. System converts all product prices to the selected currency using the exchange rates
3. All prices throughout the store are displayed in the selected currency
4. Orders are still processed in the system's default currency
