package queries

// Matches query column expressions
const (
	// MatchColumns are the column expressions for match queries
	MatchColumns = `
		m.match_id,
		m.start_time,
		m.duration,
		m.radiant_win,
		m.patch,
		md.version,
		sm.series_id,
		md.picks_bans,
		md.players_data,
		md.radiant_gold_adv,
		md.radiant_xp_adv,
		m.radiant_team_id AS radiant_team__id,
		m.dire_team_id AS dire_team__id,
		m.league_id AS league__id,
		radiant.name AS radiant_team__name,
		radiant.tag AS radiant_team__tag,
		radiant.logo_url AS radiant_team__logo_url,
		md.radiant_score AS radiant_team__score,
		md.radiant_captain AS radiant_team__captain,
		dire.name AS dire_team__name,
		dire.tag AS dire_team__tag,
		dire.logo_url AS dire_team__logo_url,
		md.dire_score AS dire_team__score,
		md.dire_captain AS dire_team__captain,
		l.name AS league__name,
		l.tier AS league__tier
	`

	// MatchSummaryColumns are the column expressions for match summary queries
	MatchSummaryColumns = `
		m.match_id,
		m.start_time,
		m.duration,
		m.radiant_win,
		sm.series_id,
		m.radiant_team_id AS radiant_team__id,
		m.dire_team_id AS dire_team__id,
		m.league_id AS league__id,
		radiant.name AS radiant_team__name,
		radiant.tag AS radiant_team__tag,
		radiant.logo_url AS radiant_team__logo_url,
		md.radiant_captain AS radiant_team__captain,
		dire.name AS dire_team__name,
		dire.tag AS dire_team__tag,
		dire.logo_url AS dire_team__logo_url,
		md.dire_captain AS dire_team__captain,
		l.name AS league__name,
		l.tier AS league__tier,
		m.radiant_heroes,
		m.dire_heroes
	`
)

// Matches query joins
const (
	// MatchJoins are the JOIN clauses for match queries
	MatchJoins = `
		LEFT JOIN matches_metadata AS md USING (match_id)
		LEFT JOIN series_match AS sm USING (match_id)
		LEFT JOIN leagues AS l USING (league_id)
		LEFT JOIN teams AS radiant ON radiant.team_id = m.radiant_team_id
		LEFT JOIN teams AS dire ON dire.team_id = m.dire_team_id
	`
)

// Matches query table expressions
const (
	// MatchesTable is the table expression for matches
	MatchesTable = "matches AS m"
)

// Matches query filters
const (
	// FilterByLeagueID filters matches by league ID
	FilterByLeagueID = "m.league_id = ?"
	
	// FilterByTeamID filters matches by team ID (either radiant or dire)
	FilterByTeamID = "m.radiant_team_id = ? OR m.dire_team_id = ?"
	
	// FilterByPlayerID filters matches by player ID (either radiant or dire)
	FilterByPlayerID = "m.radiant_players @> ? OR m.dire_players @> ?"
	
	// FilterByHeroID filters matches by hero ID (either radiant or dire)
	FilterByHeroID = "m.radiant_heroes @> ? OR m.dire_heroes @> ?"
	
	// FilterByMatchID filters matches by match ID
	FilterByMatchID = "m.match_id = ?"
)

// Matches query ordering
const (
	// OrderByMatchIDDesc orders matches by match ID descending (newest first)
	OrderByMatchIDDesc = "m.match_id DESC"
	
	// OrderByMatchIDAsc orders matches by match ID ascending (oldest first)
	OrderByMatchIDAsc = "m.match_id ASC"
)

// Cursor filters
const (
	// CursorBeforeMatchID filters for matches before a given cursor (newest sort)
	CursorBeforeMatchID = "m.match_id < ?"
	
	// CursorAfterMatchID filters for matches after a given cursor (oldest sort)
	CursorAfterMatchID = "m.match_id > ?"
)
