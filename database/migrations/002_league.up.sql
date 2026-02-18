CREATE EXTENSION IF NOT EXISTS pg_trgm;



CREATE INDEX IF NOT EXISTS leagues_name_trgm_idx 
ON leagues USING gin (name gin_trgm_ops);


CREATE INDEX IF NOT EXISTS teams_name_trgm_idx 
ON teams USING gin (name gin_trgm_ops);