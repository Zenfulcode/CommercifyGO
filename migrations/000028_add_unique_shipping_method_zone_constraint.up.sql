-- Add unique constraint to prevent duplicate shipping rates for the same method and zone combination
-- This ensures each shipping method can only have one rate per shipping zone

-- First, remove any duplicate entries, keeping only the latest one (highest ID)
WITH duplicates AS (
    SELECT id,
           ROW_NUMBER() OVER (
               PARTITION BY shipping_method_id, shipping_zone_id 
               ORDER BY id DESC
           ) as rn
    FROM shipping_rates
)
DELETE FROM shipping_rates 
WHERE id IN (
    SELECT id FROM duplicates WHERE rn > 1
);

-- Now add the unique constraint
ALTER TABLE shipping_rates 
ADD CONSTRAINT unique_shipping_method_zone 
UNIQUE (shipping_method_id, shipping_zone_id);