-- Add payment_provider column to orders table if it doesn't exist
ALTER TABLE orders ADD COLUMN IF NOT EXISTS payment_provider VARCHAR(50);
