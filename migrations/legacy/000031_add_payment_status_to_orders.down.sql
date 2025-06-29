-- Remove payment_status column from orders table
DROP INDEX IF EXISTS idx_orders_payment_status;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_status;
