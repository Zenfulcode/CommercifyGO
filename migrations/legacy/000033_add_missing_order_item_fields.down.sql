-- Remove missing fields from order_items table
DROP INDEX IF EXISTS idx_order_items_sku;
ALTER TABLE order_items DROP COLUMN IF EXISTS sku;
ALTER TABLE order_items DROP COLUMN IF EXISTS product_name;
ALTER TABLE order_items DROP COLUMN IF EXISTS weight;
