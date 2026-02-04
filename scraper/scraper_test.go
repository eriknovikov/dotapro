package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestInsertMiniBatch(t *testing.T) {
	// 1. Initialize config and DB (using your main.go logic)
	initialize()
	if err := ensureDB(); err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer DB.Close()

	// 2. Read the mini_batch.json file
	data, err := os.ReadFile("mini_batch.json")
	if err != nil {
		t.Fatalf("Failed to read mini_batch.json: %v", err)
	}

	// 3. Unmarshal into []json.RawMessage
	var rawBatch []json.RawMessage
	if err := json.Unmarshal(data, &rawBatch); err != nil {
		t.Fatalf("Failed to unmarshal json: %v", err)
	}

	// 4. Run the insertion
	fmt.Printf("Attempting to insert %d matches...\n", len(rawBatch))
	err = insertMatchBatch(rawBatch, DB)
	if err != nil {
		t.Fatalf("Insertion failed: %v", err)
	}

	fmt.Println("Success! Check your local database tables.")
}
