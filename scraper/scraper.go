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

const (
	BATCH_SIZE     = 800
	TIMES_TO_RETRY = 3
)

func ScrapeMatches(ctx context.Context, DB *bun.DB, limit int) error {
	var err error
	lastFetchedMatchId, err := fetchLastID(DB)
	if err != nil {
		return fmt.Errorf("err getting last_fetched_match_id: %w", err)
	}
	matchIds, err := fetchMatchIDs(lastFetchedMatchId, limit)
	if err != nil {
		return fmt.Errorf("err fetching match ids < last_fetched_matched_id: %w", err)
	}
	N := len(matchIds)
	if N == 0 {
		return ErrNoNewMatches
	}
	currentBatchIDs := make([]int64, 0, BATCH_SIZE)
	currentBatchMatches := make([]json.RawMessage, 0, BATCH_SIZE)
	for i := 0; i < N; i += BATCH_SIZE {
		log.Warn().Int("batchNum", i).Msg("processing batch")
		end := minInt(i+BATCH_SIZE, N)
		currentBatchIDs = matchIds[i:end]
		currentBatchMatches, err = fetchMatchBatch(currentBatchIDs)
		if err != nil {
			return fmt.Errorf("error scraping matches batch: %w", err)
		}
		odMatches, err := parseRawMatches(currentBatchMatches)
		if err != nil {
			return fmt.Errorf("failed to parse raw matches: %w", err)
		}

		if err := insertMatchBatch(odMatches, DB); err != nil {
			return fmt.Errorf("error inserting matches batch: %w", err)
		}
		if len(currentBatchIDs) == 0 {
			log.Warn().Msg("DIDNT EXPECT CURRENT BATCH ID LENGTH TO BE 0. LOGICAL ERROR, AND SHOULD NOT HAPPEN")
		}

		if len(currentBatchIDs) > 0 {
			maxID := maxInt64(currentBatchIDs)
			if err := updateLastID(ctx, DB, maxID); err != nil {
				return fmt.Errorf("error updating last fetched match id: %w", err)
			}
		}
		log.Warn().Int("batchNum", i).Msg("finished processing batch")
	}
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

func makeOpendotaRequestExplorer(query string) (*http.Response, error) {
	encodedQuery := url.PathEscape(query)
	opendotaUrl := fmt.Sprintf("https://api.opendota.com/api/explorer?sql=%s", encodedQuery)
	resp, err := http.Get(opendotaUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to make OpenDota explorer request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("OpenDota explorer request failed with status: %s", resp.Status)
	}
	return resp, nil
}

func fetchMatchIDs(lastMatchID int64, limit int) ([]int64, error) {
	resp, err := makeOpendotaRequestExplorer(queryBuilder.GetIds(lastMatchID, limit))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var opendotaResp OpendotaResponse
	if err := json.NewDecoder(resp.Body).Decode(&opendotaResp); err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(opendotaResp.Rows))
	for _, m := range opendotaResp.Rows {
		ids = append(ids, int64(m.MatchID))
	}
	return ids, nil
}

func fetchMatchBatch(matchIDs []int64) ([]json.RawMessage, error) {
	log.Warn().Int("total ids in this batch", len(matchIDs)).Msg("fetching match ids")
	resp, err := makeOpendotaRequestExplorer(queryBuilder.GetMatches(matchIDs))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var opendotaMatchResp OpendotaMatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&opendotaMatchResp); err != nil {
		return nil, fmt.Errorf("failed to decode OpenDota explorer response: %w", err)
	}

	return opendotaMatchResp.Rows, nil
}

func parseRawMatches(rawBatch []json.RawMessage) ([]ODMatch, error) {
	odMatches := make([]ODMatch, 0, len(rawBatch))
	for _, raw := range rawBatch {
		var m ODMatch
		if err := json.Unmarshal(raw, &m); err != nil {
			return nil, fmt.Errorf("failed to unmarshal match data: %w", err)
		}
		odMatches = append(odMatches, m)
	}
	return odMatches, nil
}

func extractRelatedEntities(odMatches []ODMatch) (map[int64]League, map[int64]Team, map[int64]Player, map[int64]Series) {
	leagues := make(map[int64]League)
	teams := make(map[int64]Team)
	players := make(map[int64]Player)
	seriesMap := make(map[int64]Series)

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
			log.Printf("Warning: failed to unmarshal players for match %d: %v", m.MatchID, err)
			continue
		}

		for _, p := range matchPlayers {
			if p.PlayerID > 0 {
				if _, exists := players[p.PlayerID]; !exists {
					players[p.PlayerID] = Player{
						PlayerID: p.PlayerID,
						Name:     p.Name,
					}
				}
			}
		}

		if m.RadiantTeam.Captain > 0 {
			if _, exists := players[m.RadiantTeam.Captain]; !exists {
				players[m.RadiantTeam.Captain] = Player{PlayerID: m.RadiantTeam.Captain}
			}
		}
		if m.DireTeam.Captain > 0 {
			if _, exists := players[m.DireTeam.Captain]; !exists {
				players[m.DireTeam.Captain] = Player{PlayerID: m.DireTeam.Captain}
			}
		}
	}

	return leagues, teams, players, seriesMap
}

func validateMatch(om ODMatch) error {
	var matchPlayers []ODPlayerShort
	if err := json.Unmarshal(om.Players, &matchPlayers); err != nil {
		return fmt.Errorf("failed to unmarshal players for match %d: %w", om.MatchID, err)
	}

	radiantHeroes := 0
	direHeroes := 0
	radiantPlayers := 0
	direPlayers := 0

	for _, p := range matchPlayers {
		if p.PlayerSlot < 128 {
			radiantHeroes++
			if p.PlayerID > 0 {
				radiantPlayers++
			}
		} else {
			direHeroes++
			if p.PlayerID > 0 {
				direPlayers++
			}
		}
	}

	if radiantHeroes != 5 {
		return fmt.Errorf("match %d: radiant team has %d heroes, expected 5", om.MatchID, radiantHeroes)
	}
	if direHeroes != 5 {
		return fmt.Errorf("match %d: dire team has %d heroes, expected 5", om.MatchID, direHeroes)
	}
	if radiantPlayers != 5 {
		return fmt.Errorf("match %d: radiant team has %d non-null players, expected 5", om.MatchID, radiantPlayers)
	}
	if direPlayers != 5 {
		return fmt.Errorf("match %d: dire team has %d non-null players, expected 5", om.MatchID, direPlayers)
	}

	return nil
}

func buildMatchEntities(om ODMatch) (Match, MatchMetadata, SeriesMatch) {
	m := Match{
		MatchID:    om.MatchID,
		Duration:   om.Duration,
		StartTime:  time.Unix(om.StartTime, 0),
		RadiantWin: om.RadiantWin,
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
	}
	if om.SeriesID > 0 {
		md.SeriesID = om.SeriesID
	}
	if om.RadiantTeam.Captain > 0 {
		md.RadiantCaptain = om.RadiantTeam.Captain
	}
	if om.DireTeam.Captain > 0 {
		md.DireCaptain = om.DireTeam.Captain
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

func insertRelatedEntities(ctx context.Context, tx bun.Tx, leagues map[int64]League, teams map[int64]Team, players map[int64]Player, seriesMap map[int64]Series) error {
	if len(leagues) > 0 {
		lSlice := mapsToSlice(leagues)
		if _, err := tx.NewInsert().Model(&lSlice).On("CONFLICT (league_id) DO NOTHING").Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert leagues: %w", err)
		}
	}

	if len(teams) > 0 {
		tSlice := mapsToSlice(teams)
		if _, err := tx.NewInsert().Model(&tSlice).On("CONFLICT (team_id) DO NOTHING").Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert teams: %w", err)
		}
	}

	if len(players) > 0 {
		pSlice := mapsToSlice(players)
		if _, err := tx.NewInsert().Model(&pSlice).On("CONFLICT (player_id) DO NOTHING").Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert players: %w", err)
		}
	}

	if len(seriesMap) > 0 {
		sSlice := mapsToSlice(seriesMap)
		if _, err := tx.NewInsert().Model(&sSlice).On("CONFLICT (series_id) DO NOTHING").Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert series: %w", err)
		}
	}

	return nil
}

func insertMatches(ctx context.Context, tx bun.Tx, matches []Match, metadata []MatchMetadata, seriesMatches []SeriesMatch) error {
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

func insertMatchBatch(matches []ODMatch, db *bun.DB) error {
	ctx := context.Background()

	leagues, teams, players, seriesMap := extractRelatedEntities(matches)

	var validMatches []Match
	var validMetadata []MatchMetadata
	var validSeriesMatches []SeriesMatch

	for _, om := range matches {
		if err := validateMatch(om); err != nil {
			log.Warn().Int64("match_id", om.MatchID).Err(err).Msg("Skipping match due to validation error")
			continue
		}

		m, md, sm := buildMatchEntities(om)
		validMatches = append(validMatches, m)
		validMetadata = append(validMetadata, md)
		if sm.SeriesID > 0 {
			validSeriesMatches = append(validSeriesMatches, sm)
		}
	}

	if len(validMatches) == 0 {
		log.Warn().Msg("No valid matches in batch, skipping insertion")
		return nil
	}

	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := insertRelatedEntities(ctx, tx, leagues, teams, players, seriesMap); err != nil {
			return err
		}

		return insertMatches(ctx, tx, validMatches, validMetadata, validSeriesMatches)
	})
}

func mapsToSlice[K comparable, V any](m map[K]V) []V {
	s := make([]V, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}
