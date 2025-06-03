-- +migrate Up
-- First add the new columns
ALTER TABLE orders
    ADD COLUMN customer_email VARCHAR(255),
    ADD COLUMN customer_phone VARCHAR(50),
    ADD COLUMN customer_full_name VARCHAR(255);

-- Update customer details from guest credentials for guest orders
UPDATE orders
SET 
    customer_email = guest_email,
    customer_phone = guest_phone,
    customer_full_name = guest_full_name
WHERE is_guest_order = true;

-- Update customer details from user table for non-guest orders
UPDATE orders o
SET 
    customer_email = u.email,
    customer_full_name = CONCAT(u.first_name, ' ', u.last_name)
FROM users u
WHERE o.user_id = u.id 
    AND o.is_guest_order = false
    AND o.user_id IS NOT NULL;

-- Drop the old guest columns after data migration
ALTER TABLE orders
    DROP COLUMN IF EXISTS guest_email,
    DROP COLUMN IF EXISTS guest_phone,
    DROP COLUMN IF EXISTS guest_full_name;

-- Add indexes for customer details
CREATE INDEX IF NOT EXISTS idx_orders_customer_email ON orders (customer_email);
CREATE INDEX IF NOT EXISTS idx_orders_customer_phone ON orders (customer_phone);
CREATE INDEX IF NOT EXISTS idx_orders_customer_full_name ON orders (customer_full_name);