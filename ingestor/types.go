package main

// MatchMessage represents the structure of the JSON message from SQS
type MatchMessage struct {
	MatchID            int64         `json:"match_id"`
	RadiantWin         bool          `json:"radiant_win"`
	StartTime          int64         `json:"start_time"`
	Duration           int           `json:"duration"`
	LeagueID           int           `json:"league_id"`
	LeagueName         string        `json:"league_name"`
	LeagueTier         string        `json:"league_tier"`
	RadiantTeamID      int64         `json:"radiant_team_id"`
	RadiantTeamName    string        `json:"radiant_team_name"`
	RadiantTeamTag     *string       `json:"radiant_team_tag"`
	RadiantTeamLogoURL *string       `json:"radiant_team_logo_url"`
	DireTeamID         int64         `json:"dire_team_id"`
	DireTeamName       string        `json:"dire_team_name"`
	DireTeamTag        *string       `json:"dire_team_tag"`
	DireTeamLogoURL    *string       `json:"dire_team_logo_url"`
	RadiantScore       int           `json:"radiant_score"`
	DireScore          int           `json:"dire_score"`
	RadiantCaptain     *int64        `json:"radiant_captain"`
	DireCaptain        *int64        `json:"dire_captain"`
	RadiantGoldAdv     []int         `json:"radiant_gold_adv"`
	RadiantXPAdv       []int         `json:"radiant_xp_adv"`
	Patch              string        `json:"patch"`
	Version            int           `json:"version"`
	SeriesID           *int64        `json:"series_id"`
	Players            []PlayerMatch `json:"players"`
	PicksBans          []PickBan     `json:"picks_bans"`
}

// PlayerMatch represents a player's stats in a match
type PlayerMatch struct {
	HeroID      int     `json:"hero_id"`
	AccountID   *int64  `json:"account_id"`
	PlayerSlot  int     `json:"player_slot"`
	Kills       int     `json:"kills"`
	Deaths      int     `json:"deaths"`
	Assists     int     `json:"assists"`
	GoldPerMin  int     `json:"gold_per_min"`
	XPPerMin    int     `json:"xp_per_min"`
	LastHits    int     `json:"last_hits"`
	Denies      int     `json:"denies"`
	Level       int     `json:"level"`
	Item0       int     `json:"item_0"`
	Item1       int     `json:"item_1"`
	Item2       int     `json:"item_2"`
	Item3       int     `json:"item_3"`
	Item4       int     `json:"item_4"`
	Item5       int     `json:"item_5"`
	ItemNeutral int     `json:"item_neutral"`
	Backpack0   int     `json:"backpack_0"`
	Backpack1   int     `json:"backpack_1"`
	Backpack2   int     `json:"backpack_2"`
	NetWorth    int     `json:"net_worth"`
	IsRadiant   bool    `json:"is_radiant"`
	Name        *string `json:"name"`
	Facet       int     `json:"facet"`
}

type PickBan struct {
	HeroID int  `json:"hero_id"`
	IsPick bool `json:"is_pick"`
	Order  int  `json:"order"`
	Team   int  `json:"team"`
}
