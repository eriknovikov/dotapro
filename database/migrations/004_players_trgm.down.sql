-- Remove the trigram index on players.name
DROP INDEX IF EXISTS players_name_trgm_idx;

-- Note: We do NOT drop the pg_trgm extension as it may be used by other tables
-- (leagues and teams also use pg_trgm indexes)
