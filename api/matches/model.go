package matches

import (
	"context"

	"dotapro-lambda-api/types"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type Model struct {
	DB *bun.DB
}

func NewModel(db *bun.DB) *Model { return &Model{DB: db} }

func (m *Model) GetMany(ctx context.Context, filter types.GetMatchesFilter) ([]types.MatchSummary, error) {
	var res []types.MatchSummary
	q := m.buildQueryGetMany(filter)
	if err := q.Scan(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *Model) GetOne(ctx context.Context, id int64) (*types.MatchDetail, error) {
	res := &types.MatchDetail{}
	q := m.buildQueryGetOne(id)
	log.Debug().Str("BUN_QUERY", q.String()).Send()
	if err := q.Scan(ctx, res); err != nil {
		return nil, err
	}
	return res, nil
}
