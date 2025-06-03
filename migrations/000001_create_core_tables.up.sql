-- Create core tables for Commercify e-commerce system
-- This migration consolidates all the essential tables needed for the application

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id INTEGER REFERENCES categories(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create products table with final schema (including active column and int64 prices)
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price BIGINT NOT NULL, -- stored as cents
    stock INTEGER NOT NULL DEFAULT 0,
    category_id INTEGER NOT NULL REFERENCES categories(id),
    images JSONB NOT NULL DEFAULT '[]',
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create product variants table
CREATE TABLE IF NOT EXISTS product_variants (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    sku VARCHAR(100) UNIQUE,
    price BIGINT NOT NULL, -- stored as cents
    stock INTEGER NOT NULL DEFAULT 0,
    weight DECIMAL(10, 3) NOT NULL DEFAULT 0,
    attributes JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create currency settings table
CREATE TABLE IF NOT EXISTS currencies (
    id SERIAL PRIMARY KEY,
    currency_code VARCHAR(3) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    exchange_rate DECIMAL(10, 6) NOT NULL DEFAULT 1.0,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create shipping methods table
CREATE TABLE IF NOT EXISTS shipping_methods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    base_rate BIGINT NOT NULL, -- stored as cents
    rate_per_kg BIGINT NOT NULL DEFAULT 0, -- stored as cents
    min_delivery_days INTEGER,
    max_delivery_days INTEGER,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create orders table with final schema
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    friendly_id VARCHAR(20) NOT NULL UNIQUE,
    user_id INTEGER REFERENCES users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    shipping_address JSONB NOT NULL,
    billing_address JSONB NOT NULL,
    shipping_method_id INTEGER REFERENCES shipping_methods(id),
    shipping_cost BIGINT NOT NULL DEFAULT 0, -- stored as cents
    total_amount BIGINT NOT NULL, -- stored as cents
    final_amount BIGINT NOT NULL, -- stored as cents after discounts
    discount_amount BIGINT NOT NULL DEFAULT 0, -- stored as cents
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    payment_provider VARCHAR(255),
    payment_id VARCHAR(255),
    tracking_code VARCHAR(100),
    customer_details JSONB NOT NULL DEFAULT '{}',
    discount_code VARCHAR(100),
    applied_discount JSONB,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP
);

-- Create order items table with final schema
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id),
    product_variant_id INTEGER REFERENCES product_variants(id),
    quantity INTEGER NOT NULL,
    price BIGINT NOT NULL, -- stored as cents
    subtotal BIGINT NOT NULL, -- stored as cents
    product_name VARCHAR(255) NOT NULL,
    variant_name VARCHAR(255),
    sku VARCHAR(100),
    weight DECIMAL(10, 3) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL
);

-- Create payment transactions table
CREATE TABLE IF NOT EXISTS payment_transactions (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    transaction_id VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,  -- authorize, capture, refund, cancel
    status VARCHAR(50) NOT NULL,  -- successful, failed, pending
    amount BIGINT NOT NULL, -- stored as cents
    currency VARCHAR(3) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    raw_response TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create discounts table
CREATE TABLE IF NOT EXISTS discounts (
    id SERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL, -- 'product' or 'basket'
    method VARCHAR(20) NOT NULL, -- 'percentage' or 'fixed'
    value BIGINT NOT NULL, -- percentage * 100 or cents for fixed
    min_order_value BIGINT NOT NULL DEFAULT 0, -- stored as cents
    max_discount_value BIGINT NOT NULL DEFAULT 0, -- stored as cents, 0 = no limit
    product_ids INTEGER[] DEFAULT ARRAY[]::INTEGER[],
    category_ids INTEGER[] DEFAULT ARRAY[]::INTEGER[],
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    usage_limit INTEGER NOT NULL DEFAULT 0, -- 0 = unlimited
    current_usage INTEGER NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create webhooks table
CREATE TABLE IF NOT EXISTS webhooks (
    id SERIAL PRIMARY KEY,
    url VARCHAR(500) NOT NULL,
    events TEXT[] NOT NULL,
    secret VARCHAR(255) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    action_url VARCHAR(500), -- URL for actions like payment processing
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create checkouts table (replacement for cart system)
CREATE TABLE IF NOT EXISTS checkouts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    session_id VARCHAR(255) NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    shipping_address JSONB NOT NULL DEFAULT '{}',
    billing_address JSONB NOT NULL DEFAULT '{}',
    shipping_method_id INTEGER REFERENCES shipping_methods(id) ON DELETE SET NULL,
    payment_provider VARCHAR(255),
    total_amount BIGINT NOT NULL DEFAULT 0, -- stored as cents
    shipping_cost BIGINT NOT NULL DEFAULT 0, -- stored as cents
    total_weight DECIMAL(10, 3) NOT NULL DEFAULT 0,
    customer_details JSONB NOT NULL DEFAULT '{}',
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    discount_code VARCHAR(100),
    discount_amount BIGINT NOT NULL DEFAULT 0, -- stored as cents
    final_amount BIGINT NOT NULL DEFAULT 0, -- stored as cents
    applied_discount JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_activity_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    converted_order_id INTEGER REFERENCES orders(id) ON DELETE SET NULL
);

-- Create checkout items table
CREATE TABLE IF NOT EXISTS checkout_items (
    id SERIAL PRIMARY KEY,
    checkout_id INTEGER NOT NULL REFERENCES checkouts(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id),
    product_variant_id INTEGER REFERENCES product_variants(id) ON DELETE SET NULL,
    quantity INTEGER NOT NULL,
    price BIGINT NOT NULL, -- stored as cents
    weight DECIMAL(10, 3) NOT NULL DEFAULT 0,
    product_name VARCHAR(255) NOT NULL,
    variant_name VARCHAR(255),
    sku VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create all necessary indexes
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_active ON products(active);
CREATE INDEX IF NOT EXISTS idx_product_variants_product ON product_variants(product_id);
CREATE INDEX IF NOT EXISTS idx_product_variants_sku ON product_variants(sku);

CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_friendly_id ON orders(friendly_id);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);

CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product ON order_items(product_id);
CREATE INDEX IF NOT EXISTS idx_order_items_variant ON order_items(product_variant_id);

CREATE INDEX IF NOT EXISTS idx_payment_transactions_order_id ON payment_transactions(order_id);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_transaction_id ON payment_transactions(transaction_id);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_type ON payment_transactions(type);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_status ON payment_transactions(status);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_created_at ON payment_transactions(created_at);

CREATE INDEX IF NOT EXISTS idx_discounts_code ON discounts(code);
CREATE INDEX IF NOT EXISTS idx_discounts_active ON discounts(active);
CREATE INDEX IF NOT EXISTS idx_discounts_dates ON discounts(start_date, end_date);

CREATE INDEX IF NOT EXISTS idx_checkouts_user_id ON checkouts(user_id);
CREATE INDEX IF NOT EXISTS idx_checkouts_session_id ON checkouts(session_id);
CREATE INDEX IF NOT EXISTS idx_checkouts_status ON checkouts(status);
CREATE INDEX IF NOT EXISTS idx_checkouts_expires_at ON checkouts(expires_at);
CREATE INDEX IF NOT EXISTS idx_checkouts_converted_order_id ON checkouts(converted_order_id);

CREATE INDEX IF NOT EXISTS idx_checkout_items_checkout_id ON checkout_items(checkout_id);
CREATE INDEX IF NOT EXISTS idx_checkout_items_product_id ON checkout_items(product_id);
CREATE INDEX IF NOT EXISTS idx_checkout_items_product_variant_id ON checkout_items(product_variant_id);

-- Create unique constraints
CREATE UNIQUE INDEX IF NOT EXISTS idx_checkouts_unique_active_session 
ON checkouts(session_id) 
WHERE status = 'active' AND session_id IS NOT NULL;

-- Ensure only one default currency
CREATE UNIQUE INDEX IF NOT EXISTS idx_currency_default 
ON currencies(is_default) 
WHERE is_default = true;

-- Insert default currency
INSERT INTO currencies (currency_code, name, symbol, exchange_rate, is_default, created_at, updated_at)
VALUES
('USD', 'US Dollar', '$', 1.0, true, NOW(), NOW())
('EUR', 'Euro', '€', 1.2, false, NOW(), NOW())
('GBP', 'British Pound', '£', 1.4, false, NOW(), NOW())
('DKK', 'Danish Krone', 'kr', 0.15, false, NOW(), NOW())
ON CONFLICT (currency_code) DO NOTHING;
