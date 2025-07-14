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
- **Database**: SQLite (development) / PostgreSQL (production)
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
- SQLite (for local development) or PostgreSQL 15+ (for production)
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

### Email Configuration

The application includes configurable email settings for transactional emails. Configure the following environment variables in your `.env` file:

```bash
# Email Service Configuration
EMAIL_ENABLED=true                          # Enable/disable email functionality
EMAIL_SMTP_HOST=smtp.example.com            # SMTP server hostname
EMAIL_SMTP_PORT=587                         # SMTP server port (usually 587 for TLS)
EMAIL_SMTP_USERNAME=username                # SMTP authentication username
EMAIL_SMTP_PASSWORD=password                # SMTP authentication password

# Email Addresses and Branding
EMAIL_FROM_ADDRESS=noreply@example.com      # From address for outgoing emails
EMAIL_FROM_NAME=My Store                    # From name for outgoing emails
EMAIL_ADMIN_ADDRESS=admin@example.com       # Admin email for order notifications
EMAIL_CONTACT_ADDRESS=support@example.com   # Customer support contact email (used in templates)
STORE_NAME=My Store                         # Store name displayed in email templates
```

**Email Features:**

- **Order Confirmation**: Sent when orders are placed
- **Order Shipped**: Sent when orders are marked as shipped (with optional tracking)
- **Order Notifications**: Sent to admin when new orders are received
- **Checkout Recovery**: Sent to customers who abandon their carts (if implemented)

**Template Customization:**
Email templates are located in `templates/emails/` and can be customized to match your brand.

### Database Setup

The application supports both SQLite for local development and PostgreSQL for production.

#### Option 1: SQLite (Recommended for Local Development)

SQLite is the easiest way to get started with local development:

1. Copy the local development environment file:

```bash
cp .env.local .env
```

2. Run the application:

```bash
make dev-sqlite
# or
go run cmd/api/main.go
```

The SQLite database file (`commercify.db`) will be created automatically in the project root.

#### Option 2: PostgreSQL (Production Setup)

For production or if you prefer PostgreSQL for development:

1. **Using Docker (Recommended):**

```bash
# Start PostgreSQL with Docker
make db-start

# Setup database with migrations and seed data
make dev-setup
```

2. **Manual PostgreSQL Setup:**

Create a PostgreSQL user (optional):

```bash
createuser -s newuser
```

Create a PostgreSQL database:

```bash
createdb -U newuser commercify
```

Copy and configure environment file:

```bash
cp .env.example .env
# Edit .env and set:
# DB_DRIVER=postgres
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=your_user
# DB_PASSWORD=your_password
# DB_NAME=commercify
```

Run migrations:

```bash
go run cmd/migrate/main.go -up
```

Seed the database with sample data (optional):

```bash
go run cmd/seed/main.go -all
```

#### Database Commands

The project includes helpful Make commands for database management:

```bash
# SQLite Development
make dev-sqlite              # Start with SQLite
make dev-setup-sqlite         # Setup SQLite environment
make dev-reset-sqlite         # Reset SQLite database

# PostgreSQL Development
make dev-postgres             # Start with PostgreSQL
make dev-setup               # Setup PostgreSQL environment
make dev-reset               # Reset PostgreSQL environment

# Database Operations (PostgreSQL)
make db-start                # Start PostgreSQL container
make db-stop                 # Stop PostgreSQL container
make db-logs                 # View database logs
make db-clean                # Clean database and volumes
```

### Running the Application

# Build the application

```
go build -o bin/api cmd/api/main.go
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
- `DELETE /api/checkout/items/{sku}` - Removes an item from the current checkout session using SKU.
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

## Maintenance Commands

The project includes useful maintenance commands for database cleanup and optimization:

### Checkout Cleanup

The system provides two modes for managing expired and old checkouts:

```bash
# Regular cleanup (recommended for scheduled runs)
make expire-checkouts
# or
go run ./cmd/expire-checkouts
```

This command performs the following operations:

- Marks checkouts with customer/shipping info as **abandoned** after 15 minutes of inactivity
- **Deletes** empty checkouts older than 24 hours
- **Deletes** abandoned checkouts older than 7 days
- Marks remaining expired checkouts as **expired** (legacy support)

```bash
# Force deletion (use with caution)
make force-delete-checkouts
# or
go run ./cmd/expire-checkouts -force
```

This command performs aggressive cleanup:

- **Force deletes** all expired, abandoned, and completed checkouts
- **Force deletes** checkouts older than 30 days regardless of status
- Should be used carefully as it permanently removes checkout data

### Usage Examples

```bash
# Show help and available options
./bin/expire-checkouts --help

# Regular cleanup (safe for automation)
./bin/expire-checkouts

# Force delete all expired checkouts
./bin/expire-checkouts -force
```

### Scheduling Maintenance

For production environments, consider scheduling regular cleanup:

```bash
# Example crontab entry (runs every hour)
0 * * * * /path/to/commercify/bin/expire-checkouts

# Example crontab entry for weekly force cleanup (use with caution)
0 2 * * 0 /path/to/commercify/bin/expire-checkouts -force
```
