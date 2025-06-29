-- Remove checkout_session_id column from orders table
DROP INDEX IF EXISTS idx_orders_checkout_session;
ALTER TABLE orders DROP COLUMN IF EXISTS checkout_session_id;
