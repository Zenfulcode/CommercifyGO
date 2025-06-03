-- Create checkouts table
CREATE TABLE IF NOT EXISTS checkouts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    session_id VARCHAR(255) NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    shipping_address JSONB NOT NULL DEFAULT '{}',
    billing_address JSONB NOT NULL DEFAULT '{}',
    shipping_method_id INTEGER REFERENCES shipping_methods(id) ON DELETE SET NULL,
    payment_provider VARCHAR(255),
    total_amount BIGINT NOT NULL DEFAULT 0,
    shipping_cost BIGINT NOT NULL DEFAULT 0,
    total_weight DECIMAL(10, 3) NOT NULL DEFAULT 0,
    customer_details JSONB NOT NULL DEFAULT '{}',
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    discount_code VARCHAR(100),
    discount_amount BIGINT NOT NULL DEFAULT 0,
    final_amount BIGINT NOT NULL DEFAULT 0,
    applied_discount JSONB,
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
    price BIGINT NOT NULL,
    weight DECIMAL(10, 3) NOT NULL DEFAULT 0,
    product_name VARCHAR(255) NOT NULL,
    variant_name VARCHAR(255),
    sku VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for efficient lookups
CREATE INDEX IF NOT EXISTS idx_checkouts_user_id ON checkouts(user_id);
CREATE INDEX IF NOT EXISTS idx_checkouts_session_id ON checkouts(session_id);
CREATE INDEX IF NOT EXISTS idx_checkouts_status ON checkouts(status);
CREATE INDEX IF NOT EXISTS idx_checkouts_expires_at ON checkouts(expires_at);
CREATE INDEX IF NOT EXISTS idx_checkouts_converted_order_id ON checkouts(converted_order_id);
CREATE INDEX IF NOT EXISTS idx_checkout_items_checkout_id ON checkout_items(checkout_id);
CREATE INDEX IF NOT EXISTS idx_checkout_items_product_id ON checkout_items(product_id);
CREATE INDEX IF NOT EXISTS idx_checkout_items_product_variant_id ON checkout_items(product_variant_id);