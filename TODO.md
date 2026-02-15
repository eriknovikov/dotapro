# Dotapro TODO

## Series Page UI Refactor - Filters Sidebar + Results Table

### Phase 1: Type Definitions and Setup
- [x] Rename `Params` type to `Filters` type in appropriate location
- [x] Define `Filters` TypeScript interface with fields: leagueId, teamId, sortBy
- [x] Create `filtersRef` using useRef in the series route component
- [x] Set up initial filter state structure

### Phase 2: Sidebar Component
- [x] Create `FiltersSidebar` component
- [x] Add league ID input field (text input for now)
- [x] Add team ID input field (text input for now)
- [x] Add sort by dropdown/select (newest/oldest)
- [x] Add search button component
- [x] Style sidebar with Tailwind (fixed width, proper spacing)
- [x] Connect all inputs to filtersRef (onChange handlers update ref.current)

### Phase 3: Search Functionality
- [x] Create `handleSearch` function
- [x] Read values from filtersRef on search button click
- [x] Make API request to `/series` endpoint with filter parameters
- [x] Update URL history using TanStack Router navigate with replace: true
- [x] Ensure URL update doesn't trigger component rerender
- [x] Handle loading state during API request
- [x] Handle error states

### Phase 4: Results Table Component
- [x] Create `SeriesTable` component using TanStack Table
- [x] Define table columns: id, team1, team2, score, league, date
- [x] Implement table with proper styling (Tailwind)
- [x] Add sorting capabilities (if needed)
- [x] Handle empty state (no results)
- [x] Handle loading state
- [x] Handle error state

### Phase 5: Layout Implementation
- [x] Create main layout container with flexbox
- [x] Position sidebar on the left (fixed width ~250-300px)
- [x] Position table on the right (flex-1, takes remaining space)
- [x] Set proper minimum widths for both sections
- [x] Ensure responsive design considerations
- [x] Add proper spacing between sidebar and table
- [x] Style to match Amazon-like filters + results layout

### Phase 6: Integration
- [x] Integrate FiltersSidebar into series route
- [x] Integrate SeriesTable into series route
- [x] Connect search functionality to table data refresh
- [x] Test filter updates (should not trigger rerender until search)
- [x] Test search button triggers API call and URL update
- [x] Test URL params are correctly set/updated
- [x] Test table displays results correctly

### Phase 7: Refinement
- [x] Add loading indicators
- [x] Add error messages
- [x] Improve accessibility (ARIA labels, keyboard navigation)
- [x] Add hover states and visual feedback
- [x] Test edge cases (empty filters, invalid IDs)
- [x] Verify URL params work on page load (initial state from URL)
- [x] Add clear filters button
- [x] Add keyboard navigation (Enter to search)
- [x] Pre-populate inputs from URL params

### Phase 8: Future Enhancements (NOT FOR NOW)
- [ ] Replace text inputs with searchable dropdowns
- [ ] Fetch league names and images from API
- [ ] Fetch team names and images from API
- [ ] Implement search-by-name functionality
- [ ] Add league/team selection with autocomplete
- [ ] Add visual indicators (team logos, league icons)

## Notes
- Current implementation uses direct ID inputs for filters
- Future phase will replace with dropdowns that fetch names/images from API
- Keep filter updates local (no API calls) until search button is clicked
- URL updates should be seamless without triggering rerenders
