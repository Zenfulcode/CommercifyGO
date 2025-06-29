-- Remove discount-related columns from orders table
ALTER TABLE orders DROP COLUMN IF EXISTS discount_amount;
ALTER TABLE orders DROP COLUMN IF EXISTS final_amount;
ALTER TABLE orders DROP COLUMN IF EXISTS discount_id;
ALTER TABLE orders DROP COLUMN IF EXISTS discount_code;

-- Drop discounts table
DROP TABLE IF EXISTS discounts;
