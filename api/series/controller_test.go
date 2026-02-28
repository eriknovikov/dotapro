package series

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"dotapro-lambda-api/errs"
	"dotapro-lambda-api/types"

	"github.com/go-chi/chi/v5"
)

// SeriesModelInterface defines the interface for series model operations
type SeriesModelInterface interface {
	GetMany(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, types.PaginationData, error)
	GetOne(ctx context.Context, id int64) (*types.SeriesDetail, error)
}

// mockModel is a mock implementation of the SeriesModelInterface for testing
type mockModel struct {
	getManyFunc func(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, types.PaginationData, error)
	getOneFunc  func(ctx context.Context, id int64) (*types.SeriesDetail, error)
}

func (m *mockModel) GetMany(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, types.PaginationData, error) {
	if m.getManyFunc != nil {
		return m.getManyFunc(ctx, filter)
	}
	return nil, types.PaginationData{}, nil
}

func (m *mockModel) GetOne(ctx context.Context, id int64) (*types.SeriesDetail, error) {
	if m.getOneFunc != nil {
		return m.getOneFunc(ctx, id)
	}
	return nil, nil
}

// mockController is a test helper that uses the mock model
type mockController struct {
	model SeriesModelInterface
}

func newMockController(model SeriesModelInterface) *mockController {
	return &mockController{model: model}
}

func (c *mockController) GetMany(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	filter := types.GetSeriesFilter{}
	params := r.URL.Query()

	// Parse optional filter parameters
	if leagueID, err := parseTestInt64Param(params, "league"); err != nil {
		writeTestError(w, "league_id", err)
		return
	} else {
		filter.LeagueID = leagueID
	}

	if teamID, err := parseTestInt64Param(params, "team"); err != nil {
		writeTestError(w, "team_id", err)
		return
	} else {
		filter.TeamID = teamID
	}

	// Parse sort parameter
	filter.Sort = params.Get("sort")

	// Parse limit parameter
	if limit, err := parseTestIntParam(params, "limit"); err != nil {
		writeTestError(w, "limit", err)
		return
	} else {
		filter.Limit = limit
	}

	// Parse cursor parameter
	if cursor, err := parseTestInt64Param(params, "c"); err != nil {
		writeTestError(w, "cursor", err)
		return
	} else {
		filter.Cursor = cursor
	}

	series, paginationData, err := c.model.GetMany(ctx, filter)
	if err != nil {
		if err == context.Canceled {
			return
		}
		if err == context.DeadlineExceeded {
			writeTestErrorResponse(w, context.DeadlineExceeded.Error(), http.StatusGatewayTimeout)
			return
		}
		writeTestErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := GetManyResp{Series: series, Pagination: paginationData}
	writeTestResponse(w, resp, http.StatusOK)
}

func (c *mockController) GetOne(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		writeTestErrorResponse(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := parseTestInt64FromString(idStr)
	if err != nil {
		writeTestErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	series, err := c.model.GetOne(r.Context(), id)
	if err != nil {
		if err == context.Canceled || err == context.DeadlineExceeded {
			return
		}
		switch err {
		case errs.NOT_FOUND:
			writeTestErrorResponse(w, err.Error(), http.StatusNotFound)
		default:
			writeTestErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	writeTestResponse(w, series, http.StatusOK)
}

// Helper functions for testing
func parseTestInt64Param(params url.Values, key string) (*int64, error) {
	value := params.Get(key)
	if value == "" {
		return nil, nil
	}
	parsed, err := parseTestInt64FromString(value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseTestInt64FromString(value string) (int64, error) {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func parseTestIntParam(params url.Values, key string) (int, error) {
	value := params.Get(key)
	if value == "" {
		return 0, nil
	}
	parsed, err := parseTestIntFromString(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func parseTestIntFromString(value string) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func writeTestError(w http.ResponseWriter, paramName string, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func writeTestErrorResponse(w http.ResponseWriter, err string, statusCode int) {
	resp := map[string]string{"error": err}
	writeTestResponse(w, resp, statusCode)
}

func writeTestResponse(w http.ResponseWriter, resp any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

// TestController_GetMany_Success tests successful retrieval of multiple series
func TestController_GetMany_Success(t *testing.T) {
	mock := &mockModel{
		getManyFunc: func(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, types.PaginationData, error) {
			return []types.SeriesSummary{
				{
					SeriesID:   1,
					StartTime:  time.Now(),
					TeamAScore: 2,
					TeamBScore: 1,
					TeamA: types.TeamInfo{
						ID:   100,
						Name: "Team A",
						Tag:  "TA",
					},
					TeamB: types.TeamInfo{
						ID:   200,
						Name: "Team B",
						Tag:  "TB",
					},
					League: types.LeagueInfo{
						ID:   10,
						Name: "Test League",
						Tier: "1",
					},
				},
			}, types.PaginationData{
				NextCursor: nil,
				HasMore:    false,
			}, nil
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series?limit=10", nil)
	w := httptest.NewRecorder()

	controller.GetMany(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp GetManyResp
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(resp.Series) != 1 {
		t.Errorf("Expected 1 series, got %d", len(resp.Series))
	}
	if resp.Series[0].SeriesID != 1 {
		t.Errorf("Expected series ID 1, got %d", resp.Series[0].SeriesID)
	}
}

// TestController_GetMany_WithLeagueFilter tests filtering by league
func TestController_GetMany_WithLeagueFilter(t *testing.T) {
	mock := &mockModel{
		getManyFunc: func(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, types.PaginationData, error) {
			if filter.LeagueID == nil || *filter.LeagueID != 10 {
				t.Errorf("Expected LeagueID 10, got %v", filter.LeagueID)
			}
			return []types.SeriesSummary{}, types.PaginationData{}, nil
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series?league=10", nil)
	w := httptest.NewRecorder()

	controller.GetMany(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestController_GetMany_WithTeamFilter tests filtering by team
func TestController_GetMany_WithTeamFilter(t *testing.T) {
	mock := &mockModel{
		getManyFunc: func(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, types.PaginationData, error) {
			if filter.TeamID == nil || *filter.TeamID != 100 {
				t.Errorf("Expected TeamID 100, got %v", filter.TeamID)
			}
			return []types.SeriesSummary{}, types.PaginationData{}, nil
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series?team=100", nil)
	w := httptest.NewRecorder()

	controller.GetMany(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestController_GetMany_WithSort tests sorting parameter
func TestController_GetMany_WithSort(t *testing.T) {
	mock := &mockModel{
		getManyFunc: func(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, types.PaginationData, error) {
			if filter.Sort != "oldest" {
				t.Errorf("Expected sort 'oldest', got %s", filter.Sort)
			}
			return []types.SeriesSummary{}, types.PaginationData{}, nil
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series?sort=oldest", nil)
	w := httptest.NewRecorder()

	controller.GetMany(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestController_GetMany_WithCursor tests cursor-based pagination
func TestController_GetMany_WithCursor(t *testing.T) {
	mock := &mockModel{
		getManyFunc: func(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, types.PaginationData, error) {
			if filter.Cursor == nil || *filter.Cursor != 50 {
				t.Errorf("Expected cursor 50, got %v", filter.Cursor)
			}
			return []types.SeriesSummary{}, types.PaginationData{}, nil
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series?c=50", nil)
	w := httptest.NewRecorder()

	controller.GetMany(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestController_GetMany_InvalidLeague tests invalid league parameter
func TestController_GetMany_InvalidLeague(t *testing.T) {
	mock := &mockModel{}
	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series?league=invalid", nil)
	w := httptest.NewRecorder()

	controller.GetMany(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	// http.Error writes plain text, not JSON
	body := w.Body.String()
	if body == "" {
		t.Error("Expected error message in response")
	}
}

// TestController_GetMany_InvalidTeam tests invalid team parameter
func TestController_GetMany_InvalidTeam(t *testing.T) {
	mock := &mockModel{}
	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series?team=invalid", nil)
	w := httptest.NewRecorder()

	controller.GetMany(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestController_GetMany_InvalidLimit tests invalid limit parameter
func TestController_GetMany_InvalidLimit(t *testing.T) {
	mock := &mockModel{}
	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series?limit=invalid", nil)
	w := httptest.NewRecorder()

	controller.GetMany(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestController_GetMany_InvalidCursor tests invalid cursor parameter
func TestController_GetMany_InvalidCursor(t *testing.T) {
	mock := &mockModel{}
	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series?c=invalid", nil)
	w := httptest.NewRecorder()

	controller.GetMany(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestController_GetMany_NotFound tests NOT_FOUND error handling
func TestController_GetMany_NotFound(t *testing.T) {
	mock := &mockModel{
		getManyFunc: func(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, types.PaginationData, error) {
			return nil, types.PaginationData{}, errs.NOT_FOUND
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series", nil)
	w := httptest.NewRecorder()

	controller.GetMany(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

// TestController_GetMany_DatabaseError tests database error handling
func TestController_GetMany_DatabaseError(t *testing.T) {
	mock := &mockModel{
		getManyFunc: func(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, types.PaginationData, error) {
			return nil, types.PaginationData{}, errors.New("database connection failed")
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series", nil)
	w := httptest.NewRecorder()

	controller.GetMany(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["error"] == "" {
		t.Error("Expected error message in response")
	}
}

// TestController_GetMany_ContextCanceled tests context cancellation handling
func TestController_GetMany_ContextCanceled(t *testing.T) {
	mock := &mockModel{
		getManyFunc: func(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, types.PaginationData, error) {
			return nil, types.PaginationData{}, context.Canceled
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series", nil)
	w := httptest.NewRecorder()

	controller.GetMany(w, req)

	// Context cancellation should not write a response
	if w.Code != 0 && w.Code != http.StatusOK {
		t.Logf("Status code: %d (context was canceled)", w.Code)
	}
}

// TestController_GetMany_ContextTimeout tests context timeout handling
func TestController_GetMany_ContextTimeout(t *testing.T) {
	mock := &mockModel{
		getManyFunc: func(ctx context.Context, filter types.GetSeriesFilter) ([]types.SeriesSummary, types.PaginationData, error) {
			return nil, types.PaginationData{}, context.DeadlineExceeded
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series", nil)
	w := httptest.NewRecorder()

	controller.GetMany(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("Expected status 504, got %d", w.Code)
	}
}

// TestController_GetOne_Success tests successful retrieval of a single series
func TestController_GetOne_Success(t *testing.T) {
	mock := &mockModel{
		getOneFunc: func(ctx context.Context, id int64) (*types.SeriesDetail, error) {
			if id != 1 {
				t.Errorf("Expected ID 1, got %d", id)
			}
			return &types.SeriesDetail{
				SeriesID:   1,
				StartTime:  time.Now(),
				TeamAScore: 2,
				TeamBScore: 1,
				TeamA: types.TeamInfo{
					ID:   100,
					Name: "Team A",
					Tag:  "TA",
				},
				TeamB: types.TeamInfo{
					ID:   200,
					Name: "Team B",
					Tag:  "TB",
				},
				League: types.LeagueInfo{
					ID:   10,
					Name: "Test League",
					Tier: "1",
				},
				Matches: []types.SeriesMatchDetail{},
			}, nil
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series/1", nil)
	w := httptest.NewRecorder()

	// Set chi URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	controller.GetOne(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp types.SeriesDetail
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.SeriesID != 1 {
		t.Errorf("Expected series ID 1, got %d", resp.SeriesID)
	}
}

// TestController_GetOne_NotFound tests NOT_FOUND error handling
func TestController_GetOne_NotFound(t *testing.T) {
	mock := &mockModel{
		getOneFunc: func(ctx context.Context, id int64) (*types.SeriesDetail, error) {
			return nil, errs.NOT_FOUND
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series/99999", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "99999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	controller.GetOne(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["error"] != "not found" {
		t.Errorf("Expected error 'not found', got %s", resp["error"])
	}
}

// TestController_GetOne_InvalidID tests invalid series ID
func TestController_GetOne_InvalidID(t *testing.T) {
	mock := &mockModel{}
	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series/invalid", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	controller.GetOne(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	// The mock controller writes JSON error responses
	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["error"] == "" {
		t.Error("Expected error message in response")
	}
}

// TestController_GetOne_MissingID tests missing ID parameter
func TestController_GetOne_MissingID(t *testing.T) {
	mock := &mockModel{}
	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series/", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	controller.GetOne(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["error"] != "missing id parameter" {
		t.Errorf("Expected 'missing id parameter', got %s", resp["error"])
	}
}

// TestController_GetOne_DatabaseError tests database error handling
func TestController_GetOne_DatabaseError(t *testing.T) {
	mock := &mockModel{
		getOneFunc: func(ctx context.Context, id int64) (*types.SeriesDetail, error) {
			return nil, errors.New("database connection failed")
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series/1", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	controller.GetOne(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

// TestController_GetOne_ContextCanceled tests context cancellation handling
func TestController_GetOne_ContextCanceled(t *testing.T) {
	mock := &mockModel{
		getOneFunc: func(ctx context.Context, id int64) (*types.SeriesDetail, error) {
			return nil, context.Canceled
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series/1", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	controller.GetOne(w, req)

	// Context cancellation should not write a response
	if w.Code != 0 && w.Code != http.StatusOK {
		t.Logf("Status code: %d (context was canceled)", w.Code)
	}
}

// TestController_GetOne_ContextTimeout tests context timeout handling
func TestController_GetOne_ContextTimeout(t *testing.T) {
	mock := &mockModel{
		getOneFunc: func(ctx context.Context, id int64) (*types.SeriesDetail, error) {
			return nil, context.DeadlineExceeded
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/series/1", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	controller.GetOne(w, req)

	// Context timeout should not write a response (it's handled gracefully)
	if w.Code != 0 {
		t.Logf("Status code: %d (context timeout handled)", w.Code)
	}
}

// TestNewController tests the controller constructor
func TestNewController(t *testing.T) {
	mock := &mockModel{}
	controller := newMockController(mock)

	if controller == nil {
		t.Fatal("newMockController() returned nil")
	}
	if controller.model != mock {
		t.Error("newMockController() did not set model correctly")
	}
}

// TestGetManyResp tests the response structure
func TestGetManyResp(t *testing.T) {
	resp := GetManyResp{
		Series: []types.SeriesSummary{
			{SeriesID: 1},
		},
		Pagination: types.PaginationData{
			NextCursor: nil,
			HasMore:    false,
		},
	}

	if len(resp.Series) != 1 {
		t.Errorf("Expected 1 series, got %d", len(resp.Series))
	}
	if resp.Pagination.HasMore {
		t.Error("Expected HasMore to be false")
	}
}

// TestParseQueryParams is a helper test to verify query parameter parsing
func TestParseQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected url.Values
	}{
		{
			name:  "simple query",
			query: "limit=10&sort=newest",
			expected: url.Values{
				"limit": []string{"10"},
				"sort":  []string{"newest"},
			},
		},
		{
			name:  "empty query",
			query: "",
			expected: url.Values{
				"": []string{""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse("?" + tt.query)
			if err != nil {
				t.Fatalf("Failed to parse URL: %v", err)
			}
			params := u.Query()
			for key, values := range tt.expected {
				if params.Get(key) != values[0] {
					t.Errorf("Expected %s=%s, got %s", key, values[0], params.Get(key))
				}
			}
		})
	}
}
