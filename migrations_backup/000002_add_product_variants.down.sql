-- Drop product_variants table
DROP TABLE IF EXISTS product_variants;

-- Remove has_variants column from products table
ALTER TABLE products DROP COLUMN IF EXISTS has_variants;
