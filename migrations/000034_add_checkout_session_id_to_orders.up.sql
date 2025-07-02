-- Add checkout_session_id column to orders table
ALTER TABLE orders ADD COLUMN checkout_session_id VARCHAR(255);

-- Create index on checkout_session_id for faster lookups
CREATE INDEX idx_orders_checkout_session ON orders(checkout_session_id);

-- Populate existing orders with checkout session IDs from the checkouts table
-- This links orders to their corresponding checkout sessions
UPDATE orders 
SET checkout_session_id = c.session_id
FROM checkouts c
WHERE orders.id = c.converted_order_id 
  AND c.session_id IS NOT NULL 
  AND c.session_id != '';
