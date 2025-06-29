-- Remove currency support from the database

-- Drop indexes
DROP INDEX IF EXISTS idx_product_variant_prices_currency_code;
DROP INDEX IF EXISTS idx_product_variant_prices_variant_id;
DROP INDEX IF EXISTS idx_product_prices_currency_code;
DROP INDEX IF EXISTS idx_product_prices_product_id;

-- Drop tables
DROP TABLE IF EXISTS product_variant_prices;
DROP TABLE IF EXISTS product_prices;

-- Remove default constraint on payment_transactions.currency
ALTER TABLE payment_transactions ALTER COLUMN currency DROP DEFAULT;

-- Drop currencies table
DROP TABLE IF EXISTS currencies;