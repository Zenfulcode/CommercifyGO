-- Improve product deletion by removing the problematic trigger and relying on application logic
-- The application layer will handle the business rule of ensuring products have at least one variant

-- Drop the trigger that was causing complexity with product deletion
DROP TRIGGER IF EXISTS prevent_last_variant_deletion ON product_variants;

-- Drop the trigger function as well
DROP FUNCTION IF EXISTS check_product_has_variants();

-- Add a comment explaining the new approach
COMMENT ON TABLE product_variants IS 'Product variants. Business rule "products must have at least one variant" is enforced at the application layer.';
