package series

import (
	"context"
	"errors"

	"dotapro-lambda-api/constants"
	"dotapro-lambda-api/errs"
	"dotapro-lambda-api/types"
	"dotapro-lambda-api/utils"

	"github.com/jackc/pgx/v5"
	"github.com/uptrace/bun"
)

type Model struct {
	DB *bun.DB
}

func NewModel(db *bun.DB) *Model { return &Model{DB: db} }

func (m *Model) GetMany(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, types.PaginationData, error) {
	var res []types.SeriesSummary

	// Apply default limit if not specified
	limit := filter.Limit
	if limit <= 0 {
		limit = constants.DefaultLimit
	}

	q := m.DB.NewSelect().
		ColumnExpr("s.series_id").
		ColumnExpr("s.start_time").
		ColumnExpr("s.team_a_id AS team_a__id").
		ColumnExpr("s.team_b_id AS team_b__id").
		ColumnExpr("s.league_id AS league__id").
		ColumnExpr("team_a.name AS team_a__name").
		ColumnExpr("team_a.tag AS team_a__tag").
		ColumnExpr("team_a.logo_url AS team_a__logo_url").
		ColumnExpr("team_b.name AS team_b__name").
		ColumnExpr("team_b.tag AS team_b__tag").
		ColumnExpr("team_b.logo_url AS team_b__logo_url").
		ColumnExpr("l.name AS league__name").
		ColumnExpr("l.tier AS league__tier").
		ColumnExpr("s.team_a_score").
		ColumnExpr("s.team_b_score").
		TableExpr("series AS s").
		Join("LEFT JOIN leagues AS l USING (league_id)").
		Join("LEFT JOIN teams AS team_a ON team_a.team_id = s.team_a_id").
		Join("LEFT JOIN teams AS team_b ON team_b.team_id = s.team_b_id")

	if filter.LeagueID != nil {
		q = q.Where("s.league_id = ?", *filter.LeagueID)
	}
	if filter.TeamID != nil {
		q = q.Where("s.team_a_id = ? OR s.team_b_id = ?", *filter.TeamID, *filter.TeamID)
	}

	switch filter.Sort {
	case "oldest":
		q = q.Order("s.series_id ASC")
		if filter.Cursor != nil {
			q = q.Where("s.series_id > ?", *filter.Cursor)
		}
	default:
		q = q.Order("s.series_id DESC")
		if filter.Cursor != nil {
			q = q.Where("s.series_id < ?", *filter.Cursor)
		}
	}
	q = q.Limit(limit + 1)

	if err := q.Scan(ctx, &res); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, types.PaginationData{}, errs.NOT_FOUND
		}
		return nil, types.PaginationData{}, err
	}

	// Process pagination using the helper
	results, paginationData := utils.ProcessPagination(res, limit, func(s types.SeriesSummary) int64 {
		return s.SeriesID
	})

	return results, types.PaginationData{
		NextCursor: paginationData.NextCursor,
		HasMore:    paginationData.HasMore,
	}, nil
}

func (m *Model) getMatchesForSeries(ctx context.Context, seriesID int64) ([]types.SeriesMatchDetail, error) {
	var res []types.SeriesMatchDetail
	q := m.DB.NewSelect().
		ColumnExpr("m.match_id").
		ColumnExpr("m.duration").
		ColumnExpr("m.radiant_win").
		ColumnExpr("md.picks_bans").
		ColumnExpr("md.players_data").
		ColumnExpr("md.radiant_gold_adv").
		ColumnExpr("md.radiant_xp_adv").
		ColumnExpr("md.radiant_score").
		ColumnExpr("md.dire_score").
		ColumnExpr("md.radiant_captain").
		ColumnExpr("md.dire_captain").
		TableExpr("matches AS m").
		Join("LEFT JOIN matches_metadata AS md USING (match_id)").
		Join("LEFT JOIN series_match AS sm USING (match_id)").
		Where("sm.series_id = ?", seriesID).
		Order("m.match_id ASC")
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
	q := m.DB.NewSelect().
		ColumnExpr("s.series_id").
		ColumnExpr("s.start_time").
		ColumnExpr("s.team_a_id AS team_a__id").
		ColumnExpr("s.team_b_id AS team_b__id").
		ColumnExpr("s.league_id AS league__id").
		ColumnExpr("team_a.name AS team_a__name").
		ColumnExpr("team_a.tag AS team_a__tag").
		ColumnExpr("team_a.logo_url AS team_a__logo_url").
		ColumnExpr("team_b.name AS team_b__name").
		ColumnExpr("team_b.tag AS team_b__tag").
		ColumnExpr("team_b.logo_url AS team_b__logo_url").
		ColumnExpr("l.name AS league__name").
		ColumnExpr("l.tier AS league__tier").
		ColumnExpr("s.team_a_score").
		ColumnExpr("s.team_b_score").
		TableExpr("series AS s").
		Join("LEFT JOIN leagues AS l USING (league_id)").
		Join("LEFT JOIN teams AS team_a ON team_a.team_id = s.team_a_id").
		Join("LEFT JOIN teams AS team_b ON team_b.team_id = s.team_b_id").
		Where("s.series_id = ?", id)

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
