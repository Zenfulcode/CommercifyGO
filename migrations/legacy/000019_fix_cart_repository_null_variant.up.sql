-- This migration documents the fix done to the cart repository code
-- to properly handle NULL values for product_variant_id
COMMENT ON COLUMN cart_items.product_variant_id IS 
'Reference to product_variants.id. NULL indicates this is a regular product without variants.';