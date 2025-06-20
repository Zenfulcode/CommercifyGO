-- Remove currency column from orders table
DROP INDEX IF EXISTS idx_orders_currency;
ALTER TABLE orders DROP COLUMN IF EXISTS currency;
