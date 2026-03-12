-- Enable pg_trgm extension if not already enabled
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create GIN index on players.name for trigram-based fuzzy search
-- This enables efficient word_similarity operator (%>) for player name searches
CREATE INDEX IF NOT EXISTS players_name_trgm_idx 
ON players USING gin (name gin_trgm_ops);
