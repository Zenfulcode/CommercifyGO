-- Create payment transactions table
CREATE TABLE payment_transactions (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    transaction_id VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,  -- authorize, capture, refund, cancel
    status VARCHAR(50) NOT NULL,  -- successful, failed, pending
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    raw_response TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create indexes
CREATE INDEX idx_payment_transactions_order_id ON payment_transactions(order_id);
CREATE INDEX idx_payment_transactions_transaction_id ON payment_transactions(transaction_id);
CREATE INDEX idx_payment_transactions_type ON payment_transactions(type);
CREATE INDEX idx_payment_transactions_status ON payment_transactions(status);
CREATE INDEX idx_payment_transactions_created_at ON payment_transactions(created_at);

-- Backfill payment transactions for existing orders
DO $$
DECLARE
    order_rec RECORD;
    now_time TIMESTAMP WITH TIME ZONE := NOW();
BEGIN
    -- Find orders with payment IDs that should have transaction records
    FOR order_rec IN 
        SELECT 
            id, 
            payment_id, 
            payment_provider, 
            COALESCE(final_amount, total_amount) as amount, 
            status, 
            updated_at
        FROM orders 
        WHERE payment_id IS NOT NULL 
          AND payment_id != '' 
          AND payment_provider IS NOT NULL 
          AND payment_provider != ''
    LOOP
        -- Create transaction records based on order status
        CASE order_rec.status
            -- For paid orders, create an authorization transaction
            WHEN 'paid' THEN
                INSERT INTO payment_transactions (
                    order_id, transaction_id, type, status, amount, currency, provider, 
                    metadata, created_at, updated_at
                ) VALUES (
                    order_rec.id, order_rec.payment_id, 'authorize', 'successful', 
                    order_rec.amount, 'USD', order_rec.payment_provider,
                    '{"payment_method":"credit_card"}'::JSONB, order_rec.updated_at, now_time
                );
                
            -- For captured orders, create both auth and capture transactions
            WHEN 'captured' THEN
                -- Create auth transaction (happening before capture)
                INSERT INTO payment_transactions (
                    order_id, transaction_id, type, status, amount, currency, provider, 
                    metadata, created_at, updated_at
                ) VALUES (
                    order_rec.id, order_rec.payment_id, 'authorize', 'successful', 
                    order_rec.amount, 'USD', order_rec.payment_provider,
                    '{"payment_method":"credit_card"}'::JSONB, 
                    order_rec.updated_at - INTERVAL '1 hour', now_time
                );
                
                -- Create capture transaction
                INSERT INTO payment_transactions (
                    order_id, transaction_id, type, status, amount, currency, provider, 
                    metadata, created_at, updated_at
                ) VALUES (
                    order_rec.id, order_rec.payment_id, 'capture', 'successful', 
                    order_rec.amount, 'USD', order_rec.payment_provider,
                    '{"full_capture":"true","remaining_amount":"0"}'::JSONB,
                    order_rec.updated_at, now_time
                );
                
            -- For refunded orders, create auth, capture, and refund transactions
            WHEN 'refunded' THEN
                -- Create auth transaction (happened first)
                INSERT INTO payment_transactions (
                    order_id, transaction_id, type, status, amount, currency, provider, 
                    metadata, created_at, updated_at
                ) VALUES (
                    order_rec.id, order_rec.payment_id, 'authorize', 'successful', 
                    order_rec.amount, 'USD', order_rec.payment_provider,
                    '{"payment_method":"credit_card"}'::JSONB, 
                    order_rec.updated_at - INTERVAL '2 hours', now_time
                );
                
                -- Create capture transaction (happened after auth)
                INSERT INTO payment_transactions (
                    order_id, transaction_id, type, status, amount, currency, provider, 
                    metadata, created_at, updated_at
                ) VALUES (
                    order_rec.id, order_rec.payment_id, 'capture', 'successful', 
                    order_rec.amount, 'USD', order_rec.payment_provider,
                    '{"full_capture":"true","remaining_amount":"0"}'::JSONB,
                    order_rec.updated_at - INTERVAL '1 hour', now_time
                );
                
                -- Create refund transaction
                INSERT INTO payment_transactions (
                    order_id, transaction_id, type, status, amount, currency, provider, 
                    metadata, created_at, updated_at
                ) VALUES (
                    order_rec.id, order_rec.payment_id, 'refund', 'successful', 
                    order_rec.amount, 'USD', order_rec.payment_provider,
                    (format('{"full_refund":"true","total_refunded":"%s","remaining_available":"0"}', order_rec.amount::text))::JSONB,
                    order_rec.updated_at, now_time
                );
                
            -- For cancelled orders, create auth and cancel transactions
            WHEN 'cancelled' THEN
                -- Create auth transaction (happened first)
                INSERT INTO payment_transactions (
                    order_id, transaction_id, type, status, amount, currency, provider, 
                    metadata, created_at, updated_at
                ) VALUES (
                    order_rec.id, order_rec.payment_id, 'authorize', 'successful', 
                    order_rec.amount, 'USD', order_rec.payment_provider,
                    '{"payment_method":"credit_card"}'::JSONB, 
                    order_rec.updated_at - INTERVAL '1 hour', now_time
                );
                
                -- Create cancel transaction
                INSERT INTO payment_transactions (
                    order_id, transaction_id, type, status, amount, currency, provider, 
                    metadata, created_at, updated_at
                ) VALUES (
                    order_rec.id, order_rec.payment_id, 'cancel', 'successful', 
                    0, 'USD', order_rec.payment_provider,
                    '{}'::JSONB, order_rec.updated_at, now_time
                );
                
            -- For pending_action orders, create a pending authorization transaction
            WHEN 'pending_action' THEN
                INSERT INTO payment_transactions (
                    order_id, transaction_id, type, status, amount, currency, provider, 
                    metadata, created_at, updated_at
                ) VALUES (
                    order_rec.id, order_rec.payment_id, 'authorize', 'pending', 
                    order_rec.amount, 'USD', order_rec.payment_provider,
                    '{"payment_method":"credit_card","requires_action":"true"}'::JSONB,
                    order_rec.updated_at, now_time
                );
                
            -- For any other status with payment_id, create a basic authorization record
            ELSE
                INSERT INTO payment_transactions (
                    order_id, transaction_id, type, status, amount, currency, provider, 
                    metadata, created_at, updated_at
                ) VALUES (
                    order_rec.id, order_rec.payment_id, 'authorize', 'successful', 
                    order_rec.amount, 'USD', order_rec.payment_provider,
                    '{"payment_method":"credit_card"}'::JSONB, order_rec.updated_at, now_time
                );
        END CASE;
    END LOOP;
END $$;