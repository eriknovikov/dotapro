# Section 2: Card Animations & Interactions

## Objective

Enhance SeriesCard and MatchCard components with polished hover effects, entrance animations, and micro-interactions.

---

## Critical Review Findings

**Confirmed via code inspection:**
- [`SeriesList.tsx:65-67`](ui/src/components/series/SeriesList.tsx:65-67) - Current card rendering: `{series.map(s => (<SeriesCard key={s.series_id} series={s} />))}`
- [`MatchList.tsx:56-58`](ui/src/components/matches/MatchList.tsx:56-58) - Current card rendering: `{matches.map(match => (<MatchCard key={match.match_id} match={match} />))}`
- Both use `SeriesCardSkeleton` in loading state
- MatchList already has a basic `MatchListSkeleton` function (lines 78-93) but not a proper `MatchCardSkeleton`

**Issues identified in original plan:**

1. **Conflicting Card wrapper modifications:** The original plan modifies the Card component in SeriesCard/MatchCard directly. However, these Card components are local to each file and not shared. This is acceptable but we should be careful about applying hover transforms that might conflict with existing shadows.

2. **drop-shadow with currentColor issue:** The original plan correctly identified that `drop-shadow` doesn't work with `currentColor` on `<img>` elements. The fix using `brightness-125 hover:saturate-150` is a reasonable alternative but won't give true color glow.

3. **Shimmer duplicate:** Plan 5 also covers shimmer animation. We must ensure shimmer keyframes are defined only in THIS plan and referenced by Plan 5.

4. **Grid layout mismatch:** The stagger delay approach works but we need to ensure the grid layout (`sm:grid-cols-[repeat(auto-fit,minmax(450px,1fr))]` for Series, `sm:grid-cols-[repeat(auto-fit,minmax(270px,1fr))]` for Matches) doesn't cause layout shifts during animation.

**Resolved decisions:**
- Use `brightness-125 hover:saturate-150` for logo glow (no true color glow but performant)
- Define fadeInUp keyframes AND shimmer keyframes in THIS plan (Plan 5 will reference them)
- Apply hover transforms on Card wrapper, not individual elements

---

## Components Affected

- [`ui/src/components/series/SeriesCard.tsx`](ui/src/components/series/SeriesCard.tsx:8-19) - Card component with hover
- [`ui/src/components/matches/MatchCard.tsx`](ui/src/components/matches/MatchCard.tsx:8-19) - Card component with hover
- [`ui/src/components/series/SeriesList.tsx`](ui/src/components/series/SeriesList.tsx:65-67) - Stagger animation
- [`ui/src/components/matches/MatchList.tsx`](ui/src/components/matches/MatchList.tsx:56-58) - Stagger animation
- [`ui/src/index.css`](ui/src/index.css:1) - All keyframes and animations

---

## 1. CSS Keyframes & Base Styles

Add to [`ui/src/index.css`](ui/src/index.css:1) (append before closing, after line 298):

```css
/* ============================================
   SECTION 2: Card Animations
   ============================================ */

/* Staggered fade-in-up for card entrances */
@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(20px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

/* Subtle float for hero image or decorative elements */
@keyframes float {
    0%, 100% {
        transform: translateY(0px);
    }
    50% {
        transform: translateY(-6px);
    }
}

/* Shimmer animation for skeletons - ALSO used by Plan 5 */
@keyframes shimmer {
    0% {
        background-position: -200% 0;
    }
    100% {
        background-position: 200% 0;
    }
}

/* Utility classes for card animations */
.card-entrance {
    animation: fadeInUp 0.4s ease-out forwards;
    opacity: 0;
}

.card-float {
    animation: float 3s ease-in-out infinite;
}

/* Hover depth effect - applied via group-hover on card wrapper */
.card-depth {
    transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.card-depth:hover {
    transform: translateY(-4px) scale(1.01);
    box-shadow: 
        0 20px 25px -5px hsl(0, 0%, 0% / 0.5),
        0 8px 10px -6px hsl(0, 0%, 0% / 0.4);
}

/* Shimmer skeleton utility - used by Plan 5 */
.skeleton-shimmer {
    background: linear-gradient(
        90deg,
        hsl(0, 0%, 10%) 25%,
        hsl(0, 0%, 14%) 50%,
        hsl(0, 0%, 10%) 75%
    );
    background-size: 200% 100%;
    animation: shimmer 1.5s ease-in-out infinite;
}

/* Respect prefers-reduced-motion */
@media (prefers-reduced-motion: reduce) {
    .card-entrance {
        animation: none;
        opacity: 1;
    }
    
    .card-depth:hover {
        transform: none;
    }
    
    .card-float {
        animation: none;
    }
    
    .skeleton-shimmer {
        animation: none;
    }
}
```

---

## 2. Staggered Card Entrance Animation

Modify [`ui/src/components/series/SeriesList.tsx`](ui/src/components/series/SeriesList.tsx:65-67):

### Current (lines 65-67):
```tsx
{series.map(s => (
    <SeriesCard key={s.series_id} series={s} />
))}
```

### New:
```tsx
{series.map((s, index) => (
    <SeriesCard 
        key={s.series_id} 
        series={s}
        style={{
            animationDelay: `${Math.min(index, 10) * 50}ms`,
        } as React.CSSProperties}
        className="card-entrance"
    />
))}
```

Modify [`ui/src/components/matches/MatchList.tsx`](ui/src/components/matches/MatchList.tsx:56-58):

### Current (lines 56-58):
```tsx
{matches.map(match => (
    <MatchCard key={match.match_id} match={match} />
))}
```

### New:
```tsx
{matches.map((match, index) => (
    <MatchCard 
        key={match.match_id} 
        match={match}
        style={{
            animationDelay: `${Math.min(index, 10) * 50}ms`,
        } as React.CSSProperties}
        className="card-entrance"
    />
))}
```

**Stagger Calculation:**
- `Math.min(index, 10) * 50ms` caps at 500ms delay for 11+ cards
- First 10 cards stagger up to 500ms total
- Cards beyond 10 start immediately (no additional delay)

---

## 3. Team Logo Hover Glow

Modify [`SeriesCard.tsx`](ui/src/components/series/SeriesCard.tsx:38-47) team logo section.

**Current team logo (lines 38-47):**
```tsx
{series.team_a.logo_url && (
    <img
        src={series.team_a.logo_url}
        alt={`${series.team_a.name} logo`}
        className="h-6 w-auto max-w-8 shrink-0 rounded select-none sm:h-7 sm:max-w-11"
    />
)}
```

**New:**
```tsx
{series.team_a.logo_url && (
    <img
        src={series.team_a.logo_url}
        alt={`${series.team_a.name} logo`}
        className="h-6 w-auto max-w-8 shrink-0 rounded select-none transition-all duration-300 hover:brightness-125 hover:saturate-150 sm:h-7 sm:max-w-11"
    />
)}
```

Apply same pattern to `series.team_b.logo_url`.

**For MatchCard:** Apply same `hover:brightness-125 hover:saturate-150` to team logos.

---

## 4. Card Depth Effect (Hover Transform)

Modify the Card component in [`SeriesCard.tsx`](ui/src/components/series/SeriesCard.tsx:8-19):

### Current Card (lines 8-19):
```tsx
const Card = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(({ className, ...props }, ref) => (
    <div
        ref={ref}
        className={cn(
            "bg-background-card/80 text-card-foreground relative overflow-hidden rounded-xl shadow-xl backdrop-blur-sm",
            className,
        )}
        {...props}
    >
        <div className="relative z-10 flex h-full flex-col">{props.children}</div>
    </div>
))
```

### New Card:
```tsx
const Card = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(({ className, ...props }, ref) => (
    <div
        ref={ref}
        className={cn(
            "bg-background-card/80 text-card-foreground relative overflow-hidden rounded-xl shadow-xl backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:scale-[1.01] hover:shadow-2xl",
            className,
        )}
        {...props}
    >
        <div className="relative z-10 flex h-full flex-col">{props.children}</div>
    </div>
))
```

Apply same change to [`MatchCard.tsx`](ui/src/components/matches/MatchCard.tsx:8-19) Card component.

---

## 5. Accessibility Considerations

Added `prefers-reduced-motion` media query in CSS section above to disable animations for users who prefer reduced motion.

---

## Files to Modify

| File | Changes |
|------|---------|
| [`ui/src/index.css`](ui/src/index.css:1) | Add fadeInUp, float, shimmer keyframes; card-entrance/card-depth utility classes; skeleton-shimmer class; prefers-reduced-motion |
| [`ui/src/components/series/SeriesList.tsx`](ui/src/components/series/SeriesList.tsx:65-67) | Add stagger animation to SeriesCard |
| [`ui/src/components/matches/MatchList.tsx`](ui/src/components/matches/MatchList.tsx:56-58) | Add stagger animation to MatchCard |
| [`ui/src/components/series/SeriesCard.tsx`](ui/src/components/series/SeriesCard.tsx:8-19) | Add hover depth effect to Card, logo glow |
| [`ui/src/components/matches/MatchCard.tsx`](ui/src/components/matches/MatchCard.tsx:8-19) | Add hover depth effect to Card, logo glow |

---

## Testing Checklist

- [ ] Cards fade in with stagger when list loads
- [ ] Cards lift up (`-translate-y-1 scale-[1.01]`) on hover
- [ ] Shadow increases on hover (`shadow-2xl`)
- [ ] Team logos brighten/saturate on hover
- [ ] Animations disabled when `prefers-reduced-motion: reduce`
- [ ] No layout shift during hover transitions
- [ ] Works on both Series and Matches pages
- [ ] Mobile touch devices: hover states don't stick

---

## Commit Message

```
feat(cards): add staggered entrance animations and hover effects

- Add fadeInUp keyframe animation with card-entrance utility
- Add hover depth effect (lift + shadow increase) on card wrappers
- Add team logo brightness/saturation boost on hover
- Cap stagger delay at 500ms (10 cards) to prevent long wait times
- Define shimmer keyframe for use by Plan 5
- Respect prefers-reduced-motion for accessibility
```

---

## Dependencies

None - pure CSS animations with existing Tailwind classes.

## Note for Plan 5

**Plan 5 should reference the shimmer keyframes defined here, not redefine them.**
