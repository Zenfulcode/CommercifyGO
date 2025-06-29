-- Remove action_url column from orders table
ALTER TABLE orders DROP COLUMN IF EXISTS action_url;