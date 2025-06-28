-- Add missing fields to order_items table
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS weight DECIMAL(10, 3) DEFAULT 0;
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS product_name VARCHAR(255) DEFAULT '';
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS sku VARCHAR(100) DEFAULT '';

-- Create indexes for the new columns
CREATE INDEX IF NOT EXISTS idx_order_items_sku ON order_items(sku);
