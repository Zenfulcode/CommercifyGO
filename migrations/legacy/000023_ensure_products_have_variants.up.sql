-- Migration to ensure all products have at least one variant
-- This enforces that variants are mandatory for all products

-- First, create default variants for products that don't have any variants
INSERT INTO product_variants (
    product_id,
    sku,
    price,
    currency_code,
    stock,
    attributes,
    images,
    is_default,
    created_at,
    updated_at
)
SELECT 
    p.id,
    p.product_number,  -- Use existing product_number as SKU
    p.price,
    p.currency_code,
    p.stock,
    '[{"name": "Default", "value": "Standard"}]'::jsonb,  -- Default attribute
    p.images,
    true,  -- Mark as default variant
    NOW(),
    NOW()
FROM products p
WHERE p.has_variants = false
   OR p.id NOT IN (SELECT DISTINCT product_id FROM product_variants);

-- Update products to mark them as having variants
UPDATE products 
SET has_variants = true
WHERE has_variants = false 
   OR id NOT IN (SELECT DISTINCT product_id FROM product_variants);

-- Add a constraint to ensure all products must have at least one variant
-- This will be enforced by the application layer, but we add a check here
CREATE OR REPLACE FUNCTION check_product_has_variants()
RETURNS trigger AS $$
BEGIN
    -- For INSERT operations on products, we allow it but warn that variants must be added
    IF TG_OP = 'INSERT' THEN
        RETURN NEW;
    END IF;
    
    -- For DELETE operations on product_variants, ensure at least one variant remains
    IF TG_OP = 'DELETE' THEN
        IF (SELECT COUNT(*) FROM product_variants WHERE product_id = OLD.product_id) <= 1 THEN
            RAISE EXCEPTION 'Cannot delete the last variant of a product. Products must have at least one variant.';
        END IF;
        RETURN OLD;
    END IF;
    
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- Create trigger to prevent deletion of the last variant
CREATE TRIGGER prevent_last_variant_deletion
    BEFORE DELETE ON product_variants
    FOR EACH ROW
    EXECUTE FUNCTION check_product_has_variants();

-- Add comment to document the new requirement
COMMENT ON TABLE products IS 'All products must have at least one variant. The has_variants field should always be true.';
COMMENT ON TABLE product_variants IS 'Product variants are mandatory. Every product must have at least one variant.';