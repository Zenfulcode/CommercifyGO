# User API Examples

This document provides example request bodies for the user system API endpoints.

## Public User Endpoints

### Register User

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

Example response:

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 123,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Smith",
      "role": "user",
      "created_at": "2023-05-15T10:30:45Z",
      "updated_at": "2023-05-15T10:30:45Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "",
    "expires_in": 3600
  }
}
```

**Status Codes:**

- `201 Created`: User registered successfully
- `400 Bad Request`: Invalid request body or validation error
- `409 Conflict`: Email already in use

### Login

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

Example response:

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 123,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Smith",
      "role": "user",
      "created_at": "2023-05-15T10:30:45Z",
      "updated_at": "2023-05-15T10:30:45Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "",
    "expires_in": 3600
  }
}
```

**Status Codes:**

- `200 OK`: Authentication successful
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Invalid credentials

## Authenticated User Endpoints

### Get User Profile

```plaintext
GET /api/users/me
```

Retrieve the current authenticated user's profile.

Example response:

```json
{
  "success": true,
  "data": {
    "id": 123,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Smith",
    "role": "user",
    "created_at": "2023-05-15T10:30:45Z",
    "updated_at": "2023-05-15T10:30:45Z"
  }
}
```

**Status Codes:**

- `200 OK`: Profile retrieved successfully
- `401 Unauthorized`: Not authenticated

### Update User Profile

```plaintext
PUT /api/users/me
```

Update the current authenticated user's profile.

**Request Body:**

```json
{
  "first_name": "Johnny",
  "last_name": "Smith"
}
```

Example response:

```json
{
  "success": true,
  "data": {
    "id": 123,
    "email": "user@example.com",
    "first_name": "Johnny",
    "last_name": "Smith",
    "role": "user",
    "created_at": "2023-05-15T10:30:45Z",
    "updated_at": "2023-05-16T14:22:30Z"
  }
}
```

**Status Codes:**

- `200 OK`: Profile updated successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated

### Change Password

```plaintext
PUT /api/users/me/password
```

Change the current authenticated user's password.

**Request Body:**

```json
{
  "current_password": "Password123!",
  "new_password": "NewPassword456!"
}
```

Example response:

```json
{
  "message": "Password changed successfully"
}
```

**Status Codes:**

- `200 OK`: Password changed successfully
- `400 Bad Request`: Invalid request body or current password is incorrect
- `401 Unauthorized`: Not authenticated

## Admin User Management Endpoints

### List Users

```plaintext
GET /api/admin/users
```

List all users (admin only).

**Query Parameters:**

- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10)

Example response:

```json
{
  "success": true,
  "data": [
    {
      "id": 123,
      "email": "user@example.com",
      "first_name": "Johnny",
      "last_name": "Smith",
      "role": "user",
      "created_at": "2023-05-15T10:30:45Z",
      "updated_at": "2023-05-16T14:22:30Z"
    },
    {
      "id": 124,
      "email": "seller@example.com",
      "first_name": "Sarah",
      "last_name": "Johnson",
      "role": "seller",
      "created_at": "2023-05-10T09:15:22Z",
      "updated_at": "2023-05-10T09:15:22Z"
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

- `200 OK`: Users retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (admin only)

### Get User By ID

```plaintext
GET /api/admin/users/{id}
```

Get a specific user by ID (admin only).

Example response:

```json
{
  "id": 123,
  "email": "user@example.com",
  "first_name": "Johnny",
  "last_name": "Smith",
  "role": "user",
  "created_at": "2023-05-15T10:30:45Z",
  "updated_at": "2023-05-16T14:22:30Z",
  "addresses": [
    {
      "id": 45,
      "user_id": 123,
      "name": "Home",
      "street_address": "123 Main St",
      "city": "San Francisco",
      "state": "CA",
      "postal_code": "94105",
      "country": "US",
      "is_default": true
    }
  ],
  "orders_count": 5,
  "last_order_date": "2023-05-20T16:45:22Z"
}
```

**Status Codes:**

- `200 OK`: User retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: User not found

### Update User Role

```plaintext
PUT /api/admin/users/{id}/role
```

Update a user's role (admin only).

**Request Body:**

```json
{
  "role": "seller"
}
```

Example response:

```json
{
  "id": 123,
  "email": "user@example.com",
  "first_name": "Johnny",
  "last_name": "Smith",
  "role": "seller",
  "created_at": "2023-05-15T10:30:45Z",
  "updated_at": "2023-05-22T09:12:30Z"
}
```

**Status Codes:**

- `200 OK`: User role updated successfully
- `400 Bad Request`: Invalid role
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: User not found

### Deactivate User

```plaintext
PUT /api/admin/users/{id}/deactivate
```

Deactivate a user account (admin only).

Example response:

```json
{
  "id": 123,
  "email": "user@example.com",
  "first_name": "Johnny",
  "last_name": "Smith",
  "role": "seller",
  "active": false,
  "created_at": "2023-05-15T10:30:45Z",
  "updated_at": "2023-05-22T11:45:15Z"
}
```

**Status Codes:**

- `200 OK`: User deactivated successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: User not found

### Reactivate User

```plaintext
PUT /api/admin/users/{id}/activate
```

Reactivate a user account (admin only).

Example response:

```json
{
  "id": 123,
  "email": "user@example.com",
  "first_name": "Johnny",
  "last_name": "Smith",
  "role": "seller",
  "active": true,
  "created_at": "2023-05-15T10:30:45Z",
  "updated_at": "2023-05-23T08:30:00Z"
}
```

**Status Codes:**

- `200 OK`: User reactivated successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: User not found

## Example Workflows

### User Registration and Login Flow

1. User registers with email, password, and profile information
2. System creates user account and returns JWT token
3. User can then access authenticated endpoints using the token
4. If token expires, user can log in again to get a new token

### Profile Management Flow

1. User logs in and receives a JWT token
2. User can view their profile details
3. User can update their profile information (first name, last name)
4. User can change their password by providing current and new passwords

### Admin User Management Flow

1. Admin logs in with admin credentials
2. Admin can view a list of all users in the system
3. Admin can view detailed information about any specific user
4. Admin can update a user's role (e.g., promote to seller)
5. Admin can deactivate/reactivate user accounts as needed
