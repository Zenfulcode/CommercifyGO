-- Add currency column to orders table to support multi-currency orders
ALTER TABLE orders ADD COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'USD' REFERENCES currencies(code);

-- Create index for currency lookups
CREATE INDEX idx_orders_currency ON orders(currency);

-- Update existing orders to use default currency
-- This ensures all existing orders have a valid currency
UPDATE orders SET currency = 'USD' WHERE currency IS NULL OR currency = '';
