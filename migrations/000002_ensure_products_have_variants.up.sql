-- Migration to ensure all products have at least one variant
-- This enforces that variants are mandatory for all products

-- Create default variants for products that don't have any variants
INSERT INTO product_variants (
    product_id,
    name,
    sku,
    price,
    stock,
    weight,
    attributes,
    created_at,
    updated_at
)
SELECT 
    p.id,
    p.name,
    CONCAT('SKU-', p.id), -- Generate SKU based on product ID
    p.price, -- Already stored as cents (BIGINT)
    p.stock,
    0.0, -- Default weight
    '{}'::jsonb, -- Empty attributes
    NOW(),
    NOW()
FROM products p
WHERE p.id NOT IN (SELECT DISTINCT product_id FROM product_variants);

-- Create function to ensure products always have at least one variant
CREATE OR REPLACE FUNCTION check_product_has_variants()
RETURNS trigger AS $$
BEGIN
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

-- Add comments to document the requirement
COMMENT ON
TABLE products IS 'All products must have at least one variant.';

COMMENT ON
TABLE product_variants IS 'Product variants are mandatory. Every product must have at least one variant.';