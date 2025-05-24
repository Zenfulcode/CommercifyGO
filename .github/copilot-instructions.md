# Commercify Development Guidelines

## Introduction

This document outlines the development guidelines and best practices for working with the Commercify e-commerce backend API. Following these guidelines will ensure code consistency, maintainability, and adherence to the project's architectural principles.

## Architectural Principles

Commercify follows **Clean Architecture** principles with clear separation of concerns:

1. **Domain Layer**: Contains business entities and repository interfaces
2. **Application Layer**: Contains use cases that implement business logic
3. **Infrastructure Layer**: Contains implementation of repositories and external services
4. **Interfaces Layer**: Contains API handlers and middleware

### Dependencies Rule

- Inner layers must not depend on outer layers
- Domain layer has no dependencies on other layers
- Application layer depends only on domain layer
- Infrastructure and interfaces layers depend on application and domain layers

## Code Structure

```
├── cmd/                  # Application entry points
│   ├── api/              # API server
│   ├── migrate/          # Database migration tool
│   └── seed/             # Database seeding tool
├── config/               # Configuration
├── internal/             # Internal packages
│   ├── api/              # API layer (handlers, middleware, server)
│   ├── application/      # Application layer (use cases)
│   ├── domain/           # Domain layer (entities, repositories interfaces)
│   └── infrastructure/   # Infrastructure layer (repositories impl, services)
├── migrations/           # Database migrations
├── templates/            # Email templates
└── testutil/             # Testing utilities
```

## Coding Standards

### General Guidelines

- Always format your code using `go fmt`
- Follow standard Go naming conventions
- Keep functions small and focused
- Document all exported functions, types, and constants
- Use meaningful variable and function names
- Handle all errors explicitly

### Naming Conventions

- Use **camelCase** for variable names (`userID`, `orderTotal`)
- Use **PascalCase** for exported functions, types (`ProcessOrder`, `UserRepository`)
- Use **snake_case** for file names (`order_repository.go`, `user_service.go`)
- Use **SCREAMING_SNAKE_CASE** for constants (`MAX_RETRY_ATTEMPTS`)

### Error Handling

- Return errors rather than using panic
- Use custom error types for domain-specific errors
- Wrap errors with context when crossing architectural boundaries
- Log errors at appropriate levels

```go
// Good
if err != nil {
    return nil, fmt.Errorf("fetching user profile: %w", err)
}

// Avoid
if err != nil {
    panic(err)
}
```

## Database Practices

### Migrations

To create new migrations, run the following command in the CLI:  
`migrate create -ext sql -dir migrations -seq "<name>"`

Note: Replace `<name>` with the desired migration name. Use quotes around the name if it contains spaces or special characters. For example: `"add_friendly_numbers"`.

### Queries

- Use prepared statements for all database queries
- Never concatenate user input into SQL strings
- Use transactions for operations that affect multiple tables

## API Design

### RESTful Principles

- Use proper HTTP methods (`GET`, `POST`, `PUT`, `DELETE`)
- Return appropriate HTTP status codes
- Use plural nouns for resource collections (`/products`, `/orders`)
- Use nested resources for relationships (`/products/{id}/variants`)

### Versioning

- Version APIs in the URL path (`/api/v1/products`)
- Maintain backward compatibility within a version

### Security

- Validate all input data
- Use HTTPS for all production environments
- Implement rate limiting for authentication endpoints
- Never store sensitive information in logs

## Payment System

### Multi-Provider Architecture

- Implement new payment providers by implementing the `PaymentService` interface
- Configure providers through environment variables
- Use the adapter pattern to standardize payment provider interactions
- Always test payment flows in sandbox environments before going live

## Testing

### Unit Tests

- Write unit tests for all business logic
- Use table-driven tests for testing multiple cases
- Use mocks for external dependencies

## Documentation

### Code Documentation

- Document all exported functions, types, and constants
- Use godoc format for documentation
- Provide examples for complex functions

### API Documentation

- Keep /docs/<filename>.md up to date
- Include request/response examples
- Document authentication requirements

#### API Documentation Format

All new API endpoints must be documented using the following format in RESTAPI.md:

````markdown
### [Endpoint Name]

```plaintext
[HTTP Method] [Endpoint Path]
```
````

**Request Body:**

```json
{
  // Example request body with all fields
}
```

**Response Body:**

```json
{
  // Example response body with all fields
}
```

**Status Codes:**

- `[Status Code]`: [Description]
- `[Status Code]`: [Description]
- ...

````

Example:

```markdown
### Create Discount

```plaintext
POST /api/discounts
````

**Request Body:**

```json
{
  "code": "SUMMER2025",
  "type": "basket",
  "method": "percentage",
  "value": 15.0,
  "min_order_value": 50.0,
  "max_discount_value": 30.0,
  "product_ids": [],
  "category_ids": [],
  "start_date": "2025-05-01T00:00:00Z",
  "end_date": "2025-08-31T23:59:59Z",
  "usage_limit": 500
}
```

**Response Body:**

```json
{
  "id": 7,
  "code": "SUMMER2025",
  "type": "basket",
  "method": "percentage",
  "value": 15,
  "min_order_value": 0,
  "max_discount_value": 0,
  "start_date": "0001-01-01T00:00:00Z",
  "end_date": "0001-01-01T00:00:00Z",
  "usage_limit": 0,
  "current_usage": 0,
  "active": true,
  "created_at": "2025-04-23T16:20:29.143443+02:00",
  "updated_at": "2025-04-23T16:20:29.143443+02:00"
}
```

**Status Codes:**

- `201 Created`: Discount created successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `409 Conflict`: Discount code already exists


## Commit Messages

- Use present tense ("Add feature" not "Added feature")
- First line should be a summary (max 50 chars)
- Optionally provide more detailed description after summary
- Reference issue numbers if applicable

## Email Templates

- Keep HTML email templates responsive
- Test emails on multiple email clients
- Provide both HTML and plain text versions

By following these guidelines, we'll maintain a high-quality, consistent codebase that is easy to understand, extend, and maintain.
