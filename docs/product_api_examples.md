# Product API Examples

This document provides example request bodies and responses for the product system API endpoints.

## Important Notes

- **All products must have at least one variant**: Every product in the system is required to have at least one product variant. If no variants are specified when creating a product, a default variant will be automatically created using the product's basic information.
- **SKUs are variant-specific**: All SKU-based operations (like adding items to checkout) must use variant SKUs, not product numbers.
- **Product numbers are deprecated**: While products still have product numbers for backward compatibility, all SKU lookups are now performed against variant SKUs.

## Public Product Endpoints

### List Products

`GET /api/products`

List all products with pagination.

**Query Parameters:**

- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10)

Example response:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "created_at": "2023-04-15T10:00:00Z",
      "updated_at": "2023-04-15T10:00:00Z",
      "name": "Smartphone",
      "description": "Latest smartphone model",
      "sku": "PROD-000001",
      "price": 999.99,
      "stock_quantity": 50,
      "weight": 0.35,
      "category_id": 1,
      "seller_id": 2,
      "images": ["smartphone.jpg"],
      "has_variants": true,
      "variants": [
        {
          "id": 1,
          "product_id": 1,
          "sku": "PROD-000001",
          "price": 999.99,
          "stock_quantity": 50,
          "attributes": [],
          "images": [],
          "is_default": true
        }
      ]
    },
    {
      "id": 2,
      "created_at": "2023-04-16T11:00:00Z",
      "updated_at": "2023-04-16T11:00:00Z",
      "name": "Laptop",
      "description": "Powerful laptop for professionals",
      "sku": "PROD-000002",
      "price": 1499.99,
      "stock_quantity": 25,
      "weight": 2.1,
      "category_id": 1,
      "seller_id": 2,
      "images": ["laptop.jpg"],
      "has_variants": true,
      "variants": [
        {
          "id": 1,
          "created_at": "2023-04-15T10:00:00Z",
          "updated_at": "2023-04-15T10:00:00Z",
          "product_id": 2,
          "sku": "LAPT-8GB-256",
          "price": 1499.99,
          "compare_price": 1599.99,
          "stock_quantity": 10,
          "attributes": {
            "ram": "8GB",
            "storage": "256GB",
            "color": "Silver"
          },
          "images": ["laptop_silver.jpg"],
          "is_default": true
        }
      ]
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 2
  }
}
```

**Status Codes:**

- `200 OK`: Products retrieved successfully
- `500 Internal Server Error`: Server error occurred

### Get Product

`GET /api/products/{id}`

Get details of a specific product.

**Query Parameters:**

- `currency` (optional): Currency code to display prices in (e.g., "EUR", "GBP")

Example response:

```json
{
  "success": true,
  "data": {
    "id": 2,
    "created_at": "2023-04-16T11:00:00Z",
    "updated_at": "2023-04-16T11:00:00Z",
    "name": "Laptop",
    "description": "Powerful laptop for professionals",
    "sku": "PROD-000002",
    "price": 1499.99,
    "stock_quantity": 25,
    "weight": 2.1,
    "category_id": 1,
    "seller_id": 2,
    "images": ["laptop.jpg"],
    "has_variants": true,
    "variants": [
      {
        "id": 1,
        "created_at": "2023-04-15T10:00:00Z",
        "updated_at": "2023-04-15T10:00:00Z",
        "product_id": 2,
        "sku": "LAPT-8GB-256",
        "price": 1499.99,
        "compare_price": 1599.99,
        "stock_quantity": 10,
        "attributes": {
          "ram": "8GB",
          "storage": "256GB",
          "color": "Silver"
        },
        "images": ["laptop_silver.jpg"],
        "is_default": true
      }
    ]
  }
}
```

**Status Codes:**

- `200 OK`: Product retrieved successfully
- `400 Bad Request`: Invalid product ID
- `404 Not Found`: Product not found
- `500 Internal Server Error`: Server error occurred

### Search Products

`POST /api/products/search`

Search products based on various criteria.

Request body:

```json
{
  "query": "laptop",
  "category_id": 1,
  "min_price": 1000,
  "max_price": 2000,
  "page": 1,
  "page_size": 10
}
```

Example response:

```json
{
  "success": true,
  "data": [
    {
      "id": 2,
      "created_at": "2023-04-16T11:00:00Z",
      "updated_at": "2023-04-16T11:00:00Z",
      "name": "Laptop",
      "description": "Powerful laptop for professionals",
      "sku": "PROD-000002",
      "price": 1499.99,
      "stock_quantity": 25,
      "weight": 2.1,
      "category_id": 1,
      "seller_id": 2,
      "images": ["laptop.jpg"],
      "has_variants": true,
      "variants": [
        {
          "id": 1,
          "created_at": "2023-04-15T10:00:00Z",
          "updated_at": "2023-04-15T10:00:00Z",
          "product_id": 2,
          "sku": "LAPT-8GB-256",
          "price": 1499.99,
          "compare_price": 1599.99,
          "stock_quantity": 10,
          "attributes": {
            "ram": "8GB",
            "storage": "256GB",
            "color": "Silver"
          },
          "images": ["laptop_silver.jpg"],
          "is_default": true
        }
      ]
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

- `200 OK`: Search results retrieved successfully
- `400 Bad Request`: Invalid request body
- `500 Internal Server Error`: Server error occurred

### List Categories

`GET /api/categories`

List all product categories.

Example response:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Electronics",
      "description": "Electronic devices and gadgets",
      "parent_id": null,
      "created_at": "2023-04-10T09:00:00Z",
      "updated_at": "2023-04-10T09:00:00Z"
    },
    {
      "id": 2,
      "name": "Smartphones",
      "description": "Mobile phones and smartphones",
      "parent_id": 1,
      "created_at": "2023-04-10T09:05:00Z",
      "updated_at": "2023-04-10T09:05:00Z"
    }
  ]
}
```

**Status Codes:**

- `200 OK`: Categories retrieved successfully
- `500 Internal Server Error`: Server error occurred

## Seller Product Endpoints

### Create Product

`POST /api/products`

Create a new product (seller only).

Request body:

```json
{
  "name": "New Product",
  "description": "Product description",
  "price": 199.99,
  "stock_quantity": 100,
  "weight": 1.5,
  "category_id": 1,
  "images": ["product.jpg"],
  "variants": []
}
```

**Note:** All products must have at least one variant. If no variants are provided in the request, a default variant will be automatically created using the product's basic information (price, stock) and the product number as the SKU.

Example response:

```json
{
  "success": true,
  "data": {
    "id": 4,
    "created_at": "2023-04-25T14:00:00Z",
    "updated_at": "2023-04-25T14:00:00Z",
    "name": "New Product",
    "description": "Product description",
    "sku": "PROD-000004",
    "price": 199.99,
    "stock_quantity": 100,
    "weight": 1.5,
    "category_id": 1,
    "seller_id": 2,
    "images": ["product.jpg"],
    "has_variants": true,
    "variants": [
      {
        "id": 1,
        "product_id": 4,
        "sku": "PROD-000004",
        "price": 199.99,
        "stock_quantity": 100,
        "attributes": [],
        "images": [],
        "is_default": true,
        "created_at": "2023-04-25T14:00:00Z",
        "updated_at": "2023-04-25T14:00:00Z"
      }
    ]
  }
}
```

**Status Codes:**

- `201 Created`: Product created successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `500 Internal Server Error`: Server error occurred

### Update Product

`PUT /api/products/{id}`

Update an existing product (seller only).

Request body:

```json
{
  "name": "Updated Product",
  "description": "Updated product description",
  "price": 249.99,
  "stock_quantity": 75,
  "weight": 1.6,
  "category_id": 1,
  "images": ["updated-product.jpg"]
}
```

Example response:

```json
{
  "success": true,
  "data": {
    "id": 4,
    "created_at": "2023-04-25T14:00:00Z",
    "updated_at": "2023-04-25T14:30:00Z",
    "name": "Updated Product",
    "description": "Updated product description",
    "sku": "PROD-000004",
    "price": 249.99,
    "stock_quantity": 75,
    "weight": 1.6,
    "category_id": 1,
    "seller_id": 2,
    "images": ["updated-product.jpg"],
    "has_variants": true
  }
}
```

**Status Codes:**

- `200 OK`: Product updated successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `500 Internal Server Error`: Server error occurred

### Delete Product

`DELETE /api/products/{id}`

Delete a product (seller only).

Example response:

```json
{
  "success": true,
  "message": "Product deleted successfully"
}
```

**Status Codes:**

- `200 OK`: Product deleted successfully
- `400 Bad Request`: Invalid product ID
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `500 Internal Server Error`: Server error occurred

### List Seller Products

`GET /api/products/seller`

List all products for the authenticated seller.

**Query Parameters:**

- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10)

Example response:

```json
{
  "success": true,
  "data": [
    {
      "id": 4,
      "created_at": "2023-04-25T14:00:00Z",
      "updated_at": "2023-04-25T14:30:00Z",
      "name": "Updated Product",
      "description": "Updated product description",
      "sku": "PROD-000004",
      "price": 249.99,
      "stock_quantity": 75,
      "weight": 1.6,
      "category_id": 1,
      "seller_id": 2,
      "images": ["updated-product.jpg"],
      "has_variants": true
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

- `200 OK`: Products retrieved successfully
- `401 Unauthorized`: Not authenticated
- `500 Internal Server Error`: Server error occurred

## Product Variant Endpoints

### Add Product Variant

`POST /api/products/{productId}/variants`

Add a variant to a product (seller only).

Request body:

```json
{
  "sku": "PROD-RED-M",
  "price": 29.99,
  "compare_price": 39.99,
  "stock_quantity": 10,
  "attributes": {
    "color": "Red",
    "size": "Medium"
  },
  "images": ["red-shirt.jpg"],
  "is_default": true
}
```

Example response:

```json
{
  "success": true,
  "data": {
    "id": 11,
    "created_at": "2023-04-28T15:00:00Z",
    "updated_at": "2023-04-28T15:00:00Z",
    "product_id": 3,
    "sku": "PROD-RED-M",
    "price": 29.99,
    "compare_price": 39.99,
    "stock_quantity": 10,
    "attributes": {
      "color": "Red",
      "size": "Medium"
    },
    "images": ["red-shirt.jpg"],
    "is_default": true
  }
}
```

**Status Codes:**

- `201 Created`: Variant created successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `500 Internal Server Error`: Server error occurred

### Update Product Variant

`PUT /api/products/{productId}/variants/{variantId}`

Update a product variant (seller only).

Request body:

```json
{
  "sku": "PROD-RED-M",
  "price": 24.99,
  "compare_price": 34.99,
  "stock_quantity": 15,
  "attributes": {
    "color": "Red",
    "size": "Medium"
  },
  "images": ["red-shirt-updated.jpg"],
  "is_default": true
}
```

Example response:

```json
{
  "success": true,
  "data": {
    "id": 11,
    "created_at": "2023-04-28T15:00:00Z",
    "updated_at": "2023-04-28T15:30:00Z",
    "product_id": 3,
    "sku": "PROD-RED-M",
    "price": 24.99,
    "compare_price": 34.99,
    "stock_quantity": 15,
    "attributes": {
      "color": "Red",
      "size": "Medium"
    },
    "images": ["red-shirt-updated.jpg"],
    "is_default": true
  }
}
```

**Status Codes:**

- `200 OK`: Variant updated successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `500 Internal Server Error`: Server error occurred

### Delete Product Variant

`DELETE /api/products/{productId}/variants/{variantId}`

Delete a product variant (seller only).

Example response:

```json
{
  "success": true,
  "message": "Variant deleted successfully"
}
```

**Status Codes:**

- `200 OK`: Variant deleted successfully
- `400 Bad Request`: Invalid product or variant ID
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `500 Internal Server Error`: Server error occurred

## Multi-Currency Product Management

## Multi-Currency Variant Pricing

### Set Variant Price in Specific Currency

```plaintext
POST /api/admin/variants/{variant_id}/prices
```

**Request Body:**

```json
{
  "currency_code": "DKK",
  "price": 250.0
}
```

**Response Body:**

```json
{
  "id": 1,
  "product_id": 1,
  "sku": "PROD-001-RED",
  "price": 25.0,
  "currency": "USD",
  "stock": 10,
  "attributes": [
    {
      "name": "Color",
      "value": "Red"
    }
  ],
  "images": [],
  "is_default": true,
  "created_at": "2025-06-20T10:00:00Z",
  "updated_at": "2025-06-20T10:30:00Z",
  "prices": {
    "USD": 25.0,
    "DKK": 250.0
  }
}
```

**Status Codes:**

- `200 OK`: Price set successfully
- `400 Bad Request`: Invalid request body or currency not enabled
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Variant or currency not found

### Set Multiple Variant Prices

```plaintext
PUT /api/admin/variants/{variant_id}/prices
```

**Request Body:**

```json
{
  "prices": {
    "USD": 25.0,
    "EUR": 21.25,
    "DKK": 250.0
  }
}
```

**Response Body:**

```json
{
  "id": 1,
  "product_id": 1,
  "sku": "PROD-001-RED",
  "price": 25.0,
  "currency": "USD",
  "stock": 10,
  "attributes": [
    {
      "name": "Color",
      "value": "Red"
    }
  ],
  "images": [],
  "is_default": true,
  "created_at": "2025-06-20T10:00:00Z",
  "updated_at": "2025-06-20T10:30:00Z",
  "prices": {
    "USD": 25.0,
    "EUR": 21.25,
    "DKK": 250.0
  }
}
```

**Status Codes:**

- `200 OK`: Prices set successfully
- `400 Bad Request`: Invalid request body or one or more currencies not enabled
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Variant not found

### Get Variant Prices

```plaintext
GET /api/variants/{variant_id}/prices
```

**Response Body:**

```json
{
  "prices": {
    "USD": 25.0,
    "EUR": 21.25,
    "DKK": 250.0
  }
}
```

**Status Codes:**

- `200 OK`: Prices retrieved successfully
- `404 Not Found`: Variant not found

### Remove Variant Price in Specific Currency

```plaintext
DELETE /api/admin/variants/{variant_id}/prices/{currency_code}
```

**Response Body:**

```json
{
  "id": 1,
  "product_id": 1,
  "sku": "PROD-001-RED",
  "price": 25.0,
  "currency": "USD",
  "stock": 10,
  "attributes": [
    {
      "name": "Color",
      "value": "Red"
    }
  ],
  "images": [],
  "is_default": true,
  "created_at": "2025-06-20T10:00:00Z",
  "updated_at": "2025-06-20T10:30:00Z",
  "prices": {
    "USD": 25.0,
    "EUR": 21.25
  }
}
```

**Status Codes:**

- `200 OK`: Price removed successfully
- `400 Bad Request`: Cannot remove default currency price
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Variant not found or price not set for this currency

## Benefits of Multi-Currency Pricing

### Precision and Accuracy

- **No conversion errors**: Set exact prices in each currency to avoid floating-point precision issues
- **Local pricing**: Set appropriate prices for each market without relying on exchange rate calculations
- **Price consistency**: Ensure consistent pricing across all customer touchpoints

### Checkout Integration

- **Automatic price selection**: When customers checkout in a specific currency, the system automatically uses the exact price set for that currency
- **Fallback conversion**: If no specific price is set for a currency, the system falls back to conversion from the default currency
- **Currency parameter**: API clients can specify the desired checkout currency as a parameter

### Example Use Case

A product originally priced at 25.00 USD might be set to exactly:

- 250.00 DKK (instead of 249.96 DKK from conversion)
- 21.25 EUR (instead of 21.24 EUR from conversion)
- 19.99 GBP (psychological pricing)

This eliminates the reported issue where a 250 DKK product was showing as 249.89 DKK due to conversion precision problems.
