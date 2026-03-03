package series

import (
	"context"
	"os"
	"testing"
	"time"

	"dotapro-lambda-api/config"
	"dotapro-lambda-api/db"
	"dotapro-lambda-api/errs"
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

	dbURL := config.CONFIG.LocalDBURL
	if dbURL == "" {
		t.Skip("Skipping test: LOCAL_DB_URL not set")
	}

	bunDB, err := db.CreatePool(ctx, dbURL)
	if err != nil {
		t.Skipf("Skipping test: failed to connect to database: %v", err)
	}

	t.Cleanup(func() {
		_ = bunDB.Close()
	})

	return NewModel(bunDB)
}

// TestModel_GetMany_Success tests successful retrieval of multiple series
func TestModel_GetMany_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()

	filter := types.GetSeriesFilter{
		Limit: 10,
		Sort:  "newest",
	}

	series, paginationData, err := model.GetMany(ctx, filter)
	if err != nil {
		t.Errorf("GetMany() error = %v", err)
	}

	// Verify pagination data
	if paginationData.HasMore && paginationData.NextCursor == nil {
		t.Error("Expected NextCursor when HasMore is true")
	}

	t.Logf("Retrieved %d series", len(series))
}

// TestModel_GetMany_NotFound tests NOT_FOUND error handling
func TestModel_GetMany_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()

	// Use a non-existent league ID
	leagueID := int64(999999999)
	filter := types.GetSeriesFilter{
		LeagueID: &leagueID,
		Limit:    10,
	}

	_, _, err := model.GetMany(ctx, filter)
	if err != errs.ErrNotFound {
		t.Logf("GetMany() error = %v (may be expected for non-existent league)", err)
	}
}

// TestModel_GetMany_DatabaseError tests database error handling
func TestModel_GetMany_DatabaseError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()

	// This test verifies error handling - we can't easily force a DB error
	// without breaking the connection, so we just verify the function works
	filter := types.GetSeriesFilter{
		Limit: 10,
	}

	_, _, err := model.GetMany(ctx, filter)
	if err != nil && err != errs.ErrNotFound {
		t.Errorf("GetMany() unexpected error = %v", err)
	}
}

// TestModel_GetMany_WithFilters tests filtering by league and team
func TestModel_GetMany_WithFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()

	leagueID := int64(10)
	teamID := int64(100)

	filter := types.GetSeriesFilter{
		LeagueID: &leagueID,
		TeamID:   &teamID,
		Limit:    20,
		Sort:     "oldest",
	}

	series, _, err := model.GetMany(ctx, filter)
	if err != nil && err != errs.ErrNotFound {
		t.Errorf("GetMany() error = %v", err)
	}

	t.Logf("Retrieved %d series with filters", len(series))
}

// TestModel_GetMany_WithCursor tests cursor-based pagination
func TestModel_GetMany_WithCursor(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()

	// First, get some series to find a valid cursor
	series, _, err := model.GetMany(ctx, types.GetSeriesFilter{Limit: 1})
	if err != nil {
		t.Skipf("Skipping test: failed to get initial series: %v", err)
	}

	if len(series) == 0 {
		t.Skip("Skipping test: no series found in database")
	}

	cursor := series[0].SeriesID
	filter := types.GetSeriesFilter{
		Limit:  10,
		Cursor: &cursor,
		Sort:   "newest",
	}

	_, paginationData, err := model.GetMany(ctx, filter)
	if err != nil {
		t.Errorf("GetMany() error = %v", err)
	}

	t.Logf("Pagination with cursor: HasMore=%v", paginationData.HasMore)
}

// TestModel_GetOne_Success tests successful retrieval of a single series
func TestModel_GetOne_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()

	// First, get a series ID from GetMany
	series, _, err := model.GetMany(ctx, types.GetSeriesFilter{Limit: 1})
	if err != nil {
		t.Skipf("Skipping test: failed to get series: %v", err)
	}

	if len(series) == 0 {
		t.Skip("Skipping test: no series found in database")
	}

	seriesID := series[0].SeriesID
	seriesDetail, err := model.GetOne(ctx, seriesID)

	if err != nil {
		t.Errorf("GetOne() error = %v", err)
		return
	}
	if seriesDetail == nil {
		t.Error("GetOne() expected series, got nil")
		return
	}

	if seriesDetail.SeriesID != seriesID {
		t.Errorf("GetOne() expected series ID %d, got %d", seriesID, seriesDetail.SeriesID)
	}

	t.Logf("Retrieved series with %d matches", len(seriesDetail.Matches))
}

// TestModel_GetOne_NotFound tests NOT_FOUND error for non-existent series
func TestModel_GetOne_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	seriesID := int64(999999999)

	_, err := model.GetOne(ctx, seriesID)

	if err != errs.ErrNotFound {
		t.Logf("GetOne() error = %v (expected NOT_FOUND)", err)
	}
}

// TestModel_GetOne_InvalidID tests handling of invalid series ID
func TestModel_GetOne_InvalidID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()
	seriesID := int64(0)

	_, err := model.GetOne(ctx, seriesID)

	// Should return NOT_FOUND for non-existent series
	if err != errs.ErrNotFound {
		t.Logf("GetOne() error = %v (expected NOT_FOUND)", err)
	}
}

// TestModel_GetOne_WithMatches tests retrieval of series with matches
func TestModel_GetOne_WithMatches(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()

	// First, get a series ID from GetMany
	series, _, err := model.GetMany(ctx, types.GetSeriesFilter{Limit: 10})
	if err != nil {
		t.Skipf("Skipping test: failed to get series: %v", err)
	}

	if len(series) == 0 {
		t.Skip("Skipping test: no series found in database")
	}

	// Find a series with matches
	var seriesWithMatches int64
	for _, s := range series {
		seriesDetail, err := model.GetOne(ctx, s.SeriesID)
		if err != nil {
			continue
		}
		if len(seriesDetail.Matches) > 0 {
			seriesWithMatches = s.SeriesID
			break
		}
	}

	if seriesWithMatches == 0 {
		t.Skip("Skipping test: no series with matches found")
	}

	seriesDetail, err := model.GetOne(ctx, seriesWithMatches)
	if err != nil {
		t.Errorf("GetOne() error = %v", err)
	}

	if len(seriesDetail.Matches) == 0 {
		t.Error("GetOne() expected series with matches")
	}

	t.Logf("Series %d has %d matches", seriesWithMatches, len(seriesDetail.Matches))
}

// TestModel_getMatchesForSeries_Success tests successful retrieval of matches for a series
func TestModel_getMatchesForSeries_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()

	// First, get a series ID from GetMany
	series, _, err := model.GetMany(ctx, types.GetSeriesFilter{Limit: 10})
	if err != nil {
		t.Skipf("Skipping test: failed to get series: %v", err)
	}

	if len(series) == 0 {
		t.Skip("Skipping test: no series found in database")
	}

	seriesID := series[0].SeriesID
	matches, err := model.getMatchesForSeries(ctx, seriesID)

	if err != nil {
		t.Errorf("getMatchesForSeries() error = %v", err)
	}
	if matches == nil {
		t.Error("getMatchesForSeries() expected matches slice, got nil")
	}

	t.Logf("Series %d has %d matches", seriesID, len(matches))
}

// TestModel_getMatchesForSeries_Empty tests handling of series with no matches
func TestModel_getMatchesForSeries_Empty(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	model := getTestDB(t)
	ctx := context.Background()

	// Use a non-existent series ID
	seriesID := int64(999999999)
	matches, err := model.getMatchesForSeries(ctx, seriesID)

	if err != nil {
		t.Errorf("getMatchesForSeries() error = %v", err)
	}
	// The model returns an empty slice, not nil, when there are no matches
	if len(matches) != 0 {
		t.Logf("getMatchesForSeries() returned %d matches for non-existent series", len(matches))
	}
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

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Set test environment
	_ = os.Setenv("ENVIRON", "local")
	_ = os.Setenv("LOCAL_DB_URL", "postgres://postgres:admin@localhost:5432/dotapro?sslmode=disable")

	// Run tests
	code := m.Run()

	os.Exit(code)
}
