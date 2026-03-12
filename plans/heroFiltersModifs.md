# HeroSelector Implementation Plan

## Overview

This plan details the implementation of an enhanced HeroSelector component that follows the pattern established by LeagueSelector and TeamSelector, with popular heroes by default and debounced search functionality.

## Current State Analysis

### Current HeroSelector.tsx
- Uses `CustomSelect` component (simple select, no search)
- Loads all 137 heroes from `heroes.json` at once
- Has local filtering on `displayName` and `name` fields
- Missing: "Popular items by default" pattern, debounced search

### Reference Implementation Pattern (LeagueSelector & TeamSelector)
- Custom implementation (NOT using CustomSelect)
- ContentEditable div for search input
- Debounced search (250ms via useDebounce hook)
- Conditional items: Popular items when empty, search results when typing
- Keyboard navigation (ArrowDown/ArrowUp/Enter/Escape/Tab)
- Click outside handling with proper state cleanup
- Initial value support via separate API calls

### Key Differences from League/Team Selectors
- **No backend API needed** - All hero data is in `heroes.json`
- **Search is fully frontend** - No pg_trgm backend query needed
- **Faster search** - Can use more sophisticated algorithms since data is local

---

## 1. Search Algorithm Analysis

### Options Comparison

| Algorithm | Description | Pros | Cons | Performance (137 heroes) |
|-----------|-------------|------|------|--------------------------|
| **Simple Includes** | `str.includes(query)` | Fast, simple | No fuzzy matching, exact substring only | ~0.1ms |
| **Regex** | `new RegExp(query, 'i')` | Flexible, case-insensitive | Can be slow with complex patterns | ~0.2-0.5ms |
| **Fuzzy Matching** | Levenshtein distance, Fuse.js | Handles typos, partial matches | Slower, more complex | ~1-5ms |
| **Trigram-like** | Split into 3-char sequences, match | Similar to pg_trgm, good for partial matches | More complex implementation | ~0.5-1ms |
| **Weighted Multi-field** | Score matches on displayName vs name | Better relevance | More complex | ~0.3-0.8ms |

### Recommended Approach: **Weighted Multi-field Search with Trigram-like Enhancement**

**Why this approach is optimal:**

1. **Performance**: With only 137 heroes, even complex algorithms are fast (<5ms)
2. **User Experience**: Handles partial matches, typos, and prioritizes displayName over name
3. **Consistency**: Similar to pg_trgm behavior used in League/Team selectors
4. **Maintainability**: Clear, understandable algorithm

**Algorithm Details:**

```typescript
function searchHeroes(query: string, heroes: Hero[]): Hero[] {
    if (!query.trim()) return []

    const normalizedQuery = query.toLowerCase().trim()
    const queryTrigrams = generateTrigrams(normalizedQuery)

    return heroes
        .map(hero => {
            const displayNameLower = hero.displayName.toLowerCase()
            const nameLower = hero.name.toLowerCase()

            // Calculate match score
            let score = 0

            // Exact match on displayName (highest priority)
            if (displayNameLower === normalizedQuery) score += 100
            // Exact match on name
            else if (nameLower === normalizedQuery) score += 80

            // Starts with displayName
            if (displayNameLower.startsWith(normalizedQuery)) score += 50
            // Starts with name
            else if (nameLower.startsWith(normalizedQuery)) score += 30

            // Contains displayName
            if (displayNameLower.includes(normalizedQuery)) score += 20
            // Contains name
            if (nameLower.includes(normalizedQuery)) score += 10

            // Trigram matching (partial matches)
            const displayNameTrigrams = generateTrigrams(displayNameLower)
            const trigramMatches = queryTrigrams.filter(t => displayNameTrigrams.includes(t)).length
            score += trigramMatches * 2

            return { hero, score }
        })
        .filter(result => result.score > 0)
        .sort((a, b) => b.score - a.score)
        .map(result => result.hero)
}

function generateTrigrams(text: string): string[] {
    const trigrams: string[] = []
    for (let i = 0; i < text.length - 2; i++) {
        trigrams.push(text.slice(i, i + 3))
    }
    return trigrams
}
```

---

## 2. Popular Heroes Selection

### Recommended 7 Popular Heroes

| ID | Name | DisplayName | Justification |
|----|------|-------------|---------------|
| 74 | invoker | Invoker | Most complex hero, iconic, high skill ceiling |
| 34 | tinker | Tinker | Popular mid hero, complex mechanics |
| 14 | pudge | Pudge | Most iconic hero, highest pick rate historically |
| 105 | meepo | Meepo | Unique multi-unit control, high skill ceiling |
| 22 | zuus | Zeus | Simple but iconic, high damage output |
| 1 | antimage | Anti-Mage | Classic carry, iconic design |
| 8 | juggernaut | Juggernaut | Popular carry, iconic spinning animation |

**Justification:**
- Mix of complexity levels (from simple Zeus to complex Invoker/Meepo)
- Iconic heroes that represent Dota 2's identity
- High pick/ban rates in professional matches
- Visually distinct and recognizable

---

## 3. Implementation Steps

### Step 1: Add popular_heroes to popular.json

**File:** `ui/src/assets/static_data/popular.json`

Add `popular_heroes` array to the existing structure:

```json
{
    "popular_teams": [...],
    "popular_leagues": [...],
    "popular_players": [...],
    "popular_heroes": [
        { "id": 74, "name": "Invoker" },
        { "id": 34, "name": "Tinker" },
        { "id": 14, "name": "Pudge" },
        { "id": 105, "name": "Meepo" },
        { "id": 22, "name": "Zeus" },
        { "id": 1, "name": "Anti-Mage" },
        { "id": 8, "name": "Juggernaut" }
    ]
}
```

**File:** `ui/src/lib/dotautils.ts`

Update `PopularData` interface and `getPopularData()` function:

```typescript
export interface PopularHero {
    id: number
    name: string
}

export interface PopularData {
    popular_leagues: PopularLeague[]
    popular_teams: PopularTeam[]
    popular_players: PopularPlayer[]
    popular_heroes: PopularHero[]
}
```

---

### Step 2: Rewrite HeroSelector.tsx following LeagueSelector/TeamSelector pattern

**File:** `ui/src/components/HeroSelector.tsx`

Complete rewrite following the reference pattern with these key differences:

1. **No API calls** - All data from `heroes.json`
2. **No loading states** - Data is immediately available
3. **No spinner** - No async operations

**New Component Structure:**

```typescript
// Imports
import { useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState } from "react"
import { useDebounce } from "../hooks"
import { cn, getPopularData } from "../lib"
import heroesData from "@/assets/static_data/heroes.json"

// Types
interface HeroSelectorProps {
    onSelect: (heroId: number | undefined) => void
    initialValue?: number
    className?: string
    id?: string
    "aria-label"?: string
    inputClassName?: string
}

interface Hero {
    id: number
    name: string
    displayName: string
}

// Constants
const DEBOUNCE_MS = 250
const NAVIGATION_KEYS = ["ArrowDown", "ArrowUp", "Enter"]

// Icons (ClearIcon, ChevronIcon - same as reference)

// Component
export function HeroSelector({ ... }: HeroSelectorProps) {
    // State
    const [isOpen, setIsOpen] = useState(false)
    const [inputValue, setInputValue] = useState("")
    const [highlightedIndex, setHighlightedIndex] = useState(-1)

    // Refs
    const containerRef = useRef<HTMLDivElement>(null)
    const inputRef = useRef<HTMLDivElement>(null)
    const listRef = useRef<HTMLUListElement>(null)
    const itemRefs = useRef<(HTMLLIElement | null)[]>([])
    const isSelectingRef = useRef(false)
    const isExternalUpdateRef = useRef(false)
    const selectedHeroIdRef = useRef<number | undefined>(undefined)
    const selectedHeroNameRef = useRef<string | undefined>(undefined)
    const hasSetInitialValueRef = useRef(false)

    // Data (no API calls - all from heroes.json)
    const debouncedQuery = useDebounce(inputValue, DEBOUNCE_MS)

    // Convert heroes.json to array
    const allHeroes: Hero[] = useMemo(() => {
        return Object.values(heroesData).map(hero => ({
            id: hero.id,
            name: hero.name,
            displayName: hero.displayName,
        }))
    }, [])

    // Search results (frontend search)
    const searchResults = useMemo(() => {
        if (!debouncedQuery.trim()) return []
        return searchHeroes(debouncedQuery, allHeroes)
    }, [debouncedQuery, allHeroes])

    // Computed values
    const items = useMemo(() => {
        if (debouncedQuery.trim().length === 0) {
            return getPopularData().popular_heroes.map(h => ({
                id: h.id,
                name: h.name,
                displayName: allHeroes.find(hero => hero.id === h.id)?.displayName || h.name
            }))
        }
        return searchResults
    }, [debouncedQuery, searchResults, allHeroes])

    // ... rest follows LeagueSelector/TeamSelector pattern
}
```

---

### Step 3: Implement Search Algorithm

**File:** `ui/src/components/HeroSelector.tsx`

Add the search algorithm function inside the component file:

```typescript
// ============================================================================
// Search Algorithm
// ============================================================================

function searchHeroes(query: string, heroes: Hero[]): Hero[] {
    if (!query.trim()) return []

    const normalizedQuery = query.toLowerCase().trim()
    const queryTrigrams = generateTrigrams(normalizedQuery)

    return heroes
        .map(hero => {
            const displayNameLower = hero.displayName.toLowerCase()
            const nameLower = hero.name.toLowerCase()

            // Calculate match score
            let score = 0

            // Exact match on displayName (highest priority)
            if (displayNameLower === normalizedQuery) score += 100
            // Exact match on name
            else if (nameLower === normalizedQuery) score += 80

            // Starts with displayName
            if (displayNameLower.startsWith(normalizedQuery)) score += 50
            // Starts with name
            else if (nameLower.startsWith(normalizedQuery)) score += 30

            // Contains displayName
            if (displayNameLower.includes(normalizedQuery)) score += 20
            // Contains name
            if (nameLower.includes(normalizedQuery)) score += 10

            // Trigram matching (partial matches)
            const displayNameTrigrams = generateTrigrams(displayNameLower)
            const trigramMatches = queryTrigrams.filter(t => displayNameTrigrams.includes(t)).length
            score += trigramMatches * 2

            return { hero, score }
        })
        .filter(result => result.score > 0)
        .sort((a, b) => b.score - a.score)
        .map(result => result.hero)
}

function generateTrigrams(text: string): string[] {
    const trigrams: string[] = []
    for (let i = 0; i < text.length - 2; i++) {
        trigrams.push(text.slice(i, i + 3))
    }
    return trigrams
}
```

---

### Step 4: Add Keyboard Navigation

**File:** `ui/src/components/HeroSelector.tsx`

Follow the exact pattern from LeagueSelector/TeamSelector:

```typescript
const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
        // Open dropdown on navigation keys when closed
        if (!isOpen) {
            if (NAVIGATION_KEYS.includes(e.key as string)) {
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

### Step 5: Handle Initial Values

**File:** `ui/src/components/HeroSelector.tsx`

Since all hero data is local, initial value handling is simpler:

```typescript
// Set initial value when component mounts or initialValue changes
useLayoutEffect(() => {
    if (initialValue !== undefined && !hasSetInitialValueRef.current) {
        // Find hero in allHeroes (since we have all data locally)
        const initialHero = allHeroes.find(h => h.id === initialValue)

        if (initialHero) {
            isExternalUpdateRef.current = true
            setInputValue(initialHero.displayName)
            selectedHeroIdRef.current = initialValue
            selectedHeroNameRef.current = initialHero.displayName
            hasSetInitialValueRef.current = true
        }
    } else if (initialValue === undefined && selectedHeroIdRef.current !== undefined) {
        // Clear the input when initialValue becomes undefined
        setInputValue("")
        selectedHeroIdRef.current = undefined
        selectedHeroNameRef.current = undefined
        hasSetInitialValueRef.current = false
    }
}, [initialValue, allHeroes])
```

---

### Step 6: Test and Verify

See Testing Checklist below.

---

## 4. Code Structure

### New HeroSelector Component Structure

```typescript
// ============================================================================
// Imports
// ============================================================================
import { useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState } from "react"
import { useDebounce } from "../hooks"
import { cn, getPopularData } from "../lib"
import heroesData from "@/assets/static_data/heroes.json"

// ============================================================================
// Types
// ============================================================================
interface HeroSelectorProps { ... }
interface Hero { ... }

// ============================================================================
// Constants
// ============================================================================
const DEBOUNCE_MS = 250
const NAVIGATION_KEYS = ["ArrowDown", "ArrowUp", "Enter"]

// ============================================================================
// Icons
// ============================================================================
const ClearIcon = () => { ... }
const ChevronIcon = ({ isOpen }: { isOpen: boolean }) => { ... }

// ============================================================================
// Search Algorithm
// ============================================================================
function searchHeroes(query: string, heroes: Hero[]): Hero[] { ... }
function generateTrigrams(text: string): string[] { ... }

// ============================================================================
// Component
// ============================================================================
export function HeroSelector({ ... }: HeroSelectorProps) {
    // ---------------------------------------------------------------------------
    // State
    // ---------------------------------------------------------------------------
    const [isOpen, setIsOpen] = useState(false)
    const [inputValue, setInputValue] = useState("")
    const [highlightedIndex, setHighlightedIndex] = useState(-1)

    // ---------------------------------------------------------------------------
    // Refs
    // ---------------------------------------------------------------------------
    const containerRef = useRef<HTMLDivElement>(null)
    const inputRef = useRef<HTMLDivElement>(null)
    const listRef = useRef<HTMLUListElement>(null)
    const itemRefs = useRef<(HTMLLIElement | null)[]>([])
    const isSelectingRef = useRef(false)
    const isExternalUpdateRef = useRef(false)
    const selectedHeroIdRef = useRef<number | undefined>(undefined)
    const selectedHeroNameRef = useRef<string | undefined>(undefined)
    const hasSetInitialValueRef = useRef(false)

    // ---------------------------------------------------------------------------
    // Data (no API calls - all from heroes.json)
    // ---------------------------------------------------------------------------
    const debouncedQuery = useDebounce(inputValue, DEBOUNCE_MS)
    const allHeroes = useMemo(() => { ... }, [])
    const searchResults = useMemo(() => { ... }, [debouncedQuery, allHeroes])

    // ---------------------------------------------------------------------------
    // Computed values
    // ---------------------------------------------------------------------------
    const items = useMemo(() => { ... }, [debouncedQuery, searchResults, allHeroes])
    const listboxId = `${id}-listbox`
    const hasQuery = debouncedQuery.trim().length > 0
    const showClearButton = inputValue
    const showEmptyState = items.length === 0

    // ---------------------------------------------------------------------------
    // Helper functions
    // ---------------------------------------------------------------------------
    const clearSelectedHero = useCallback(() => { ... }, [onSelect])
    const shouldClearHeroOnClose = useCallback(() => { ... }, [inputValue])
    const closeDropdown = useCallback(() => { ... }, [shouldClearHeroOnClose, clearSelectedHero])

    // ---------------------------------------------------------------------------
    // Effects
    // ---------------------------------------------------------------------------
    // Set initial value
    useLayoutEffect(() => { ... }, [initialValue, allHeroes])

    // Close dropdown when clicking outside
    useEffect(() => { ... }, [closeDropdown])

    // Reset highlighted index when items change
    useLayoutEffect(() => { ... }, [items])

    // Scroll highlighted item into view
    useEffect(() => { ... }, [highlightedIndex])

    // Sync div content with inputValue
    useLayoutEffect(() => { ... }, [inputValue])

    // Ensure text content is set when ref becomes available
    useLayoutEffect(() => { ... }, [inputValue])

    // ---------------------------------------------------------------------------
    // Event handlers
    // ---------------------------------------------------------------------------
    const handleSelect = useCallback((hero: Hero) => { ... }, [onSelect])
    const handleClear = useCallback(() => { ... }, [onSelect])
    const handleInputFocus = useCallback(() => { ... }, [isOpen])
    const handleInputChange = useCallback((e: React.FormEvent<HTMLDivElement>) => { ... }, [])
    const handleInputPaste = useCallback((e: React.ClipboardEvent<HTMLDivElement>) => { ... }, [])
    const handleInputBlur = useCallback(() => { ... }, [inputValue, onSelect])
    const handleItemMouseEnter = useCallback((index: number) => { ... }, [])
    const handleItemClick = useCallback((e: React.MouseEvent, hero: Hero) => { ... }, [handleSelect])

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
    const renderHeroItem = (hero: Hero, index: number) => { ... }
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
| `allHeroes` | `Hero[]` | All heroes from heroes.json |
| `searchResults` | `Hero[]` | Filtered heroes based on search |
| `items` | `Hero[]` | Either popular heroes or search results |

### Items Computation Logic

```typescript
const items = useMemo(() => {
    if (debouncedQuery.trim().length === 0) {
        // Show popular heroes when no search query
        return getPopularData().popular_heroes.map(h => ({
            id: h.id,
            name: h.name,
            displayName: allHeroes.find(hero => hero.id === h.id)?.displayName || h.name
        }))
    }
    // Show search results when typing
    return searchResults
}, [debouncedQuery, searchResults, allHeroes])
```

---

## 5. File Changes Summary

| File | Changes |
|------|---------|
| `ui/src/assets/static_data/popular.json` | Add `popular_heroes` array with 7 popular heroes |
| `ui/src/lib/dotautils.ts` | Add `PopularHero` interface and update `PopularData` interface |
| `ui/src/components/HeroSelector.tsx` | Complete rewrite following LeagueSelector/TeamSelector pattern |

---

## 6. Testing Checklist

### Functional Testing

- [ ] **Popular Heroes Display**
  - [ ] When dropdown opens with empty input, shows 7 popular heroes
  - [ ] Popular heroes are displayed in correct order
  - [ ] Each popular hero shows correct displayName

- [ ] **Search Functionality**
  - [ ] Typing in search input filters heroes
  - [ ] Search is debounced (250ms)
  - [ ] Search works on displayName field
  - [ ] Search works on name field
  - [ ] Partial matches work (e.g., "inv" matches "Invoker")
  - [ ] Case-insensitive search works
  - [ ] Empty search results show "No heroes found" message

- [ ] **Hero Selection**
  - [ ] Clicking a hero selects it
  - [ ] Selected hero displayName appears in input
  - [ ] `onSelect` callback is called with correct hero_id
  - [ ] Dropdown closes after selection
  - [ ] Input is focused after selection

- [ ] **Clear Functionality**
  - [ ] Clear button appears when hero is selected
  - [ ] Clicking clear button removes selection
  - [ ] `onSelect` callback is called with `undefined`
  - [ ] Input is cleared
  - [ ] Popular heroes are shown again

- [ ] **Keyboard Navigation**
  - [ ] ArrowDown moves highlight down
  - [ ] ArrowUp moves highlight up
  - [ ] Enter selects highlighted hero
  - [ ] Escape closes dropdown
  - [ ] Tab closes dropdown
  - [ ] Navigation keys open dropdown when closed

- [ ] **Click Outside**
  - [ ] Clicking outside dropdown closes it
  - [ ] State is properly cleaned up on close

- [ ] **Initial Value Handling**
  - [ ] When `initialValue` is provided, hero displayName is shown
  - [ ] When `initialValue` changes, input updates
  - [ ] When `initialValue` is `undefined`, input is cleared

- [ ] **Accessibility**
  - [ ] Component has proper ARIA labels
  - [ ] Keyboard navigation works
  - [ ] Screen reader announces changes

### Visual Testing

- [ ] Dropdown styling matches LeagueSelector/TeamSelector
- [ ] Highlighted item has correct background color
- [ ] Icons (Clear, Chevron) display correctly
- [ ] Placeholder text shows when input is empty
- [ ] Responsive design works on mobile

### Performance Testing

- [ ] Search is fast (<10ms for 137 heroes)
- [ ] No unnecessary re-renders
- [ ] Debounce works correctly

### Integration Testing

- [ ] Works correctly in MatchFilters
- [ ] Works correctly in SeriesFilters
- [ ] URL params update correctly on selection
- [ ] Page refresh preserves selected hero

---

## 7. Implementation Notes

### Key Differences from Reference Implementation

1. **No API calls** - All data is local from `heroes.json`
2. **No loading states** - Data is immediately available
3. **No spinner** - No async operations
4. **Simpler initial value handling** - No need to fetch hero name from API

### Dependencies

- `useDebounce` hook (already exists)
- `getPopularData` function (needs update)
- `heroes.json` (already exists)

### Browser Compatibility

- ContentEditable is widely supported
- All modern browsers support the required features

---

## 8. Next Steps

After this plan is approved:

1. Implement Step 1: Add popular_heroes to popular.json
2. Implement Step 2: Rewrite HeroSelector.tsx
3. Implement Step 3: Add search algorithm
4. Implement Step 4: Add keyboard navigation
5. Implement Step 5: Handle initial values
6. Implement Step 6: Test and verify

Each step can be implemented and tested independently.
