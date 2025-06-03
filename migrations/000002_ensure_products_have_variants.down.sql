-- Remove product variant constraints and triggers

DROP TRIGGER IF EXISTS prevent_last_variant_deletion ON product_variants;
DROP FUNCTION IF EXISTS check_product_has_variants();

-- Remove comments
COMMENT ON TABLE products IS NULL;
COMMENT ON TABLE product_variants IS NULL;
