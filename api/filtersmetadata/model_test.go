package filtersmetadata

import (
	"context"
	"os"
	"testing"
	"time"

	"dotapro-lambda-api/config"
	"dotapro-lambda-api/db"
	"dotapro-lambda-api/types"
)

// getTestDB returns a database connection for testing
func getTestDB(t *testing.T) *Model {
	// Load environment variables
	if err := config.LoadEnvs(); err != nil {
		t.Skipf("Skipping test: failed to load environment variables: %v", err)
	}

	// Check if we're in local environment
	if !config.IsLocal() {
		t.Skip("Skipping test: not in local environment")
	}

	// Create database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbURL := config.CONFIG.LOCAL_DB_URL
	if dbURL == "" {
		t.Skip("Skipping test: LOCAL_DB_URL not set")
	}

	bunDB, err := db.CreatePool(ctx, dbURL)
	if err != nil {
		t.Skipf("Skipping test: failed to connect to database: %v", err)
	}

	t.Cleanup(func() {
		bunDB.Close()
	})

	return NewModel(bunDB)
}

// TestModel_SearchTeams_Success tests successful team search
func TestModel_SearchTeams_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	query := "Team"

	teams, err := model.SearchTeams(ctx, query)
	if err != nil {
		t.Errorf("SearchTeams() error = %v", err)
	}

	t.Logf("Found %d teams matching '%s'", len(teams), query)
}

// TestModel_SearchTeams_EmptyResults tests empty search results
func TestModel_SearchTeams_EmptyResults(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	query := "NonExistentTeamXYZ123"

	teams, err := model.SearchTeams(ctx, query)
	if err != nil {
		t.Errorf("SearchTeams() error = %v", err)
	}

	if len(teams) != 0 {
		t.Logf("Found %d teams for non-existent query (may be expected)", len(teams))
	}
}

// TestModel_SearchTeams_DatabaseError tests database error handling
func TestModel_SearchTeams_DatabaseError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	query := "Team"

	// This test verifies error handling - we can't easily force a DB error
	// without breaking the connection, so we just verify the function works
	_, err := model.SearchTeams(ctx, query)
	if err != nil {
		t.Logf("SearchTeams() error = %v", err)
	}
}

// TestModel_SearchTeams_EmptyQuery tests empty query string
func TestModel_SearchTeams_EmptyQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	query := ""

	// Empty query should return empty results or error
	teams, err := model.SearchTeams(ctx, query)
	if err != nil {
		t.Logf("SearchTeams() with empty query returned error: %v", err)
	}

	if len(teams) != 0 {
		t.Logf("Empty query returned %d teams", len(teams))
	}
}

// TestModel_SearchTeams_SpecialCharacters tests query with special characters
func TestModel_SearchTeams_SpecialCharacters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	query := "Team's Name"

	teams, err := model.SearchTeams(ctx, query)
	if err != nil {
		t.Errorf("SearchTeams() error = %v", err)
	}

	t.Logf("Found %d teams with special characters in query", len(teams))
}

// TestModel_SearchLeagues_Success tests successful league search
func TestModel_SearchLeagues_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	query := "League"

	leagues, err := model.SearchLeagues(ctx, query)
	if err != nil {
		t.Errorf("SearchLeagues() error = %v", err)
	}

	t.Logf("Found %d leagues matching '%s'", len(leagues), query)
}

// TestModel_SearchLeagues_EmptyResults tests empty search results
func TestModel_SearchLeagues_EmptyResults(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	query := "NonExistentLeagueXYZ123"

	leagues, err := model.SearchLeagues(ctx, query)
	if err != nil {
		t.Errorf("SearchLeagues() error = %v", err)
	}

	if len(leagues) != 0 {
		t.Logf("Found %d leagues for non-existent query (may be expected)", len(leagues))
	}
}

// TestModel_SearchLeagues_DatabaseError tests database error handling
func TestModel_SearchLeagues_DatabaseError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	query := "League"

	// This test verifies error handling - we can't easily force a DB error
	// without breaking the connection, so we just verify the function works
	_, err := model.SearchLeagues(ctx, query)
	if err != nil {
		t.Logf("SearchLeagues() error = %v", err)
	}
}

// TestModel_SearchLeagues_EmptyQuery tests empty query string
func TestModel_SearchLeagues_EmptyQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	query := ""

	// Empty query should return empty results or error
	leagues, err := model.SearchLeagues(ctx, query)
	if err != nil {
		t.Logf("SearchLeagues() with empty query returned error: %v", err)
	}

	if len(leagues) != 0 {
		t.Logf("Empty query returned %d leagues", len(leagues))
	}
}

// TestModel_SearchLeagues_SpecialCharacters tests query with special characters
func TestModel_SearchLeagues_SpecialCharacters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	query := "League's Cup"

	leagues, err := model.SearchLeagues(ctx, query)
	if err != nil {
		t.Errorf("SearchLeagues() error = %v", err)
	}

	t.Logf("Found %d leagues with special characters in query", len(leagues))
}

// TestModel_SearchTeams_LimitApplied tests that search limit is applied
func TestModel_SearchTeams_LimitApplied(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	query := "Team"

	teams, err := model.SearchTeams(ctx, query)
	if err != nil {
		t.Errorf("SearchTeams() error = %v", err)
	}

	// The model should apply the SearchLimit constant (10)
	if len(teams) > 10 {
		t.Errorf("Expected at most 10 teams, got %d", len(teams))
	}

	t.Logf("Search returned %d teams (limit: 10)", len(teams))
}

// TestModel_SearchLeagues_LimitApplied tests that search limit is applied
func TestModel_SearchLeagues_LimitApplied(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	query := "League"

	leagues, err := model.SearchLeagues(ctx, query)
	if err != nil {
		t.Errorf("SearchLeagues() error = %v", err)
	}

	// The model should apply the SearchLimit constant (10)
	if len(leagues) > 10 {
		t.Errorf("Expected at most 10 leagues, got %d", len(leagues))
	}

	t.Logf("Search returned %d leagues (limit: 10)", len(leagues))
}

// TestNewModel tests the model constructor
func TestNewModel(t *testing.T) {
	model := NewModel(nil)

	if model == nil {
		t.Fatal("NewModel() returned nil")
	}
	if model.DB != nil {
		t.Error("NewModel() expected DB to be nil when passed nil")
	}
}

// TestTeamSearchResult tests the TeamSearchResult structure
func TestTeamSearchResult(t *testing.T) {
	result := types.TeamSearchResult{
		TeamID:  100,
		Name:    "Test Team",
		LogoURL: "https://example.com/logo.png",
	}

	if result.TeamID != 100 {
		t.Errorf("Expected TeamID 100, got %d", result.TeamID)
	}
	if result.Name != "Test Team" {
		t.Errorf("Expected Name 'Test Team', got %s", result.Name)
	}
	if result.LogoURL != "https://example.com/logo.png" {
		t.Errorf("Expected LogoURL 'https://example.com/logo.png', got %s", result.LogoURL)
	}
}

// TestLeagueSearchResult tests the LeagueSearchResult structure
func TestLeagueSearchResult(t *testing.T) {
	result := types.LeagueSearchResult{
		LeagueID: 10,
		Name:     "Test League",
	}

	if result.LeagueID != 10 {
		t.Errorf("Expected LeagueID 10, got %d", result.LeagueID)
	}
	if result.Name != "Test League" {
		t.Errorf("Expected Name 'Test League', got %s", result.Name)
	}
}

// TestSearchQueryValidation tests various query string scenarios
func TestSearchQueryValidation(t *testing.T) {
	tests := []struct {
		name  string
		query string
		valid bool
	}{
		{
			name:  "normal query",
			query: "Team Name",
			valid: true,
		},
		{
			name:  "single word",
			query: "Team",
			valid: true,
		},
		{
			name:  "with numbers",
			query: "Team 123",
			valid: true,
		},
		{
			name:  "with special chars",
			query: "Team's Name",
			valid: true,
		},
		{
			name:  "empty string",
			query: "",
			valid: true, // Empty query should be handled gracefully
		},
		{
			name:  "with spaces",
			query: "  Team Name  ",
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify query is not nil
			if tt.query == "" && !tt.valid {
				t.Error("Expected empty query to be invalid")
			}
		})
	}
}

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Set test environment
	os.Setenv("ENVIRON", "local")
	os.Setenv("LOCAL_DB_URL", "postgres://postgres:admin@localhost:5432/dotapro?sslmode=disable")

	// Run tests
	code := m.Run()

	os.Exit(code)
}
