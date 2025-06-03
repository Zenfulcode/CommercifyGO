-- Remove default currency columns from products and product_variants tables

-- Drop indexes
DROP INDEX IF EXISTS idx_products_currency_code;
DROP INDEX IF EXISTS idx_product_variants_currency_code;

-- Remove currency_code columns
ALTER TABLE products DROP COLUMN IF EXISTS currency_code;
ALTER TABLE product_variants DROP COLUMN IF EXISTS currency_code;