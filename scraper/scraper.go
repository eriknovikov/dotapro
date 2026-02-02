package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"

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
	BATCH_SIZE      = 100
	TOTAL_TO_SCRAPE = 1000
	TIMES_TO_RETRY  = 3
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

func insertMatchBatch(matches []json.RawMessage, db *bun.DB) error {
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

func makeOpendotaRequest(query string) (*http.Response, error) {
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

func fetchMatchIDs(lastMatchID int64) ([]int64, error) {
	resp, err := makeOpendotaRequest(queryBuilder.GetIds(lastMatchID, TOTAL_TO_SCRAPE))
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
	resp, err := makeOpendotaRequest(queryBuilder.GetMatches(matchIDs))
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
