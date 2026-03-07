package filtersmetadata

import (
	"context"
	"dotapro-lambda-api/constants"
	"dotapro-lambda-api/queries"
	"dotapro-lambda-api/types"

	"github.com/uptrace/bun"
)

type Model struct {
	DB *bun.DB
}

func NewModel(db *bun.DB) *Model { return &Model{DB: db} }

// SearchTeams searches for teams by name using PostgreSQL's word_similarity operator.
// Returns up to SearchLimit results ordered by similarity.
func (m *Model) SearchTeams(ctx context.Context, query string) ([]types.TeamSearchResult, error) {
	var results []types.TeamSearchResult

	err := m.DB.NewSelect().
		ColumnExpr(queries.TeamSearchColumns).
		TableExpr(queries.TeamsTable).
		Where(queries.FilterByNameSimilarity, query).
		OrderExpr(queries.OrderByNameSimilarity, query).
		Limit(constants.SearchLimit).
		Scan(ctx, &results)

	return results, err
}

// SearchLeagues searches for leagues by name using PostgreSQL's word_similarity operator.
// Returns up to SearchLimit results ordered by similarity.
func (m *Model) SearchLeagues(ctx context.Context, query string) ([]types.LeagueSearchResult, error) {
	var results []types.LeagueSearchResult

	err := m.DB.NewSelect().
		ColumnExpr(queries.LeagueSearchColumns).
		TableExpr(queries.LeaguesTable).
		Where(queries.FilterByNameSimilarity, query).
		OrderExpr(queries.OrderByNameSimilarity, query).
		Limit(constants.SearchLimit).
		Scan(ctx, &results)

	return results, err
}

// GetTeamName retrieves the name of a team by its ID
func (m *Model) GetTeamName(ctx context.Context, id int64) (map[string]string, error) {
	var result struct {
		Name string `bun:"name"`
	}
	
	err := m.DB.NewSelect().
		Column("name").
		TableExpr("teams").
		Where("team_id = ?", id).
		Scan(ctx, &result)
	
	if err != nil {
		return nil, err
	}
	
	return map[string]string{"name": result.Name}, nil
}

// GetLeagueName retrieves the name of a league by its ID
func (m *Model) GetLeagueName(ctx context.Context, id int64) (map[string]string, error) {
	var result struct {
		Name string `bun:"name"`
	}
	
	err := m.DB.NewSelect().
		Column("name").
		TableExpr("leagues").
		Where("league_id = ?", id).
		Scan(ctx, &result)
	
	if err != nil {
		return nil, err
	}
	
	return map[string]string{"name": result.Name}, nil
}
