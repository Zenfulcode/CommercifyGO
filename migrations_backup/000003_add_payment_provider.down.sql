-- Remove payment_provider column from orders table
ALTER TABLE orders DROP COLUMN IF EXISTS payment_provider;
