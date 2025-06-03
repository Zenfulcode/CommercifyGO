-- Add default currency columns to products and product_variants tables

-- Add currency_code column to products table (initially nullable)
ALTER TABLE products
ADD COLUMN currency_code VARCHAR(3) REFERENCES currencies(code);

-- Update existing products with default currency
UPDATE products
SET currency_code = (SELECT code FROM currencies WHERE is_default = true LIMIT 1);

-- Make currency_code NOT NULL for products
ALTER TABLE products
ALTER COLUMN currency_code SET NOT NULL;

-- Add currency_code column to product_variants table (initially nullable)
ALTER TABLE product_variants
ADD COLUMN currency_code VARCHAR(3) REFERENCES currencies(code);

-- Update existing variants with their product's currency
UPDATE product_variants pv
SET currency_code = p.currency_code
FROM products p
WHERE pv.product_id = p.id;

-- Make currency_code NOT NULL for product_variants
ALTER TABLE product_variants
ALTER COLUMN currency_code SET NOT NULL;

-- Create indexes for better query performance
CREATE INDEX idx_products_currency_code ON products(currency_code);
CREATE INDEX idx_product_variants_currency_code ON product_variants(currency_code);