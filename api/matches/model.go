package matches

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

func (m *Model) GetMany(ctx context.Context, filter types.GetMatchesFilter) ([]types.MatchSummary, types.PaginationData, error) {
	var res []types.MatchSummary

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	q := m.DB.NewSelect().
		ColumnExpr("m.match_id").
		ColumnExpr("m.start_time").
		ColumnExpr("m.duration").
		ColumnExpr("m.radiant_win").
		ColumnExpr("sm.series_id").
		ColumnExpr("m.radiant_team_id AS radiant_team__id").
		ColumnExpr("m.dire_team_id AS dire_team__id").
		ColumnExpr("m.league_id AS league__id").
		ColumnExpr("radiant.name AS radiant_team__name").
		ColumnExpr("radiant.tag AS radiant_team__tag").
		ColumnExpr("radiant.logo_url AS radiant_team__logo_url").
		ColumnExpr("dire.name AS dire_team__name").
		ColumnExpr("dire.tag AS dire_team__tag").
		ColumnExpr("dire.logo_url AS dire_team__logo_url").
		ColumnExpr("l.name AS league__name").
		ColumnExpr("l.tier AS league__tier").
		ColumnExpr("m.radiant_heroes").
		ColumnExpr("m.dire_heroes").
		TableExpr("matches AS m").
		Join("LEFT JOIN matches_metadata AS md USING (match_id)").
		Join("LEFT JOIN series_match AS sm USING (match_id)").
		Join("LEFT JOIN leagues AS l USING (league_id)").
		Join("LEFT JOIN teams AS radiant ON radiant.team_id = m.radiant_team_id").
		Join("LEFT JOIN teams AS dire ON dire.team_id = m.dire_team_id")

	if filter.LeagueID != nil {
		q = q.Where("m.league_id = ?", *filter.LeagueID)
	}
	if filter.TeamID != nil {
		q = q.Where("m.radiant_team_id = ? OR m.dire_team_id = ?", *filter.TeamID, *filter.TeamID)
	}
	if filter.PlayerID != nil {
		q = q.Where("m.radiant_players @> ? OR m.dire_players @> ?", []int64{*filter.PlayerID}, []int64{*filter.PlayerID})
	}
	if filter.HeroID != nil {
		q = q.Where("m.radiant_heroes @> ? OR m.dire_heroes @> ?", []int64{*filter.HeroID}, []int64{*filter.HeroID})
	}

	switch filter.Sort {
	case "oldest":
		q = q.Order("m.match_id ASC")
	default:
		q = q.Order("m.match_id DESC")
	}

	q = q.Limit(limit).Offset(offset)

	paginationQuery := m.DB.NewSelect().
		ColumnExpr("COUNT(*)").
		TableExpr("matches AS m").
		Join("LEFT JOIN matches_metadata AS md USING (match_id)").
		Join("LEFT JOIN series_match AS sm USING (match_id)").
		Join("LEFT JOIN leagues AS l USING (league_id)").
		Join("LEFT JOIN teams AS radiant ON radiant.team_id = m.radiant_team_id").
		Join("LEFT JOIN teams AS dire ON dire.team_id = m.dire_team_id")

	if filter.LeagueID != nil {
		paginationQuery = paginationQuery.Where("m.league_id = ?", *filter.LeagueID)
	}
	if filter.TeamID != nil {
		paginationQuery = paginationQuery.Where("m.radiant_team_id = ? OR m.dire_team_id = ?", *filter.TeamID, *filter.TeamID)
	}
	if filter.PlayerID != nil {
		paginationQuery = paginationQuery.Where("m.radiant_players @> ? OR m.dire_players @> ?", []int64{*filter.PlayerID}, []int64{*filter.PlayerID})
	}
	if filter.HeroID != nil {
		paginationQuery = paginationQuery.Where("m.radiant_heroes @> ? OR m.dire_heroes @> ?", []int64{*filter.HeroID}, []int64{*filter.HeroID})
	}

	if err := q.Scan(ctx, &res); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, types.PaginationData{}, errs.NOT_FOUND
		}
		return nil, types.PaginationData{}, err
	}

	var totalCount int
	if err := paginationQuery.Scan(ctx, &totalCount); err != nil {
		return nil, types.PaginationData{}, err
	}

	totalPages := totalCount / limit
	if totalCount%limit != 0 {
		totalPages++
	}

	return res, types.PaginationData{
		PageSize:    limit,
		TotalPages:  totalPages,
		CurrentPage: page,
	}, nil
}

func (m *Model) GetOne(ctx context.Context, id int64) (*types.MatchDetail, error) {
	res := &types.MatchDetail{}
	q := m.DB.NewSelect().
		ColumnExpr("m.match_id").
		ColumnExpr("m.start_time").
		ColumnExpr("m.duration").
		ColumnExpr("m.radiant_win").
		ColumnExpr("m.patch").
		ColumnExpr("md.version").
		ColumnExpr("sm.series_id").
		ColumnExpr("md.picks_bans").
		ColumnExpr("md.players_data").
		ColumnExpr("md.radiant_gold_adv").
		ColumnExpr("md.radiant_xp_adv").
		ColumnExpr("m.radiant_team_id AS radiant_team__id").
		ColumnExpr("m.dire_team_id AS dire_team__id").
		ColumnExpr("m.league_id AS league__id").
		ColumnExpr("radiant.name AS radiant_team__name").
		ColumnExpr("radiant.tag AS radiant_team__tag").
		ColumnExpr("radiant.logo_url AS radiant_team__logo_url").
		ColumnExpr("md.radiant_score AS radiant_team__score").
		ColumnExpr("md.radiant_captain AS radiant_team__captain").
		ColumnExpr("dire.name AS dire_team__name").
		ColumnExpr("dire.tag AS dire_team__tag").
		ColumnExpr("dire.logo_url AS dire_team__logo_url").
		ColumnExpr("md.dire_score AS dire_team__score").
		ColumnExpr("md.dire_captain AS dire_team__captain").
		ColumnExpr("l.name AS league__name").
		ColumnExpr("l.tier AS league__tier").
		TableExpr("matches AS m").
		Join("LEFT JOIN matches_metadata AS md USING (match_id)").
		Join("LEFT JOIN series_match AS sm USING (match_id)").
		Join("LEFT JOIN leagues AS l USING (league_id)").
		Join("LEFT JOIN teams AS radiant ON radiant.team_id = m.radiant_team_id").
		Join("LEFT JOIN teams AS dire ON dire.team_id = m.dire_team_id").
		Where("m.match_id = ?", id)

	log.Debug().Str("BUN_QUERY", q.String()).Send()
	if err := q.Scan(ctx, res); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NOT_FOUND
		}
		return nil, err
	}
	return res, nil
}
