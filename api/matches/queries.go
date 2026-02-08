package matches

import (
	"github.com/uptrace/bun"
	"dotapro-lambda-api/types"
)

var EMPTY_MATCH_DETAIL *types.MatchDetail

func (m *Model) buildQueryGetOne(id int64) *bun.SelectQuery {
	return m.DB.NewSelect().
		Model(EMPTY_MATCH_DETAIL).
		ColumnExpr("m.match_id").
		ColumnExpr("m.start_time").
		ColumnExpr("m.duration").
		ColumnExpr("m.radiant_win").
		ColumnExpr("m.patch").
		ColumnExpr("md.version").
		ColumnExpr("m.radiant_heroes").
		ColumnExpr("m.dire_heroes").
		ColumnExpr("m.radiant_players").
		ColumnExpr("m.dire_players").
		ColumnExpr("md.series_id").
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
		Join("LEFT JOIN matches_metadata AS md ON md.match_id = m.match_id").
		Join("LEFT JOIN leagues AS l ON l.league_id = m.league_id").
		Join("LEFT JOIN teams AS radiant ON radiant.team_id = m.radiant_team_id").
		Join("LEFT JOIN teams AS dire ON dire.team_id = m.dire_team_id").
		Where("m.match_id = ?", id)
}
