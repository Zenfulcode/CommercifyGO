-- First drop the indexes
DROP INDEX IF EXISTS idx_orders_customer_email;
DROP INDEX IF EXISTS idx_orders_customer_phone;
DROP INDEX IF EXISTS idx_orders_customer_full_name;

-- Add back the guest columns
ALTER TABLE orders
    ADD COLUMN guest_email VARCHAR(255),
    ADD COLUMN guest_phone VARCHAR(50),
    ADD COLUMN guest_full_name VARCHAR(255);

-- Restore guest data from customer details for guest orders
UPDATE orders
SET 
    guest_email = customer_email,
    guest_phone = customer_phone,
    guest_full_name = customer_full_name
WHERE is_guest_order = true;

-- Drop the new customer detail columns
ALTER TABLE orders
    DROP COLUMN IF EXISTS customer_email,
    DROP COLUMN IF EXISTS customer_phone,
    DROP COLUMN IF EXISTS customer_full_name; 