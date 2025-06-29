-- Migration to change money fields from DECIMAL to BIGINT (int64)
-- This stores monetary values as cents instead of dollars to avoid floating point issues

-- Order table
ALTER TABLE orders
    ALTER COLUMN shipping_cost TYPE BIGINT USING (shipping_cost * 100)::BIGINT,
    ALTER COLUMN total_amount TYPE BIGINT USING (total_amount * 100)::BIGINT,
    ALTER COLUMN final_amount TYPE BIGINT USING (final_amount * 100)::BIGINT,
    ALTER COLUMN discount_amount TYPE BIGINT USING (discount_amount * 100)::BIGINT;

-- OrderItem table
ALTER TABLE order_items
    ALTER COLUMN price TYPE BIGINT USING (price * 100)::BIGINT,
    ALTER COLUMN subtotal TYPE BIGINT USING (subtotal * 100)::BIGINT;

-- PaymentTransaction table
ALTER TABLE payment_transactions
    ALTER COLUMN amount TYPE BIGINT USING (amount * 100)::BIGINT;

-- Discounts table
ALTER TABLE discounts
    ALTER COLUMN min_order_value TYPE BIGINT USING (min_order_value * 100)::BIGINT,
    ALTER COLUMN max_discount_value TYPE BIGINT USING (max_discount_value * 100)::BIGINT;

-- ShippingRate table
ALTER TABLE shipping_rates
    ALTER COLUMN base_rate TYPE BIGINT USING (base_rate * 100)::BIGINT,
    ALTER COLUMN min_order_value TYPE BIGINT USING (min_order_value * 100)::BIGINT,
    ALTER COLUMN free_shipping_threshold TYPE BIGINT USING (free_shipping_threshold * 100)::BIGINT;

-- WeightBasedRate table
ALTER TABLE weight_based_rates
    ALTER COLUMN rate TYPE BIGINT USING (rate * 100)::BIGINT;

-- ValueBasedRate table
ALTER TABLE value_based_rates
    ALTER COLUMN rate TYPE BIGINT USING (rate * 100)::BIGINT,
    ALTER COLUMN min_order_value TYPE BIGINT USING (min_order_value * 100)::BIGINT,
    ALTER COLUMN max_order_value TYPE BIGINT USING (max_order_value * 100)::BIGINT;

-- Products table
ALTER TABLE products
    ALTER COLUMN price TYPE BIGINT USING (price * 100)::BIGINT;

-- ProductVariants table
ALTER TABLE product_variants
    ALTER COLUMN price TYPE BIGINT USING (price * 100)::BIGINT;