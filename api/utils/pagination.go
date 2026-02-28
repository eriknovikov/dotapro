package utils

import "dotapro-lambda-api/constants"

// CursorExtractor is a function that extracts the cursor ID from a result item
type CursorExtractor[T any] func(T) int64

// ProcessPagination processes a paginated result set and returns pagination metadata.
// It handles the "fetch one extra" pattern for cursor-based pagination.
//
// Parameters:
//   - results: The fetched results (may include one extra item)
//   - limit: The requested limit
//   - extractCursor: A function to extract the cursor ID from a result item
//
// Returns:
//   - trimmed: The results trimmed to the requested limit
//   - pagination: Pagination metadata including next cursor and has_more flag
func ProcessPagination[T any](results []T, limit int, extractCursor CursorExtractor[T]) ([]T, PaginationData) {
	// Apply default limit if not specified
	if limit <= 0 {
		limit = constants.DefaultLimit
	}

	var nextCursor *int64
	hasMore := false

	// Check if we have more results than requested
	if len(results) > limit {
		hasMore = true
		results = results[:limit]
		
		// Extract cursor from the last item
		if len(results) > 0 {
			lastID := extractCursor(results[len(results)-1])
			nextCursor = &lastID
		}
	}

	return results, PaginationData{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}
}

// PaginationData represents pagination metadata
type PaginationData struct {
	NextCursor *int64 `json:"nc"`       // ID for next page
	HasMore    bool   `json:"has_more"` // Whether there are more results
}
