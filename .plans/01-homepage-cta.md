# Section 1: Homepage CTA Update

## Objective

Update the homepage (`ui/src/routes/index.tsx`) to replace the "About dotapro" button with "View Matches" button.

---

## Critical Review Findings

**Confirmed via code inspection:**

- [`ui/src/routes/index.tsx:34-41`](ui/src/routes/index.tsx:34-41) - Current CTA buttons confirmed (lines were offset from original plan's 24-42)
- [`ui/src/components/ui/button.tsx:11-31`](ui/src/components/ui/button.tsx:11-31) - Button variants confirmed: `primary`, `secondary`, `outline`, `destructive`, `destructive-outline`, `cool-outline`, `ghost`, `link`, `default`
- Routes exist: `/series`, `/matches`, `/about`

**Issue identified:**

- Plan suggested `variant="secondary"` but the current secondary variant uses `bg-secondary` (gold) which is correct for differentiation. No conflicts.

**Decision:** Use `variant="cool-outline"` for the Matches button to match the styling used elsewhere in the app (Series page pagination buttons use `cool-outline`), creating consistency.

---

## Changes

### 1. Update CTA Buttons in [`ui/src/routes/index.tsx`](ui/src/routes/index.tsx:34-41)

**Current state (lines 34-41):**

```tsx
<Button
    asChild
    variant="outline"
    size="sm"
    className="px-6 text-sm sm:px-6 sm:text-sm md:px-6 md:text-sm lg:px-8 lg:text-base"
>
    <Link to="/about">About dotapro</Link>
</Button>
```

**New state:**

```tsx
<Button
    asChild
    variant="cool-outline"
    size="sm"
    className="px-6 text-sm sm:px-6 sm:text-sm md:px-6 md:text-sm lg:px-8 lg:text-base"
>
    <Link to="/matches">View Matches</Link>
</Button>
```

---

## Implementation Notes

1. **Button variant choice:** `cool-outline` is used because:
   - Primary (red) is for the main CTA "View Series"
   - Cool-outline has gradient hover effect matching app aesthetic
   - Consistent with pagination buttons on Series/Matches pages

2. **No layout changes needed** - existing gap classes `gap-2 sm:gap-3 md:gap-4` are appropriate

3. **Sizing consistency** - keep same `size="sm"` and padding classes for visual consistency

---

## Files to Modify

| File                                                       | Change                                                                                                                                         |
| ---------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| [`ui/src/routes/index.tsx`](ui/src/routes/index.tsx:34-41) | Change second Button `variant="outline"` to `variant="cool-outline"`, `to="/about"` to `to="/matches"`, text "About dotapro" to "View Matches" |

---

## Testing Checklist

- [ ] View Series button still works → navigates to `/series`
- [ ] View Matches button works → navigates to `/matches`
- [ ] Both buttons have consistent sizing and spacing
- [ ] Responsive behavior maintained on mobile (1 column on small screens)
- [ ] No console errors
- [ ] Visual hover effects work on new button

---

## Commit Message

```
feat(homepage): add View Matches CTA, remove About button

- Replace "About dotapro" button with "View Matches"
- Use cool-outline variant for visual consistency with app
- Navigate to /matches route (already exists)
```
