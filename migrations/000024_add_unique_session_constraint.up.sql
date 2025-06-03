-- Add unique constraint to prevent multiple active checkouts for the same session_id
-- This constraint ensures only one active checkout can exist per session_id
-- We use a partial unique index to only apply the constraint to active checkouts

-- First, clean up any duplicate active checkouts (keeping the most recent one)
DELETE FROM checkouts 
WHERE id NOT IN (
    SELECT DISTINCT ON (session_id) id 
    FROM checkouts 
    WHERE status = 'active' AND session_id IS NOT NULL
    ORDER BY session_id, created_at DESC
) AND status = 'active' AND session_id IS NOT NULL;

-- Create a partial unique index on session_id for active checkouts
CREATE UNIQUE INDEX idx_checkouts_unique_active_session 
ON checkouts(session_id) 
WHERE status = 'active' AND session_id IS NOT NULL;