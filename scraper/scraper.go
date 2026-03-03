package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

// ScraperConfig holds configuration for the scraper
type ScraperConfig struct {
	BatchSize   int // default 800
	MaxRetries  int
	HTTPTimeout time.Duration
	DBTimeout   time.Duration
}

// DefaultScraperConfig returns a sensible default scraper configuration
func DefaultScraperConfig() ScraperConfig {
	return ScraperConfig{
		BatchSize:   800,
		MaxRetries:  3,
		HTTPTimeout: 30 * time.Second,
		DBTimeout:   5 * time.Second,
	}
}

// Global HTTP client with connection pooling and timeout
var httpClient = &http.Client{
	Timeout: DefaultScraperConfig().HTTPTimeout,
	Transport: &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
		DisableKeepAlives:   false,
	},
}

type ResponseMatch struct {
	MatchID float64 `json:"match_id"`
}

type OpendotaResponse struct {
	Rows []ResponseMatch `json:"rows"`
}

type OpendotaMatchResponse struct {
	Rows []json.RawMessage `json:"rows"`
}

var queryBuilder = &QueryBuilder{}
var validator = NewValidator()

// BatchProcessor holds reusable allocations to reduce memory pressure
type BatchProcessor struct {
	// Maps for related entities
	leagues      map[int64]League
	teams        map[int64]Team
	players      map[int64]Player
	seriesMap    map[int64]Series
	seriesScores map[int64]SeriesScore

	// Slices for valid matches
	validMatches       []Match
	validMetadata      []MatchMetadata
	validSeriesMatches []SeriesMatch

	// Slice for parsed matches
	odMatches []ODMatch

	// Slice for match IDs
	ids []int64
}

// NewBatchProcessor creates a new BatchProcessor with pre-allocated capacity
func NewBatchProcessor(config ScraperConfig) *BatchProcessor {
	return &BatchProcessor{
		leagues:            make(map[int64]League, config.BatchSize),
		teams:              make(map[int64]Team, config.BatchSize*2),
		players:            make(map[int64]Player, config.BatchSize*10),
		seriesMap:          make(map[int64]Series, config.BatchSize),
		seriesScores:       make(map[int64]SeriesScore, config.BatchSize),
		validMatches:       make([]Match, 0, config.BatchSize),
		validMetadata:      make([]MatchMetadata, 0, config.BatchSize),
		validSeriesMatches: make([]SeriesMatch, 0, config.BatchSize),
		odMatches:          make([]ODMatch, 0, config.BatchSize),
		ids:                make([]int64, 0, config.BatchSize),
	}
}

// clear resets all maps and slices for reuse
func (bp *BatchProcessor) clear() {
	// Clear maps
	for k := range bp.leagues {
		delete(bp.leagues, k)
	}
	for k := range bp.teams {
		delete(bp.teams, k)
	}
	for k := range bp.players {
		delete(bp.players, k)
	}
	for k := range bp.seriesMap {
		delete(bp.seriesMap, k)
	}
	for k := range bp.seriesScores {
		delete(bp.seriesScores, k)
	}

	// Reset slices
	bp.validMatches = bp.validMatches[:0]
	bp.validMetadata = bp.validMetadata[:0]
	bp.validSeriesMatches = bp.validSeriesMatches[:0]
	bp.odMatches = bp.odMatches[:0]
	bp.ids = bp.ids[:0]
}

// ScrapeMatches is the main entry point for scraping matches
func ScrapeMatches(ctx context.Context, DB *bun.DB, maxBatches int) error {
	config := DefaultScraperConfig()

	// Simple counters for logging
	var matchesInserted int
	var errorCount int

	// Check for context cancellation early
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context cancelled before scraping started: %w", err)
	}

	lastFetchedMatchId, err := fetchLastID(DB)
	if err != nil {
		return fmt.Errorf("error getting last_fetched_match_id: %w", err)
	}

	matchesToFetchLimit := maxBatches * config.BatchSize
	matchIds, err := fetchMatchIDs(ctx, lastFetchedMatchId, matchesToFetchLimit)
	if err != nil {
		errorCount++
		return fmt.Errorf("error fetching match ids: %w", err)
	}

	N := len(matchIds)
	if N == 0 {
		return ErrNoNewMatches
	}

	processor := NewBatchProcessor(config)
	defer processor.clear()

	for i := 0; i < N; i += config.BatchSize {
		// Check for context cancellation before each batch
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("scraping cancelled: %w", err)
		}

		batchNum := i/config.BatchSize + 1
		end := minInt(i+config.BatchSize, N)
		currentBatchIDs := matchIds[i:end]

		if len(currentBatchIDs) == 0 {
			continue
		}

		currentBatchMatches, err := fetchODMatches(ctx, currentBatchIDs)
		if err != nil {
			errorCount++
			return fmt.Errorf("error scraping matches batch %d: %w", batchNum, err)
		}

		if err := processBatch(ctx, currentBatchMatches, DB, processor, config, &matchesInserted, &errorCount); err != nil {
			errorCount++
			return fmt.Errorf("error processing matches batch %d: %w", batchNum, err)
		}

		matchesInserted += len(currentBatchIDs)

		maxID := maxInt64(currentBatchIDs)
		if err := updateLastID(ctx, DB, maxID); err != nil {
			errorCount++
			return fmt.Errorf("error updating last fetched match id: %w", err)
		}

		processor.clear()
	}

	// Log final metrics
	log.Info().
		Int("matches_inserted", matchesInserted).
		Int("error_count", errorCount).
		Msg("scraping complete")

	return nil
}

func minInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func maxInt64(slice []int64) int64 {
	if len(slice) == 0 {
		return 0
	}
	return slices.Max(slice)
}

// makeOpendotaRequestExplorer makes a request to the OpenDota explorer API
func makeOpendotaRequestExplorer(ctx context.Context, query string) (*http.Response, error) {
	encodedQuery := url.PathEscape(query)
	opendotaUrl := fmt.Sprintf("https://api.opendota.com/api/explorer?sql=%s", encodedQuery)

	req, err := http.NewRequestWithContext(ctx, "GET", opendotaUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make OpenDota explorer request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("OpenDota explorer request failed with status: %s", resp.Status)
	}

	return resp, nil
}

// fetchMatchIDs fetches match IDs from OpenDota with retry logic
func fetchMatchIDs(ctx context.Context, lastMatchID int64, limit int) ([]int64, error) {
	retryConfig := RetryConfig{
		MaxAttempts: DefaultScraperConfig().MaxRetries,
		BaseDelay:   500 * time.Millisecond,
		MaxDelay:    10 * time.Second,
		Multiplier:  2.0,
	}

	var ids []int64
	var fetchErr error

	err := RetryWithBackoff(ctx, retryConfig, func(ctx context.Context) error {
		resp, err := makeOpendotaRequestExplorer(ctx, queryBuilder.GetIds(lastMatchID, limit))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var opendotaResp OpendotaResponse
		if err := json.NewDecoder(resp.Body).Decode(&opendotaResp); err != nil {
			return fmt.Errorf("failed to decode OpenDota response: %w", err)
		}

		ids = make([]int64, 0, len(opendotaResp.Rows))
		for _, m := range opendotaResp.Rows {
			ids = append(ids, int64(m.MatchID))
		}

		return nil
	})

	if err != nil {
		fetchErr = fmt.Errorf("failed to fetch match IDs after retries: %w", err)
		return nil, fetchErr
	}

	return ids, nil
}

// fetchODMatches fetches match details from OpenDota with retry logic
func fetchODMatches(ctx context.Context, matchIDs []int64) ([]json.RawMessage, error) {
	retryConfig := RetryConfig{
		MaxAttempts: DefaultScraperConfig().MaxRetries,
		BaseDelay:   500 * time.Millisecond,
		MaxDelay:    10 * time.Second,
		Multiplier:  2.0,
	}

	var matches []json.RawMessage
	var fetchErr error

	err := RetryWithBackoff(ctx, retryConfig, func(ctx context.Context) error {
		resp, err := makeOpendotaRequestExplorer(ctx, queryBuilder.GetMatches(matchIDs))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var opendotaMatchResp OpendotaMatchResponse
		if err := json.NewDecoder(resp.Body).Decode(&opendotaMatchResp); err != nil {
			return fmt.Errorf("failed to decode OpenDota match response: %w", err)
		}

		matches = opendotaMatchResp.Rows
		return nil
	})

	if err != nil {
		fetchErr = fmt.Errorf("failed to fetch match details after retries: %w", err)
		return nil, fetchErr
	}

	return matches, nil
}

// parseRawMatches parses raw JSON messages into ODMatch structs
func parseRawMatches(rawBatch []json.RawMessage, odMatches []ODMatch) ([]ODMatch, error) {
	odMatches = odMatches[:0]
	for i, raw := range rawBatch {
		var m ODMatch
		if err := json.Unmarshal(raw, &m); err != nil {
			return nil, fmt.Errorf("failed to unmarshal match data at index %d: %w", i, err)
		}
		odMatches = append(odMatches, m)
	}
	return odMatches, nil
}

// extractRelatedEntities extracts leagues, teams, players, and series from matches
func extractRelatedEntities(odMatches []ODMatch, leagues map[int64]League, teams map[int64]Team, players map[int64]Player, seriesMap map[int64]Series) {
	for _, m := range odMatches {
		if m.League.ID > 0 {
			leagues[m.League.ID] = League{
				LeagueID: m.League.ID,
				Name:     m.League.Name,
				Tier:     m.League.Tier,
			}
		}

		if m.RadiantTeam.ID > 0 {
			teams[m.RadiantTeam.ID] = Team{
				TeamID:  m.RadiantTeam.ID,
				Name:    m.RadiantTeam.Name,
				Tag:     m.RadiantTeam.Tag,
				LogoURL: m.RadiantTeam.LogoURL,
			}
		}

		if m.DireTeam.ID > 0 {
			teams[m.DireTeam.ID] = Team{
				TeamID:  m.DireTeam.ID,
				Name:    m.DireTeam.Name,
				Tag:     m.DireTeam.Tag,
				LogoURL: m.DireTeam.LogoURL,
			}
		}

		if m.SeriesID > 0 {
			s := Series{
				SeriesID:  m.SeriesID,
				StartTime: time.Unix(m.StartTime, 0),
			}
			if m.League.ID > 0 {
				s.LeagueID = m.League.ID
			}
			if m.RadiantTeam.ID > 0 {
				s.TeamAID = m.RadiantTeam.ID
			}
			if m.DireTeam.ID > 0 {
				s.TeamBID = m.DireTeam.ID
			}
			seriesMap[m.SeriesID] = s
		}

		var matchPlayers []ODPlayerShort
		if err := json.Unmarshal(m.Players, &matchPlayers); err != nil {
			continue
		}

		for _, p := range matchPlayers {
			if p.PlayerID > 0 && p.Name != "" {
				if _, exists := players[p.PlayerID]; !exists {
					players[p.PlayerID] = Player{
						PlayerID: p.PlayerID,
						Name:     p.Name,
					}
				}
			}
		}
	}
}

// buildMatchEntities builds Match, MatchMetadata, and SeriesMatch entities from ODMatch
func buildMatchEntities(om ODMatch, players map[int64]Player) (Match, MatchMetadata, SeriesMatch) {
	m := Match{
		MatchID:    om.MatchID,
		Duration:   om.Duration,
		StartTime:  time.Unix(om.StartTime, 0),
		RadiantWin: om.RadiantWin,
		Patch:      om.Patch,
	}
	if om.League.ID > 0 {
		m.LeagueID = om.League.ID
	}
	if om.RadiantTeam.ID > 0 {
		m.RadiantTeamID = om.RadiantTeam.ID
	}
	if om.DireTeam.ID > 0 {
		m.DireTeamID = om.DireTeam.ID
	}

	var matchPlayers []ODPlayerShort
	json.Unmarshal(om.Players, &matchPlayers)

	for _, p := range matchPlayers {
		if p.PlayerSlot < 128 {
			m.RadiantHeroes = append(m.RadiantHeroes, p.HeroID)
			m.RadiantPlayers = append(m.RadiantPlayers, p.PlayerID)
		} else {
			m.DireHeroes = append(m.DireHeroes, p.HeroID)
			m.DirePlayers = append(m.DirePlayers, p.PlayerID)
		}
	}

	md := MatchMetadata{
		MatchID:        om.MatchID,
		PicksBans:      om.PicksBans,
		PlayersData:    om.Players,
		RadiantGoldAdv: om.RadiantGoldAdv,
		RadiantXPAdv:   om.RadiantXPAdv,
		RadiantScore:   om.RadiantTeam.Score,
		DireScore:      om.DireTeam.Score,
		Version:        om.Version,
	}
	if om.RadiantTeam.Captain != nil {
		if _, exists := players[*om.RadiantTeam.Captain]; exists {
			md.RadiantCaptain = om.RadiantTeam.Captain
		}
	}
	if om.DireTeam.Captain != nil {
		if _, exists := players[*om.DireTeam.Captain]; exists {
			md.DireCaptain = om.DireTeam.Captain
		}
	}

	var sm SeriesMatch
	if om.SeriesID > 0 {
		sm = SeriesMatch{
			SeriesID: om.SeriesID,
			MatchID:  om.MatchID,
		}
	}

	return m, md, sm
}

// insertRelatedEntities inserts leagues, teams, players, and series into the database
func insertRelatedEntities(ctx context.Context, tx bun.Tx, leagues map[int64]League, teams map[int64]Team, players map[int64]Player, seriesMap map[int64]Series) error {
	if len(leagues) > 0 {
		lSlice := mapsToSlice(leagues)
		// Validate leagues before insertion
		for i := range lSlice {
			if err := validator.ValidateLeague(lSlice[i]); err != nil {
				continue
			}
		}
		if _, err := tx.NewInsert().Model(&lSlice).On("CONFLICT (league_id) DO NOTHING").Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert leagues: %w", err)
		}
	}

	if len(teams) > 0 {
		tSlice := mapsToSlice(teams)
		// Validate teams before insertion
		for i := range tSlice {
			if err := validator.ValidateTeam(tSlice[i]); err != nil {
				continue
			}
		}
		if _, err := tx.NewInsert().Model(&tSlice).On("CONFLICT (team_id) DO NOTHING").Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert teams: %w", err)
		}
	}

	if len(players) > 0 {
		pSlice := mapsToSlice(players)
		// Validate players before insertion
		for i := range pSlice {
			if err := validator.ValidatePlayer(pSlice[i]); err != nil {
				continue
			}
		}
		if _, err := tx.NewInsert().Model(&pSlice).On("CONFLICT (player_id) DO NOTHING").Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert players: %w", err)
		}
	}

	if len(seriesMap) > 0 {
		sSlice := mapsToSlice(seriesMap)
		// Validate series before insertion
		for i := range sSlice {
			if err := validator.ValidateSeries(sSlice[i]); err != nil {
				continue
			}
		}
		if _, err := tx.NewInsert().Model(&sSlice).On("CONFLICT (series_id) DO NOTHING").Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert series: %w", err)
		}
	}

	return nil
}

// insertBatch inserts matches, metadata, and series matches into the database
func insertBatch(ctx context.Context, tx bun.Tx, matches []Match, metadata []MatchMetadata, seriesMatches []SeriesMatch) error {
	// Validate matches and metadata before insertion
	validCount, _, _ := validator.ValidateBatch(matches, metadata)

	if validCount == 0 {
		return nil
	}

	if _, err := tx.NewInsert().Model(&matches).On("CONFLICT (match_id) DO NOTHING").Exec(ctx); err != nil {
		return fmt.Errorf("failed to insert matches: %w", err)
	}

	if _, err := tx.NewInsert().Model(&metadata).On("CONFLICT (match_id) DO NOTHING").Exec(ctx); err != nil {
		return fmt.Errorf("failed to insert match metadata: %w", err)
	}

	if len(seriesMatches) > 0 {
		if _, err := tx.NewInsert().Model(&seriesMatches).On("CONFLICT (series_id, match_id) DO NOTHING").Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert series matches: %w", err)
		}
	}

	return nil
}

// calculateSeriesScores calculates series scores from match results
func calculateSeriesScores(odMatches []ODMatch, scores map[int64]SeriesScore) {
	for _, m := range odMatches {
		if m.SeriesID == 0 {
			continue
		}

		score, exists := scores[m.SeriesID]
		if !exists {
			score = SeriesScore{SeriesID: m.SeriesID}
		}

		var winningTeamID int64
		if m.RadiantWin {
			winningTeamID = m.RadiantTeam.ID
		} else {
			winningTeamID = m.DireTeam.ID
		}

		switch winningTeamID {
		case m.RadiantTeam.ID:
			score.TeamAWins++
		case m.DireTeam.ID:
			score.TeamBWins++
		default:
		}

		scores[m.SeriesID] = score
	}
}

// updateSeriesScores updates series scores in the database
func updateSeriesScores(ctx context.Context, tx bun.Tx, scores map[int64]SeriesScore) error {
	if len(scores) == 0 {
		return nil
	}

	for seriesID, score := range scores {
		_, err := tx.NewUpdate().
			Model(&Series{}).
			Set("team_a_score = team_a_score + ?", score.TeamAWins).
			Set("team_b_score = team_b_score + ?", score.TeamBWins).
			Where("series_id = ?", seriesID).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to update series score for series_id %d: %w", seriesID, err)
		}
	}

	return nil
}

// processBatch processes a batch of raw match data
func processBatch(ctx context.Context, rawBatch []json.RawMessage, db *bun.DB, bp *BatchProcessor, config ScraperConfig, matchesInserted *int, errorCount *int) error {
	odMatches, err := parseRawMatches(rawBatch, bp.odMatches)
	if err != nil {
		return fmt.Errorf("failed to parse raw matches: %w", err)
	}

	// Validate OD matches before processing
	for _, m := range odMatches {
		if err := validator.ValidateODMatch(m); err != nil {
			continue
		}
	}

	extractRelatedEntities(odMatches, bp.leagues, bp.teams, bp.players, bp.seriesMap)

	bp.validMatches = bp.validMatches[:0]
	bp.validMetadata = bp.validMetadata[:0]
	bp.validSeriesMatches = bp.validSeriesMatches[:0]

	for _, om := range odMatches {
		m, md, sm := buildMatchEntities(om, bp.players)
		bp.validMatches = append(bp.validMatches, m)
		bp.validMetadata = append(bp.validMetadata, md)
		if sm.SeriesID > 0 {
			bp.validSeriesMatches = append(bp.validSeriesMatches, sm)
		}
	}

	if len(bp.validMatches) == 0 {
		return nil
	}

	err = db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := insertRelatedEntities(ctx, tx, bp.leagues, bp.teams, bp.players, bp.seriesMap); err != nil {
			return err
		}

		if err := insertBatch(ctx, tx, bp.validMatches, bp.validMetadata, bp.validSeriesMatches); err != nil {
			return err
		}

		calculateSeriesScores(odMatches, bp.seriesScores)
		if err := updateSeriesScores(ctx, tx, bp.seriesScores); err != nil {
			return err
		}

		return nil
	})

	// Clear series scores after processing to free memory
	for k := range bp.seriesScores {
		delete(bp.seriesScores, k)
	}

	return err
}

// mapsToSlice converts a map to a slice of values
func mapsToSlice[K comparable, V any](m map[K]V) []V {
	s := make([]V, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}
