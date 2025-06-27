-- Remove product_variant_id column from order_items table
DROP INDEX IF EXISTS idx_order_items_variant;
ALTER TABLE order_items DROP COLUMN IF EXISTS product_variant_id;
