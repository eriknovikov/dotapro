# Plan: Make `/series/$id` Endpoint Responsive

## Overview
Transform the series detail page to be mobile-friendly with no horizontal scrolling, following the same UX philosophy used in the `/series` endpoint redesign.

## Design Decisions

### 1. SeriesHeader Component
**Reference:** [`ui/src/components/SeriesCard.tsx`](../ui/src/components/SeriesCard.tsx) (lines 42-94)

**Mobile Layout (< md breakpoint):**
```
┌─────────────────────────┐
│ ← Back to Series        │
├─────────────────────────┤
│ [Logo] Team A (2)       │
│                         │
│      ⚔️                 │
│                         │
│ (1) Team B [Logo]      │
├─────────────────────────┤
│ League Name             │
│ ~2 hours ago            │
└─────────────────────────┘
```

**Desktop Layout (≥ md breakpoint):**
- Keep current layout unchanged
- Teams side-by-side with centered score
- League, tier, and time on one line

**Implementation Details:**
- Use `hidden md:flex` for desktop layout
- Use `flex md:hidden` for mobile layout
- Import `Swords` icon from `lucide-react`
- Mobile team row: Logo (left) → Team Name (centered) → Score (in parentheses, right)
- VS icon centered between teams on mobile
- League info: Only show league name and formatted start time on mobile (no tier badge)
- Team logos: `h-8` on mobile (vs `h-16 md:h-24` on desktop)
- Team names: `text-base` on mobile (vs `text-xl md:text-3xl` on desktop)

---

### 2. GameTabContent Component
**File:** [`ui/src/components/series/GameTabContent.tsx`](../ui/src/components/series/GameTabContent.tsx)

**Mobile Layout (< sm breakpoint):**
```
┌─────────────────────────┐
│ Match ID                │
│ 🕐 45:32                │
│                         │
│              [Copy ID]  │
│           [OpenDota →]  │
│              [Stratz →] │
└─────────────────────────┘
```

**Implementation Details:**
- Make the card container `relative` on mobile
- Position buttons in top-right corner using absolute positioning
- Stack buttons vertically with minimal spacing
- Button container: `absolute top-4 right-4 flex flex-col gap-2`
- Match info section: `pr-24` (padding-right to avoid overlap with buttons)
- On desktop (≥ sm): Keep current `flex-row` layout
- Buttons should be compact on mobile: `size="sm"` or smaller

---

### 3. TabsList (in $id.tsx)
**File:** [`ui/src/routes/series/$id.tsx`](../ui/src/routes/series/$id.tsx)

**Mobile Layout:**
```
┌─────────────────────────┐
│ [G1] [G2] [G3] [G4] >  │
└─────────────────────────┘
```

**Implementation Details:**
- Change tab labels conditionally based on breakpoint
- Mobile: `G{i + 1}` (e.g., "G1", "G2")
- Desktop: `Game {i + 1}` (e.g., "Game 1", "Game 2")
- Keep existing `max-w-full overflow-x-auto overflow-y-hidden`
- Consider smaller padding on mobile tabs

---

### 4. PlayersTable Component
**Status:** DEFERRED - Will be analyzed and implemented separately in a future task.

---

## Files to Modify

| File | Changes | Priority |
|------|---------|----------|
| [`ui/src/components/series/SeriesHeader.tsx`](../ui/src/components/series/SeriesHeader.tsx) | Add mobile layout with VS icon, scores in parentheses, simplified league info | 1 |
| [`ui/src/components/series/GameTabContent.tsx`](../ui/src/components/series/GameTabContent.tsx) | Absolute position buttons in top-right on mobile | 2 |
| [`ui/src/routes/series/$id.tsx`](../ui/src/routes/series/$id.tsx) | Compact tab labels (G1, G2, etc.) | 3 |

---

## Implementation Order

1. **SeriesHeader** - Establishes mobile-first pattern, most visible change
2. **GameTabContent** - Improves action accessibility
3. **TabsList** - Quick text change

---

## Breakpoints

Following Tailwind conventions:
- `sm`: 640px (small phones)
- `md`: 768px (tablets, large phones)
- `lg`: 1024px (laptops)

---

## Key Patterns to Follow

### Hidden/Block Pattern (from SeriesCard)
```tsx
{/* Mobile only */}
<div className="flex md:hidden">
  {/* Mobile content */}
</div>

{/* Desktop only */}
<div className="hidden md:flex">
  {/* Desktop content */}
</div>
```

### Responsive Text Sizes
```tsx
className="text-base md:text-xl lg:text-3xl"
```

### Responsive Spacing
```tsx
className="gap-2 md:gap-4 lg:gap-6"
```

---

## Potential Issues to Watch For

1. **Long team names** - Use `truncate` class to prevent overflow
2. **Missing logos** - Handle gracefully with conditional rendering
3. **Button overlap** - Ensure proper padding on match info section
4. **Tab overflow** - Test with many games (5+)
5. **Touch targets** - Ensure buttons are at least 44px tall on mobile

---

## Testing Checklist

- [ ] SeriesHeader displays correctly on mobile (320px+)
- [ ] SeriesHeader displays correctly on tablet (768px+)
- [ ] SeriesHeader displays correctly on desktop (1024px+)
- [ ] Back button works on all screen sizes
- [ ] VS icon centers properly between teams on mobile
- [ ] Scores in parentheses don't cause overflow
- [ ] League info shows only name + time on mobile
- [ ] GameTabContent buttons don't overlap match info
- [ ] Buttons are accessible and tappable on mobile
- [ ] Tab labels change from "Game N" to "GN" on mobile
- [ ] Tabs scroll horizontally when many games exist
- [ ] No horizontal scroll on the main page

---

## Dependencies

- `lucide-react` - For `Swords` icon (already imported in SeriesCard)
- Tailwind CSS v4 - For responsive utilities
- Existing components: `Button`, `Tabs`, `TabsList`, `TabsTrigger`, `TabsContent`

---

## Notes

- The `formatRelativeTime` function is already imported from `@/lib`
- The `Swords` icon is already used in `SeriesCard.tsx`
- Follow the existing code style and patterns from `SeriesCard.tsx`
- Use `cn()` utility for conditional class merging
