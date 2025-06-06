-- Remove unique constraint from shipping_rates table

ALTER TABLE shipping_rates 
DROP CONSTRAINT IF EXISTS unique_shipping_method_zone;