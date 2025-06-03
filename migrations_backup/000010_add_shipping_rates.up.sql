-- Create shipping methods table
CREATE TABLE IF NOT EXISTS shipping_methods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    estimated_delivery_days INT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create shipping rate rules table
CREATE TABLE IF NOT EXISTS shipping_zones (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    countries JSONB NOT NULL DEFAULT '[]',
    states JSONB NOT NULL DEFAULT '[]',
    zip_codes JSONB NOT NULL DEFAULT '[]',
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create shipping rates table to connect methods with rules
CREATE TABLE IF NOT EXISTS shipping_rates (
    id SERIAL PRIMARY KEY,
    shipping_method_id INT NOT NULL REFERENCES shipping_methods(id) ON DELETE CASCADE,
    shipping_zone_id INT NOT NULL REFERENCES shipping_zones(id) ON DELETE CASCADE,
    base_rate DECIMAL(10, 2) NOT NULL,
    min_order_value DECIMAL(10, 2) DEFAULT 0,
    free_shipping_threshold DECIMAL(10, 2) DEFAULT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create weight-based rates table
CREATE TABLE IF NOT EXISTS weight_based_rates (
    id SERIAL PRIMARY KEY,
    shipping_rate_id INT NOT NULL REFERENCES shipping_rates(id) ON DELETE CASCADE,
    min_weight DECIMAL(10, 2) NOT NULL DEFAULT 0,
    max_weight DECIMAL(10, 2) NOT NULL,
    rate DECIMAL(10, 2) NOT NULL
);

-- Create order value-based surcharges/discounts
CREATE TABLE IF NOT EXISTS value_based_rates (
    id SERIAL PRIMARY KEY,
    shipping_rate_id INT NOT NULL REFERENCES shipping_rates(id) ON DELETE CASCADE,
    min_order_value DECIMAL(10, 2) NOT NULL DEFAULT 0,
    max_order_value DECIMAL(10, 2) NOT NULL,
    rate DECIMAL(10, 2) NOT NULL
);

-- Add shipping_method_id, shipping_cost, weight to the orders table
ALTER TABLE orders ADD COLUMN IF NOT EXISTS shipping_method_id INT REFERENCES shipping_methods(id);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS shipping_cost DECIMAL(10, 2) DEFAULT 0;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS total_weight DECIMAL(10, 2) DEFAULT 0;

-- Add weight field to products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS weight DECIMAL(10, 2) DEFAULT 0;

-- Create indexes
CREATE INDEX idx_shipping_rates_method_id ON shipping_rates(shipping_method_id);
CREATE INDEX idx_shipping_rates_zone_id ON shipping_rates(shipping_zone_id);
CREATE INDEX idx_weight_based_rates_shipping_rate_id ON weight_based_rates(shipping_rate_id);
CREATE INDEX idx_value_based_rates_shipping_rate_id ON value_based_rates(shipping_rate_id);