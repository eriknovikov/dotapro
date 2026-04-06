# Section 3: Visual Effects (Hero Parallax, Background)

## Objective

Enhance the homepage with parallax scrolling effects on the hero image, background animations, and Dota 2-themed visual polish.

---

## Critical Review Findings

**Confirmed via code inspection:**
- [`ui/src/routes/index.tsx:47-67`](ui/src/routes/index.tsx:47-67) - Hero image section with 3D transform already exists
- Current hero uses `perspective-[2000px]` and `transform-[rotateX(-18deg)_rotateY(7deg)_rotateZ(-0.75deg)]`
- Hero image is wrapped in a container with `overflow-hidden`

**Issues identified in original plan:**

1. **3D transform conflict:** The hero already has a complex 3D transform. Adding parallax (translateY based on scroll) could conflict with or break the existing 3D effect. 
   - **Resolution:** Instead of scroll-based parallax, implement a subtle "float" animation that simulates depth without conflicting with the 3D transform.

2. **Radar/scanline effect:** Could be too distracting or gimmicky. 
   - **Resolution:** Make it very subtle (low opacity) and only show occasionally.

3. **Performance concern with canvas noise:** Canvas-based noise texture at 10fps could be CPU intensive on mobile.
   - **Resolution:** Use CSS gradient animation instead (much lighter).

4. **CSS gradient mesh conflicts:** Background already has `hsl(0, 0%, 4%)` base. Adding gradient mesh might clash.
   - **Resolution:** Apply gradient to a pseudo-element with low opacity so it adds atmosphere without overriding.

5. **CRITICAL FIX - float keyframe:** The original plan defined `@keyframes float` which DUPLICATES the one in Plan 2.
   - **Resolution:** DO NOT define float here - reference the one from Plan 2's CSS.

**Resolved decisions:**
- Use CSS gradient animation (Option A) for background
- Use float animation instead of scroll parallax - but reference float from Plan 2, don't redefine
- Radar effect: subtle, low opacity, periodic
- Apply effects to specific elements, not global background

---

## Files to Modify

- [`ui/src/routes/index.tsx`](ui/src/routes/index.tsx:47-67) - Main hero section
- [`ui/src/index.css`](ui/src/index.css:1) - New keyframes and effects (NOTE: Do NOT define float here - Plan 2 defines it)

---

## 1. CSS Keyframes & Base Styles

Add to [`ui/src/index.css`](ui/src/index.css:1) (append before closing, after Plan 2's CSS):

```css
/* ============================================
   SECTION 3: Visual Effects
   ============================================ */

/* NOTE: @keyframes float is defined in Plan 2's CSS section */
/* Do NOT redefine it here - reference .card-float class instead */

/* Animated gradient mesh background - subtle atmospheric effect */
@keyframes gradientShift {
    0%, 100% {
        background-position: 0% 50%;
        opacity: 0.3;
    }
    25% {
        background-position: 50% 100%;
        opacity: 0.4;
    }
    50% {
        background-position: 100% 50%;
        opacity: 0.3;
    }
    75% {
        background-position: 50% 0%;
        opacity: 0.4;
    }
}

/* Gradient mesh overlay - applied to hero container */
.bg-gradient-mesh {
    position: relative;
}

.bg-gradient-mesh::before {
    content: "";
    position: absolute;
    inset: 0;
    background: linear-gradient(
        135deg,
        hsl(0, 10%, 8%) 0%,
        hsl(0, 5%, 6%) 25%,
        hsl(0, 10%, 8%) 50%,
        hsl(0, 5%, 6%) 75%,
        hsl(0, 10%, 8%) 100%
    );
    background-size: 400% 400%;
    animation: gradientShift 20s ease infinite;
    opacity: 0.5;
    z-index: 0;
    pointer-events: none;
}

/* Scanline effect for Dota/esports aesthetic */
@keyframes scanline {
    0% {
        transform: translateY(-100%);
        opacity: 0;
    }
    5% {
        opacity: 0.15;
    }
    95% {
        opacity: 0.15;
    }
    100% {
        transform: translateY(100vh);
        opacity: 0;
    }
}

.scanline-effect {
    position: relative;
    overflow: hidden;
}

.scanline-effect::after {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 100px;
    background: linear-gradient(
        to bottom,
        transparent,
        hsl(0, 84%, 50% / 0.1),
        transparent
    );
    animation: scanline 8s linear infinite;
    pointer-events: none;
    z-index: 5;
}

/* Radiant/Dire glow effects */
@keyframes radiantGlow {
    0%, 100% {
        box-shadow: 0 0 20px hsl(38, 92%, 50% / 0.3);
    }
    50% {
        box-shadow: 0 0 40px hsl(38, 92%, 50% / 0.5);
    }
}

@keyframes direGlow {
    0%, 100% {
        box-shadow: 0 0 20px hsl(0, 84%, 50% / 0.3);
    }
    50% {
        box-shadow: 0 0 40px hsl(0, 84%, 50% / 0.5);
    }
}

.radiant-glow {
    animation: radiantGlow 3s ease-in-out infinite;
}

.dire-glow {
    animation: direGlow 3s ease-in-out infinite;
}

/* Floating Dota icons animation */
@keyframes iconFloat {
    0%, 100% {
        transform: translateY(0) rotate(0deg);
    }
    25% {
        transform: translateY(-10px) rotate(5deg);
    }
    75% {
        transform: translateY(5px) rotate(-5deg);
    }
}

.icon-float {
    animation: iconFloat 4s ease-in-out infinite;
}

/* Hero image subtle float - uses float keyframe from Plan 2 */
.hero-float {
    animation: float 6s ease-in-out infinite;
}

/* NOTE: The @keyframes float is defined in Plan 2's CSS section.
   Do NOT redefine it here. The .hero-float class references it. */

/* Respect prefers-reduced-motion */
@media (prefers-reduced-motion: reduce) {
    .hero-float {
        animation: none;
    }
    
    .bg-gradient-mesh::before {
        animation: none;
        opacity: 0;
    }
    
    .scanline-effect::after {
        animation: none;
        opacity: 0;
    }
    
    .radiant-glow,
    .dire-glow {
        animation: none;
    }
    
    .icon-float {
        animation: none;
    }
}
```

---

## 2. Update Homepage Hero Section

Modify [`ui/src/routes/index.tsx`](ui/src/routes/index.tsx:47-67):

### Current (lines 47-67):
```tsx
{/* Hero Image with Linear-style 3D effect - Full width outside container */}
<div className="sm:perspective:[2000px] flex w-full items-center justify-center overflow-hidden px-2 py-1 perspective-[2000px] sm:py-2 md:py-8 lg:py-12">
    <div className="relative w-full transform-[rotateX(-18deg)_rotateY(7deg)_rotateZ(-0.75deg)] transition-shadow duration-300 ease-in-out hover:shadow-[0_0_0_1px_rgba(0,0,0,.06),0_16px_24px_rgba(0,0,0,.08),0_40px_60px_rgba(0,0,0,.06),0_80px_100px_rgba(0,0,0,.06),0_160px_200px_rgba(0,0,0,.04),0_400px_320px_rgba(0,0,0,.04)] sm:w-auto sm:transform-[rotateX(-25deg)_rotateY(10deg)_rotateZ(-1deg)]">
        {/* Ambient Glow (Light Bleed) */}
        <div
            className="bg-primary-500/10 absolute -inset-2 opacity-30 blur-3xl sm:-inset-4"
            aria-hidden="true"
        />

        {/* Main Image Container */}
        <div className="bg-background-card relative overflow-hidden rounded-xl shadow-[0_0_0_1px_rgba(0,0,0,.04),0_4px_6px_rgba(0,0,0,.02),0_10px_15px_rgba(0,0,0,.015),0_20px_25px_rgba(0,0,0,.015),0_40px_50px_rgba(0,0,0,.01),0_100px_80px_rgba(0,0,0,.01)]">
            <img
                src="/hero.webp"
                alt="Series Screenshot"
                width="2525"
                height="1222"
                className="block h-auto w-full sm:max-w-[80vw] md:max-w-[90vw] lg:max-w-[80vw]"
            />
        </div>
    </div>
</div>
```

### New:
```tsx
{/* Hero Image with Linear-style 3D effect - Full width outside container */}
<div className="sm:perspective:[2000px] bg-gradient-mesh relative flex w-full items-center justify-center overflow-hidden px-2 py-1 perspective-[2000px] sm:py-2 md:py-8 lg:py-12">
    {/* Scanline effect overlay */}
    <div className="scanline-effect absolute inset-0 z-10 pointer-events-none" aria-hidden="true" />
    
    <div className="hero-float relative w-full transform-[rotateX(-18deg)_rotateY(7deg)_rotateZ(-0.75deg)] transition-shadow duration-300 ease-in-out hover:shadow-[0_0_0_1px_rgba(0,0,0,.06),0_16px_24px_rgba(0,0,0,.08),0_40px_60px_rgba(0,0,0,.06),0_80px_100px_rgba(0,0,0,.06),0_160px_200px_rgba(0,0,0,.04),0_400px_320px_rgba(0,0,0,.04)] sm:w-auto sm:transform-[rotateX(-25deg)_rotateY(10deg)_rotateZ(-1deg)]">
        {/* Ambient Glow (Light Bleed) - with Radiant gold tint */}
        <div
            className="bg-primary-500/10 absolute -inset-2 opacity-30 blur-3xl sm:-inset-4"
            aria-hidden="true"
        />

        {/* Main Image Container */}
        <div className="bg-background-card relative overflow-hidden rounded-xl shadow-[0_0_0_1px_rgba(0,0,0,.04),0_4px_6px_rgba(0,0,0,.02),0_10px_15px_rgba(0,0,0,.015),0_20px_25px_rgba(0,0,0,.015),0_40px_50px_rgba(0,0,0,.01),0_100px_80px_rgba(0,0,0,.01)]">
            <img
                src="/hero.webp"
                alt="Series Screenshot"
                width="2525"
                height="1222"
                className="block h-auto w-full sm:max-w-[80vw] md:max-w-[90vw] lg:max-w-[80vw]"
            />
        </div>
    </div>
</div>
```

**Changes made:**
1. Added `bg-gradient-mesh` class to container for subtle animated background
2. Added `scanline-effect` overlay div for subtle scanline
3. Added `hero-float` class to the inner transform container for gentle floating motion

---

## 3. Floating Dota Icons (Optional Enhancement)

If desired, add floating Dota-themed icons near the hero section. This is purely decorative.

**Decision:** Deferred to future enhancement. The current effects are sufficient for initial implementation.

---

## Files to Modify

| File | Changes |
|------|---------|
| [`ui/src/index.css`](ui/src/index.css:1) | Add gradientShift, scanline, glow keyframes and utility classes. NOTE: Do NOT add float keyframe - use from Plan 2 |
| [`ui/src/routes/index.tsx`](ui/src/routes/index.tsx:47-67) | Add bg-gradient-mesh, scanline-effect, hero-float classes |

---

## Testing Checklist

- [ ] Background has subtle animated gradient effect
- [ ] Scanline sweeps across hero section periodically (every 8s)
- [ ] Hero image has subtle floating motion (not scroll-based)
- [ ] Hover shadow enhancement still works
- [ ] All animations disabled when `prefers-reduced-motion: reduce`
- [ ] No performance issues on mobile
- [ ] No layout conflicts with 3D transform
- [ ] Effects are subtle, not distracting
- [ ] No duplicate keyframe definitions (float comes from Plan 2)

---

## Commit Message

```
feat(homepage): add visual effects - gradient mesh, scanline, hero float

- Add subtle animated gradient mesh background to hero section
- Add periodic scanline effect for esports aesthetic
- Use float animation from Plan 2 (avoid duplicate keyframe)
- Respect prefers-reduced-motion for accessibility
```

---

## Dependencies

Plan 2 must be implemented first (for the `float` keyframe).

## Coordination Note

**IMPORTANT:** This plan references `@keyframes float` from Plan 2's CSS. Do NOT define it here.
