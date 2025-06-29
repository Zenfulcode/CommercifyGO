-- Restore the original constraint configuration
-- First, drop the unique index if it exists
DROP INDEX IF EXISTS cart_items_unique_product_variant;

-- Restore the original unique constraint
ALTER TABLE cart_items ADD CONSTRAINT cart_items_cart_id_product_id_key 
UNIQUE (cart_id, product_id);

-- Drop and recreate the original foreign key constraint
ALTER TABLE cart_items DROP CONSTRAINT IF EXISTS fk_cart_items_product_variant;
ALTER TABLE cart_items 
ADD CONSTRAINT fk_cart_items_product_variant 
FOREIGN KEY (product_variant_id) 
REFERENCES product_variants(id) 
ON DELETE SET NULL;