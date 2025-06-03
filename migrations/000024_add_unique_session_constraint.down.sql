-- Remove the unique constraint on session_id for active checkouts
DROP INDEX IF EXISTS idx_checkouts_unique_active_session;