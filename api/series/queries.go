package series

import (
	"dotapro-lambda-api/types"

	"github.com/uptrace/bun"
)

func (m *Model) buildQueryGetMany(filter types.GetSeriesFilter) *bun.SelectQuery {
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
	default:
		q = q.Order("s.series_id DESC")
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	return q.Limit(limit).Offset(offset)
}

func (m *Model) buildQueryGetOne(id int64) *bun.SelectQuery {
	return m.DB.NewSelect().
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
}

func (m *Model) buildQueryGetMatchesForSeries(seriesID int64) *bun.SelectQuery {
	return m.DB.NewSelect().
		ColumnExpr("m.match_id").
		ColumnExpr("m.start_time").
		ColumnExpr("m.duration").
		ColumnExpr("m.radiant_win").
		ColumnExpr("m.patch").
		ColumnExpr("md.version").
		ColumnExpr("md.picks_bans").
		ColumnExpr("md.players_data").
		ColumnExpr("md.radiant_gold_adv").
		ColumnExpr("md.radiant_xp_adv").
		ColumnExpr("m.radiant_heroes").
		ColumnExpr("m.dire_heroes").
		ColumnExpr("m.radiant_players").
		ColumnExpr("m.dire_players").
		ColumnExpr("md.radiant_score").
		ColumnExpr("md.dire_score").
		ColumnExpr("md.radiant_captain").
		ColumnExpr("md.dire_captain").
		TableExpr("matches AS m").
		Join("LEFT JOIN matches_metadata AS md USING (match_id)").
		Where("md.series_id = ?", seriesID).
		Order("m.match_id ASC")
}
