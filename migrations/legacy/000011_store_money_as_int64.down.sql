-- Migration to revert money fields from INT (cents) back to DECIMAL
-- Create temporary columns with _decimal suffix

-- Products table
ALTER TABLE products ADD COLUMN price_decimal DECIMAL(10, 2);
UPDATE products SET price_decimal = price::DECIMAL / 100;

-- Product variants table
ALTER TABLE product_variants ADD COLUMN price_decimal DECIMAL(10, 2);
UPDATE product_variants SET 
    price_decimal = price::DECIMAL / 100;

-- Orders table
ALTER TABLE orders ADD COLUMN total_amount_decimal DECIMAL(10, 2);
ALTER TABLE orders ADD COLUMN shipping_cost_decimal DECIMAL(10, 2);
ALTER TABLE orders ADD COLUMN discount_amount_decimal DECIMAL(10, 2);
ALTER TABLE orders ADD COLUMN final_amount_decimal DECIMAL(10, 2);
UPDATE orders SET 
    total_amount_decimal = total_amount::DECIMAL / 100,
    shipping_cost_decimal = shipping_cost::DECIMAL / 100,
    discount_amount_decimal = discount_amount::DECIMAL / 100,
    final_amount_decimal = final_amount::DECIMAL / 100;

-- Order items table
ALTER TABLE order_items ADD COLUMN price_decimal DECIMAL(10, 2);
ALTER TABLE order_items ADD COLUMN subtotal_decimal DECIMAL(10, 2);
UPDATE order_items SET 
    price_decimal = price::DECIMAL / 100,
    subtotal_decimal = subtotal::DECIMAL / 100;

-- Shipping rates table
ALTER TABLE shipping_rates ADD COLUMN base_rate_decimal DECIMAL(10, 2);
ALTER TABLE shipping_rates ADD COLUMN min_order_value_decimal DECIMAL(10, 2);
ALTER TABLE shipping_rates ADD COLUMN free_shipping_threshold_decimal DECIMAL(10, 2);
UPDATE shipping_rates SET 
    base_rate_decimal = base_rate::DECIMAL / 100,
    min_order_value_decimal = min_order_value::DECIMAL / 100,
    free_shipping_threshold_decimal = CASE WHEN free_shipping_threshold IS NOT NULL THEN free_shipping_threshold::DECIMAL / 100 ELSE NULL END;

-- Weight-based rates table
ALTER TABLE weight_based_rates ADD COLUMN rate_decimal DECIMAL(10, 2);
UPDATE weight_based_rates SET rate_decimal = rate::DECIMAL / 100;

-- Value-based rates table
ALTER TABLE value_based_rates ADD COLUMN min_order_value_decimal DECIMAL(10, 2);
ALTER TABLE value_based_rates ADD COLUMN max_order_value_decimal DECIMAL(10, 2);
ALTER TABLE value_based_rates ADD COLUMN rate_decimal DECIMAL(10, 2);
UPDATE value_based_rates SET 
    min_order_value_decimal = min_order_value::DECIMAL / 100,
    max_order_value_decimal = max_order_value::DECIMAL / 100,
    rate_decimal = rate::DECIMAL / 100;

-- Discounts table
ALTER TABLE discounts ADD COLUMN min_order_value_decimal DECIMAL(10, 2);
ALTER TABLE discounts ADD COLUMN max_discount_value_decimal DECIMAL(10, 2);
UPDATE discounts SET 
    min_order_value_decimal = min_order_value::DECIMAL / 100,
    max_discount_value_decimal = max_discount_value::DECIMAL / 100;

-- Payment transactions table
ALTER TABLE payment_transactions ADD COLUMN amount_decimal DECIMAL(10, 2);
UPDATE payment_transactions SET amount_decimal = amount::DECIMAL / 100;

-- Now drop the int columns and rename the decimal ones
-- Products
ALTER TABLE products DROP COLUMN IF EXISTS price;
ALTER TABLE products RENAME COLUMN price_decimal TO price;

-- Product variants
ALTER TABLE product_variants DROP COLUMN price;
ALTER TABLE product_variants RENAME COLUMN price_decimal TO price;

-- Orders
ALTER TABLE orders DROP COLUMN total_amount;
ALTER TABLE orders DROP COLUMN shipping_cost;
ALTER TABLE orders DROP COLUMN discount_amount;
ALTER TABLE orders DROP COLUMN final_amount;
ALTER TABLE orders RENAME COLUMN total_amount_decimal TO total_amount;
ALTER TABLE orders RENAME COLUMN shipping_cost_decimal TO shipping_cost;
ALTER TABLE orders RENAME COLUMN discount_amount_decimal TO discount_amount;
ALTER TABLE orders RENAME COLUMN final_amount_decimal TO final_amount;

-- Order items
ALTER TABLE order_items DROP COLUMN price;
ALTER TABLE order_items DROP COLUMN subtotal;
ALTER TABLE order_items RENAME COLUMN price_decimal TO price;
ALTER TABLE order_items RENAME COLUMN subtotal_decimal TO subtotal;

-- Shipping rates
ALTER TABLE shipping_rates DROP COLUMN base_rate;
ALTER TABLE shipping_rates DROP COLUMN min_order_value;
ALTER TABLE shipping_rates DROP COLUMN free_shipping_threshold;
ALTER TABLE shipping_rates RENAME COLUMN base_rate_decimal TO base_rate;
ALTER TABLE shipping_rates RENAME COLUMN min_order_value_decimal TO min_order_value;
ALTER TABLE shipping_rates RENAME COLUMN free_shipping_threshold_decimal TO free_shipping_threshold;

-- Weight-based rates
ALTER TABLE weight_based_rates DROP COLUMN rate;
ALTER TABLE weight_based_rates RENAME COLUMN rate_decimal TO rate;

-- Value-based rates
ALTER TABLE value_based_rates DROP COLUMN min_order_value;
ALTER TABLE value_based_rates DROP COLUMN max_order_value;
ALTER TABLE value_based_rates DROP COLUMN rate;
ALTER TABLE value_based_rates RENAME COLUMN min_order_value_decimal TO min_order_value;
ALTER TABLE value_based_rates RENAME COLUMN max_order_value_decimal TO max_order_value;
ALTER TABLE value_based_rates RENAME COLUMN rate_decimal TO rate;

-- Discounts
ALTER TABLE discounts DROP COLUMN min_order_value;
ALTER TABLE discounts DROP COLUMN max_discount_value;
ALTER TABLE discounts RENAME COLUMN min_order_value_decimal TO min_order_value;
ALTER TABLE discounts RENAME COLUMN max_discount_value_decimal TO max_discount_value;

-- Payment transactions
ALTER TABLE payment_transactions DROP COLUMN amount;
ALTER TABLE payment_transactions RENAME COLUMN amount_decimal TO amount;
