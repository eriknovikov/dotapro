package main

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// TestProcessMessage_MarshalingAndInsert reads sample match JSON, processes it into the DB,
// then cleans up rows for that specific match.
func TestProcessMessage_MarshalingAndInsert(t *testing.T) {
	t.Helper()

	_ = godotenv.Load("../.env")

	// Read sample JSON (represents the SQS message body)
	body, err := os.ReadFile("../scraper/samplematches/1.json")
	if err != nil {
		t.Fatalf("failed to read sample json: %v", err)
	}

	// Unmarshal once to get match_id for cleanup
	var mm MatchMessage
	if err := json.Unmarshal(body, &mm); err != nil {
		t.Fatalf("failed to unmarshal sample json: %v", err)
	}

	ctx := context.Background()

	pool, err := getPostgresClient()
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	// Ensure cleanup of the inserted match rows
	t.Cleanup(func() {
		if mm.MatchID == 0 {
			return
		}
		// Remove rows for this specific match. CASCADE removes matches_metadata.
		_, _ = pool.Exec(ctx, "DELETE FROM user_feeds WHERE match_id = $1", mm.MatchID)
		_, _ = pool.Exec(ctx, "DELETE FROM matches WHERE match_id = $1", mm.MatchID)
	})

	// Process the message (this runs the full insertion flow in a transaction)
	if err := processMessage(ctx, pool, string(body)); err != nil {
		t.Fatalf("processMessage failed: %v", err)
	}

	// Basic existence check: the match should now exist
	var exists bool
	if err := pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM matches WHERE match_id = $1)", mm.MatchID).Scan(&exists); err != nil {
		t.Fatalf("query exists failed: %v", err)
	}
	if !exists {
		t.Fatalf("expected match %d to exist after insertion", mm.MatchID)
	}
}
