package main

import (
	"context"
	"encoding/json"
	"os"
	"testing"
)

func TestQueryBuilder(t *testing.T) {
	ctx := context.Background()
	matchIds, err := fetchMatchIDs(ctx, 8673330227, 100)
	if err != nil {
		t.Fatalf("err fetching match ids < last_fetched_matched_id: %v", err)
	}
	b, err := fetchODMatches(ctx, matchIds)
	if err != nil {
		t.Fatal(err)
	}
	f, _ := os.Create("batch.json")
	defer f.Close()
	json.NewEncoder(f).Encode(b)
}
