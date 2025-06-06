-- Migration to fix has_variants field values based on actual variant count
-- This updates existing products to have correct has_variants values

-- Update has_variants to true for products that have more than one variant
UPDATE products 
SET has_variants = true, updated_at = NOW()
WHERE id IN (
    SELECT p.id 
    FROM products p
    JOIN (
        SELECT product_id, COUNT(*) as variant_count
        FROM product_variants
        GROUP BY product_id
        HAVING COUNT(*) > 1
    ) v ON p.id = v.product_id
);

-- Update has_variants to false for products that have only one variant
UPDATE products 
SET has_variants = false, updated_at = NOW()
WHERE id IN (
    SELECT p.id 
    FROM products p
    JOIN (
        SELECT product_id, COUNT(*) as variant_count
        FROM product_variants
        GROUP BY product_id
        HAVING COUNT(*) = 1
    ) v ON p.id = v.product_id
);

-- Ensure products without any variants are set to false (should not happen with current system)
UPDATE products 
SET has_variants = false, updated_at = NOW()
WHERE id NOT IN (
    SELECT DISTINCT product_id 
    FROM product_variants
);