# PlayerSelector Implementation Plan

## Overview

This plan details the implementation of an enhanced PlayerSelector component that follows the pattern established by LeagueSelector and TeamSelector, with popular players by default and debounced search functionality using the API.

## Current State Analysis

### Current PlayerSelector.tsx
- Uses a regular `<input>` element (NOT ContentEditable div)
- Has debounced search (250ms via useDebounce hook)
- Fetches players from API (`searchPlayers` function)
- Has loading states and spinner
- Has initial value support via `getPlayerName` API call
- Shows popular players when input is empty
- Has keyboard navigation (ArrowDown/ArrowUp/Enter/Escape/Tab)
- Missing: ContentEditable div pattern (like TeamSelector/LeagueSelector)

### Reference Implementation Pattern (LeagueSelector & TeamSelector)
- Custom implementation (NOT using CustomSelect)
- **ContentEditable div** for search input (NOT regular input)
- Debounced search (250ms via useDebounce hook)
- Conditional items: Popular items when empty, search results when typing
- Keyboard navigation (ArrowDown/ArrowUp/Enter/Escape/Tab)
- Click outside handling with proper state cleanup
- Initial value support via separate API calls
- Loading states and spinner for API calls
- Proper cursor position handling in ContentEditable div

### Key Differences from HeroSelector
- **API calls required** - Players are fetched from `/filtersmetadata/players` endpoint
- **Loading states needed** - Spinner shown while fetching
- **Initial value API call** - Need to fetch player name by ID via `/filtersmetadata/player` endpoint
- **Uses pg_trgm backend search** - Similar to leagues/teams, not frontend search

---

## 1. Popular Players Selection

### Current Popular Players (in popular.json)

| ID | Name | Notes |
|----|------|-------|
| 86745914 | Clement Ivanov | Puppey |
| 154714040 | Jesse Vainikka | JerAx |
| 87278757 | Amer Al-Barkawi | Miracle- |
| 88271537 | Kuro Takhasomi | KuroKy |
| 84772440 | Lasse Aukusti | MATUMBAMAN |
| 111620041 | Syed Sumail | SumaiL |
| 94135151 | Topias Taavitsainen | Topson |
| 1295838 | Danil Ishutin | Dendi |

### Recommended Popular Players Update

**User will provide IDs for:**
- Nisha
- Topson (already exists: 94135151)
- bzm
- Other popular players as needed

**Note:** The user will provide the specific player IDs for the popular players section. The current `popular_players` array in `popular.json` should be updated with the new IDs provided by the user.

---

## 2. Implementation Steps

### Step 1: Update popular.json with new popular players

**File:** `ui/src/assets/static_data/popular.json`

Update the `popular_players` array with new IDs provided by the user:

```json
{
    "popular_teams": [...],
    "popular_leagues": [...],
    "popular_players": [
        { "id": <NISHA_ID>, "name": "Nisha" },
        { "id": 94135151, "name": "Topias Taavitsainen" },
        { "id": <BZM_ID>, "name": "bzm" },
        // ... other popular players as provided by user
    ],
    "popular_heroes": [...]
}
```

**Note:** User will provide the actual IDs for Nisha, bzm, and other players.

---

### Step 2: Rewrite PlayerSelector.tsx following TeamSelector/LeagueSelector pattern

**File:** `ui/src/components/PlayerSelector.tsx`

Complete rewrite following the reference pattern with these key characteristics:

1. **Use ContentEditable div** instead of regular `<input>` element
2. **API calls for search** - `searchPlayers` function
3. **Loading states** - Spinner shown while fetching
4. **Initial value API call** - `getPlayerName` function
5. **Popular players by default** - From `popular.json`

**New Component Structure:**

```typescript
// Imports
import { useQuery } from "@tanstack/react-query"
import { useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState } from "react"
import { getPlayerName, searchPlayers } from "../api"
import { useDebounce } from "../hooks"
import { cn, getPopularData } from "../lib"
import { Spinner } from "./Spinner"

// Types
interface PlayerSelectorProps {
    onSelect: (playerId: number | undefined) => void
    initialValue?: number
    className?: string
    id?: string
    "aria-label"?: string
    inputClassName?: string
}

interface Player {
    id: number
    name: string
}

// Constants
const DEBOUNCE_MS = 250
const STALE_TIME_MS = 5000
const NAVIGATION_KEYS = ["ArrowDown", "ArrowUp", "Enter"]

// Icons (ClearIcon, ChevronIcon - same as reference)

// Component
export function PlayerSelector({ ... }: PlayerSelectorProps) {
    // State
    const [isOpen, setIsOpen] = useState(false)
    const [inputValue, setInputValue] = useState("")
    const [highlightedIndex, setHighlightedIndex] = useState(-1)
    const [hasFetchedInitialValue, setHasFetchedInitialValue] = useState(false)

    // Refs
    const containerRef = useRef<HTMLDivElement>(null)
    const inputRef = useRef<HTMLDivElement>(null)
    const listRef = useRef<HTMLUListElement>(null)
    const itemRefs = useRef<(HTMLLIElement | null)[]>([])
    const isSelectingRef = useRef(false)
    const isExternalUpdateRef = useRef(false)
    const selectedPlayerIdRef = useRef<number | undefined>(undefined)
    const selectedPlayerNameRef = useRef<string | undefined>(undefined)
    const hasSetInitialValueRef = useRef(false)

    // Data fetching
    const debouncedQuery = useDebounce(inputValue, DEBOUNCE_MS)

    const {
        data: searchResults = [],
        isLoading,
        isFetching,
    } = useQuery({
        queryKey: ["players", "search", debouncedQuery],
        queryFn: async ({ signal }) => {
            if (!debouncedQuery.trim()) return []
            return searchPlayers(debouncedQuery, signal)
        },
        enabled: debouncedQuery.trim().length > 0,
        staleTime: STALE_TIME_MS,
    })

    // Fetch player name by ID when initialValue is provided
    const { data: playerNameData, isLoading: isPlayerNameLoading } = useQuery({
        queryKey: ["player", "name", initialValue],
        queryFn: async ({ signal }) => {
            if (initialValue === undefined) return null
            return getPlayerName(initialValue, signal)
        },
        enabled: initialValue !== undefined && inputValue === "" && !hasFetchedInitialValue,
        staleTime: 10 * 60 * 1000, // 10 minutes
    })

    // ... rest follows TeamSelector/LeagueSelector pattern
}
```

---

### Step 3: Implement ContentEditable Input

**File:** `ui/src/components/PlayerSelector.tsx`

Replace the regular `<input>` element with a ContentEditable div:

```typescript
const renderInput = () => (
    <div className="relative">
        {inputValue === "" && (
            <span className="text-foreground-muted pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-sm select-none">
                Type to search players...
            </span>
        )}
        <div
            ref={inputRef}
            id={id}
            contentEditable
            suppressContentEditableWarning
            onInput={handleInputChange}
            onPaste={handleInputPaste}
            onBlur={handleInputBlur}
            onFocus={handleInputFocus}
            onKeyDown={handleKeyDown}
            tabIndex={0}
            aria-label={ariaLabel || "Search players"}
            aria-autocomplete="list"
            aria-controls={listboxId}
            aria-expanded={isOpen}
            aria-activedescendant={highlightedIndex >= 0 ? `${listboxId}-option-${highlightedIndex}` : undefined}
            role="combobox"
            className={cn(
                "min-h-10 w-full rounded-lg border px-3 py-2 pr-10 text-sm",
                "bg-background border-border text-foreground",
                "focus-visible:ring-2 focus-visible:ring-red-800/50 focus-visible:outline-none",
                "transition-all duration-200 ease-in-out",
                "wrap-break-word whitespace-pre-wrap",
                inputClassName,
            )}
        />
    </div>
)
```

---

### Step 4: Add Input Event Handlers

**File:** `ui/src/components/PlayerSelector.tsx`

Implement the ContentEditable-specific event handlers:

```typescript
const handleInputChange = useCallback((e: React.FormEvent<HTMLDivElement>) => {
    const text = e.currentTarget.textContent || ""
    setInputValue(text)
    if (!isExternalUpdateRef.current) {
        setIsOpen(true)
    }
}, [])

const handleInputPaste = useCallback((e: React.ClipboardEvent<HTMLDivElement>) => {
    e.preventDefault()
    const text = e.clipboardData.getData("text/plain")
    document.execCommand("insertText", false, text)
}, [])

const handleInputBlur = useCallback(() => {
    // Ensure inputValue is in sync with div content on blur
    if (inputRef.current) {
        const text = inputRef.current.textContent || ""
        // If input is cleared, remove the filter entirely
        if (text === "" && selectedPlayerIdRef.current !== undefined) {
            onSelect(undefined)
            selectedPlayerIdRef.current = undefined
            selectedPlayerNameRef.current = undefined
        }
        // Only update state if value actually changed to avoid unnecessary re-renders
        if (text !== inputValue) {
            setInputValue(text)
        }
    }
}, [inputValue, onSelect])

const handleInputFocus = useCallback(() => {
    if (!isSelectingRef.current && !isOpen && !isExternalUpdateRef.current) {
        setIsOpen(true)
    }
}, [isOpen])
```

---

### Step 5: Add Cursor Position Handling

**File:** `ui/src/components/PlayerSelector.tsx`

Add the cursor position handling for ContentEditable div:

```typescript
// Sync div content with inputValue when it changes externally (e.g., from selection)
useLayoutEffect(() => {
    if (inputRef.current) {
        // Always clear the div content when inputValue is empty to ensure placeholder shows
        if (inputValue === "") {
            inputRef.current.textContent = ""
            isExternalUpdateRef.current = false
            return
        }

        if (isExternalUpdateRef.current) {
            // Save the current selection/cursor position
            const selection = window.getSelection()
            const range = selection && selection.rangeCount > 0 ? selection.getRangeAt(0) : null
            const cursorOffset = range?.startOffset ?? 0

            inputRef.current.textContent = inputValue

            // Restore the cursor position if it was within the input
            if (selection && range && inputRef.current.contains(range.commonAncestorContainer)) {
                const newRange = document.createRange()
                const textNode = inputRef.current.firstChild as Text
                if (textNode) {
                    const newOffset = Math.min(cursorOffset, textNode.length)
                    newRange.setStart(textNode, newOffset)
                    newRange.setEnd(textNode, newOffset)
                    selection.removeAllRanges()
                    selection.addRange(newRange)
                }
            }

            isExternalUpdateRef.current = false
        }
    }
}, [inputValue])

// Ensure text content is set when ref becomes available (for initial value on page refresh)
useLayoutEffect(() => {
    if (inputRef.current && inputValue && inputRef.current.textContent !== inputValue) {
        inputRef.current.textContent = inputValue
    }
}, [inputValue, playerNameData])
```

---

### Step 6: Add Keyboard Navigation

**File:** `ui/src/components/PlayerSelector.tsx`

Follow the exact pattern from TeamSelector/LeagueSelector:

```typescript
const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
        // Open dropdown on navigation keys when closed
        if (!isOpen) {
            if (NAVIGATION_KEYS.includes(e.key as string)) {
                // Prevent newline insertion when pressing Enter to open dropdown
                if (e.key === "Enter") {
                    e.preventDefault()
                    e.stopPropagation()
                }
                setIsOpen(true)
            }
            return
        }

        switch (e.key) {
            case "ArrowDown":
                e.preventDefault()
                e.stopPropagation()
                setHighlightedIndex(prev => (prev < items.length - 1 ? prev + 1 : prev))
                break
            case "ArrowUp":
                e.preventDefault()
                e.stopPropagation()
                setHighlightedIndex(prev => (prev > 0 ? prev - 1 : 0))
                break
            case "Enter":
                e.preventDefault()
                e.stopPropagation()
                if (highlightedIndex >= 0 && items[highlightedIndex]) {
                    handleSelect(items[highlightedIndex])
                }
                break
            case "Escape":
                e.stopPropagation()
                closeDropdown()
                inputRef.current?.focus()
                break
            case "Tab":
                e.stopPropagation()
                closeDropdown()
                break
        }
    },
    [isOpen, items, highlightedIndex, handleSelect, closeDropdown],
)
```

---

### Step 7: Handle Initial Values

**File:** `ui/src/components/PlayerSelector.tsx`

Since player data is fetched from API, initial value handling requires API call:

```typescript
// Set initial value when component mounts or initialValue changes
useLayoutEffect(() => {
    if (initialValue !== undefined && !hasSetInitialValueRef.current) {
        // First try to find in items (popular players or search results)
        const initialItem = items.find(p => p.id === initialValue)

        if (initialItem) {
            // Always update if we have the item, regardless of current inputValue
            isExternalUpdateRef.current = true
            setInputValue(initialItem.name)
            selectedPlayerIdRef.current = initialValue
            selectedPlayerNameRef.current = initialItem.name
            hasSetInitialValueRef.current = true
        } else if (playerNameData && !isPlayerNameLoading) {
            // If not found in items but we have the name from API, use that
            isExternalUpdateRef.current = true
            setInputValue(playerNameData.name)
            selectedPlayerIdRef.current = initialValue
            selectedPlayerNameRef.current = playerNameData.name
            hasSetInitialValueRef.current = true
        }
    } else if (initialValue === undefined && selectedPlayerIdRef.current !== undefined) {
        // Clear the input when initialValue becomes undefined
        setInputValue("")
        selectedPlayerIdRef.current = undefined
        selectedPlayerNameRef.current = undefined
        setHasFetchedInitialValue(false)
        hasSetInitialValueRef.current = false
    }
}, [initialValue, items, playerNameData, isPlayerNameLoading])
```

---

### Step 8: Add Render Helpers

**File:** `ui/src/components/PlayerSelector.tsx`

Add the render helper functions:

```typescript
const renderInputSuffix = () => (
    <div className="absolute top-1/2 right-2 flex -translate-y-1/2 items-center gap-1">
        {showSpinner && <Spinner size="sm" aria-label="Loading players" />}
        {showClearButton && (
            <button
                type="button"
                onClick={handleClear}
                className="text-foreground-muted hover:text-foreground transition-colors"
                aria-label="Clear search"
            >
                <ClearIcon />
            </button>
        )}
        {!showSpinner && !showClearButton && (
            <button
                type="button"
                onClick={() => {
                    setIsOpen(!isOpen)
                    inputRef.current?.focus()
                }}
                className="text-foreground-muted hover:text-foreground transition-colors"
                aria-label={isOpen ? "Close dropdown" : "Open dropdown"}
                aria-expanded={isOpen}
            >
                <ChevronIcon isOpen={isOpen} />
            </button>
        )}
    </div>
)

const renderEmptyState = () => (
    <li role="option" aria-selected="false" className="text-foreground-muted px-4 py-3 text-sm">
        No players found
    </li>
)

const renderPlayerItem = (player: Player, index: number) => (
    <li
        key={player.id}
        ref={el => {
            itemRefs.current[index] = el
        }}
        id={`${listboxId}-option-${index}`}
        role="option"
        aria-selected={highlightedIndex === index}
        onClick={e => handleItemClick(e, player)}
        onMouseEnter={() => handleItemMouseEnter(index)}
        className={cn(
            "flex cursor-pointer items-center gap-3 px-4 py-3 text-sm",
            highlightedIndex === index && "bg-red-800/20",
        )}
    >
        <span className="truncate">{player.name}</span>
    </li>
)

const renderDropdown = () => (
    <ul
        ref={listRef}
        id={listboxId}
        role="listbox"
        aria-label={ariaLabel || "Players"}
        className="border-border-accent bg-background-card absolute top-full z-50 mt-1 max-h-72 w-full overflow-auto rounded-lg border shadow-xl"
    >
        {showEmptyState ? renderEmptyState() : items.map(renderPlayerItem)}
    </ul>
)
```

---

### Step 9: Test and Verify

See Testing Checklist below.

---

## 3. Code Structure

### New PlayerSelector Component Structure

```typescript
// ============================================================================
// Imports
// ============================================================================
import { useQuery } from "@tanstack/react-query"
import { useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState } from "react"
import { getPlayerName, searchPlayers } from "../api"
import { useDebounce } from "../hooks"
import { cn, getPopularData } from "../lib"
import { Spinner } from "./Spinner"

// ============================================================================
// Types
// ============================================================================
interface PlayerSelectorProps { ... }
interface Player { ... }

// ============================================================================
// Constants
// ============================================================================
const DEBOUNCE_MS = 250
const STALE_TIME_MS = 5000
const NAVIGATION_KEYS = ["ArrowDown", "ArrowUp", "Enter"]

// ============================================================================
// Icons
// ============================================================================
const ClearIcon = () => { ... }
const ChevronIcon = ({ isOpen }: { isOpen: boolean }) => { ... }

// ============================================================================
// Component
// ============================================================================
export function PlayerSelector({ ... }: PlayerSelectorProps) {
    // ---------------------------------------------------------------------------
    // State
    // ---------------------------------------------------------------------------
    const [isOpen, setIsOpen] = useState(false)
    const [inputValue, setInputValue] = useState("")
    const [highlightedIndex, setHighlightedIndex] = useState(-1)
    const [hasFetchedInitialValue, setHasFetchedInitialValue] = useState(false)

    // ---------------------------------------------------------------------------
    // Refs
    // ---------------------------------------------------------------------------
    const containerRef = useRef<HTMLDivElement>(null)
    const inputRef = useRef<HTMLDivElement>(null)
    const listRef = useRef<HTMLUListElement>(null)
    const itemRefs = useRef<(HTMLLIElement | null)[]>([])
    const isSelectingRef = useRef(false)
    const isExternalUpdateRef = useRef(false)
    const selectedPlayerIdRef = useRef<number | undefined>(undefined)
    const selectedPlayerNameRef = useRef<string | undefined>(undefined)
    const hasSetInitialValueRef = useRef(false)

    // ---------------------------------------------------------------------------
    // Data fetching
    // ---------------------------------------------------------------------------
    const debouncedQuery = useDebounce(inputValue, DEBOUNCE_MS)
    const { data: searchResults, isLoading, isFetching } = useQuery({ ... })
    const { data: playerNameData, isLoading: isPlayerNameLoading } = useQuery({ ... })

    // ---------------------------------------------------------------------------
    // Computed values
    // ---------------------------------------------------------------------------
    const items = useMemo(() => { ... }, [debouncedQuery, searchResults])
    const listboxId = `${id}-listbox`
    const hasQuery = debouncedQuery.trim().length > 0
    const showSpinner = isFetching && hasQuery
    const showClearButton = inputValue && !showSpinner
    const showEmptyState = items.length === 0 && !isLoading && !isFetching

    // ---------------------------------------------------------------------------
    // Helper functions
    // ---------------------------------------------------------------------------
    const clearSelectedPlayer = useCallback(() => { ... }, [onSelect])
    const shouldClearPlayerOnClose = useCallback(() => { ... }, [inputValue])
    const closeDropdown = useCallback(() => { ... }, [shouldClearPlayerOnClose, clearSelectedPlayer])

    // ---------------------------------------------------------------------------
    // Effects
    // ---------------------------------------------------------------------------
    // Set initial value
    useLayoutEffect(() => { ... }, [initialValue, items, playerNameData, isPlayerNameLoading])

    // Mark initial value as fetched
    useEffect(() => { ... }, [playerNameData, hasFetchedInitialValue])

    // Close dropdown when clicking outside
    useEffect(() => { ... }, [closeDropdown])

    // Reset highlighted index when items change
    useLayoutEffect(() => { ... }, [items])

    // Scroll highlighted item into view
    useEffect(() => { ... }, [highlightedIndex])

    // Sync div content with inputValue
    useLayoutEffect(() => { ... }, [inputValue])

    // Ensure text content is set when ref becomes available
    useLayoutEffect(() => { ... }, [inputValue, playerNameData])

    // ---------------------------------------------------------------------------
    // Event handlers
    // ---------------------------------------------------------------------------
    const handleSelect = useCallback((player: Player) => { ... }, [onSelect])
    const handleClear = useCallback(() => { ... }, [onSelect])
    const handleInputFocus = useCallback(() => { ... }, [isOpen])
    const handleInputChange = useCallback((e: React.FormEvent<HTMLDivElement>) => { ... }, [])
    const handleInputPaste = useCallback((e: React.ClipboardEvent<HTMLDivElement>) => { ... }, [])
    const handleInputBlur = useCallback(() => { ... }, [inputValue, onSelect])
    const handleItemMouseEnter = useCallback((index: number) => { ... }, [])
    const handleItemClick = useCallback((e: React.MouseEvent, player: Player) => { ... }, [handleSelect])

    // ---------------------------------------------------------------------------
    // Keyboard navigation
    // ---------------------------------------------------------------------------
    const handleKeyDown = useCallback((e: React.KeyboardEvent) => { ... }, [isOpen, items, highlightedIndex, handleSelect, closeDropdown])

    // ---------------------------------------------------------------------------
    // Render helpers
    // ---------------------------------------------------------------------------
    const renderInput = () => { ... }
    const renderInputSuffix = () => { ... }
    const renderEmptyState = () => { ... }
    const renderPlayerItem = (player: Player, index: number) => { ... }
    const renderDropdown = () => { ... }

    // ---------------------------------------------------------------------------
    // Render
    // ---------------------------------------------------------------------------
    return (
        <div ref={containerRef} className={cn("relative w-full", className)}>
            <div className="relative">
                {renderInput()}
                {renderInputSuffix()}
                {isOpen && renderDropdown()}
            </div>
        </div>
    )
}
```

### State Variables Summary

| Variable | Type | Purpose |
|----------|------|---------|
| `isOpen` | `boolean` | Controls dropdown visibility |
| `inputValue` | `string` | Current search input text |
| `highlightedIndex` | `number` | Index of currently highlighted item |
| `debouncedQuery` | `string` | Debounced version of inputValue |
| `searchResults` | `Player[]` | Search results from API |
| `playerNameData` | `{ name: string } \| null` | Player name from API for initial value |
| `items` | `Player[]` | Either popular players or search results |

### Items Computation Logic

```typescript
const items = useMemo(() => {
    if (debouncedQuery.trim().length === 0) {
        // Show popular players when no search query
        return getPopularData().popular_players.map(p => ({
            id: p.id,
            name: p.name,
        }))
    }
    // Show search results when typing
    return (searchResults || []).map(p => ({ id: p.player_id, name: p.name }))
}, [debouncedQuery, searchResults])
```

---

## 4. File Changes Summary

| File | Changes |
|------|---------|
| `ui/src/assets/static_data/popular.json` | Update `popular_players` array with new IDs (Nisha, Topson, bzm, etc.) |
| `ui/src/components/PlayerSelector.tsx` | Complete rewrite following TeamSelector/LeagueSelector pattern with ContentEditable div |

---

## 5. Testing Checklist

### Functional Testing

- [ ] **Popular Players Display**
  - [ ] When dropdown opens with empty input, shows popular players
  - [ ] Popular players are displayed in correct order
  - [ ] Each popular player shows correct name

- [ ] **Search Functionality**
  - [ ] Typing in search input filters players
  - [ ] Search is debounced (250ms)
  - [ ] Search works on player name field
  - [ ] Partial matches work (e.g., "top" matches "Topson")
  - [ ] Case-insensitive search works
  - [ ] Empty search results show "No players found" message
  - [ ] Spinner shows while fetching search results

- [ ] **Player Selection**
  - [ ] Clicking a player selects it
  - [ ] Selected player name appears in input
  - [ ] `onSelect` callback is called with correct player_id
  - [ ] Dropdown closes after selection
  - [ ] Input is focused after selection

- [ ] **Clear Functionality**
  - [ ] Clear button appears when player is selected
  - [ ] Clicking clear button removes selection
  - [ ] `onSelect` callback is called with `undefined`
  - [ ] Input is cleared
  - [ ] Popular players are shown again

- [ ] **Keyboard Navigation**
  - [ ] ArrowDown moves highlight down
  - [ ] ArrowUp moves highlight up
  - [ ] Enter selects highlighted player
  - [ ] Escape closes dropdown
  - [ ] Tab closes dropdown
  - [ ] Navigation keys open dropdown when closed
  - [ ] Cursor position is preserved in ContentEditable div

- [ ] **Click Outside**
  - [ ] Clicking outside dropdown closes it
  - [ ] State is properly cleaned up on close

- [ ] **Initial Value Handling**
  - [ ] When `initialValue` is provided, player name is shown
  - [ ] When `initialValue` is in popular players, name is shown immediately
  - [ ] When `initialValue` is not in popular players, name is fetched from API
  - [ ] When `initialValue` changes, input updates
  - [ ] When `initialValue` is `undefined`, input is cleared
  - [ ] Spinner shows while fetching initial player name

- [ ] **ContentEditable Behavior**
  - [ ] Placeholder text shows when input is empty
  - [ ] Text can be typed and edited
  - [ ] Paste works correctly
  - [ ] Cursor position is preserved after selection
  - [ ] Blur event properly syncs div content with state

- [ ] **Accessibility**
  - [ ] Component has proper ARIA labels
  - [ ] Keyboard navigation works
  - [ ] Screen reader announces changes
  - [ ] Focus management is correct

### Visual Testing

- [ ] Dropdown styling matches TeamSelector/LeagueSelector
- [ ] Highlighted item has correct background color
- [ ] Icons (Clear, Chevron) display correctly
- [ ] Placeholder text shows when input is empty
- [ ] Spinner displays correctly during loading
- [ ] Responsive design works on mobile

### Performance Testing

- [ ] Search is fast (API response time)
- [ ] No unnecessary re-renders
- [ ] Debounce works correctly
- [ ] ContentEditable cursor handling is smooth

### Integration Testing

- [ ] Works correctly in MatchFilters
- [ ] Works correctly in SeriesFilters
- [ ] URL params update correctly on selection
- [ ] Page refresh preserves selected player
- [ ] API calls are made correctly

---

## 6. Implementation Notes

### Key Differences from HeroSelector

1. **API calls required** - Players are fetched from `/filtersmetadata/players` endpoint
2. **Loading states needed** - Spinner shown while fetching
3. **Initial value API call** - Need to fetch player name by ID via `/filtersmetadata/player` endpoint
4. **Uses pg_trgm backend search** - Similar to leagues/teams, not frontend search
5. **ContentEditable div** - Uses ContentEditable div pattern like TeamSelector/LeagueSelector

### Dependencies

- `useDebounce` hook (already exists)
- `getPopularData` function (already exists)
- `searchPlayers` API function (already exists)
- `getPlayerName` API function (already exists)
- `Spinner` component (already exists)

### Browser Compatibility

- ContentEditable is widely supported
- All modern browsers support the required features

### API Endpoints

- **Search Players:** `GET /filtersmetadata/players?q={query}`
- **Get Player Name:** `GET /filtersmetadata/player?id={id}`

### Popular Players

The user will provide the specific player IDs for the popular players section. The current `popular_players` array in `popular.json` should be updated with the new IDs provided by the user.

---

## 7. Next Steps

After this plan is approved:

1. **Get popular player IDs from user** - User will provide IDs for Nisha, Topson, bzm, etc.
2. Implement Step 1: Update popular.json with new popular players
3. Implement Step 2: Rewrite PlayerSelector.tsx with ContentEditable div
4. Implement Step 3: Implement ContentEditable Input
5. Implement Step 4: Add Input Event Handlers
6. Implement Step 5: Add Cursor Position Handling
7. Implement Step 6: Add Keyboard Navigation
8. Implement Step 7: Handle Initial Values
9. Implement Step 8: Add Render Helpers
10. Implement Step 9: Test and verify

Each step can be implemented and tested independently.
