-- This migration ensures the checkouts table exists and has a session_id column
-- First, make sure the checkouts table exists
CREATE TABLE IF NOT EXISTS checkouts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    shipping_address JSONB NOT NULL DEFAULT '{}',
    billing_address JSONB NOT NULL DEFAULT '{}',
    shipping_method_id INTEGER REFERENCES shipping_methods(id) ON DELETE SET NULL,
    payment_provider VARCHAR(255),
    total_amount BIGINT NOT NULL DEFAULT 0,
    shipping_cost BIGINT NOT NULL DEFAULT 0,
    total_weight DECIMAL(10, 3) NOT NULL DEFAULT 0,
    customer_details JSONB NOT NULL DEFAULT '{}',
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    discount_code VARCHAR(100),
    discount_amount BIGINT NOT NULL DEFAULT 0,
    final_amount BIGINT NOT NULL DEFAULT 0,
    applied_discount JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_activity_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL DEFAULT (NOW() + INTERVAL '1 hour'),
    completed_at TIMESTAMP,
    converted_order_id INTEGER REFERENCES orders(id) ON DELETE SET NULL
);

-- Now, check if the session_id column exists. If not, add it.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'checkouts' 
        AND column_name = 'session_id'
    ) THEN
        ALTER TABLE checkouts ADD COLUMN session_id VARCHAR(255) NULL;
        CREATE INDEX IF NOT EXISTS idx_checkouts_session_id ON checkouts(session_id);
    END IF;
END $$;

-- Add other indices if they don't exist
CREATE INDEX IF NOT EXISTS idx_checkouts_user_id ON checkouts(user_id);
CREATE INDEX IF NOT EXISTS idx_checkouts_status ON checkouts(status);