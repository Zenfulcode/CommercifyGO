-- Add currency support to the database

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

-- Create product_prices table to store prices in different currencies
CREATE TABLE IF NOT EXISTS product_prices (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    currency_code VARCHAR(3) NOT NULL REFERENCES currencies(code) ON DELETE CASCADE,
    price BIGINT NOT NULL, -- stored in cents/smallest currency unit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(product_id, currency_code)
);

-- Create product_variant_prices table to store variant prices in different currencies
CREATE TABLE IF NOT EXISTS product_variant_prices (
    id SERIAL PRIMARY KEY,
    variant_id INT NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    currency_code VARCHAR(3) NOT NULL REFERENCES currencies(code) ON DELETE CASCADE,
    price BIGINT NOT NULL, -- stored in cents/smallest currency unit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(variant_id, currency_code)
);

-- Add default currency column to payment_transactions
ALTER TABLE payment_transactions ALTER COLUMN currency SET DEFAULT 'USD';

-- Create indexes for better query performance
CREATE INDEX idx_product_prices_product_id ON product_prices(product_id);
CREATE INDEX idx_product_prices_currency_code ON product_prices(currency_code);
CREATE INDEX idx_product_variant_prices_variant_id ON product_variant_prices(variant_id);
CREATE INDEX idx_product_variant_prices_currency_code ON product_variant_prices(currency_code);

-- Insert default currencies
INSERT INTO currencies (code, name, symbol, exchange_rate, is_default, is_enabled)
VALUES 
('USD', 'US Dollar', '$', 1.0, true, true),
('EUR', 'Euro', '€', 0.85, false, true),
('DKK', 'Danish Krone', 'kr', 0.15, false, true),
('GBP', 'British Pound', '£', 0.75, false, true),
('JPY', 'Japanese Yen', '¥', 110.0, false, true),
('CAD', 'Canadian Dollar', 'CA$', 1.25, false, true)
ON CONFLICT (code) DO NOTHING;