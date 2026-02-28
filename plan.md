# Responsive Design Plan - `/series` Page

## Overview
Make the `/series` page responsive for mobile and tablet devices with a bottom sheet filter modal and adaptive grid layout.

## Scope
- **Target Page**: `/series` (series list page only)
- **Devices**: Mobile (<640px) and Tablet (640px-1024px)
- **Out of Scope**: `/series/$id` (series detail page) - to be addressed later

---

## Changes Required

### 1. FiltersSidebar - Bottom Sheet Modal

**Current State**: Slide-in from left sidebar with fixed width (288px)

**New Behavior**:
- Mobile & Tablet: Bottom sheet slides up from bottom
- Desktop: Keep current fixed sidebar behavior
- Height: ~80% of viewport
- Backdrop overlay: Click outside to close
- Trigger: FAB button (bottom-right) on mobile/tablet

**Files to Modify**:
- `ui/src/components/FiltersSidebar.tsx`
- `ui/src/routes/series/index.tsx`

**Implementation Details**:

#### FiltersSidebar.tsx Changes
```tsx
// Current: Fixed sidebar from left
className="fixed top-16 left-0 z-40 w-72 h-[calc(100vh-4rem)]"

// New: Bottom sheet on mobile/tablet, sidebar on desktop
className={`
    fixed z-50
    ${isMobileOpen ? "translate-y-0" : "translate-y-full"}
    md:fixed md:top-16 md:left-0 md:translate-x-0 md:translate-y-0
    md:w-72 md:h-[calc(100vh-4rem)]
    bottom-0 left-0 right-0
    h-[80vh]
    bg-background border-t border-border
    rounded-t-2xl md:rounded-none
    shadow-2xl
    transition-transform duration-300 ease-in-out
`}
```

#### Backdrop Overlay
```tsx
// Add backdrop overlay for mobile/tablet
{isMobileOpen && (
    <div
        className="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 md:hidden"
        onClick={onMobileClose}
        aria-hidden="true"
    />
)}
```

#### Handle/Close Button
```tsx
// Add drag handle at top of bottom sheet
<div className="flex items-center justify-center pt-3 pb-2 md:hidden">
    <div className="w-12 h-1.5 bg-foreground-muted/30 rounded-full" />
</div>
```

---

### 2. SeriesList Grid - Full Width When Filters Closed

**Current State**: Grid has `ml-72` margin to account for fixed sidebar

**New Behavior**:
- Mobile & Tablet: Full width (no margin)
- Desktop: Keep `ml-72` margin for sidebar

**Files to Modify**:
- `ui/src/routes/series/index.tsx`

**Implementation Details**:

```tsx
// Current: Fixed margin for sidebar
<main className="flex-1 py-6 ml-72 overflow-y-auto">

// New: Responsive margin
<main className="flex-1 py-6 md:ml-72 overflow-y-auto">
```

---

### 3. SeriesList Grid - Adaptive Grid with 400px Min Width

**Current State**: `grid-cols-[repeat(auto-fit,minmax(450px,1fr))]`

**New Behavior**:
- Keep `minmax(400px,1fr)` for now
- Later: Reduce to `minmax(300px,1fr)` with SeriesCard improvements

**Files to Modify**:
- `ui/src/components/SeriesList.tsx`

**Implementation Details**:

```tsx
// Current
<div className="grid grid-cols-[repeat(auto-fit,minmax(450px,1fr))] gap-6">

// New: 400px min width
<div className="grid grid-cols-[repeat(auto-fit,minmax(400px,1fr))] gap-4 sm:gap-6">
```

---

### 4. SeriesCard - Responsive Layout Improvements

**Current State**: Fixed widths, horizontal layout that may overflow on small screens

**New Behavior**:
- Add responsive classes to wrap team names
- Adjust button width on mobile
- Ensure logos don't overflow

**Files to Modify**:
- `ui/src/components/SeriesCard.tsx`

**Implementation Details**:

```tsx
// Teams and Score section
<div className="flex items-center justify-between gap-2 sm:gap-4 border-b border-border mb-3 pb-3">
    {/* Team A */}
    <div className="flex items-center min-w-0 flex-1">
        <div className="flex min-w-0 gap-1 sm:gap-2 items-center">
            {series.team_a.logo_url && (
                <img
                    src={series.team_a.logo_url}
                    className="h-6 sm:h-7 w-auto max-w-8 sm:max-w-11 rounded shrink-0 select-none"
                />
            )}
            <span className="font-bold text-foreground text-sm sm:text-md font-shantell truncate">
                {series.team_a.name}
            </span>
        </div>
    </div>

    {/* Score */}
    <div className="flex items-center justify-center shrink-0 select-none px-1">
        <span className="font-bold text-xs text-foreground-muted">{series.team_a_score}</span>
        <span className="mx-1 sm:mx-2 text-gray-400">-</span>
        <span className="font-bold text-xs text-foreground-muted">{series.team_b_score}</span>
    </div>

    {/* Team B */}
    <div className="flex items-center justify-end min-w-0 flex-1">
        <div className="flex min-w-0 gap-1 sm:gap-2 items-center justify-end">
            <span className="font-bold text-foreground text-sm sm:text-md text-end font-shantell truncate">
                {series.team_b.name}
            </span>
            {series.team_b.logo_url && (
                <img
                    src={series.team_b.logo_url}
                    className="h-6 sm:h-7 w-auto max-w-8 sm:max-w-11 rounded shrink-0 select-none"
                />
            )}
        </div>
    </div>
</div>

// Button section
<div className="flex justify-end items-center mt-auto">
    <Button
        variant="cool-outline"
        size="sm"
        className="w-24 sm:w-32 group/btn text-xs"
        onClick={() => navigate({ to: `/series/${series.series_id}` })}
    >
        View series
        <div className="flex justify-center items-center w-5 sm:w-6 h-4">
            <Eye className="size-3 transition-all duration-300 group-hover/btn:size-5" />
        </div>
    </Button>
</div>
```

---

### 5. FAB Button - Filter Trigger

**Current State**: Fixed bottom-right button on mobile only

**New Behavior**:
- Show on mobile AND tablet
- Hide on desktop (sidebar is always visible)

**Files to Modify**:
- `ui/src/routes/series/index.tsx`

**Implementation Details**:

```tsx
// Current: md:hidden
<div className="md:hidden fixed bottom-4 right-4 z-50">

// New: Show on mobile and tablet, hide on desktop
<div className="lg:hidden fixed bottom-4 right-4 z-50">
```

---

## File Changes Summary

| File | Changes |
|------|---------|
| `ui/src/components/FiltersSidebar.tsx` | Convert to bottom sheet on mobile/tablet, add backdrop, handle, responsive classes |
| `ui/src/routes/series/index.tsx` | Update main content margin, FAB visibility |
| `ui/src/components/SeriesList.tsx` | Update grid min-width to 400px |
| `ui/src/components/SeriesCard.tsx` | Add responsive classes for team names, logos, button |

---

## Implementation Order

1. **FiltersSidebar.tsx** - Convert to bottom sheet
2. **series/index.tsx** - Update layout and FAB
3. **SeriesList.tsx** - Update grid min-width
4. **SeriesCard.tsx** - Add responsive classes

---

## Testing Checklist

- [ ] Mobile (<640px): Bottom sheet opens/closes correctly
- [ ] Mobile (<640px): Grid takes full width
- [ ] Mobile (<640px): Cards display correctly with 400px min-width
- [ ] Mobile (<640px): Team names truncate properly
- [ ] Tablet (640px-1024px): Bottom sheet opens/closes correctly
- [ ] Tablet (640px-1024px): Grid takes full width
- [ ] Tablet (640px-1024px): Cards display correctly
- [ ] Desktop (>1024px): Sidebar remains fixed
- [ ] Desktop (>1024px): Grid has proper margin
- [ ] Desktop (>1024px): FAB button is hidden
- [ ] Backdrop click closes bottom sheet
- [ ] Escape key closes bottom sheet
- [ ] Smooth transitions on all breakpoints

---

## Future Enhancements (Out of Scope)

1. Reduce SeriesCard min-width to 300px with more aggressive wrapping
2. Add swipe gesture to close bottom sheet
3. Add drag handle for bottom sheet
4. Implement `/series/$id` page responsiveness
5. Add skeleton loading states for mobile
