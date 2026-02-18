package filtersmetadata

import (
	"context"

	"dotapro-lambda-api/types"

	"github.com/uptrace/bun"
)

type Model struct {
	DB *bun.DB
}

func NewModel(db *bun.DB) *Model { return &Model{DB: db} }

func (m *Model) SearchTeams(ctx context.Context, query string) ([]types.TeamSearchResult, error) {
	var results []types.TeamSearchResult

	err := m.DB.NewSelect().
		ColumnExpr("team_id").
		ColumnExpr("name").
		ColumnExpr("logo_url").
		TableExpr("teams").
		Where("name %> ?", query).
		OrderExpr("word_similarity(name, ?) DESC", query).
		Limit(10).
		Scan(ctx, &results)

	return results, err
}

func (m *Model) SearchLeagues(ctx context.Context, query string) ([]types.LeagueSearchResult, error) {
	var results []types.LeagueSearchResult

	err := m.DB.NewSelect().
		ColumnExpr("league_id").
		ColumnExpr("name").
		TableExpr("leagues").
		Where("name %> ?", query).
		OrderExpr("word_similarity(name, ?) DESC", query).
		Limit(10).
		Scan(ctx, &results)

	return results, err
}
