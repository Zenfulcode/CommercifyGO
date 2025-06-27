-- Add product_variant_id column to order_items table
ALTER TABLE order_items ADD COLUMN product_variant_id INTEGER REFERENCES product_variants(id);

-- Create index for the new column
CREATE INDEX idx_order_items_variant ON order_items(product_variant_id);
