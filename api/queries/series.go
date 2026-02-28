package queries

// Series query column expressions
const (
	// SeriesColumns are the column expressions for series queries
	SeriesColumns = `
		s.series_id,
		s.start_time,
		s.team_a_id AS team_a__id,
		s.team_b_id AS team_b__id,
		s.league_id AS league__id,
		team_a.name AS team_a__name,
		team_a.tag AS team_a__tag,
		team_a.logo_url AS team_a__logo_url,
		team_b.name AS team_b__name,
		team_b.tag AS team_b__tag,
		team_b.logo_url AS team_b__logo_url,
		l.name AS league__name,
		l.tier AS league__tier,
		s.team_a_score,
		s.team_b_score
	`
)

// Series query joins
const (
	// SeriesJoins are the JOIN clauses for series queries
	SeriesJoins = `
		LEFT JOIN leagues AS l USING (league_id)
		LEFT JOIN teams AS team_a ON team_a.team_id = s.team_a_id
		LEFT JOIN teams AS team_b ON team_b.team_id = s.team_b_id
	`
)

// Series query table expressions
const (
	// SeriesTable is the table expression for series
	SeriesTable = "series AS s"
)

// Series query filters
const (
	// FilterSeriesByLeagueID filters series by league ID
	FilterSeriesByLeagueID = "s.league_id = ?"
	
	// FilterSeriesByTeamID filters series by team ID (either team A or team B)
	FilterSeriesByTeamID = "s.team_a_id = ? OR s.team_b_id = ?"
	
	// FilterSeriesBySeriesID filters series by series ID
	FilterSeriesBySeriesID = "s.series_id = ?"
)

// Series query ordering
const (
	// OrderBySeriesIDDesc orders series by series ID descending (newest first)
	OrderBySeriesIDDesc = "s.series_id DESC"
	
	// OrderBySeriesIDAsc orders series by series ID ascending (oldest first)
	OrderBySeriesIDAsc = "s.series_id ASC"
)

// Cursor filters
const (
	// CursorBeforeSeriesID filters for series before a given cursor (newest sort)
	CursorBeforeSeriesID = "s.series_id < ?"
	
	// CursorAfterSeriesID filters for series after a given cursor (oldest sort)
	CursorAfterSeriesID = "s.series_id > ?"
)

// Series match query column expressions
const (
	// SeriesMatchColumns are the column expressions for series match queries
	SeriesMatchColumns = `
		m.match_id,
		m.duration,
		m.radiant_win,
		md.picks_bans,
		md.players_data,
		md.radiant_gold_adv,
		md.radiant_xp_adv,
		md.radiant_score,
		md.dire_score,
		md.radiant_captain,
		md.dire_captain
	`
)

// Series match query joins
const (
	// SeriesMatchJoins are the JOIN clauses for series match queries
	SeriesMatchJoins = `
		LEFT JOIN matches_metadata AS md USING (match_id)
		LEFT JOIN series_match AS sm USING (match_id)
	`
)

// Series match query table expressions
const (
	// SeriesMatchTable is the table expression for series matches
	SeriesMatchTable = "matches AS m"
)

// Series match query filters
const (
	// FilterMatchesBySeriesID filters matches by series ID
	FilterMatchesBySeriesID = "sm.series_id = ?"
)

// Series match query ordering
const (
	// OrderMatchesByMatchIDAsc orders matches by match ID ascending
	OrderMatchesByMatchIDAsc = "m.match_id ASC"
)
