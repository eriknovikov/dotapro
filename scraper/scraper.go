package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

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
	BATCH_SIZE     = 100
	TIMES_TO_RETRY = 3
)

func ScrapeMatches(matchIds []int64, DB *bun.DB) error {
	N := len(matchIds)
	currentBatchIDs := make([]int64, 0, BATCH_SIZE)
	currentBatchMatches := make([]json.RawMessage, 0, BATCH_SIZE)
	var err error
	for i := 0; i < N; i += BATCH_SIZE {
		end := minInt(i+BATCH_SIZE, N)
		currentBatchIDs = matchIds[i:end]
		currentBatchMatches, err = fetchMatchBatch(currentBatchIDs)
		if err != nil {
			return fmt.Errorf("error scraping matches batch: %w", err)
		}

		if err := insertMatchBatch(currentBatchMatches, DB); err != nil {
			return fmt.Errorf("error inserting matches batch: %w", err)
		}
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

func insertMatchBatch(rawBatch []json.RawMessage, db *bun.DB) error {
	ctx := context.Background()

	var odMatches []ODMatch
	leagues := make(map[int64]League)
	teams := make(map[int64]Team)
	players := make(map[int64]Player)
	seriesMap := make(map[int64]Series)

	for _, raw := range rawBatch {
		var m ODMatch
		if err := json.Unmarshal(raw, &m); err != nil {
			return err
		}
		odMatches = append(odMatches, m)

		leagues[m.League.ID] = League{LeagueID: m.League.ID, Name: m.League.Name, Tier: m.League.Tier}
		teams[m.RadiantTeam.ID] = Team{TeamID: m.RadiantTeam.ID, Name: m.RadiantTeam.Name, Tag: m.RadiantTeam.Tag, LogoURL: m.RadiantTeam.LogoURL}
		teams[m.DireTeam.ID] = Team{TeamID: m.DireTeam.ID, Name: m.DireTeam.Name, Tag: m.DireTeam.Tag, LogoURL: m.DireTeam.LogoURL}

		if m.SeriesID != 0 {
			seriesMap[m.SeriesID] = Series{
				SeriesID:  m.SeriesID,
				LeagueID:  m.League.ID,
				TeamAID:   m.RadiantTeam.ID,
				TeamBID:   m.DireTeam.ID,
				StartTime: time.Unix(m.StartTime, 0),
			}
		}
	var matchPlayers []ODPlayerShort
	json.Unmarshal(m.Players, &matchPlayers)
	for _, p := range matchPlayers {
		if p.PlayerID > 0 {
			players[p.PlayerID] = Player{PlayerID: p.PlayerID, Name: p.Name}
		}
	}
	// Add captains to players map if they're not already there
	if m.RadiantTeam.Captain != nil && *m.RadiantTeam.Captain > 0 {
		if _, exists := players[*m.RadiantTeam.Captain]; !exists {
			players[*m.RadiantTeam.Captain] = Player{PlayerID: *m.RadiantTeam.Captain}
		}
	}
	if m.DireTeam.Captain != nil && *m.DireTeam.Captain > 0 {
		if _, exists := players[*m.DireTeam.Captain]; !exists {
			players[*m.DireTeam.Captain] = Player{PlayerID: *m.DireTeam.Captain}
		}
	}
}


	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if len(leagues) > 0 {
			lSlice := mapsToSlice(leagues)
			if _, err := tx.NewInsert().Model(&lSlice).On("CONFLICT (league_id) DO NOTHING").Exec(ctx); err != nil {
				return err
			}
		}
		if len(teams) > 0 {
			tSlice := mapsToSlice(teams)
			if _, err := tx.NewInsert().Model(&tSlice).On("CONFLICT (team_id) DO NOTHING").Exec(ctx); err != nil {
				return err
			}
		}
		if len(players) > 0 {
			pSlice := mapsToSlice(players)
			if _, err := tx.NewInsert().Model(&pSlice).On("CONFLICT (player_id) DO NOTHING").Exec(ctx); err != nil {
				return err
			}
		}
		if len(seriesMap) > 0 {
			sSlice := mapsToSlice(seriesMap)
			if _, err := tx.NewInsert().Model(&sSlice).On("CONFLICT (series_id) DO NOTHING").Exec(ctx); err != nil {
				return err
			}
		}

		var finalMatches []Match
		var finalMetadata []MatchMetadata
		var finalSeriesMatches []SeriesMatch

		for _, om := range odMatches {
			m := Match{
				MatchID:       om.MatchID,
				LeagueID:      om.League.ID,
				RadiantTeamID: om.RadiantTeam.ID,
				DireTeamID:    om.DireTeam.ID,
				Duration:      om.Duration,
				StartTime:     time.Unix(om.StartTime, 0),
				RadiantWin:    om.RadiantWin,
			}

			var ps []ODPlayerShort
			json.Unmarshal(om.Players, &ps)
			for _, p := range ps {
				if p.PlayerSlot < 128 {
					m.RadiantHeroes = append(m.RadiantHeroes, p.HeroID)
					m.RadiantPlayers = append(m.RadiantPlayers, p.PlayerID)
				} else {
					m.DireHeroes = append(m.DireHeroes, p.HeroID)
					m.DirePlayers = append(m.DirePlayers, p.PlayerID)
				}
			}
			finalMatches = append(finalMatches, m)

			finalMetadata = append(finalMetadata, MatchMetadata{
				MatchID:        om.MatchID,
				SeriesID:       om.SeriesID,
				RadiantCaptain: om.RadiantTeam.Captain,
				DireCaptain:    om.DireTeam.Captain,
				PicksBans:      om.PicksBans,
				PlayersData:    om.Players,
				RadiantGoldAdv: om.RadiantGoldAdv,
				RadiantXPAdv:   om.RadiantXPAdv,
				RadiantScore:   om.RadiantTeam.Score,
				DireScore:      om.DireTeam.Score,
			})

			if om.SeriesID != 0 {
				finalSeriesMatches = append(finalSeriesMatches, SeriesMatch{SeriesID: om.SeriesID, MatchID: om.MatchID})
			}
		}

		if _, err := tx.NewInsert().Model(&finalMatches).On("CONFLICT (match_id) DO NOTHING").Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewInsert().Model(&finalMetadata).On("CONFLICT (match_id) DO NOTHING").Exec(ctx); err != nil {
			return err
		}
		if len(finalSeriesMatches) > 0 {
			if _, err := tx.NewInsert().Model(&finalSeriesMatches).On("CONFLICT (series_id, match_id) DO NOTHING").Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

// Helper to convert map values to slice for Bun
func mapsToSlice[K comparable, V any](m map[K]V) []V {
	s := make([]V, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}
