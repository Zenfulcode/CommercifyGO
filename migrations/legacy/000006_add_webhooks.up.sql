-- Create webhooks table
CREATE TABLE IF NOT EXISTS webhooks (
    id SERIAL PRIMARY KEY,
    provider VARCHAR(50) NOT NULL, -- e.g., 'mobilepay', 'stripe', etc.
    external_id VARCHAR(255), -- ID assigned by the provider
    url VARCHAR(255) NOT NULL,
    events JSONB NOT NULL, -- Array of event types this webhook is registered for
    secret VARCHAR(255), -- Webhook secret for verification
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on provider for faster lookups
CREATE INDEX IF NOT EXISTS idx_webhooks_provider ON webhooks (provider);