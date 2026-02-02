package main

import (
	"os"
	"testing"
)

func TestQueryBuilder(t *testing.T) {
	testMatchesIds := []int64{8673330227}
	q := queryBuilder.GetMatches(testMatchesIds)

	if err := os.WriteFile("./od.query.sql", []byte(q), 0644); err != nil {
		t.Fatalf("Failed to write query file: %v", err)
	}
	// dumb change
	// change #2
	t.Log("Query written to ./od.query.sql")
}
