# Product API Examples

This document provides example request bodies and responses for the product system API endpoints.

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
      "has_variants": false
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
  "has_variants": false
}
```

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
    "has_variants": false
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
    "has_variants": false
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
      "has_variants": false
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

### Setting Product Currency Prices

When creating or updating products and their variants, you can specify prices in multiple currencies using the `currency_prices` array property. Each entry in this array should include:

- `currency_code`: The three-letter ISO code of the currency (e.g., "USD", "EUR", "GBP")
- `price`: The price in the specified currency
- `compare_price` (optional): The compare price (original/before discount price) in the specified currency

The system always requires a price in the default currency, and additional currency prices are optional. If a currency price is not specified for a particular currency, the system will automatically convert the price from the default currency using the current exchange rate when needed.

### Retrieving Products with Specific Currency Prices

When retrieving products, you can specify a currency code in the query parameters to get prices in that currency:

```
GET /api/products/1?currency=EUR
```

This will return the product with prices in euros, either using the explicitly set euro prices or converting from the default currency if no specific euro prices are set.

## Example Workflow

### Product Management Flow (Seller)

1. Seller creates a base product through the seller interface
2. If the product has variants, seller adds variants with different attributes (color, size, etc.)
3. Seller can update product information or variant details as needed
4. Seller can manage inventory levels for products and variants
5. Seller can deactivate or delete products when they're no longer available

### Product Shopping Flow (Customer)

1. Customers browse products by category or use the search function
2. Customers can view detailed product information including available variants
3. When adding to checkout, customers select specific variants if the product has them
4. Products and variants are displayed with current inventory levels
5. Out-of-stock products or variants can be marked as unavailable for purchase
