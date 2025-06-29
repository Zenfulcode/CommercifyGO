-- Add product_variant_id column to cart_items table
ALTER TABLE cart_items ADD COLUMN IF NOT EXISTS product_variant_id INTEGER;

-- Add foreign key constraint
ALTER TABLE cart_items 
ADD CONSTRAINT fk_cart_items_product_variant 
FOREIGN KEY (product_variant_id) 
REFERENCES product_variants(id) 
ON DELETE SET NULL;

-- Add index for faster lookups
CREATE INDEX IF NOT EXISTS idx_cart_items_product_variant_id ON cart_items(product_variant_id);
