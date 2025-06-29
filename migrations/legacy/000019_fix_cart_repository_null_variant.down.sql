-- No action needed for rollback as we're just documenting the change
-- to the cart repository code
COMMENT ON COLUMN cart_items.product_variant_id IS NULL;