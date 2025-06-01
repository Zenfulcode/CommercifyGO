-- Migration to disable cart tables since they've been replaced by the checkout system
-- This migration preserves existing data but prevents new operations on cart tables

-- Create a temporary table to archive existing cart data for reference
CREATE TABLE cart_archive AS 
SELECT * FROM carts;

CREATE TABLE cart_items_archive AS
SELECT * FROM cart_items;

-- Create triggers to prevent inserts/updates to cart tables
CREATE OR REPLACE FUNCTION prevent_cart_operations()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Cart operations are disabled. Please use the checkout system instead.';
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create triggers on carts table
CREATE TRIGGER prevent_cart_insert
BEFORE INSERT ON carts
FOR EACH ROW EXECUTE FUNCTION prevent_cart_operations();

CREATE TRIGGER prevent_cart_update
BEFORE UPDATE ON carts
FOR EACH ROW EXECUTE FUNCTION prevent_cart_operations();

-- Create triggers on cart_items table
CREATE TRIGGER prevent_cart_items_insert
BEFORE INSERT ON cart_items
FOR EACH ROW EXECUTE FUNCTION prevent_cart_operations();

CREATE TRIGGER prevent_cart_items_update
BEFORE UPDATE ON cart_items
FOR EACH ROW EXECUTE FUNCTION prevent_cart_operations();

-- Comment the tables to indicate they're deprecated
COMMENT ON TABLE carts IS 'DEPRECATED: This table has been replaced by the checkout system. Use checkouts instead.';
COMMENT ON TABLE cart_items IS 'DEPRECATED: This table has been replaced by the checkout system. Use checkout_items instead.';

-- Create a view to make cart data accessible through the checkout system if needed
CREATE VIEW legacy_carts AS
SELECT 
    c.id, 
    c.user_id,
    c.session_id,
    c.created_at,
    c.updated_at
FROM carts c;

CREATE VIEW legacy_cart_items AS
SELECT 
    ci.id,
    ci.cart_id,
    ci.product_id,
    ci.product_variant_id,
    ci.quantity,
    ci.created_at,
    ci.updated_at
FROM cart_items ci;

-- Add indexes on the archive tables to maintain query performance if needed
CREATE INDEX idx_cart_archive_user_id ON cart_archive(user_id);
CREATE INDEX idx_cart_items_archive_cart_id ON cart_items_archive(cart_id);