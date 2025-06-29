-- Commercify E-commerce Database Schema
-- This is a consolidated migration that replaces all previous migrations
-- with a clean, up-to-date schema without legacy components

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
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    parent_id INTEGER REFERENCES categories(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create currencies table
CREATE TABLE IF NOT EXISTS currencies (
    code VARCHAR(3) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    exchange_rate DECIMAL(16, 6) NOT NULL DEFAULT 1.0,
    is_default BOOLEAN NOT NULL DEFAULT false,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Insert default USD currency
INSERT INTO currencies (code, name, symbol, exchange_rate, is_default, is_enabled)
VALUES
 ('USD', 'US Dollar', '$', 1.0, true, true),
 ('DKK', 'Danish Krone', 'dKr', 6.54, false, true),
 ('EUR', 'Euro', '€', 0.92, false, true)
ON CONFLICT (code) DO NOTHING;

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id INTEGER NOT NULL REFERENCES categories(id),
    images JSONB NOT NULL DEFAULT '[]',
    has_variants BOOLEAN NOT NULL DEFAULT false,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create product_variants table
CREATE TABLE IF NOT EXISTS product_variants (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(100) NOT NULL UNIQUE,
    stock INTEGER NOT NULL DEFAULT 0,
    attributes JSONB NOT NULL,
    images JSONB NOT NULL DEFAULT '[]',
    is_default BOOLEAN NOT NULL DEFAULT false,
    weight DECIMAL(10, 3) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create product_prices table for multi-currency support
CREATE TABLE IF NOT EXISTS product_prices (
    id SERIAL PRIMARY KEY,
    variant_id INTEGER NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    currency_code VARCHAR(3) NOT NULL REFERENCES currencies(code) ON DELETE CASCADE,
    price BIGINT NOT NULL, -- stored in cents/smallest currency unit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(variant_id, currency_code)
);

-- Create payment_providers table
CREATE TABLE IF NOT EXISTS payment_providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    config JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create discounts table
CREATE TABLE IF NOT EXISTS discounts (
    id SERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL, -- 'product', 'basket'
    method VARCHAR(20) NOT NULL, -- 'percentage', 'fixed'
    value DECIMAL(10, 2) NOT NULL,
    min_order_value BIGINT NOT NULL DEFAULT 0, -- stored in cents
    max_discount_value BIGINT NOT NULL DEFAULT 0, -- stored in cents
    product_ids JSONB NOT NULL DEFAULT '[]',
    category_ids JSONB NOT NULL DEFAULT '[]',
    start_date TIMESTAMP NOT NULL DEFAULT NOW(),
    end_date TIMESTAMP NOT NULL DEFAULT NOW(),
    usage_limit INTEGER NOT NULL DEFAULT 0,
    current_usage INTEGER NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create webhooks table
CREATE TABLE IF NOT EXISTS webhooks (
    id SERIAL PRIMARY KEY,
    event VARCHAR(100) NOT NULL,
    url VARCHAR(255) NOT NULL,
    secret VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create shipping_zones table
CREATE TABLE IF NOT EXISTS shipping_zones (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    countries JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create shipping_methods table
CREATE TABLE IF NOT EXISTS shipping_methods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    zone_id INTEGER NOT NULL REFERENCES shipping_zones(id) ON DELETE CASCADE,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(name, zone_id)
);

-- Create shipping_rates table
CREATE TABLE IF NOT EXISTS shipping_rates (
    id SERIAL PRIMARY KEY,
    method_id INTEGER NOT NULL REFERENCES shipping_methods(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL, -- 'flat', 'weight_based', 'value_based'
    base_rate BIGINT NOT NULL DEFAULT 0, -- stored in cents
    min_order_value BIGINT NOT NULL DEFAULT 0, -- stored in cents
    free_shipping_threshold BIGINT NOT NULL DEFAULT 0, -- stored in cents
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create weight_based_rates table
CREATE TABLE IF NOT EXISTS weight_based_rates (
    id SERIAL PRIMARY KEY,
    rate_id INTEGER NOT NULL REFERENCES shipping_rates(id) ON DELETE CASCADE,
    min_weight DECIMAL(10, 3) NOT NULL,
    max_weight DECIMAL(10, 3) NOT NULL,
    rate BIGINT NOT NULL, -- stored in cents
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create value_based_rates table
CREATE TABLE IF NOT EXISTS value_based_rates (
    id SERIAL PRIMARY KEY,
    rate_id INTEGER NOT NULL REFERENCES shipping_rates(id) ON DELETE CASCADE,
    min_order_value BIGINT NOT NULL, -- stored in cents
    max_order_value BIGINT NOT NULL, -- stored in cents
    rate BIGINT NOT NULL, -- stored in cents
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create checkouts table (replaces deprecated carts)
CREATE TABLE IF NOT EXISTS checkouts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    checkout_session_id VARCHAR(255) UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    shipping_address JSONB NOT NULL DEFAULT '{}',
    billing_address JSONB NOT NULL DEFAULT '{}',
    customer_details JSONB NOT NULL DEFAULT '{}',
    shipping_method_id INTEGER REFERENCES shipping_methods(id) ON DELETE SET NULL,
    payment_provider VARCHAR(255),
    currency VARCHAR(3) NOT NULL REFERENCES currencies(code),
    discount_code VARCHAR(100),
    discount_amount BIGINT NOT NULL DEFAULT 0, -- stored in cents
    shipping_cost BIGINT NOT NULL DEFAULT 0, -- stored in cents
    subtotal BIGINT NOT NULL DEFAULT 0, -- stored in cents
    final_amount BIGINT NOT NULL DEFAULT 0, -- stored in cents
    total_weight DECIMAL(10, 3) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_activity_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    converted_order_id INTEGER REFERENCES orders(id) ON DELETE SET NULL
);

-- Create checkout_items table
CREATE TABLE IF NOT EXISTS checkout_items (
    id SERIAL PRIMARY KEY,
    checkout_id INTEGER NOT NULL REFERENCES checkouts(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id),
    product_variant_id INTEGER REFERENCES product_variants(id) ON DELETE SET NULL,
    quantity INTEGER NOT NULL,
    price BIGINT NOT NULL, -- stored in cents
    weight DECIMAL(10, 3) NOT NULL DEFAULT 0,
    product_name VARCHAR(255) NOT NULL,
    variant_name VARCHAR(255),
    sku VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE,
    user_id INTEGER REFERENCES users(id),
    checkout_session_id VARCHAR(255),
    subtotal BIGINT NOT NULL, -- stored in cents
    shipping_cost BIGINT NOT NULL DEFAULT 0, -- stored in cents
    final_amount BIGINT NOT NULL, -- stored in cents
    discount_amount BIGINT NOT NULL DEFAULT 0, -- stored in cents
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    payment_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    shipping_address JSONB NOT NULL,
    billing_address JSONB NOT NULL,
    customer_details JSONB NOT NULL DEFAULT '{}',
    currency VARCHAR(3) NOT NULL REFERENCES currencies(code),
    payment_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP
);

-- Create order_items table
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id),
    product_variant_id INTEGER REFERENCES product_variants(id) ON DELETE SET NULL,
    quantity INTEGER NOT NULL,
    price BIGINT NOT NULL, -- stored in cents
    subtotal BIGINT NOT NULL, -- stored in cents
    product_name VARCHAR(255) NOT NULL,
    variant_name VARCHAR(255),
    sku VARCHAR(100),
    created_at TIMESTAMP NOT NULL
);

-- Create payment_transactions table
CREATE TABLE IF NOT EXISTS payment_transactions (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id),
    payment_provider VARCHAR(100) NOT NULL,
    transaction_id VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL, -- stored in cents
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(50) NOT NULL,
    provider_response JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_active ON products(active);
CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
CREATE INDEX idx_product_variants_sku ON product_variants(sku);
CREATE INDEX idx_product_variant_prices_variant_id ON product_prices(variant_id);
CREATE INDEX idx_product_variant_prices_currency ON product_prices(currency_code);
CREATE INDEX idx_discounts_code ON discounts(code);
CREATE INDEX idx_discounts_active ON discounts(active);
CREATE INDEX idx_checkouts_user_id ON checkouts(user_id);
CREATE INDEX idx_checkouts_session_id ON checkouts(session_id);
CREATE INDEX idx_checkouts_status ON checkouts(status);
CREATE INDEX idx_checkouts_expires_at ON checkouts(expires_at);
CREATE INDEX idx_checkouts_converted_order_id ON checkouts(converted_order_id);
CREATE INDEX idx_checkout_items_checkout_id ON checkout_items(checkout_id);
CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_payment_status ON orders(payment_status);
CREATE INDEX idx_orders_checkout_session ON orders(checkout_session_id);
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_payment_transactions_order ON payment_transactions(order_id);
CREATE INDEX idx_payment_transactions_provider ON payment_transactions(payment_provider);

-- Triggers to automatically update timestamps
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_timestamp BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_categories_timestamp BEFORE UPDATE ON categories FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_currencies_timestamp BEFORE UPDATE ON currencies FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_products_timestamp BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_product_variants_timestamp BEFORE UPDATE ON product_variants FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_product_prices_timestamp BEFORE UPDATE ON product_prices FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_product_variant_prices_timestamp BEFORE UPDATE ON product_prices FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_payment_providers_timestamp BEFORE UPDATE ON payment_providers FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_discounts_timestamp BEFORE UPDATE ON discounts FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_webhooks_timestamp BEFORE UPDATE ON webhooks FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_shipping_zones_timestamp BEFORE UPDATE ON shipping_zones FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_shipping_methods_timestamp BEFORE UPDATE ON shipping_methods FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_shipping_rates_timestamp BEFORE UPDATE ON shipping_rates FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_weight_based_rates_timestamp BEFORE UPDATE ON weight_based_rates FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_value_based_rates_timestamp BEFORE UPDATE ON value_based_rates FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_checkouts_timestamp BEFORE UPDATE ON checkouts FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_checkout_items_timestamp BEFORE UPDATE ON checkout_items FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_orders_timestamp BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER update_payment_transactions_timestamp BEFORE UPDATE ON payment_transactions FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- Function to generate friendly IDs
CREATE OR REPLACE FUNCTION generate_friendly_id()
RETURNS TRIGGER AS $$
DECLARE
    counter INTEGER := 1;
    friendly_id VARCHAR(50);
BEGIN
    -- Generate friendly ID based on table name and sequence
    IF TG_TABLE_NAME = 'orders' THEN
        LOOP
            friendly_id := 'ord-' || LPAD(counter::text, 6, '0');
            EXIT WHEN NOT EXISTS (SELECT 1 FROM orders WHERE order_number = friendly_id);
            counter := counter + 1;
        END LOOP;
    END IF;
    
    NEW.friendly_id := friendly_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for friendly ID generation
CREATE TRIGGER generate_product_friendly_id BEFORE INSERT ON products FOR EACH ROW WHEN (NEW.friendly_id IS NULL) EXECUTE FUNCTION generate_friendly_id();
CREATE TRIGGER generate_order_friendly_id BEFORE INSERT ON orders FOR EACH ROW WHEN (NEW.friendly_id IS NULL) EXECUTE FUNCTION generate_friendly_id();

-- Constraints to ensure data integrity
ALTER TABLE products ADD CONSTRAINT chk_products_price_positive CHECK (price >= 0);
ALTER TABLE product_variants ADD CONSTRAINT chk_variants_price_positive CHECK (price >= 0);
ALTER TABLE checkouts ADD CONSTRAINT chk_checkouts_amounts_positive CHECK (subtotal >= 0 AND shipping_cost >= 0 AND discount_amount >= 0 AND final_amount >= 0);
ALTER TABLE orders ADD CONSTRAINT chk_orders_amounts_positive CHECK (subtotal >= 0 AND shipping_cost >= 0 AND discount_amount >= 0 AND final_amount >= 0);
ALTER TABLE order_items ADD CONSTRAINT chk_order_items_amounts_positive CHECK (price >= 0 AND subtotal >= 0 AND quantity > 0);
ALTER TABLE checkout_items ADD CONSTRAINT chk_checkout_items_amounts_positive CHECK (price >= 0 AND quantity > 0);

-- Comments for documentation
COMMENT ON TABLE users IS 'User accounts for customers and administrators';
COMMENT ON TABLE categories IS 'Product categories with hierarchical support';
COMMENT ON TABLE currencies IS 'Supported currencies with exchange rates';
COMMENT ON TABLE products IS 'Core product catalog';
COMMENT ON TABLE product_variants IS 'Product variations (size, color, etc.)';
COMMENT ON TABLE product_prices IS 'Multi-currency pricing for product variants';
COMMENT ON TABLE checkouts IS 'Shopping cart replacement with expiration and guest support';
COMMENT ON TABLE checkout_items IS 'Items in checkout sessions';
COMMENT ON TABLE orders IS 'Completed purchase orders';
COMMENT ON TABLE order_items IS 'Items within completed orders';
COMMENT ON TABLE discounts IS 'Discount codes and promotional offers';
COMMENT ON TABLE payment_transactions IS 'Payment processing transaction logs';
COMMENT ON TABLE shipping_zones IS 'Geographic shipping zones';
COMMENT ON TABLE shipping_methods IS 'Available shipping methods per zone';
COMMENT ON TABLE shipping_rates IS 'Shipping cost calculation rules';
COMMENT ON TABLE webhooks IS 'External webhook integrations';

COMMENT ON COLUMN product_variants.price IS 'Price stored in cents to avoid floating point precision issues';
COMMENT ON COLUMN checkouts.session_id IS 'Session ID for guest checkouts, nullable for authenticated users';
COMMENT ON COLUMN checkout_items.product_variant_id IS 'Reference to product_variants.id. NULL indicates this is a regular product without variants';
COMMENT ON COLUMN orders.checkout_session_id IS 'Links order back to the checkout session that created it';
COMMENT ON COLUMN order_items.product_variant_id IS 'Reference to product_variants.id. NULL indicates this is a regular product without variants';
