package types

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

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
	LogoURL       string `bun:"logo_url"`
}

type Player struct {
	bun.BaseModel `bun:"table:players"`
	PlayerID      int64  `bun:"player_id,pk"`
	Name          string `bun:"name,notnull"`
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
	Patch          string    `bun:"patch"`
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
	Version        int             `bun:"version"`
}

type SeriesMatch struct {
	bun.BaseModel `bun:"table:series_match"`
	SeriesID      int64 `bun:"series_id,pk"`
	MatchID       int64 `bun:"match_id,pk"`
}

type TeamInfo struct {
	ID      int64  `json:"id" bun:"id"`
	Name    string `json:"name" bun:"name"`
	Tag     string `json:"tag" bun:"tag"`
	LogoURL string `json:"logo_url" bun:"logo_url"`
	Score   int    `json:"score" bun:"score"`
	Captain int64  `json:"captain" bun:"captain"`
}

type LeagueInfo struct {
	ID   int64  `json:"id" bun:"id"`
	Name string `json:"name" bun:"name"`
	Tier string `json:"tier" bun:"tier"`
}

type MatchDetail struct {
	MatchID        int64       `json:"match_id" bun:"match_id"`
	StartTime      time.Time   `json:"start_time" bun:"start_time"`
	Duration       int         `json:"duration" bun:"duration"`
	RadiantWin     bool        `json:"radiant_win" bun:"radiant_win"`
	Patch          string      `json:"patch" bun:"patch"`
	Version        int         `json:"version" bun:"version"`
	RadiantTeam    TeamInfo    `json:"radiant_team" bun:"radiant_team"`
	DireTeam       TeamInfo    `json:"dire_team" bun:"dire_team"`
	League         LeagueInfo  `json:"league" bun:"league"`
	SeriesID       int64       `json:"series_id" bun:"series_id"`
	PicksBans      json.RawMessage `json:"picks_bans" bun:"picks_bans"`
	PlayersData    json.RawMessage `json:"players_data" bun:"players_data"`
	RadiantGoldAdv []int32     `json:"radiant_gold_adv" bun:"radiant_gold_adv"`
	RadiantXPAdv   []int32     `json:"radiant_xp_adv" bun:"radiant_xp_adv"`
	RadiantHeroes  []int64     `json:"radiant_heroes" bun:"radiant_heroes"`
	DireHeroes     []int64     `json:"dire_heroes" bun:"dire_heroes"`
	RadiantPlayers []int64     `json:"radiant_players" bun:"radiant_players"`
	DirePlayers    []int64     `json:"dire_players" bun:"dire_players"`
}
