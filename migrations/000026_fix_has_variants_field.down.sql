-- Rollback migration for fixing has_variants field
-- This migration doesn't have a specific rollback since it fixes data integrity
-- The previous state was incorrect, so rolling back would restore incorrect data

-- If needed, you could reset all products to has_variants = true (previous default behavior)
-- UPDATE products SET has_variants = true WHERE id IN (SELECT DISTINCT product_id FROM product_variants);