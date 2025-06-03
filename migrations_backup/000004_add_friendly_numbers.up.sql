-- Add order_number column to orders table
ALTER TABLE orders ADD COLUMN IF NOT EXISTS order_number VARCHAR(50) UNIQUE;

-- Add product_number column to products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS product_number VARCHAR(50) UNIQUE;

-- Update existing orders with order numbers
UPDATE orders SET order_number = 'ORD-' || to_char(created_at, 'YYYYMMDD') || '-' || LPAD(id::text, 6, '0') WHERE order_number IS NULL;

-- Update existing products with product numbers
UPDATE products SET product_number = 'PROD-' || LPAD(id::text, 6, '0') WHERE product_number IS NULL;
