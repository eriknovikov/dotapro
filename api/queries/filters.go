package queries

// Filters query column expressions
const (
	// TeamSearchColumns are the column expressions for team search queries
	TeamSearchColumns = `
		team_id,
		name,
		logo_url
	`

	// LeagueSearchColumns are the column expressions for league search queries
	LeagueSearchColumns = `
		league_id,
		name
	`
)

// Filters query table expressions
const (
	// TeamsTable is the table expression for teams
	TeamsTable = "teams"

	// LeaguesTable is the table expression for leagues
	LeaguesTable = "leagues"
)

// Filters query filters
const (
	// FilterByNameSimilarity filters by name similarity using PostgreSQL's word_similarity operator
	FilterByNameSimilarity = "name %> ?"
)

// Filters query ordering
const (
	// OrderByNameSimilarity orders by word similarity descending
	OrderByNameSimilarity = "word_similarity(name, ?) DESC"
)

// Filters query limits
const (
	// SearchLimit is the limit for search results
	SearchLimit = 10
)
