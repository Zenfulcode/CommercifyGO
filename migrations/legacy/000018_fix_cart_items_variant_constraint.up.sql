-- First, drop the existing constraint
ALTER TABLE cart_items DROP CONSTRAINT IF EXISTS fk_cart_items_product_variant;

-- Re-add the constraint with the proper NULL handling
ALTER TABLE cart_items 
ADD CONSTRAINT fk_cart_items_product_variant 
FOREIGN KEY (product_variant_id) 
REFERENCES product_variants(id) 
ON DELETE SET NULL;

-- Make sure the unique constraint on cart_items allows for NULL variant_id
ALTER TABLE cart_items DROP CONSTRAINT IF EXISTS cart_items_cart_id_product_id_key;

-- Add a new unique constraint that allows NULL variant_id
-- This ensures each product or product variant combination is unique per cart
-- Using a partial index to handle NULL values properly
CREATE UNIQUE INDEX cart_items_unique_product_variant 
ON cart_items (cart_id, product_id, COALESCE(product_variant_id, 0));

-- Add a comment explaining the constraint
COMMENT ON CONSTRAINT fk_cart_items_product_variant ON cart_items IS 
'Foreign key constraint to product_variants table. Allows NULL for products without variants.';