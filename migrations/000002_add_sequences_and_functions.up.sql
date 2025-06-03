-- Migration to add sequences and utility functions
-- This adds support for friendly order numbers and other utility functions

-- Create sequence for friendly order numbers
CREATE SEQUENCE IF NOT EXISTS order_friendly_sequence START 1001;

-- Function to generate friendly order IDs
CREATE OR REPLACE FUNCTION generate_friendly_order_id() RETURNS VARCHAR(20) AS $$
DECLARE
    next_val BIGINT;
    friendly_id VARCHAR(20);
BEGIN
    SELECT nextval('order_friendly_sequence') INTO next_val;
    friendly_id := 'ORD-' || LPAD(next_val::text, 6, '0');
    RETURN friendly_id;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to auto-generate friendly_id for orders
CREATE OR REPLACE FUNCTION set_order_friendly_id()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.friendly_id IS NULL OR NEW.friendly_id = '' THEN
        NEW.friendly_id := generate_friendly_order_id();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger on orders table
CREATE TRIGGER trigger_set_order_friendly_id
    BEFORE INSERT ON orders
    FOR EACH ROW
    EXECUTE FUNCTION set_order_friendly_id();

-- Function to update timestamps
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updated_at on all relevant tables
CREATE TRIGGER trigger_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_product_variants_updated_at
    BEFORE UPDATE ON product_variants
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_currency_settings_updated_at
    BEFORE UPDATE ON currency_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_shipping_methods_updated_at
    BEFORE UPDATE ON shipping_methods
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_payment_transactions_updated_at
    BEFORE UPDATE ON payment_transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_discounts_updated_at
    BEFORE UPDATE ON discounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_webhooks_updated_at
    BEFORE UPDATE ON webhooks
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_checkouts_updated_at
    BEFORE UPDATE ON checkouts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_checkout_items_updated_at
    BEFORE UPDATE ON checkout_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();
