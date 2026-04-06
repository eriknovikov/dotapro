# Section 4: Active Filter Chips

## Objective

Add active filter chips to the Series and Matches filter sidebar, allowing users to quickly see and remove individual filters.

---

## Critical Review Findings

**Confirmed via code inspection:**
- [`ui/src/components/series/SeriesFilters.tsx`](ui/src/components/series/SeriesFilters.tsx:1) - Complex component with mobile drawer, uses `CustomSelect` dropdowns
- [`ui/src/components/matches/MatchFilters.tsx`](ui/src/components/matches/MatchFilters.tsx:1) - Simpler inline filter layout using `RadixSelect`
- [`ui/src/components/ui/badge.tsx`](ui/src/components/ui/badge.tsx:1) - Badge uses CVA with variants: `default`, `secondary`, `destructive`, `outline`

**Key observations:**
- SeriesFilters is a complex mobile drawer component (lines 107-265)
- MatchFilters is simpler inline grid layout (lines 50-93)
- SeriesFilters already has an `itemType` prop that can be "series" or "matches" (line 15)
- Both components use URL search params for state via `useNavigate`
- SeriesFilters already has `handleClear` function (lines 76-93) that handles conditional clearing based on itemType

**Issues identified in original plan:**

1. **Integration complexity:** SeriesFilters is complex with mobile drawer, sidebar, and multiple sections. Simply adding ActiveFiltersBar to the top might break the mobile drawer animation or layout.
   - **Resolution:** Add ActiveFiltersBar INSIDE the sidebar/drawer content area, not outside. Position it at the top of the filter content.

2. **Lookup data for filter names:** To show human-readable names (e.g., "League: DPC 2024"), we need lookup data. The original plan suggested heroes.json but didn't address leagues/teams.
   - **Resolution:** Use existing selectors (LeagueSelector, TeamSelector) that already handle name display. We'll need to expose selected option labels or use a lookup map.

3. **MatchFilters is separate:** Unlike SeriesFilters, MatchFilters is a standalone component without the drawer pattern. The ActiveFiltersBar should work for both but may need different integration.

4. **Badge variant naming conflict:** Original plan suggested `variant="filter"` but this could conflict with existing patterns. Use `variant="active"` instead.

5. **Handler conflict:** Plan 4 suggested adding `handleClearAll` but SeriesFilters already has `handleClear` that works correctly. We should reuse the existing handler, not create a new one.

**Resolved decisions:**
- Place ActiveFiltersBar inside the filter content area, below the mobile header
- Use `variant="active"` for badge (distinct from other badge uses)
- For name lookups: pass pre-loaded name maps or derive from selector state
- Keep implementation simple - focus on SeriesFilters first, then MatchFilters
- Reuse existing `handleClear` function instead of creating `handleClearAll`

---

## Current State

### SeriesFilters (lines 107-265):
```tsx
<aside className="... fixed z-50 ...">
    <div className="h-full overflow-y-auto p-6">
        {/* Mobile close button */}
        <div className="mb-6 flex items-center justify-between lg:hidden">
            ...
        </div>
        
        {/* Desktop title */}
        <div className="mb-6 hidden items-center gap-2 lg:flex">
            ...
        </div>
        
        <div className="space-y-12">
            {/* League Filter */}
            {/* Team Filter */}
            {/* etc */}
        </div>
    </div>
</aside>
```

### MatchFilters (lines 50-93):
```tsx
<div className="bg-background-accent mb-6 rounded-lg p-4">
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        {/* Inline filter controls */}
    </div>
</div>
```

---

## Design Goal

Show active filters as dismissible chips above the filter controls:

```
┌─────────────────────────────────────────┐
│  Active Filters                         │
│  ┌──────────┐ ┌──────────┐ ┌────────┐  │
│  │ League X │ │ Team Y   │ │ Hero Z │  │
│  │    ×     │ │    ×     │ │   ×    │  │
│  └──────────┘ └──────────┘ └────────┘  │
│                                         │
│  League: [Dropdown]                      │
│  Team: [Dropdown]                       │
│  ...                                    │
└─────────────────────────────────────────┘
```

---

## 1. Create ActiveFiltersBar Component

Create `ui/src/components/ActiveFiltersBar.tsx`:

```tsx
import { X } from "lucide-react"
import { Badge } from "./ui/badge"
import { Button } from "./ui/button"

interface ActiveFiltersBarProps {
    filters: {
        league?: number
        team?: number
        player?: number
        hero?: number
        sort?: string
    }
    onRemove: (key: "league" | "team" | "hero" | "player") => void
    onClearAll: () => void
    labels?: {
        league?: string
        team?: string
        player?: string
        hero?: string
    }
}

export function ActiveFiltersBar({
    filters,
    onRemove,
    onClearAll,
    labels = {},
}: ActiveFiltersBarProps) {
    const activeFilters = Object.entries(filters).filter(
        ([, value]) => value !== undefined && value !== ""
    )

    if (activeFilters.length === 0) return null

    const getLabel = (key: string, value: string | number) => {
        const label = labels[key as keyof typeof labels]
        if (label) return label
        // Fallback display
        switch (key) {
            case "league":
                return `League #${value}`
            case "team":
                return `Team #${value}`
            case "hero":
                return `Hero #${value}`
            case "player":
                return `Player #${value}`
            case "sort":
                return value === "newest" ? "Newest first" : "Oldest first"
            default:
                return String(value)
        }
    }

    return (
        <div className="mb-4 space-y-2">
            <div className="flex items-center justify-between">
                <span className="text-xs font-medium uppercase tracking-wide text-foreground-muted">
                    Active Filters
                </span>
                <Button
                    variant="ghost"
                    size="sm"
                    onClick={onClearAll}
                    className="h-auto p-0 text-xs text-foreground-muted hover:text-foreground"
                >
                    Clear all
                </Button>
            </div>

            <div className="flex flex-wrap gap-2">
                {activeFilters.map(([key, value]) => (
                    <Badge
                        key={key}
                        variant="active"
                        className="pl-2 pr-1"
                    >
                        <span className="mr-1 capitalize">{key}:</span>
                        <span className="font-normal">
                            {getLabel(key, value as string | number)}
                        </span>
                        <button
                            onClick={() => onRemove(key as "league" | "team" | "hero" | "player")}
                            className="ml-1 rounded-full p-0.5 transition-colors hover:bg-primary-500/20"
                            aria-label={`Remove ${key} filter`}
                        >
                            <X className="h-3 w-3" />
                        </button>
                    </Badge>
                ))}
            </div>
        </div>
    )
}
```

---

## 2. Add Active Badge Variant

Update [`ui/src/components/ui/badge.tsx`](ui/src/components/ui/badge.tsx:10-16):

### Current badge variants (lines 10-16):
```tsx
variants: {
    variant: {
        default: "border-transparent bg-primary text-primary-foreground shadow hover:bg-primary/80",
        secondary: "border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80",
        destructive: "border-transparent bg-destructive text-destructive-foreground shadow hover:bg-destructive/80",
        outline: "text-foreground",
    },
},
```

### New - add active variant:
```tsx
variants: {
    variant: {
        default: "border-transparent bg-primary text-primary-foreground shadow hover:bg-primary/80",
        secondary: "border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80",
        destructive: "border-transparent bg-destructive text-destructive-foreground shadow hover:bg-destructive/80",
        outline: "text-foreground",
        // NEW: Active filter chip variant
        active: "border-primary-500/30 bg-primary-500/10 text-foreground hover:bg-primary-500/20",
    },
},
```

---

## 3. Update SeriesFilters Component

Modify [`ui/src/components/series/SeriesFilters.tsx`](ui/src/components/series/SeriesFilters.tsx):

**Key integration points:**
1. Add `ActiveFiltersBar` inside the scrollable content area
2. Pass `filters` object and navigation handlers
3. Extract selected labels from selectors for display

### Add import (around line 2):
```tsx
import { ActiveFiltersBar } from "@/components/ActiveFiltersBar"
```

**Note:** You'll need to export ActiveFiltersBar from components/index.ts after creating it.

### Add handlers for individual filter removal (after line 93, before return):

The existing `handleClear` function (lines 76-93) already handles clearing ALL filters properly based on `itemType`. We just need to add individual filter removal:

```tsx
// Handler for removing individual filter
const handleRemoveFilter = (key: "league" | "team" | "hero" | "player") => {
    const newFilters = { ...filters }
    delete newFilters[key]
    navigate({
        to: ".",
        search: newFilters,
    })
}
```

### Add ActiveFiltersBar inside the aside content (around line 143, before `<div className="space-y-12">`):

Find the section around line 143:
```tsx
{/* Desktop title */}
<div className="mb-6 hidden items-center gap-2 lg:flex">
    ...
</div>

<div className="space-y-12">
```

Add ActiveFiltersBar after the desktop title but before space-y-12:
```tsx
{/* Desktop title */}
<div className="mb-6 hidden items-center gap-2 lg:flex">
    ...
</div>

{/* Active Filters Bar */}
<ActiveFiltersBar
    filters={filters}
    onRemove={handleRemoveFilter}
    onClearAll={handleClear}
/>

<div className="space-y-12">
```

---

## 4. Update MatchFilters Component

Modify [`ui/src/components/matches/MatchFilters.tsx`](ui/src/components/matches/MatchFilters.tsx):

MatchFilters is simpler with inline layout. We can add ActiveFiltersBar above the grid.

### Add import (around line 15):
```tsx
import { ActiveFiltersBar } from "@/components/ActiveFiltersBar"
```

### Add handlers (around line 48, before return):
```tsx
// Handler for removing individual filter
const handleRemoveFilter = (key: "league" | "team" | "hero" | "player") => {
    updateFilters({ [key]: undefined })
}

// Handler for clearing all filters
const handleClearAll = () => {
    updateFilters({
        league: undefined,
        team: undefined,
        hero: undefined,
        player: undefined,
        sort: undefined,
        limit: undefined,
    })
}
```

### Add ActiveFiltersBar inside the component (around line 51, before the grid):
```tsx
return (
    <div className="bg-background-accent mb-6 rounded-lg p-4">
        {/* Active Filters Bar */}
        <ActiveFiltersBar
            filters={search}
            onRemove={handleRemoveFilter}
            onClearAll={handleClearAll}
        />
        
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
            ...
        </div>
    </div>
)
```

---

## 5. Files to Create

| File | Purpose |
|------|---------|
| `ui/src/components/ActiveFiltersBar.tsx` | Reusable active filter chips component |

---

## 6. Files to Modify

| File | Changes |
|------|---------|
| `ui/src/components/ActiveFiltersBar.tsx` | NEW - Create this component |
| `ui/src/components/ui/badge.tsx` | Add `active` variant |
| `ui/src/components/series/SeriesFilters.tsx` | Add ActiveFiltersBar, handleRemoveFilter |
| `ui/src/components/matches/MatchFilters.tsx` | Add ActiveFiltersBar, handlers |
| `ui/src/components/index.ts` | Export ActiveFiltersBar |

---

## Testing Checklist

- [ ] ActiveFiltersBar appears when any filter is set
- [ ] Filter chips show key and value (e.g., "league: DPC 2024" or fallback "league: 12345")
- [ ] Clicking × removes that specific filter
- [ ] "Clear all" removes all filters (uses existing handleClear for SeriesFilters)
- [ ] Filters work correctly with SeriesFilters (mobile drawer)
- [ ] Filters work correctly with MatchFilters (inline)
- [ ] URL updates correctly when removing/clearing filters
- [ ] No layout breaking in mobile drawer
- [ ] Empty state (no active filters) shows nothing

---

## Commit Message

```
feat(filters): add active filter chips display

- Create ActiveFiltersBar component with dismissible chips
- Add 'active' badge variant for filter chips
- Integrate with SeriesFilters (drawer) and MatchFilters (inline)
- Allow individual filter removal and clear all
```

---

## Dependencies

None - uses existing Badge and Button components.
