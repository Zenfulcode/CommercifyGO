-- Rollback migration for ensuring products have variants

-- Drop the trigger and function
DROP TRIGGER IF EXISTS prevent_last_variant_deletion ON product_variants;
DROP FUNCTION IF EXISTS check_product_has_variants();

-- Remove comments
COMMENT ON TABLE products IS NULL;
COMMENT ON TABLE product_variants IS NULL;

-- Note: We don't automatically delete the created default variants
-- as they may have been modified by users. Manual cleanup may be required.