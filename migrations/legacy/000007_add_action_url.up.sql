-- Add action_url column to orders table
ALTER TABLE orders ADD COLUMN IF NOT EXISTS action_url TEXT;