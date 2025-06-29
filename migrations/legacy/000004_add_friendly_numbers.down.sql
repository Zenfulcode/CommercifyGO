-- Remove order_number column from orders table
ALTER TABLE orders DROP COLUMN IF EXISTS order_number;

-- Remove product_number column from products table
ALTER TABLE products DROP COLUMN IF EXISTS product_number;
