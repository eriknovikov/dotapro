package series

import (
	"context"
	"errors"

	"dotapro-lambda-api/errs"
	"dotapro-lambda-api/types"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type Model struct {
	DB *bun.DB
}

func NewModel(db *bun.DB) *Model { return &Model{DB: db} }

func (m *Model) GetMany(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, error) {
	var res []types.SeriesSummary
	q := m.buildQueryGetMany(filter)
	if err := q.Scan(ctx, &res); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NOT_FOUND
		}
		return nil, err
	}
	return res, nil
}

func (m *Model) getMatchesForSeries(ctx context.Context, seriesID int64) ([]types.SeriesMatchDetail, error) {
	var res []types.SeriesMatchDetail
	q := m.buildQueryGetMatchesForSeries(seriesID)
	log.Debug().Str("BUN_QUERY", q.String()).Send()
	if err := q.Scan(ctx, &res); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []types.SeriesMatchDetail{}, nil
		}
		return nil, err
	}
	return res, nil
}

func (m *Model) GetOne(ctx context.Context, id int64) (*types.SeriesDetail, error) {
	res := &types.SeriesDetail{}
	q := m.buildQueryGetOne(id)
	log.Debug().Str("BUN_QUERY", q.String()).Send()
	if err := q.Scan(ctx, res); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NOT_FOUND
		}
		return nil, err
	}

	matches, err := m.getMatchesForSeries(ctx, id)
	if err != nil {
		return nil, err
	}

	res.Matches = matches
	return res, nil
}
