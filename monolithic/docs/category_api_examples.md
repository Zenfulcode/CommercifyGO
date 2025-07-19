# Category API Examples

This document provides example request bodies for the category system API endpoints.

## Public Category Endpoints

### List Categories

```plaintext
GET /api/categories
```

Get all categories in the system.

**Response Body:**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Electronics",
      "description": "Electronic devices and accessories",
      "parent_id": null,
      "created_at": "2025-07-07T10:30:45Z",
      "updated_at": "2025-07-07T10:30:45Z"
    },
    {
      "id": 2,
      "name": "Smartphones",
      "description": "Mobile phones and accessories",
      "parent_id": 1,
      "created_at": "2025-07-07T10:30:45Z",
      "updated_at": "2025-07-07T10:30:45Z"
    },
    {
      "id": 3,
      "name": "Laptops",
      "description": "Laptop computers and accessories",
      "parent_id": 1,
      "created_at": "2025-07-07T10:30:45Z",
      "updated_at": "2025-07-07T10:30:45Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 3,
    "total": 3
  }
}
```

**Status Codes:**

- `200 OK`: Categories retrieved successfully
- `500 Internal Server Error`: Failed to retrieve categories

### Get Category

```plaintext
GET /api/categories/{id}
```

Get a category by ID.

**Path Parameters:**

- `id` (required): Category ID

**Response Body:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices and accessories",
    "parent_id": null,
    "created_at": "2025-07-07T10:30:45Z",
    "updated_at": "2025-07-07T10:30:45Z"
  }
}
```

**Status Codes:**

- `200 OK`: Category retrieved successfully
- `404 Not Found`: Category not found
- `500 Internal Server Error`: Failed to retrieve category

### Get Child Categories

```plaintext
GET /api/categories/{id}/children
```

Get child categories of a parent category.

**Path Parameters:**

- `id` (required): Parent category ID

**Response Body:**

```json
{
  "success": true,
  "data": [
    {
      "id": 2,
      "name": "Smartphones",
      "description": "Mobile phones and accessories",
      "parent_id": 1,
      "created_at": "2025-07-07T10:30:45Z",
      "updated_at": "2025-07-07T10:30:45Z"
    },
    {
      "id": 3,
      "name": "Laptops",
      "description": "Laptop computers and accessories",
      "parent_id": 1,
      "created_at": "2025-07-07T10:30:45Z",
      "updated_at": "2025-07-07T10:30:45Z"
    }
  ]
}
```

**Status Codes:**

- `200 OK`: Child categories retrieved successfully
- `404 Not Found`: Parent category not found
- `500 Internal Server Error`: Failed to retrieve child categories

## Admin Category Endpoints

All admin category endpoints require authentication and admin role.

### Create Category

```plaintext
POST /api/admin/categories
```

Create a new category (admin only).

**Request Body:**

```json
{
  "name": "Gaming",
  "description": "Gaming devices and accessories",
  "parent_id": 1
}
```

**Response Body:**

```json
{
  "success": true,
  "message": "Category created successfully",
  "data": {
    "id": 4,
    "name": "Gaming",
    "description": "Gaming devices and accessories",
    "parent_id": 1,
    "created_at": "2025-07-07T10:30:45Z",
    "updated_at": "2025-07-07T10:30:45Z"
  }
}
```

**Status Codes:**

- `201 Created`: Category created successfully
- `400 Bad Request`: Invalid request body or validation error
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

### Update Category

```plaintext
PUT /api/admin/categories/{id}
```

Update an existing category (admin only).

**Path Parameters:**

- `id` (required): Category ID

**Request Body:**

```json
{
  "name": "Gaming Consoles",
  "description": "Gaming consoles and accessories",
  "parent_id": 1
}
```

**Response Body:**

```json
{
  "success": true,
  "message": "Category updated successfully",
  "data": {
    "id": 4,
    "name": "Gaming Consoles",
    "description": "Gaming consoles and accessories",
    "parent_id": 1,
    "created_at": "2025-07-07T10:30:45Z",
    "updated_at": "2025-07-07T10:35:15Z"
  }
}
```

**Status Codes:**

- `200 OK`: Category updated successfully
- `400 Bad Request`: Invalid request body or validation error
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Category not found

### Delete Category

```plaintext
DELETE /api/admin/categories/{id}
```

Delete a category (admin only).

**Path Parameters:**

- `id` (required): Category ID

**Response Body:**

```json
{
  "success": true,
  "message": "Category deleted successfully"
}
```

**Status Codes:**

- `200 OK`: Category deleted successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Category not found
- `409 Conflict`: Cannot delete category with existing products or subcategories
