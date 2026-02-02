DROP TABLE IF EXISTS user_feeds CASCADE;
DROP TABLE IF EXISTS user_feed_filters CASCADE;
DROP TABLE IF EXISTS matches_metadata CASCADE;
DROP TABLE IF EXISTS matches CASCADE;
DROP TABLE IF EXISTS series CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS players CASCADE;
DROP TABLE IF EXISTS teams CASCADE;
DROP TABLE IF EXISTS leagues CASCADE;
DROP TABLE IF EXISTS series_match CASCADE;


-- Create tables in dependency order (no foreign keys first, then dependent tables)

-- 1. Base tables with no foreign key dependencies
CREATE TABLE leagues (
	league_id BIGINT PRIMARY KEY,
	name TEXT NOT NULL,
	tier TEXT
);

CREATE TABLE teams (
	team_id BIGINT PRIMARY KEY,
	name TEXT,
	tag TEXT,
	logo_url TEXT
);

CREATE TABLE players (
	account_id BIGINT PRIMARY KEY,
	name TEXT,
	profile_img TEXT
);

CREATE TABLE users (
	user_id UUID PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	username TEXT UNIQUE NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
);

-- 2. Tables that depend on leagues and teams
CREATE TABLE series (
	series_id BIGINT PRIMARY KEY,
	league_id BIGINT REFERENCES leagues(league_id),
	team_a_id BIGINT REFERENCES teams(team_id),
	team_b_id BIGINT REFERENCES teams(team_id),
	start_time TIMESTAMP NOT NULL,
	team_a_score SMALLINT DEFAULT 0,
	team_b_score SMALLINT DEFAULT 0
);

CREATE TABLE series_match (
	series_id BIGINT,
	match_id BIGINT,
	PRIMARY KEY (series_id, match_id)
);

-- 3. Matches table (depends on leagues, teams, and series)
CREATE TABLE matches (
	match_id BIGINT PRIMARY KEY,
	league_id BIGINT REFERENCES leagues(league_id),
	series_id BIGINT NOT NULL,
	radiant_team_id BIGINT REFERENCES teams(team_id),
	dire_team_id BIGINT REFERENCES teams(team_id),
	radiant_heroes BIGINT[],
	dire_heroes BIGINT[],
	radiant_players BIGINT[],
	dire_players BIGINT[],
	start_time TIMESTAMP NOT NULL,
	radiant_win BOOLEAN NOT NULL
);


CREATE TABLE matches_metadata (
	match_id BIGINT PRIMARY KEY REFERENCES matches(match_id) ON DELETE CASCADE,
	series_id BIGINT REFERENCES series(series_id),
	radiant_captain BIGINT REFERENCES players(account_id),
	dire_captain BIGINT REFERENCES players(account_id),
	duration INTEGER NOT NULL,
	picks_bans JSONB,
	players_data JSONB,
	radiant_gold_adv INTEGER[],
	radiant_xp_adv INTEGER[],
	radiant_score INTEGER NOT NULL,
	dire_score INTEGER NOT NULL
);



-- Indexes
CREATE INDEX leagues_league_id_idx ON leagues USING btree (league_id);
CREATE INDEX players_account_id_idx ON players USING btree (account_id);
CREATE INDEX teams_team_id_idx ON teams USING btree (team_id);

CREATE INDEX matches_dire_heroes_idx ON matches USING gin (dire_heroes);
CREATE INDEX matches_dire_players_idx ON matches USING gin (dire_players);
CREATE INDEX matches_dire_team_id_idx ON matches USING btree (dire_team_id);
CREATE INDEX matches_league_id_idx ON matches USING btree (league_id);
CREATE INDEX matches_radiant_heroes_idx ON matches USING gin (radiant_heroes);
CREATE INDEX matches_radiant_players_idx ON matches USING gin (radiant_players);
CREATE INDEX matches_radiant_team_id_idx ON matches USING btree (radiant_team_id);
CREATE INDEX matches_series_id_idx ON matches USING btree (series_id);
CREATE INDEX matches_start_time_idx ON matches USING btree (start_time);

CREATE INDEX series_league_id_idx ON series USING btree (league_id);
CREATE INDEX series_start_time_idx ON series USING btree (start_time);
CREATE INDEX series_team_a_id_idx ON series USING btree (team_a_id);
CREATE INDEX series_team_b_id_idx ON series USING btree (team_b_id);



CREATE OR REPLACE FUNCTION after_match_insert()
RETURNS TRIGGER AS $$
DECLARE 
	team_a_won BOOLEAN;
	team_a_id BIGINT;
	--NEW = Match Record
BEGIN
	INSERT INTO series (series_id, league_id, team_a_id, team_b_id, start_time) VALUES (NEW.series_id, NEW.league_id, NEW.radiant_team_id, NEW.dire_team_id, NEW.start_time) ON CONFLICT DO NOTHING;
	SELECT series.team_a_id INTO team_a_id FROM series WHERE series.series_id = NEW.series_id;
	INSERT INTO series_match VALUES(NEW.series_id, NEW.match_id);
	team_a_won := FALSE;
	IF (team_a_id = NEW.radiant_team_id AND NEW.radiant_win) OR (team_a_id = NEW.dire_team_id AND NOT NEW.radiant_win)
		THEN team_a_won := TRUE;
	END IF;

	IF team_a_won THEN
		UPDATE series SET team_a_score = team_a_score + 1 WHERE series_id = NEW.series_id;
	ELSE
		UPDATE series SET team_b_score = team_b_score + 1 WHERE series_id = NEW.series_id;
	END IF;
	RETURN NULL;
END;
$$ LANGUAGE plpgsql;



CREATE TRIGGER after_match_insert_trigger 
AFTER INSERT ON matches
FOR EACH ROW EXECUTE FUNCTION after_match_insert();

