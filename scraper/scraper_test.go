//go:build insert
// +build insert

package main

import (
	"context"
	"testing"
)

func TestInsert(t *testing.T) {
	initialize()
	if err := ensureDB(); err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer DB.Close()

	if err := ScrapeMatches(context.Background(), DB, 1*800); err != nil {
		t.Fatal(err)
	}
	t.Log("success. check db")
}
