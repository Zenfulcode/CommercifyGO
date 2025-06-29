-- Drop indexes
DROP INDEX IF EXISTS idx_shipping_rates_method_id;
DROP INDEX IF EXISTS idx_shipping_rates_zone_id;
DROP INDEX IF EXISTS idx_weight_based_rates_shipping_rate_id;
DROP INDEX IF EXISTS idx_value_based_rates_shipping_rate_id;

-- Remove columns from products and orders tables
ALTER TABLE products DROP COLUMN IF EXISTS weight;
ALTER TABLE orders DROP COLUMN IF EXISTS shipping_method_id;
ALTER TABLE orders DROP COLUMN IF EXISTS shipping_cost;
ALTER TABLE orders DROP COLUMN IF EXISTS total_weight;

-- Drop tables in reverse order of creation to avoid foreign key constraint issues
DROP TABLE IF EXISTS value_based_rates;
DROP TABLE IF EXISTS weight_based_rates;
DROP TABLE IF EXISTS shipping_rates;
DROP TABLE IF EXISTS shipping_zones;
DROP TABLE IF EXISTS shipping_methods;