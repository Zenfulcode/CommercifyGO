# Commercify Backend API

A robust, scalable e-commerce backend API built with Go, following clean architecture principles and best practices.

## Features

- **User Management**: Registration, authentication, profile management
- **Product Management**: CRUD operations, categories, variants, search
- **Checkout System**: Session-based checkout, add/update/remove items, apply discounts, shipping methods
- **Order Processing**: Create orders, payment processing, order status tracking
- **Payment Integration**: Support for multiple payment providers (Stripe, MobilePay, etc.)
- **Email Notifications**: Order confirmations, status updates

> **New**: The Checkout System now uses a session-based approach with cookies. See [Checkout Session Documentation](/docs/checkout_session.md) for details.

## Technology Stack

- **Language**: Go 1.20+
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Payment Processing**: Stripe, MobilePay
- **Email**: SMTP integration

## Project Structure

The project follows clean architecture principles with clear separation of concerns:

```
├── cmd/ # Application entry points
│ ├── api/ # API server
│ ├── migrate/ # Database migration tool
│ └── seed/ # Database seeding tool
├── config/ # Configuration
├── internal/ # Internal packages
│ ├── api/ # API layer (handlers, middleware, server)
│ ├── application/ # Application layer (use cases)
│ ├── domain/ # Domain layer (entities, repositories interfaces)
│ └── infrastructure/ # Infrastructure layer (repositories implementation, services)
├── migrations/ # Database migrations
├── templates/ # Email templates
└── testutil/ # Testing utilities
```

## Setup and Installation

### Prerequisites

- Go 1.20+
- PostgreSQL 15 (Only tested on v15)
- Docker (optional)

### Docker Setup

For a quick start with Docker Compose:

1. Clone the repository:

```bash
git clone https://github.com/zenfulcode/commercifygo.git
cd commercify
```

2. Start the services using Docker Compose:

```bash
docker-compose up -d
```

This will start:

- PostgreSQL database
- Commercify API server

3. Run database migrations (First startup also migrates automatically):

```bash
docker-compose exec api /app/commercify-migrate -up
```

4. Seed the database with sample data (optional):

```bash
docker-compose exec api /app/commercify-seed -all
```

5. Access the API at `http://localhost:6091`

6. To stop the services:

```bash
docker-compose down
```

### Environment Variables

Create a `.env` file in the root directory by copying the `.env.example`

```bash
cp .env.example .env
```

### Database Setup

1. Create a PostgreSQL user (optional):

```bash
createuser -s newuser
```

2. Create a PostgreSQL database:

```bash
createdb -U newuser commercify
```

3. Run migrations:

```bash
go run cmd/migrate/main.go -up
```

4. Seed the database with sample data (optional):

```bash
go run cmd/seed/main.go -all
```

### Running the Application

# Build the application

```
go build -o commercify cmd/api/main.go
```

# Run the application

```bash
./commercify
```

Or simply:

```bash
go run cmd/api/main.go
```

## API Documentation

### Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

### API Endpoints

#### Checkout

- `GET /api/checkout` - Retrieves the current checkout session for a user. If no checkout exists, a new one will be created.
- `POST /api/checkout/items` - Adds a product item to the current checkout session.
- `PUT /api/checkout/items/{productId}` Updates the quantity or variant of an item in the current checkout.
- `DELETE /api/checkout/items/{productId}` - Removes an item from the current checkout session.
- `DELETE /api/checkout` - Removes all items from the current checkout session.
- `PUT /api/checkout/shipping-addres` - Sets the shipping address for the current checkout.
- `PUT /api/checkout/billing-address` - Sets the billing address for the current checkout.

#### Users

- `POST /api/users/register` - Register a new user
- `POST /api/users/login` - Login and get JWT token
- `GET /api/users/me` - Get current user profile
- `PUT /api/users/me` - Update user profile
- `PUT /api/users/me/password` - Change password
- `GET /api/admin/users` - List all users (admin only)
- `GET /api/admin/users/{id}` - Get user by ID (admin only)
- `PUT /api/admin/users/{id}/role` - Update user role (admin only)
- `PUT /api/admin/users/{id}/deactivate` - Deactivate user (admin only)
- `PUT /api/admin/users/{id}/activate` - Reactivate user (admin only)

#### Products

- `GET /api/admin/products` - List products with pagination
- `GET /api/products/{id}` - Get product details
- `GET /api/products/search` - Search products
- `GET /api/categories` - List product categories
- `POST /api/admin/products` - Create product
- `PUT /api/admin/products/{id}` - Update product
- `DELETE /api/admin/products/{id}` - Delete product

#### Product Variants

- `POST /api/admin/products/{productId}/variants` - Add variant
- `PUT /api/admin/products/{productId}/variants/{variantId}` - Update variant
- `DELETE /api/admin/products/{productId}/variants/{variantId}` - Delete variant

#### Orders

- `POST /api/orders` - Create order for authenticated user
- `GET /api/orders/{id}` - Get order details
- `GET /api/orders` - List user's orders
- `POST /api/orders/{id}/payment` - Process payment for user order
- `POST /api/orders/{id}/discounts` - Apply discount to order
- `DELETE /api/orders/{id}/discounts` - Remove discount from order
- `GET /api/admin/orders` - List all orders (admin only)
- `PUT /api/admin/orders/{id}/status` - Update order status (admin only)

#### Payment

- `GET /api/payment/providers` - Get available payment providers
- `POST /api/admin/payments/{paymentId}/capture` - Capture payment (admin only)
- `POST /api/admin/payments/{paymentId}/cancel` - Cancel payment (admin only)
- `POST /api/admin/payments/{paymentId}/refund` - Refund payment (admin only)

#### Shipping

- `POST /api/shipping/options` - Calculate shipping options for address and order
- `POST /api/shipping/rates/{id}/cost` - Calculate cost for specific shipping rate
- `POST /api/admin/shipping/methods` - Create shipping method (admin only)
- `PUT /api/admin/shipping/methods/{id}` - Update shipping method (admin only)
- `POST /api/admin/shipping/zones` - Create shipping zone (admin only)
- `PUT /api/admin/shipping/zones/{id}` - Update shipping zone (admin only)
- `POST /api/admin/shipping/rates` - Create shipping rate (admin only)
- `PUT /api/admin/shipping/rates/{id}` - Update shipping rate (admin only)
- `POST /api/admin/shipping/rates/weight` - Add weight-based rate (admin only)
- `POST /api/admin/shipping/rates/value` - Add value-based rate (admin only)

#### Discounts

- `GET /api/discounts` - List active discounts
- `POST /api/orders/{id}/discounts` - Apply discount to order
- `DELETE /api/orders/{id}/discounts` - Remove discount from order
- `POST /api/admin/discounts` - Create discount (admin only)
- `PUT /api/admin/discounts/{id}` - Update discount (admin only)
- `DELETE /api/admin/discounts/{id}` - Delete discount (admin only)
- `GET /api/admin/discounts` - List all discounts (admin only)

#### Webhooks

- `POST /api/webhooks/stripe` - Stripe webhook endpoint
- `POST /api/webhooks/mobilepay` - MobilePay webhook endpoint

## Database Schema

The database consists of the following tables:

### Users and Authentication

- `users` - User accounts and authentication information

### Products

- `categories` - Product categories with hierarchical structure
- `products` - Product information including name, description, price, and stock
- `product_variants` - Variations of products with different attributes (size, color, etc.)

### Orders

- `orders` - Customer orders with status, amounts, and addresses
- `order_items` - Individual items in orders

### Payments

- `payment_transactions` - Record of payment attempts, successes, and failures

### Discounts

- `discounts` - Promotion codes with various discount types and rules

### Shipping

- `shipping_methods` - Available shipping methods
- `shipping_zones` - Geographic shipping zones
- `shipping_rates` - Shipping pricing based on weight, value, or other factors

### Webhooks

- `webhooks` - Configuration for external service webhook endpoints

## Development

### Running Tests

# Run all tests

```bash
go test ./...
```

# Run tests with coverage

```bash
go test -cover ./...
```

### Adding Migrations

# Create a new migration

Install the [golang-migrate/migrate](https://github.com/golang-migrate/migrate) tool

### Homebrew

```bash
brew install migrate
```

using the cli tool to generate migration files, then you can use the following command, which creates both the files in the right format:

```bash
migrate create -ext sql -dir migrations -seq add_friendly_numbers
```

otherwise you can create them manually using:

```bash
touch migrations/[sequence]_[migration_name].up.sql
touch migrations/[sequence]_[migration_name].down.sql
```

Where `migrations` is the migrations folder, the `sequence` is the 6 digits in front and `migration_name` is a short description

## Multi-Provider Payment System

The application supports multiple payment providers through a flexible payment service architecture:

- **Stripe**: Credit card payments
- **MobilePay** Mobile payments
- **Mock**: Test payment provider for development

Payment providers can be enabled/disabled through configuration, and new providers can be added by implementing the `PaymentService` interface.

## User-Friendly Identifiers

The system uses user-friendly identifiers for better readability:

- **Order Numbers**: Format `ORD-YYYYMMDD-000001` (date-based with sequential numbering)
- **Product Numbers**: Format `PROD-000001` (sequential numbering)

These identifiers make it easier to reference orders and products in the UI and customer communications.

## Payment Provider Implementations

### Stripe Payment Provider

Commercify implements Stripe as a payment provider following Clean Architecture principles. The implementation:

- Supports credit card payments using Stripe's Payment Intents API
- Handles 3D Secure authentication flows
- Provides webhook integration for asynchronous event handling
- Manages payment lifecycle (authorize, capture, refund, cancel)

#### Stripe Setup

1. Create a Stripe account at [stripe.com](https://stripe.com)
2. Get your API keys from the Stripe Dashboard
3. Set the following environment variables:

```
STRIPE_ENABLED=true
STRIPE_SECRET_KEY=sk_test_your_key
STRIPE_PUBLIC_KEY=pk_test_your_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_signing_secret
STRIPE_PAYMENT_DESCRIPTION=Commercify Store Purchase
```

#### Webhook Configuration

To handle asynchronous payment events (3D Secure authentication, payment success/failure):

1. Create a webhook endpoint in your Stripe Dashboard
2. Point it to: `https://your-domain.com/api/webhooks/stripe`
3. Select the following events to listen for:
   - `payment_intent.succeeded`
   - `payment_intent.payment_failed`
   - `payment_intent.canceled`
   - `payment_intent.requires_action`
   - `payment_intent.processing`
   - `payment_intent.amount_capturable_updated`
   - `charge.succeeded`
   - `charge.failed`
   - `charge.refunded`
   - `charge.dispute.created`
   - `charge.dispute.closed`
4. Copy the signing secret and set it as `STRIPE_WEBHOOK_SECRET` in your environment

#### Payment Flows

Commercify supports several payment flows with Stripe:

**Direct Payment**
Payment is authorized and captured immediately.

**Authorization and Capture**
Payment is first authorized, then captured later when the order is fulfilled.

**3D Secure Authentication**
When required by the bank, customers will be redirected to complete 3D Secure authentication.

#### Testing Stripe Integration

Use Stripe's test cards for development:

- `4242 4242 4242 4242` - Successful payment
- `4000 0000 0000 3220` - 3D Secure authentication required
- `4000 0000 0000 9995` - Payment declined

For more test card numbers, visit [Stripe's testing documentation](https://stripe.com/docs/testing).
