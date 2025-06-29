-- Add payment_status column to orders table
ALTER TABLE orders ADD COLUMN IF NOT EXISTS payment_status VARCHAR(20) NOT NULL DEFAULT 'pending';

-- Create index for payment_status
CREATE INDEX IF NOT EXISTS idx_orders_payment_status ON orders(payment_status);

-- Update existing orders to have proper payment_status based on their current status
-- Orders with status 'paid', 'shipped', 'completed' should have payment_status 'captured'
-- Orders with status 'cancelled' should have payment_status 'cancelled'
-- Orders with status 'pending' should have payment_status 'pending'
UPDATE orders 
SET payment_status = 
    CASE 
        WHEN status IN ('paid', 'shipped', 'completed') THEN 'captured'
        WHEN status = 'cancelled' THEN 'cancelled'
        ELSE 'pending'
    END;
