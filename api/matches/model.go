package matches

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"dotapro-lambda-api/types"
)
type Model struct {
	DB *bun.DB
}

func NewModel(db *bun.DB) *Model { return &Model{DB: db} }

func (m *Model) GetOne(ctx context.Context, id int64) (*types.MatchDetail, error) {
	res := &types.MatchDetail{}
	q := m.buildQueryGetOne(id)
	log.Debug().Str("BUN_QUERY", q.String()).Send()
	if err := q.Scan(ctx, res); err != nil {
		return nil, err
	}
	return res, nil
}
