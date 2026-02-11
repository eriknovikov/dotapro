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
	player_id BIGINT PRIMARY KEY,
	name TEXT NOT NULL
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

CREATE TABLE matches (
	match_id BIGINT PRIMARY KEY,
	league_id BIGINT REFERENCES leagues(league_id),
	radiant_team_id BIGINT REFERENCES teams(team_id),
	dire_team_id BIGINT REFERENCES teams(team_id),
	radiant_heroes BIGINT[],
	dire_heroes BIGINT[],
	radiant_players BIGINT[],
	dire_players BIGINT[],
	duration INTEGER NOT NULL,
	start_time TIMESTAMP NOT NULL,
	radiant_win BOOLEAN NOT NULL,
	patch TEXT
);

CREATE TABLE matches_metadata (
	match_id BIGINT PRIMARY KEY REFERENCES matches(match_id) ON DELETE CASCADE,
	radiant_captain BIGINT REFERENCES players(player_id),
	dire_captain BIGINT REFERENCES players(player_id),
	picks_bans JSONB,
	players_data JSONB,
	radiant_gold_adv INTEGER[],
	radiant_xp_adv INTEGER[],
	radiant_score INTEGER NOT NULL,
	dire_score INTEGER NOT NULL,
	version INTEGER
);

CREATE TABLE scraper_metadata (
	id SMALLINT PRIMARY KEY CHECK(id=1), --singleton
	last_fetched_match_id BIGINT NOT NULL,
	updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX series_match_match_id_idx ON series_match USING btree (match_id);