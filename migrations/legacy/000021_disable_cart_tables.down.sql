-- Rollback migration to re-enable cart tables
-- This will remove the triggers and restore the cart functionality

-- Drop the triggers preventing cart operations
DROP TRIGGER IF EXISTS prevent_cart_insert ON carts;
DROP TRIGGER IF EXISTS prevent_cart_update ON carts;
DROP TRIGGER IF EXISTS prevent_cart_items_insert ON cart_items;
DROP TRIGGER IF EXISTS prevent_cart_items_update ON cart_items;

-- Drop the trigger function
DROP FUNCTION IF EXISTS prevent_cart_operations();

-- Remove comments on tables
COMMENT ON TABLE carts IS '';
COMMENT ON TABLE cart_items IS '';

-- Drop the legacy views
DROP VIEW IF EXISTS legacy_carts;
DROP VIEW IF EXISTS legacy_cart_items;

-- Drop the archive tables if they are no longer needed
-- Note: You might want to keep these for historical data
-- DROP TABLE IF EXISTS cart_archive;
-- DROP TABLE IF EXISTS cart_items_archive;