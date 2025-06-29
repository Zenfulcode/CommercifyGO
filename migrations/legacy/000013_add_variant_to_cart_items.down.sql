-- Remove foreign key constraint
ALTER TABLE cart_items DROP CONSTRAINT IF EXISTS fk_cart_items_product_variant;

-- Remove index
DROP INDEX IF EXISTS idx_cart_items_product_variant_id;

-- Remove product_variant_id column
ALTER TABLE cart_items DROP COLUMN IF EXISTS product_variant_id;
