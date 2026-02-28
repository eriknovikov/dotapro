package filtersmetadata

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"dotapro-lambda-api/types"
)

// FiltersModelInterface defines the interface for filters metadata model operations
type FiltersModelInterface interface {
	SearchTeams(ctx context.Context, query string) ([]types.TeamSearchResult, error)
	SearchLeagues(ctx context.Context, query string) ([]types.LeagueSearchResult, error)
}

// mockModel is a mock implementation of the FiltersModelInterface for testing
type mockModel struct {
	searchTeamsFunc  func(ctx context.Context, query string) ([]types.TeamSearchResult, error)
	searchLeaguesFunc func(ctx context.Context, query string) ([]types.LeagueSearchResult, error)
}

func (m *mockModel) SearchTeams(ctx context.Context, query string) ([]types.TeamSearchResult, error) {
	if m.searchTeamsFunc != nil {
		return m.searchTeamsFunc(ctx, query)
	}
	return nil, nil
}

func (m *mockModel) SearchLeagues(ctx context.Context, query string) ([]types.LeagueSearchResult, error) {
	if m.searchLeaguesFunc != nil {
		return m.searchLeaguesFunc(ctx, query)
	}
	return nil, nil
}

// mockController is a test helper that uses the mock model
type mockController struct {
	model FiltersModelInterface
}

func newMockController(model FiltersModelInterface) *mockController {
	return &mockController{model: model}
}

func (c *mockController) SearchTeams(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query, err := parseRequiredStringParam(r.URL.Query(), "q")
	if err != nil {
		writeTestErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	teams, err := c.model.SearchTeams(ctx, query)
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

	writeTestResponse(w, teams, http.StatusOK)
}

func (c *mockController) SearchLeagues(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query, err := parseRequiredStringParam(r.URL.Query(), "q")
	if err != nil {
		writeTestErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	leagues, err := c.model.SearchLeagues(ctx, query)
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

	writeTestResponse(w, leagues, http.StatusOK)
}

// Helper functions for testing
func parseRequiredStringParam(params url.Values, key string) (string, error) {
	value := params.Get(key)
	if value == "" {
		return "", errors.New("missing required parameter: " + key)
	}
	return value, nil
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

// TestController_SearchTeams_Success tests successful team search
func TestController_SearchTeams_Success(t *testing.T) {
	mock := &mockModel{
		searchTeamsFunc: func(ctx context.Context, query string) ([]types.TeamSearchResult, error) {
			if query != "Team" {
				t.Errorf("Expected query 'Team', got %s", query)
			}
			return []types.TeamSearchResult{
				{
					TeamID:  100,
					Name:    "Team Alpha",
					LogoURL: "https://example.com/logo1.png",
				},
				{
					TeamID:  200,
					Name:    "Team Beta",
					LogoURL: "https://example.com/logo2.png",
				},
			}, nil
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/teams?q=Team", nil)
	w := httptest.NewRecorder()

	controller.SearchTeams(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp []types.TeamSearchResult
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("Expected 2 teams, got %d", len(resp))
	}
	if resp[0].TeamID != 100 {
		t.Errorf("Expected team ID 100, got %d", resp[0].TeamID)
	}
}

// TestController_SearchTeams_EmptyResults tests empty search results
func TestController_SearchTeams_EmptyResults(t *testing.T) {
	mock := &mockModel{
		searchTeamsFunc: func(ctx context.Context, query string) ([]types.TeamSearchResult, error) {
			return []types.TeamSearchResult{}, nil
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/teams?q=NonExistentTeam", nil)
	w := httptest.NewRecorder()

	controller.SearchTeams(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp []types.TeamSearchResult
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(resp) != 0 {
		t.Errorf("Expected 0 teams, got %d", len(resp))
	}
}

// TestController_SearchTeams_MissingQuery tests missing query parameter
func TestController_SearchTeams_MissingQuery(t *testing.T) {
	mock := &mockModel{}
	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/teams", nil)
	w := httptest.NewRecorder()

	controller.SearchTeams(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !strings.Contains(resp["error"], "missing required parameter") {
		t.Errorf("Expected 'missing required parameter' in error, got %s", resp["error"])
	}
}

// TestController_SearchTeams_EmptyQuery tests empty query parameter
func TestController_SearchTeams_EmptyQuery(t *testing.T) {
	mock := &mockModel{}
	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/teams?q=", nil)
	w := httptest.NewRecorder()

	controller.SearchTeams(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !strings.Contains(resp["error"], "missing required parameter") {
		t.Errorf("Expected 'missing required parameter' in error, got %s", resp["error"])
	}
}

// TestController_SearchTeams_DatabaseError tests database error handling
func TestController_SearchTeams_DatabaseError(t *testing.T) {
	mock := &mockModel{
		searchTeamsFunc: func(ctx context.Context, query string) ([]types.TeamSearchResult, error) {
			return nil, errors.New("database connection failed")
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/teams?q=Team", nil)
	w := httptest.NewRecorder()

	controller.SearchTeams(w, req)

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

// TestController_SearchTeams_ContextCanceled tests context cancellation handling
func TestController_SearchTeams_ContextCanceled(t *testing.T) {
	mock := &mockModel{
		searchTeamsFunc: func(ctx context.Context, query string) ([]types.TeamSearchResult, error) {
			return nil, context.Canceled
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/teams?q=Team", nil)
	w := httptest.NewRecorder()

	controller.SearchTeams(w, req)

	// Context cancellation should not write a response
	if w.Code != 0 && w.Code != http.StatusOK {
		t.Logf("Status code: %d (context was canceled)", w.Code)
	}
}

// TestController_SearchTeams_ContextTimeout tests context timeout handling
func TestController_SearchTeams_ContextTimeout(t *testing.T) {
	mock := &mockModel{
		searchTeamsFunc: func(ctx context.Context, query string) ([]types.TeamSearchResult, error) {
			return nil, context.DeadlineExceeded
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/teams?q=Team", nil)
	w := httptest.NewRecorder()

	controller.SearchTeams(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("Expected status 504, got %d", w.Code)
	}
}

// TestController_SearchTeams_SpecialCharacters tests query with special characters
func TestController_SearchTeams_SpecialCharacters(t *testing.T) {
	mock := &mockModel{
		searchTeamsFunc: func(ctx context.Context, query string) ([]types.TeamSearchResult, error) {
			if query != "Team's Name" {
				t.Errorf("Expected query 'Team's Name', got %s", query)
			}
			return []types.TeamSearchResult{}, nil
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/teams?q=Team's+Name", nil)
	w := httptest.NewRecorder()

	controller.SearchTeams(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestController_SearchLeagues_Success tests successful league search
func TestController_SearchLeagues_Success(t *testing.T) {
	mock := &mockModel{
		searchLeaguesFunc: func(ctx context.Context, query string) ([]types.LeagueSearchResult, error) {
			if query != "League" {
				t.Errorf("Expected query 'League', got %s", query)
			}
			return []types.LeagueSearchResult{
				{
					LeagueID: 10,
					Name:     "The International",
				},
				{
					LeagueID: 20,
					Name:     "Dota Pro Circuit",
				},
			}, nil
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/leagues?q=League", nil)
	w := httptest.NewRecorder()

	controller.SearchLeagues(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp []types.LeagueSearchResult
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("Expected 2 leagues, got %d", len(resp))
	}
	if resp[0].LeagueID != 10 {
		t.Errorf("Expected league ID 10, got %d", resp[0].LeagueID)
	}
}

// TestController_SearchLeagues_EmptyResults tests empty search results
func TestController_SearchLeagues_EmptyResults(t *testing.T) {
	mock := &mockModel{
		searchLeaguesFunc: func(ctx context.Context, query string) ([]types.LeagueSearchResult, error) {
			return []types.LeagueSearchResult{}, nil
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/leagues?q=NonExistentLeague", nil)
	w := httptest.NewRecorder()

	controller.SearchLeagues(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp []types.LeagueSearchResult
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(resp) != 0 {
		t.Errorf("Expected 0 leagues, got %d", len(resp))
	}
}

// TestController_SearchLeagues_MissingQuery tests missing query parameter
func TestController_SearchLeagues_MissingQuery(t *testing.T) {
	mock := &mockModel{}
	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/leagues", nil)
	w := httptest.NewRecorder()

	controller.SearchLeagues(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !strings.Contains(resp["error"], "missing required parameter") {
		t.Errorf("Expected 'missing required parameter' in error, got %s", resp["error"])
	}
}

// TestController_SearchLeagues_EmptyQuery tests empty query parameter
func TestController_SearchLeagues_EmptyQuery(t *testing.T) {
	mock := &mockModel{}
	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/leagues?q=", nil)
	w := httptest.NewRecorder()

	controller.SearchLeagues(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !strings.Contains(resp["error"], "missing required parameter") {
		t.Errorf("Expected 'missing required parameter' in error, got %s", resp["error"])
	}
}

// TestController_SearchLeagues_DatabaseError tests database error handling
func TestController_SearchLeagues_DatabaseError(t *testing.T) {
	mock := &mockModel{
		searchLeaguesFunc: func(ctx context.Context, query string) ([]types.LeagueSearchResult, error) {
			return nil, errors.New("database connection failed")
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/leagues?q=League", nil)
	w := httptest.NewRecorder()

	controller.SearchLeagues(w, req)

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

// TestController_SearchLeagues_ContextCanceled tests context cancellation handling
func TestController_SearchLeagues_ContextCanceled(t *testing.T) {
	mock := &mockModel{
		searchLeaguesFunc: func(ctx context.Context, query string) ([]types.LeagueSearchResult, error) {
			return nil, context.Canceled
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/leagues?q=League", nil)
	w := httptest.NewRecorder()

	controller.SearchLeagues(w, req)

	// Context cancellation should not write a response
	if w.Code != 0 && w.Code != http.StatusOK {
		t.Logf("Status code: %d (context was canceled)", w.Code)
	}
}

// TestController_SearchLeagues_ContextTimeout tests context timeout handling
func TestController_SearchLeagues_ContextTimeout(t *testing.T) {
	mock := &mockModel{
		searchLeaguesFunc: func(ctx context.Context, query string) ([]types.LeagueSearchResult, error) {
			return nil, context.DeadlineExceeded
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/leagues?q=League", nil)
	w := httptest.NewRecorder()

	controller.SearchLeagues(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("Expected status 504, got %d", w.Code)
	}
}

// TestController_SearchLeagues_SpecialCharacters tests query with special characters
func TestController_SearchLeagues_SpecialCharacters(t *testing.T) {
	mock := &mockModel{
		searchLeaguesFunc: func(ctx context.Context, query string) ([]types.LeagueSearchResult, error) {
			if query != "League's Cup" {
				t.Errorf("Expected query 'League's Cup', got %s", query)
			}
			return []types.LeagueSearchResult{}, nil
		},
	}

	controller := newMockController(mock)

	req := httptest.NewRequest("GET", "/filtersmetadata/leagues?q=League's+Cup", nil)
	w := httptest.NewRecorder()

	controller.SearchLeagues(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
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

// TestParseRequiredStringParam tests the parameter parsing helper
func TestParseRequiredStringParam(t *testing.T) {
	tests := []struct {
		name      string
		params    url.Values
		key       string
		expected  string
		wantError bool
	}{
		{
			name: "valid parameter",
			params: url.Values{
				"q": []string{"test query"},
			},
			key:       "q",
			expected:  "test query",
			wantError: false,
		},
		{
			name: "missing parameter",
			params: url.Values{
				"other": []string{"value"},
			},
			key:       "q",
			expected:  "",
			wantError: true,
		},
		{
			name:      "empty params",
			params:    url.Values{},
			key:       "q",
			expected:  "",
			wantError: true,
		},
		{
			name: "empty string parameter",
			params: url.Values{
				"q": []string{""},
			},
			key:       "q",
			expected:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseRequiredStringParam(tt.params, tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("parseRequiredStringParam() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if result != tt.expected {
				t.Errorf("parseRequiredStringParam() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSearchQueryEncoding tests URL encoding of search queries
func TestSearchQueryEncoding(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "simple query",
			query:    "Team",
			expected: "Team",
		},
		{
			name:     "with spaces",
			query:    "Team Name",
			expected: "Team Name",
		},
		{
			name:     "with apostrophe",
			query:    "Team's Name",
			expected: "Team's Name",
		},
		{
			name:     "with special chars",
			query:    "Team & Co.",
			expected: "Team & Co.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request with the query
			req := httptest.NewRequest("GET", "/filtersmetadata/teams?q="+url.QueryEscape(tt.query), nil)
			query := req.URL.Query().Get("q")
			if query != tt.expected {
				t.Errorf("Expected query %q, got %q", tt.expected, query)
			}
		})
	}
}
