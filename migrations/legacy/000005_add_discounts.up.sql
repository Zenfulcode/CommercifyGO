-- Create discounts table
CREATE TABLE IF NOT EXISTS discounts (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL, -- 'basket' or 'product'
    method VARCHAR(20) NOT NULL, -- 'fixed' or 'percentage'
    value DECIMAL(10, 2) NOT NULL,
    min_order_value DECIMAL(10, 2) NOT NULL DEFAULT 0,
    max_discount_value DECIMAL(10, 2) NOT NULL DEFAULT 0,
    product_ids JSONB NOT NULL DEFAULT '[]',
    category_ids JSONB NOT NULL DEFAULT '[]',
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    usage_limit INTEGER NOT NULL DEFAULT 0,
    current_usage INTEGER NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Add discount-related columns to orders table
ALTER TABLE orders
ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(10, 2) NOT NULL DEFAULT 0;

ALTER TABLE orders
ADD COLUMN IF NOT EXISTS final_amount DECIMAL(10, 2);

ALTER TABLE orders
ADD COLUMN IF NOT EXISTS discount_id INTEGER REFERENCES discounts (id);

ALTER TABLE orders
ADD COLUMN IF NOT EXISTS discount_code VARCHAR(50);

-- Update existing orders to set final_amount equal to total_amount
UPDATE orders
SET
    final_amount = orders.total_amount
WHERE
    final_amount IS NULL;

-- Create indexes
CREATE INDEX idx_discounts_code ON discounts (code);

CREATE INDEX idx_discounts_active ON discounts (active);

CREATE INDEX idx_discounts_dates ON discounts (start_date, end_date);