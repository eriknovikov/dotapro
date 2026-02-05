package main

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

// --- DB Models ---

type League struct {
	bun.BaseModel `bun:"table:leagues"`
	LeagueID      int64  `bun:"league_id,pk"`
	Name          string `bun:"name,notnull"`
	Tier          string `bun:"tier"`
}

type Team struct {
	bun.BaseModel `bun:"table:teams"`
	TeamID        int64  `bun:"team_id,pk"`
	Name          string `bun:"name"`
	Tag           string `bun:"tag"`
	LogoURL       string `bun:"logo_url"` // Empty string represents NULL
}

type Player struct {
	bun.BaseModel `bun:"table:players"`
	PlayerID      int64  `bun:"player_id,pk"`
	Name          string `bun:"name"`       // Empty string represents NULL
	ProfileImg    string `bun:"profile_img"` // Empty string represents NULL
}

type Series struct {
	bun.BaseModel `bun:"table:series"`
	SeriesID      int64     `bun:"series_id,pk"`
	LeagueID      int64     `bun:"league_id,nullzero"`
	TeamAID       int64     `bun:"team_a_id,nullzero"`
	TeamBID       int64     `bun:"team_b_id,nullzero"`
	StartTime     time.Time `bun:"start_time,notnull"`
	TeamAScore    int16     `bun:"team_a_score,default:0"`
	TeamBScore    int16     `bun:"team_b_score,default:0"`
}

type Match struct {
	bun.BaseModel  `bun:"table:matches"`
	MatchID        int64     `bun:"match_id,pk"`
	LeagueID       int64     `bun:"league_id,nullzero"`
	RadiantTeamID  int64     `bun:"radiant_team_id,nullzero"`
	DireTeamID     int64     `bun:"dire_team_id,nullzero"`
	RadiantHeroes  []int64   `bun:"radiant_heroes,array"`
	DireHeroes     []int64   `bun:"dire_heroes,array"`
	RadiantPlayers []int64   `bun:"radiant_players,array"`
	DirePlayers    []int64   `bun:"dire_players,array"`
	Duration       int       `bun:"duration,notnull"`
	StartTime      time.Time `bun:"start_time,notnull"`
	RadiantWin     bool      `bun:"radiant_win,notnull"`
}

type MatchMetadata struct {
	bun.BaseModel  `bun:"table:matches_metadata"`
	MatchID        int64           `bun:"match_id,pk"`
	SeriesID       int64           `bun:"series_id,nullzero"`
	RadiantCaptain int64           `bun:"radiant_captain,nullzero"`
	DireCaptain    int64           `bun:"dire_captain,nullzero"`
	PicksBans      json.RawMessage `bun:"picks_bans,type:jsonb"`
	PlayersData    json.RawMessage `bun:"players_data,type:jsonb"`
	RadiantGoldAdv []int32         `bun:"radiant_gold_adv,array"`
	RadiantXPAdv   []int32         `bun:"radiant_xp_adv,array"`
	RadiantScore   int             `bun:"radiant_score,notnull"`
	DireScore      int             `bun:"dire_score,notnull"`
}

type SeriesMatch struct {
	bun.BaseModel `bun:"table:series_match"`
	SeriesID      int64 `bun:"series_id,pk"`
	MatchID       int64 `bun:"match_id,pk"`
}

type ScraperMetadata struct {
	bun.BaseModel        `bun:"table:scraper_metadata"`
	ID                   int16  `bun:"id,pk"`
	LastFetchedMatchID   int64  `bun:"last_fetched_match_id,notnull"`
}

// --- Source Structs (OpenDota Mapping) ---

type ODMatch struct {
	MatchID        int64           `json:"match_id"`
	RadiantWin     bool            `json:"radiant_win"`
	StartTime      int64           `json:"start_time"`
	Duration       int             `json:"duration"`
	SeriesID       int64           `json:"series_id"`
	RadiantGoldAdv []int32         `json:"radiant_gold_adv"`
	RadiantXPAdv   []int32         `json:"radiant_xp_adv"`
	League         ODLeague        `json:"league"`
	RadiantTeam    ODTeam          `json:"radiant_team"`
	DireTeam       ODTeam          `json:"dire_team"`
	Players        json.RawMessage `json:"players"`
	PicksBans      json.RawMessage `json:"picks_bans"`
}

type ODLeague struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Tier string `json:"tier"`
}

type ODTeam struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	LogoURL string `json:"logo_url"` // Empty string represents NULL
	Score   int    `json:"score"`
	Captain int64  `json:"captain"`  // 0 represents NULL
}

type ODPlayerShort struct {
	PlayerID   int64  `json:"player_id"`
	HeroID     int64  `json:"hero_id"`
	PlayerSlot int    `json:"player_slot"`
	Name       string `json:"name"` // Empty string represents NULL
}
