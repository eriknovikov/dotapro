# LeagueSelector Component - Implementation Plan

## Overview
Create a custom searchable league selector component using Radix UI, Tailwind CSS, and Downshift. The component will display static popular leagues by default and fetch search results from the backend when the user types.

## Requirements

### Functional Requirements
1. **Default State**: Show "Any League" as the default selected value
2. **Empty Query**: When search query is empty, display static popular leagues from `static_data.json`
3. **Search**: When user types, debounce for 250ms, then fetch leagues from backend using pg_trgm
4. **Results Display**: Swap dropdown items between static and search results
5. **"Any League" Option**: Always show "Any League" option (unless already selected)
6. **Loading State**: Show loading indicator while fetching
7. **Empty State**: Show "No leagues found" when search returns no results

### Tech Stack
- **Radix UI**: For accessible primitives (if needed)
- **Tailwind CSS**: For styling
- **Downshift**: For combobox/select functionality
- **TanStack Query**: For data fetching and caching
- **TypeScript**: For type safety

## Backend API Analysis

### Endpoint
- **URL**: `GET /filtersmetadata/leagues?q={query}`
- **Controller**: [`api/filtersmetadata/controller.go`](api/filtersmetadata/controller.go:48) - `SearchLeagues`
- **Model**: [`api/filtersmetadata/model.go`](api/filtersmetadata/model.go:33) - `SearchLeagues`

### Request Parameters
- `q` (required): Search query string
- Returns 400 if query is empty

### Response Type
```typescript
type LeagueSearchResult = {
    league_id: number
    name: string
}
```

### Backend Behavior
- Uses PostgreSQL pg_trgm extension with `name %> ?` operator (word similarity)
- Orders by `word_similarity(name, ?) DESC`
- Limits to 10 results
- 5 second timeout

## Frontend Data Sources

### Static Data
- **File**: [`ui/src/assets/static_data.json`](ui/src/assets/static_data.json:1)
- **Key**: `popular_leagues`
- **Structure**:
```typescript
{
    id: number
    name: string
}[]
```

### API Function
- **File**: [`ui/src/api/api.tsx`](ui/src/api/api.tsx:1)
- **Function**: `searchLeagues(query: string, signal: AbortSignal): Promise<LeagueSearchResult[]>`
- **Type**: `LeagueSearchResult`

### Existing Hooks
- **File**: [`ui/src/hooks/useDebounce.ts`](ui/src/hooks/useDebounce.ts:1)
- **Function**: `useDebounce<T>(value: T, delay: number): T`

## Component Architecture

### State Management
```typescript
interface LeagueSelectorState {
    // Input value (what user types)
    inputValue: string

    // Debounced search query
    debouncedQuery: string

    // Currently selected league
    selectedLeague: League | null

    // Is dropdown open
    isOpen: boolean

    // Is menu controlled (opened by user interaction)
    isMenuControlled: boolean
}
```

### Data Flow
```
User Types → inputValue updates → useDebounce(250ms) → debouncedQuery updates
                                                          ↓
                                                    TanStack Query fetches
                                                          ↓
                                                    searchResults update
                                                          ↓
                                                    items array recomputed
                                                          ↓
                                                    Dropdown renders items
```

### Items Array Logic
```typescript
const items = useMemo(() => {
    const anyLeague = { id: -1, name: "Any League" }

    if (debouncedQuery.trim().length === 0) {
        // Show static popular leagues
        return [anyLeague, ...popularLeagues]
    } else {
        // Show search results
        return [anyLeague, ...searchResults.map(l => ({
            id: l.league_id,
            name: l.name
        }))]
    }
}, [debouncedQuery, searchResults, popularLeagues])
```

## Component Structure

### Props Interface
```typescript
interface LeagueSelectorProps {
    onSelect: (leagueId: number | undefined) => void
    initialValue?: number
    className?: string
    id?: string
    "aria-label"?: string
}
```

### Component Layout
```
┌─────────────────────────────────────┐
│  [Input Field]              [Clear] │
│  ┌───────────────────────────────┐   │
│  │ Any League                    │   │
│  │ The International 2025        │   │
│  │ DreamLeague Season 26        │   │
│  │ ...                           │   │
│  └───────────────────────────────┘   │
└─────────────────────────────────────┘
```

### Key Features

#### 1. Input Field
- Controlled input with `inputValue` state
- Placeholder: "Search leagues..."
- Focus opens dropdown
- Typing triggers search

#### 2. Clear Button
- Shows when `inputValue` is not empty
- Clears input and resets selection to "Any League"

#### 3. Dropdown Menu
- Controlled by Downshift's `useSelect` hook
- Shows items based on search state
- Highlights hovered item
- Closes on item selection or click outside

#### 4. Loading State
- Shows Spinner component while fetching
- Positioned in the dropdown or inline

#### 5. Empty State
- Shows "No leagues found" when search returns empty
- Only shows when query is not empty

## Implementation Steps

### Step 1: Setup and Imports
```typescript
import { useSelect } from "downshift"
import { useQuery } from "@tanstack/react-query"
import { useState, useMemo, useRef, useEffect } from "react"
import { searchLeagues, type LeagueSearchResult } from "../api/api"
import { useDebounce } from "../hooks/useDebounce"
import { cn } from "../lib/utils"
import { Spinner } from "./ui/spinner"
import popularData from "../assets/static_data.json"
```

### Step 2: Define Constants and Types
```typescript
const ANY_LEAGUE_ID = -1
const ANY_LEAGUE_NAME = "Any League"

const anyLeagueItem = { id: ANY_LEAGUE_ID, name: ANY_LEAGUE_NAME }

interface League {
    id: number
    name: string
}
```

### Step 3: Initialize State
```typescript
const [inputValue, setInputValue] = useState("")
const [isMenuControlled, setIsMenuControlled] = useState(false)
const debouncedQuery = useDebounce(inputValue, 250)
const inputRef = useRef<HTMLInputElement>(null)
```

### Step 4: Fetch Search Results
```typescript
const { data: searchResults = [], isLoading, isFetching } = useQuery({
    queryKey: ["leagues", "search", debouncedQuery],
    queryFn: async ({ signal }) => {
        if (!debouncedQuery.trim()) return []
        return searchLeagues(debouncedQuery, signal)
    },
    enabled: debouncedQuery.trim().length > 0,
    staleTime: 5000,
})
```

### Step 5: Compute Items Array
```typescript
const items = useMemo(() => {
    if (debouncedQuery.trim().length === 0) {
        // Show static popular leagues
        return [anyLeagueItem, ...popularData.popular_leagues]
    } else {
        // Show search results
        return [anyLeagueItem, ...searchResults.map(l => ({
            id: l.league_id,
            name: l.name
        }))]
    }
}, [debouncedQuery, searchResults])
```

### Step 6: Initialize Downshift
```typescript
const initialItem = initialValue
    ? items.find(l => l.id === initialValue) || anyLeagueItem
    : anyLeagueItem

const {
    isOpen,
    selectedItem,
    getToggleButtonProps,
    getMenuProps,
    highlightedIndex,
    getItemProps,
    reset,
    openMenu,
    closeMenu,
} = useSelect({
    items,
    initialSelectedItem: initialItem,
    onSelectedItemChange: ({ selectedItem }) => {
        onSelect(selectedItem?.id === ANY_LEAGUE_ID ? undefined : selectedItem?.id)
        setInputValue(selectedItem?.name || "")
        setIsMenuControlled(false)
    },
})
```

### Step 7: Handle Input Changes
```typescript
const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value

    // If typing on "Any League", clear it first
    if (inputValue === ANY_LEAGUE_NAME && !isMenuControlled) {
        setInputValue(value)
        setIsMenuControlled(true)
    } else {
        setInputValue(value)
    }

    // Open menu if not already open
    if (!isOpen) {
        openMenu()
    }
}
```

### Step 8: Handle Input Focus
```typescript
const handleInputFocus = () => {
    if (!isOpen) {
        openMenu()
    }
}
```

### Step 9: Handle Clear
```typescript
const handleClear = () => {
    setInputValue("")
    setIsMenuControlled(false)
    reset()
    onSelect(undefined)
}
```

### Step 10: Sync Input Value with Selected Item
```typescript
useEffect(() => {
    if (selectedItem && !isMenuControlled && inputValue !== selectedItem.name) {
        setInputValue(selectedItem.name)
    }
}, [selectedItem, isMenuControlled])
```

### Step 11: Render Component
```typescript
return (
    <div className={cn("relative w-full", className)}>
        {/* Input Field */}
        <div className="relative">
            <input
                ref={inputRef}
                type="text"
                id={id}
                value={inputValue}
                onChange={handleInputChange}
                onFocus={handleInputFocus}
                placeholder="Search leagues..."
                aria-label={ariaLabel}
                className={/* Tailwind classes */}
            />

            {/* Clear/Icon Button */}
            <div className="absolute right-2 top-1/2 -translate-y-1/2">
                {isLoading || isFetching ? (
                    <Spinner size="sm" />
                ) : inputValue ? (
                    <button onClick={handleClear}>Clear</button>
                ) : (
                    <ChevronIcon />
                )}
            </div>

            {/* Hidden Toggle Button for Downshift */}
            <button {...getToggleButtonProps()} className="hidden" />
        </div>

        {/* Dropdown Menu */}
        <ul {...getMenuProps()} className={/* Tailwind classes */}>
            {isOpen && (
                <>
                    {items.length === 0 && !isLoading && (
                        <li>No leagues found</li>
                    )}
                    {items
                        .filter(league => league.id !== selectedItem?.id)
                        .map((league, index) => (
                            <li
                                key={league.id}
                                {...getItemProps({ item: league, index })}
                                className={/* Tailwind classes */}
                            >
                                {league.name}
                            </li>
                        ))}
                </>
            )}
        </ul>
    </div>
)
```

## Styling Considerations

### Theme Colors (from ARCHITECTURE.md)
- **Dire Red**: `hsl(0,84%,50%)`
- **Radiant Gold**: `hsl(38,92%,50%)`
- **Deep Dark**: `hsl(0,0%,4%)`

### Tailwind Classes
- Input: `h-10 w-full rounded-lg border px-3 py-2 text-sm bg-background border-border text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-800/50`
- Dropdown: `absolute z-50 mt-1.5 max-h-72 w-full overflow-auto rounded-lg border border-border-accent bg-background-card shadow-xl`
- Item: `cursor-pointer px-4 py-3 text-sm hover:bg-red-800/20`
- Highlighted: `bg-red-800/20`

## Edge Cases and Considerations

### 1. Initial Value Handling
- If `initialValue` is provided, find matching item in items array
- If not found, default to "Any League"

### 2. Query Empty State
- When `debouncedQuery` is empty, show static popular leagues
- When user clears input, reset to "Any League"

### 3. "Any League" Always Visible
- Always include `anyLeagueItem` at the start of items array
- Filter out if already selected (to avoid duplicate)

### 4. Loading State
- Show spinner while fetching
- Don't show spinner if query is empty

### 5. Empty Results
- Show "No leagues found" only when:
  - Query is not empty
  - Not loading
  - Results array is empty

### 6. Keyboard Navigation
- Downshift handles arrow keys, Enter, Escape
- Ensure focus management works correctly

### 7. Click Outside
- Downshift handles click outside to close menu
- Should not reset selection

## Testing Checklist

- [ ] Default shows "Any League" selected
- [ ] Empty query shows popular leagues
- [ ] Typing triggers search after 250ms
- [ ] Search results replace popular leagues
- [ ] "Any League" always appears (unless selected)
- [ ] Loading state shows spinner
- [ ] Empty state shows message
- [ ] Clear button resets to "Any League"
- [ ] Click outside closes menu
- [ ] Keyboard navigation works
- [ ] Initial value is set correctly
- [ ] onSelect callback fires correctly

## Files to Modify/Create

### Existing Files to Modify
1. `ui/src/components/LeagueSelector.tsx` - Complete rewrite

### Existing Files (No Changes Needed)
1. `ui/src/hooks/useDebounce.ts` - Already exists
2. `ui/src/api/api.tsx` - Already has `searchLeagues` function
3. `ui/src/assets/static_data.json` - Already has popular leagues

### New Files (None)
All functionality can be implemented in the existing component file.
