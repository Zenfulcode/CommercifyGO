-- Add active field to products table
ALTER TABLE products ADD COLUMN active BOOLEAN NOT NULL DEFAULT false;