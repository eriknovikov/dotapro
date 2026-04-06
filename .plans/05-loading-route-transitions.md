# Section 5: Loading States & Route Transitions

## Objective

Improve loading states with animated skeleton shimmer effects and add smooth route transitions between pages.

---

## Critical Review Findings

**Confirmed via code inspection:**
- [`ui/src/components/Skeleton.tsx:7-19`](ui/src/components/Skeleton.tsx:7-19) - Current skeleton uses `animate-pulse`
- [`ui/src/components/Skeleton.tsx:23-66`](ui/src/components/Skeleton.tsx:23-66) - `SeriesCardSkeleton` exists but uses `animate-pulse`
- [`ui/src/components/matches/MatchList.tsx:78-93`](ui/src/components/matches/MatchList.tsx:78-93) - `MatchListSkeleton` exists but is basic (not proper `MatchCardSkeleton`)
- [`ui/src/routes/__root.tsx`](ui/src/routes/__root.tsx:1) - Root layout with existing structure

**Key observations:**
- TanStack Router v1 has built-in transition support via CSS class approach
- **Shimmer keyframes are now defined in Plan 2** - do NOT redefine them
- `MatchCardSkeleton` doesn't exist as a standalone component - need to create it
- `SeriesCardSkeleton` and `Skeleton` are exported from `components/index.ts:15`

**Issues identified in original plan:**

1. **CRITICAL FIX - Duplicate shimmer keyframes:** The original plan defined shimmer keyframes. Plan 2 now defines them. Plan 5 should reference them, NOT redefine.
   - **Resolution:** Remove shimmer keyframe definition from this plan. Reference `.skeleton-shimmer` class from Plan 2.

2. **TanStack Router transition approach:** The original plan suggested creating a `PageTransition` component and restructuring `__root.tsx`. This was too aggressive and would break the existing good structure.
   - **Resolution:** Use simpler approach - add `page-enter` class to existing content div without restructuring.

3. **MatchCardSkeleton:** The original plan showed adding a proper MatchCardSkeleton but didn't note it needs to be exported.
   - **Resolution:** Create proper `MatchCardSkeleton` component in Skeleton.tsx AND export it from components/index.ts.

4. **Existing skeletons use animate-pulse:** Both SeriesCardSkeleton and the inline MatchListSkeleton use `animate-pulse`. Need to update to shimmer.

5. **Missing prefers-reduced-motion for shimmer:** Plan 2 added some prefers-reduced-motion but not for skeleton shimmer.

**Resolved decisions:**
- Reference shimmer keyframes from Plan 2 (do not redefine)
- Use minimal changes to __root.tsx (add class only, don't restructure)
- Create proper MatchCardSkeleton in Skeleton.tsx
- Export MatchCardSkeleton from components/index.ts
- Add prefers-reduced-motion for .skeleton-shimmer

---

## Files to Modify

- [`ui/src/components/Skeleton.tsx`](ui/src/components/Skeleton.tsx:1) - Update skeletons to use shimmer, add MatchCardSkeleton
- [`ui/src/components/series/SeriesList.tsx`](ui/src/components/series/SeriesList.tsx:1) - Already uses SeriesCardSkeleton (no change needed if skeleton updated)
- [`ui/src/components/matches/MatchList.tsx`](ui/src/components/matches/MatchList.tsx:1) - Use new MatchCardSkeleton
- [`ui/src/components/index.ts`](ui/src/components/index.ts:15) - Export MatchCardSkeleton
- [`ui/src/routes/__root.tsx`](ui/src/routes/__root.tsx:1) - Add page-enter class (minimal change)
- [`ui/src/index.css`](ui/src/index.css:1) - Add page-enter keyframes and loading bar (NO shimmer - use from Plan 2)

---

## 1. CSS for Route Transitions

Add to [`ui/src/index.css`](ui/src/index.css:1) (append before closing):

```css
/* ============================================
   SECTION 5: Loading States & Transitions
   ============================================ */

/* NOTE: @keyframes shimmer and .skeleton-shimmer are defined in Plan 2 */
/* Do NOT redefine them here - use the existing class */

/* Route transition animations */
@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.page-enter {
    animation: fadeIn 0.3s ease-out;
}

/* Loading bar for route changes */
@keyframes loadingBar {
    0% {
        width: 0%;
        opacity: 1;
    }
    50% {
        width: 70%;
        opacity: 1;
    }
    100% {
        width: 100%;
        opacity: 0;
    }
}

.loading-bar {
    position: fixed;
    top: 0;
    left: 0;
    height: 2px;
    background: linear-gradient(90deg, hsl(0, 84%, 50%), hsl(38, 92%, 50%));
    animation: loadingBar 0.6s ease-out forwards;
    z-index: 9999;
}

/* Respect prefers-reduced-motion for skeletons (supplement to Plan 2) */
@media (prefers-reduced-motion: reduce) {
    .skeleton-shimmer {
        animation: none;
        background: hsl(0, 0%, 10%);
    }
}
```

---

## 2. Update Skeleton Component

Modify [`ui/src/components/Skeleton.tsx`](ui/src/components/Skeleton.tsx:7-19):

### Current Skeleton function (lines 7-19):
```tsx
export function Skeleton({ className, variant = "default", ...props }: SkeletonProps) {
    return (
        <div
            className={cn(
                "bg-background-accent animate-pulse rounded-md",
                variant === "card" && "bg-background-card/80 rounded-xl",
                variant === "text" && "h-4 w-3/4",
                variant === "button" && "h-8 w-32 rounded-md",
                className,
            )}
            {...props}
        />
    )
}
```

### New Skeleton function:
```tsx
export function Skeleton({ className, variant = "default", ...props }: SkeletonProps) {
    return (
        <div
            className={cn(
                "skeleton-shimmer rounded-md",
                variant === "card" && "rounded-xl",
                variant === "text" && "h-4 w-3/4",
                variant === "button" && "h-8 w-32 rounded-md",
                className,
            )}
            {...props}
        />
    )
}
```

**Changes:**
- Removed `bg-background-accent` (shimmer gradient provides color)
- Removed `animate-pulse`
- Added `skeleton-shimmer` class (defined in Plan 2)

---

## 3. Update SeriesCardSkeleton

Modify [`ui/src/components/Skeleton.tsx`](ui/src/components/Skeleton.tsx:23-66):

### Current (lines 23-66):
```tsx
export function SeriesCardSkeleton() {
    return (
        <div className="group bg-background-card/80 text-card-foreground relative overflow-hidden rounded-xl shadow-xl backdrop-blur-sm">
            <div className="relative z-10 flex h-full flex-col p-4">
                {/* Teams and Score */}
                <div className="border-border mb-3 flex items-center justify-between gap-4 border-b pb-3">
                    {/* Team A */}
                    <div className="flex min-w-0 flex-1 items-center">
                        <div className="flex min-w-0 items-center gap-1">
                            <Skeleton className="h-7 w-11 shrink-0 rounded" variant="default" />
                            <Skeleton className="h-5 w-24" variant="text" />
                        </div>
                    </div>

                    {/* Score */}
                    <div className="flex shrink-0 items-center justify-center">
                        <Skeleton className="h-4 w-4" variant="default" />
                        <span className="mx-2 text-gray-400">-</span>
                        <Skeleton className="h-4 w-4" variant="default" />
                    </div>

                    {/* Team B */}
                    <div className="flex min-w-0 flex-1 items-center justify-end">
                        <div className="flex min-w-0 items-center justify-end gap-1">
                            <Skeleton className="h-5 w-24" variant="text" />
                            <Skeleton className="h-7 w-11 shrink-0 rounded" variant="default" />
                        </div>
                    </div>
                </div>

                {/* League and Date */}
                <div className="text-foreground-muted flex flex-1 flex-col gap-2 text-xs">
                    <Skeleton className="h-3 w-32" variant="default" />
                    <Skeleton className="h-3 w-24" variant="default" />
                </div>

                {/* Button */}
                <div className="mt-auto flex items-center justify-end">
                    <Skeleton className="h-8 w-32 rounded-md" variant="button" />
                </div>
            </div>
        </div>
    )
}
```

### New:
```tsx
export function SeriesCardSkeleton() {
    return (
        <div className="group bg-background-card/80 text-card-foreground relative overflow-hidden rounded-xl shadow-xl backdrop-blur-sm">
            <div className="relative z-10 flex h-full flex-col p-4">
                {/* Teams and Score */}
                <div className="border-border mb-3 flex items-center justify-between gap-4 border-b pb-3">
                    {/* Team A */}
                    <div className="flex min-w-0 flex-1 items-center">
                        <div className="flex min-w-0 items-center gap-1">
                            <Skeleton className="h-7 w-11 shrink-0 rounded" variant="default" />
                            <Skeleton className="h-5 w-24" variant="text" />
                        </div>
                    </div>

                    {/* Score */}
                    <div className="flex shrink-0 items-center justify-center">
                        <Skeleton className="h-4 w-4" variant="default" />
                        <span className="text-foreground-muted mx-2">-</span>
                        <Skeleton className="h-4 w-4" variant="default" />
                    </div>

                    {/* Team B */}
                    <div className="flex min-w-0 flex-1 items-center justify-end">
                        <div className="flex min-w-0 items-center justify-end gap-1">
                            <Skeleton className="h-5 w-24" variant="text" />
                            <Skeleton className="h-7 w-11 shrink-0 rounded" variant="default" />
                        </div>
                    </div>
                </div>

                {/* League and Date */}
                <div className="text-foreground-muted flex flex-1 flex-col gap-2 text-xs">
                    <Skeleton className="h-3 w-32" variant="default" />
                    <Skeleton className="h-3 w-24" variant="default" />
                </div>

                {/* Button */}
                <div className="mt-auto flex items-center justify-end">
                    <Skeleton className="h-8 w-32 rounded-md" variant="button" />
                </div>
            </div>
        </div>
    )
}
```

**Changes:**
- Fixed `text-gray-400` to `text-foreground-muted` for consistency
- Skeleton variants now use shimmer (via updated Skeleton component)

---

## 4. Create MatchCardSkeleton

Add to [`ui/src/components/Skeleton.tsx`](ui/src/components/Skeleton.tsx) after `SeriesCardSkeleton` (after line 66):

```tsx
// MatchCard skeleton matching the MatchCard layout
export function MatchCardSkeleton() {
    return (
        <div className="bg-background-card/80 text-card-foreground relative min-w-67.5 overflow-hidden rounded-xl shadow-xl backdrop-blur-sm">
            <div className="relative z-10 flex h-full flex-col p-4">
                {/* Radiant Team */}
                <div className="mb-2 flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        <Skeleton className="h-8 w-10 rounded" variant="default" />
                        <Skeleton className="h-5 w-28" variant="text" />
                    </div>
                </div>

                {/* Radiant Heroes */}
                <div className="mb-2 flex gap-1">
                    {Array.from({ length: 5 }).map((_, i) => (
                        <Skeleton key={`rad-${i}`} className="h-6 w-8" variant="default" />
                    ))}
                </div>

                {/* VS Divider */}
                <div className="my-3 flex justify-center">
                    <Skeleton className="h-6 w-6 rounded-full" variant="default" />
                </div>

                {/* Dire Team */}
                <div className="mb-2 flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        <Skeleton className="h-8 w-10 rounded" variant="default" />
                        <Skeleton className="h-5 w-24" variant="text" />
                    </div>
                </div>

                {/* Dire Heroes */}
                <div className="mb-3 flex gap-1">
                    {Array.from({ length: 5 }).map((_, i) => (
                        <Skeleton key={`dire-${i}`} className="h-6 w-8" variant="default" />
                    ))}
                </div>

                {/* Separator */}
                <div className="border-border mb-3 border-t" />

                {/* League and Date */}
                <div className="text-foreground-muted mb-2 flex flex-col gap-2 text-xs">
                    <Skeleton className="h-3 w-32" variant="default" />
                    <Skeleton className="h-3 w-24" variant="default" />
                </div>

                {/* Button */}
                <div className="mt-auto flex justify-end">
                    <Skeleton className="h-8 w-28 rounded-md" variant="button" />
                </div>
            </div>
        </div>
    )
}
```

---

## 5. Export MatchCardSkeleton

Update [`ui/src/components/index.ts`](ui/src/components/index.ts:15):

### Current (line 15):
```ts
export { SeriesCardSkeleton, Skeleton } from "./Skeleton"
```

### New:
```ts
export { MatchCardSkeleton, SeriesCardSkeleton, Skeleton } from "./Skeleton"
```

---

## 6. Update MatchListSkeleton

Modify [`ui/src/components/matches/MatchList.tsx`](ui/src/components/matches/MatchList.tsx:78-93) to use the new MatchCardSkeleton:

### Current (lines 78-93):
```tsx
function MatchListSkeleton({ count = 9 }: { count?: number }) {
    return (
        <div className="w-full px-2 md:px-12">
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-[repeat(auto-fit,minmax(270px,1fr))] sm:gap-6">
                {Array.from({ length: count }).map((_, i) => (
                    <div key={i} className="bg-background-card/80 min-w-67.5 rounded-xl p-4 shadow-xl">
                        <Skeleton className="mb-2 h-6 w-3/4" />
                        <Skeleton className="mb-4 h-4 w-1/2" />
                        <Skeleton className="mb-1 h-4 w-full" />
                        <Skeleton className="h-4 w-2/3" />
                    </div>
                ))}
            </div>
        </div>
    )
}
```

### New:
```tsx
function MatchListSkeleton({ count = 9 }: { count?: number }) {
    return (
        <div className="w-full px-2 md:px-12">
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-[repeat(auto-fit,minmax(270px,1fr))] sm:gap-6">
                {Array.from({ length: count }).map((_, i) => (
                    <MatchCardSkeleton key={i} />
                ))}
            </div>
        </div>
    )
}
```

### Add import at top of MatchList.tsx:
```tsx
import { Button, EmptyState, ErrorState, MatchCardSkeleton } from ".."
```

---

## 7. Update __root.tsx for Route Transitions

Modify [`ui/src/routes/__root.tsx`](ui/src/routes/__root.tsx:1) - **MINIMAL CHANGE**:

### Current (lines 6-18):
```tsx
const RootLayout = () => {
    return (
        <ErrorBoundary>
            <div className="bg-background flex min-h-screen flex-col overflow-x-hidden">
                <Navbar />
                <div className="flex-1 pt-14 sm:pt-16">
                    <Outlet />
                </div>
                <Footer />
                {import.meta.env.DEV && <TanStackRouterDevtools />}
            </div>
        </ErrorBoundary>
    )
}
```

### New:
```tsx
const RootLayout = () => {
    return (
        <ErrorBoundary>
            <div className="bg-background flex min-h-screen flex-col overflow-x-hidden">
                <Navbar />
                <div className="page-enter flex-1 pt-14 sm:pt-16">
                    <Outlet />
                </div>
                <Footer />
                {import.meta.env.DEV && <TanStackRouterDevtools />}
            </div>
        </ErrorBoundary>
    )
}
```

**Changes:**
- Added `page-enter` class to the content div (the one containing `<Outlet />`)
- No restructuring - just added one class name

---

## Files to Create/Modify

### Modify:

| File | Changes |
|------|---------|
| [`ui/src/components/index.ts`](ui/src/components/index.ts:15) | Export `MatchCardSkeleton` |
| [`ui/src/components/Skeleton.tsx`](ui/src/components/Skeleton.tsx:7-19) | Update Skeleton to use shimmer, add MatchCardSkeleton |
| [`ui/src/components/matches/MatchList.tsx`](ui/src/components/matches/MatchList.tsx:78-93) | Use MatchCardSkeleton in MatchListSkeleton, add import |
| [`ui/src/routes/__root.tsx`](ui/src/routes/__root.tsx:1) | Add page-enter class (minimal change) |
| [`ui/src/index.css`](ui/src/index.css:1) | Add page-enter, loading-bar keyframes. NOTE: NO shimmer - use from Plan 2 |

---

## Testing Checklist

- [ ] SeriesCardSkeleton shows shimmer animation (not pulse)
- [ ] MatchCardSkeleton shows shimmer animation
- [ ] Skeleton variants (text, button, card) all use shimmer
- [ ] MatchList uses MatchCardSkeleton (proper layout)
- [ ] SeriesList skeleton (SeriesCardSkeleton) updated
- [ ] Route transitions have fade-in animation on page content
- [ ] No layout shift during skeleton to content transition
- [ ] Animations respect prefers-reduced-motion (including skeleton-shimmer)
- [ ] No duplicate keyframe definitions

---

## Commit Message

```
feat(loading): add shimmer skeletons and route transitions

- Replace animate-pulse with shimmer animation on all skeletons
- Reference shimmer from Plan 2 (avoid duplicate keyframe)
- Create MatchCardSkeleton component with proper layout
- Export MatchCardSkeleton from components/index.ts
- Update MatchListSkeleton to use MatchCardSkeleton
- Add page-enter animation class for route transitions
- Add prefers-reduced-motion for skeleton-shimmer
```

---

## Dependencies

- **Plan 2 must be implemented first** (for shimmer keyframe and .skeleton-shimmer class)

## Coordination Note

**IMPORTANT:** This plan references `@keyframes shimmer` and `.skeleton-shimmer` from Plan 2. Do NOT define them here.
