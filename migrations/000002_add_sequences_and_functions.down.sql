-- Drop all triggers
DROP TRIGGER IF EXISTS trigger_users_updated_at ON users;
DROP TRIGGER IF EXISTS trigger_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS trigger_products_updated_at ON products;
DROP TRIGGER IF EXISTS trigger_product_variants_updated_at ON product_variants;
DROP TRIGGER IF EXISTS trigger_currency_settings_updated_at ON currency_settings;
DROP TRIGGER IF EXISTS trigger_shipping_methods_updated_at ON shipping_methods;
DROP TRIGGER IF EXISTS trigger_orders_updated_at ON orders;
DROP TRIGGER IF EXISTS trigger_payment_transactions_updated_at ON payment_transactions;
DROP TRIGGER IF EXISTS trigger_discounts_updated_at ON discounts;
DROP TRIGGER IF EXISTS trigger_webhooks_updated_at ON webhooks;
DROP TRIGGER IF EXISTS trigger_checkouts_updated_at ON checkouts;
DROP TRIGGER IF EXISTS trigger_checkout_items_updated_at ON checkout_items;

-- Drop trigger on orders
DROP TRIGGER IF EXISTS trigger_set_order_friendly_id ON orders;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at();
DROP FUNCTION IF EXISTS set_order_friendly_id();
DROP FUNCTION IF EXISTS generate_friendly_order_id();

-- Drop sequence
DROP SEQUENCE IF EXISTS order_friendly_sequence;
