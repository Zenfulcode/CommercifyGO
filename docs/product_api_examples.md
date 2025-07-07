# Product API Examples

This document provides example request bodies and responses for the product system API endpoints.

## Important Notes

- **All products must have at least one variant**: Every product in the system is required to have at least one product variant. If no variants are specified when creating a product, a default variant will be automatically created using the product's basic information.
- **SKUs are variant-specific**: All SKU-based operations (like adding items to checkout) must use variant SKUs, not product numbers.
- **Product numbers are deprecated**: While products still have product numbers for backward compatibility, all SKU lookups are now performed against variant SKUs.

## Public Product Endpoints

### Get Product

```plaintext
GET /api/products/{productId}
```

Get a product by ID.

**Path Parameters:**

- `productId` (required): Product ID

**Response Body:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Smartphone",
    "description": "Latest smartphone model",
    "currency": "USD",
    "price": 999.99,
    "sku": "PROD-001",
    "total_stock": 50,
    "category": "Electronics",
    "category_id": 1,
    "images": ["https://example.com/smartphone.jpg"],
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
          "color": "black"
        },
        "images": ["https://example.com/variant1.jpg"],
        "is_default": true,
        "weight": 0.35,
        "price": 999.99,
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

**Status Codes:**

- `200 OK`: Product retrieved successfully
- `404 Not Found`: Product not found

### Search Products

```plaintext
GET /api/products/search
```

Search for products with optional filters.

**Query Parameters:**

- `query` (string, optional): Search term
- `category_id` (number, optional): Filter by category ID
- `min_price` (number, optional): Minimum price filter
- `max_price` (number, optional): Maximum price filter
- `currency` (string, optional): Currency code (default: USD)
- `active_only` (boolean, optional): Show only active products (default: true)
- `page` (number, optional): Page number (default: 1)
- `page_size` (number, optional): Items per page (default: 10)

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

## Admin Product Endpoints

All admin product endpoints require authentication and admin role.

### List Products

```plaintext
GET /api/admin/products
```

List all products (admin only).

**Query Parameters:**

- `page` (number, optional): Page number (default: 1)
- `page_size` (number, optional): Items per page (default: 10)
- `query` (string, optional): Search term
- `category_id` (number, optional): Filter by category ID
- `min_price` (number, optional): Minimum price filter
- `max_price` (number, optional): Maximum price filter
- `currency` (string, optional): Currency code
- `active_only` (boolean, optional): Show only active products

**Status Codes:**

- `200 OK`: Products retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

### Create Product

```plaintext
POST /api/admin/products
```

Create a new product (admin only).

**Request Body:**

```json
{
  "name": "New Product",
  "description": "Product description",
  "currency": "USD",
  "category_id": 1,
  "images": ["https://example.com/product.jpg"],
  "active": true,
  "variants": [
    {
      "sku": "PROD-004-L",
      "stock": 100,
      "attributes": [
        { "name": "size", "value": "L" },
        { "name": "color", "value": "blue" }
      ],
      "images": ["https://example.com/variant1.jpg"],
      "is_default": true,
      "weight": 1.5,
      "price": 199.99
    }
  ]
}
```

**Note:** All products must have at least one variant. If no variants are provided in the request, a default variant will be automatically created.

**Response Body:**

```json
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "id": 4,
    "name": "New Product",
    "description": "Product description",
    "currency": "USD",
    "price": 199.99,
    "sku": "PROD-004-L",
    "total_stock": 100,
    "category": "Electronics",
    "category_id": 1,
    "images": ["https://example.com/product.jpg"],
    "has_variants": true,
    "active": true,
    "variants": [
      {
        "id": 1,
        "product_id": 4,
        "variant_name": "Size L",
        "sku": "PROD-004-L",
        "stock": 100,
        "attributes": {
          "size": "L",
          "color": "blue"
        },
        "images": ["https://example.com/variant1.jpg"],
        "is_default": true,
        "weight": 1.5,
        "price": 199.99,
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

**Status Codes:**

- `201 Created`: Product created successfully
- `400 Bad Request`: Invalid request body or validation error
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

### Update Product

```plaintext
PUT /api/admin/products/{productId}
```

Update an existing product (admin only).

**Path Parameters:**

- `productId` (required): Product ID

**Request Body:**

```json
{
  "name": "Updated Product Name",
  "description": "Updated description",
  "currency": "USD",
  "category_id": 2,
  "images": ["https://example.com/updated-product.jpg"],
  "active": true
}
```

**Status Codes:**

- `200 OK`: Product updated successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Product not found

### Delete Product

```plaintext
DELETE /api/admin/products/{productId}
```

Delete a product (admin only).

**Path Parameters:**

- `productId` (required): Product ID

**Status Codes:**

- `200 OK`: Product deleted successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Product not found
- `409 Conflict`: Cannot delete product with existing orders

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

All variant endpoints require authentication and admin role.

### Add Product Variant

```plaintext
POST /api/admin/products/{productId}/variants
```

Add a variant to a product (admin only).

**Path Parameters:**

- `productId` (required): Product ID

**Request Body:**

```json
{
  "sku": "PROD-004-M",
  "stock": 10,
  "attributes": [
    { "name": "color", "value": "Red" },
    { "name": "size", "value": "Medium" }
  ],
  "images": ["https://example.com/red-variant.jpg"],
  "is_default": false,
  "weight": 1.2,
  "price": 29.99
}
```

**Response Body:**

```json
{
  "success": true,
  "message": "Variant added successfully",
  "data": {
    "id": 11,
    "product_id": 4,
    "variant_name": "Color Red, Size Medium",
    "sku": "PROD-004-M",
    "stock": 10,
    "attributes": {
      "color": "Red",
      "size": "Medium"
    },
    "images": ["https://example.com/red-variant.jpg"],
    "is_default": false,
    "weight": 1.2,
    "price": 29.99,
    "currency": "USD",
    "created_at": "2025-07-07T10:30:45Z",
    "updated_at": "2025-07-07T10:30:45Z"
  }
}
```

**Status Codes:**

- `201 Created`: Variant created successfully
- `400 Bad Request`: Invalid request body or validation error
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Product not found
- `409 Conflict`: Variant with this SKU already exists

### Update Product Variant

```plaintext
PUT /api/admin/products/{productId}/variants/{variantId}
```

Update a product variant (admin only).

**Path Parameters:**

- `productId` (required): Product ID
- `variantId` (required): Variant ID

**Request Body:**

```json
{
  "sku": "PROD-004-M-UPDATED",
  "stock": 15,
  "attributes": [
    { "name": "color", "value": "Dark Red" },
    { "name": "size", "value": "Medium" }
  ],
  "images": ["https://example.com/dark-red-variant.jpg"],
  "is_default": false,
  "weight": 1.3,
  "price": 24.99
}
```

**Status Codes:**

- `200 OK`: Variant updated successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Product or variant not found
- `409 Conflict`: Variant with this SKU already exists

### Delete Product Variant

```plaintext
DELETE /api/admin/products/{productId}/variants/{variantId}
```

Delete a product variant (admin only).

**Path Parameters:**

- `productId` (required): Product ID
- `variantId` (required): Variant ID

**Status Codes:**

- `200 OK`: Variant deleted successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Product or variant not found
- `409 Conflict`: Cannot delete the last variant of a product or variant with existing orders

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
