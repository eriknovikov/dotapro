package matches

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type Match struct {
	MatchID        int64     `bun:"match_id,bigint,pk"`
	LeagueID       int64     `bun:"league_id,bigint"`
	RadiantTeamID  int64     `bun:"radiant_team_id,bigint"`
	DireTeamID     int64     `bun:"dire_team_id,bigint"`
	RadiantHeroes  []int64   `bun:"radiant_heroes,bigint[]"`
	DireHeroes     []int64   `bun:"dire_heroes,bigint[]"`
	RadiantPlayers []int64   `bun:"radiant_players,bigint[]"`
	DirePlayers    []int64   `bun:"dire_players,bigint[]"`
	Duration       int       `bun:"duration,int"`
	StartTime      time.Time `bun:"start_time,timestamp"`
	RadiantWin     bool      `bun:"radiant_win,bool"`
}

type MatchMetadata struct {
	MatchID       int64   `bun:"match_id,bigint,pk"`
	SeriesID      int64   `bun:"series_id,bigint"`
	RadiantCaptain int64  `bun:"radiant_captain,bigint"`
	DireCaptain   int64   `bun:"dire_captain,bigint"`
	PicksBans     []byte  `bun:"picks_bans,jsonb"`
	PlayersData   []byte  `bun:"players_data,jsonb"`
	RadiantGoldAdv []int64 `bun:"radiant_gold_adv,bigint[]"`
	RadiantXPAdv  []int64 `bun:"radiant_xp_adv,bigint[]"`
	RadiantScore  int     `bun:"radiant_score,int"`
	DireScore     int     `bun:"dire_score,int"`
}
type Model struct {
	DB *bun.DB
}

func NewModel(db *bun.DB) *Model { return &Model{DB: db} }

func (m *Model) GetOne(ctx context.Context, id int64) (*Match, error) {
	res := &Match{}
	q := m.buildQueryGetOne(id)
	log.Debug().Str("BUN_QUERY", q.String()).Send()
	if err := q.Scan(ctx, res); err != nil {
		return nil, err
	}
	return res, nil
}
