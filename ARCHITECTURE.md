# Dotapro Architecture

## Overview

Dota 2 match data aggregation system. Scrapes from OpenDota API, stores in PostgreSQL, provides REST API. React frontend for visualization.

## Infrastructure

AWS RDS for database, Lambda for API and scraper, S3 + CloudFront for frontend, Cognito for auth. RDS in public VPC for cost efficiency.

## Tech Stack

- API: Go 1.25
- DB: PostgreSQL 16
- ORM: Bun + pgx
- Migrations: golang-migrate
- Frontend: React + Tailwind v4 + Vite
  - TanStack Router for routing
  - TanStack Query for data fetching
  - TanStack Table for data tables
  - shadcn/ui for UI components
- Logging: zerolog

## Directory Structure

- api/ - REST API Lambda
  - matches/ - Match endpoints
  - config/ - Environment config
  - db/ - Database connection
  - errs/ - Error handling
- scraper/ - Scraper Lambda
  - main.go - Lambda entry point
  - scraper.go - Scraping logic
  - types.go - Data models
  - querybuilder.go - SQL for OpenDota Explorer
  - utils.go - DB setup, SSM integration
  - config/ - Environment config
- database/ - Database schemas
  - migrations/ - Migration files
  - queries/ - Raw SQL queries
  - makefile - Migration commands
- ui/ - React + Vite app
  - src/
    - components/
      - ui/ - Shadcn/ui components (add with `npx shadcn@latest add <component>`)
      - FiltersSidebar.tsx
      - SeriesList.tsx
    - lib/
      - utils.ts - Utility functions (cn() for class merging)
    - routes/
      - __root.tsx
      - index.tsx
      - about.tsx
    - index.css - Global styles, design tokens, and CSS variables

## App Flow

Users request from React UI via CloudFront to API Gateway, routes to API Lambda, queries RDS, returns JSON. Scraper Lambda runs every minute via EventBridge, queries OpenDota for new matches, processes and inserts to RDS.

## Core Data

See database/ directory for schemas and queries.

## API Endpoints

See api/docs.yaml for full API documentation.

- GET /matches - List matches with filters (league, team, player, hero, sort by newest/oldest, cursor-based pagination with `c` parameter)
- GET /matches/:id - Match details (id, league, radiant, dire, etc)
- GET /series - Series list by team, league (id, team1, team2, score, league, date, cursor-based pagination with `c` parameter)
- GET /series/:id - Matches in a series

## Frontend Design System

### Overview

The design system is built on top of Tailwind CSS v4 and provides:
- **Design Tokens**: Consistent values for colors, spacing, typography, shadows, and transitions
- **Reusable Components**: Pre-built UI components that follow the design system
- **Utility Functions**: Helper functions for common styling patterns

**Theme**: Dota 2 Inspired Dark Theme
- **Primary**: Rich Red (Dire faction) - `hsl(0, 84%, 50%)`
- **Secondary**: Golden Yellow (Radiant faction) - `hsl(38, 92%, 50%)`
- **Background**: Deep Dark - `hsl(0, 0%, 4%)`

### Design Tokens

#### Colors

**Primary Colors - Rich Red (Dire Faction)**
- `primary-50` to `primary-950`: Full red color scale for primary actions and accents
- Primary base: `hsl(0, 84%, 50%)` - Rich Red

**Secondary Colors - Golden Yellow (Radiant Faction)**
- `secondary-50` to `secondary-950`: Full golden yellow color scale for secondary actions and highlights
- Secondary base: `hsl(38, 92%, 50%)` - Golden Yellow

**Semantic Colors - Dark Theme**
- `background`: Deep dark backgrounds (`hsl(0, 0%, 4%)`)
- `background-muted`: Muted backgrounds (`hsl(0, 0%, 7%)`)
- `background-accent`: Accent backgrounds (`hsl(0, 0%, 10%)`)
- `background-card`: Card backgrounds (`hsl(0, 0%, 8%)`)
- `foreground`: Light text (`hsl(0, 0%, 98%)`)
- `foreground-muted`: Muted text (`hsl(0, 0%, 64%)`)
- `foreground-subtle`: Subtle text (`hsl(0, 0%, 45%)`)
- `border`: Default borders (`hsl(0, 0%, 18%)`)
- `border-muted`: Muted borders (`hsl(0, 0%, 12%)`)
- `border-accent`: Accent borders (`hsl(0, 0%, 25%)`)

**Status Colors**
- `success`: Green scale for success states (`hsl(142, 71%, 45%)` base)
- `warning`: Golden yellow scale for warnings (`hsl(38, 92%, 50%)` base)
- `error`: Red scale for errors (`hsl(0, 84%, 50%)` base)

#### Spacing
- Custom spacing values: `18` (4.5rem), `88` (22rem), `128` (32rem)
- Standard Tailwind spacing: `xs` (0.25rem), `sm` (0.5rem), `md` (1rem), `lg` (1.5rem), `xl` (2rem)

#### Typography
- Font sizes: `xxs` (0.625rem) in addition to standard Tailwind sizes
- Font weights: `font-medium`, `font-semibold`, `font-bold`

#### Border Radius
- `sm` (0.25rem), `md` (0.375rem), `lg` (0.5rem), `xl` (0.75rem), `2xl` (1rem), `3xl` (1.5rem)

#### Shadows
- `sm`, `DEFAULT`, `md`, `lg`, `xl`, `2xl`, `inner`, `none`

#### Transitions
- Durations: `150`, `200`, `300` (ms)
- Timing functions: `ease-in-out`

### Reusable Components

#### Button (shadcn/ui)

A versatile button component with multiple variants and sizes.

```tsx
import { Button } from "@/components/ui/button"

// Primary button (Rich Red)
<Button onClick={handleClick}>Click me</Button>

// Secondary button (Golden Yellow)
<Button variant="secondary">Cancel</Button>

// Outline button
<Button variant="outline">Outline</Button>

// Ghost button
<Button variant="ghost">Ghost</Button>

// Destructive button
<Button variant="destructive">Delete</Button>

// Different sizes
<Button size="sm">Small</Button>
<Button size="default">Default</Button>
<Button size="lg">Large</Button>
<Button size="icon"><Icon /></Button>
```

**Props:**
- `variant`: `"default" | "destructive" | "outline" | "secondary" | "ghost" | "link"` (default: `"default"`)
- `size`: `"default" | "sm" | "lg" | "icon"` (default: `"default"`)
- All standard HTML button attributes

#### Input (shadcn/ui)

A text input component.

```tsx
import { Input } from "@/components/ui/input"

// Basic input
<Input type="email" placeholder="Enter your email" />

// Disabled input
<Input disabled placeholder="Disabled input" />

// With custom styling
<Input className="bg-background-muted" placeholder="Custom styled" />
```

**Props:**
- All standard HTML input attributes
- Use `className` for custom styling

#### Select (shadcn/ui)

A select dropdown component.

```tsx
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"

<Select>
    <SelectTrigger>
        <SelectValue placeholder="Select a country" />
    </SelectTrigger>
    <SelectContent>
        <SelectItem value="us">United States</SelectItem>
        <SelectItem value="ca">Canada</SelectItem>
    </SelectContent>
</Select>
```

#### Card (shadcn/ui)

A card container with header, content, and footer sections.

```tsx
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from "@/components/ui/card"

<Card>
    <CardHeader>
        <CardTitle>Card Title</CardTitle>
        <CardDescription>Card description</CardDescription>
    </CardHeader>
    <CardContent>
        <p>Card content goes here</p>
    </CardContent>
    <CardFooter>
        <Button>Action</Button>
    </CardFooter>
</Card>
```

#### Badge (shadcn/ui)

A small badge/tag for displaying status or labels.

```tsx
import { Badge } from "@/components/ui/badge"

<Badge>Default</Badge>
<Badge variant="secondary">Secondary</Badge>
<Badge variant="outline">Outline</Badge>
<Badge variant="destructive">Destructive</Badge>
```

**Props:**
- `variant`: `"default" | "secondary" | "destructive" | "outline"` (default: `"default"`)
- All standard HTML span attributes

#### Spinner (Custom)

A loading spinner animation.

```tsx
import { Spinner } from "@/components/ui/spinner"

<Spinner />
<Spinner size="lg" />
```

**Props:**
- `size`: `"sm" | "md" | "lg" | "xl"` (default: `"md"`)
- All standard HTML div attributes

#### EmptyState (Custom)

A component to display when there's no data.

```tsx
import { EmptyState } from "@/components/ui/empty-state"

<EmptyState
    title="No results found"
    description="Try adjusting your filters to find what you're looking for."
    action={<Button>Clear Filters</Button>}
/>
```

**Props:**
- `title`: `string` (default: `"No data found"`)
- `description`: `string`
- `action`: `ReactNode` - Action button/element
- `size`: `"sm" | "md" | "lg"` (default: `"md"`)
- `icon`: `ReactNode` - Icon/illustration
- All standard HTML div attributes

#### ErrorState (Custom)

A component to display error states.

```tsx
import { ErrorState } from "@/components/ui/error-state"

<ErrorState
    error={new Error("Failed to load data")}
    action={<Button onClick={retry}>Retry</Button>}
/>
```

**Props:**
- `error`: `Error | string` - Error object or message
- `title`: `string` (default: `"Error"`)
- `action`: `ReactNode` - Action button/element
- `size`: `"sm" | "md" | "lg"` (default: `"md"`)
- All standard HTML div attributes

### Utility Functions

#### cn()

A utility function to merge Tailwind CSS classes using clsx.

```tsx
import { cn } from "@/lib/utils"

const className = cn(
    "px-4 py-2",
    isActive && "bg-primary-500",
    "text-foreground"
)
```

### Using Design Tokens in Tailwind Classes

#### Colors

Use the semantic color tokens for consistent theming:

```tsx
// Backgrounds
<div className="bg-background">Main background</div>
<div className="bg-background-muted">Muted background</div>
<div className="bg-background-card">Card background</div>
<div className="bg-background-accent">Accent background</div>

// Text
<p className="text-foreground">Default text</p>
<p className="text-foreground-muted">Muted text</p>
<p className="text-foreground-subtle">Subtle text</p>

// Primary (Rich Red)
<button className="bg-primary-500 text-white">Primary action</button>
<button className="bg-primary-600 hover:bg-primary-700">Hover state</button>

// Secondary (Golden Yellow)
<button className="bg-secondary-500 text-foreground">Secondary action</button>

// Borders
<div className="border-border">Default border</div>
<div className="border-border-muted">Muted border</div>
<div className="border-border-accent">Accent border</div>

// Status
<div className="text-success-500">Success message</div>
<div className="text-warning-500">Warning message</div>
<div className="text-error-500">Error message</div>
```

#### Spacing

Use the custom spacing tokens:

```tsx
<div className="p-18">Custom padding (4.5rem)</div>
<div className="w-88">Custom width (22rem)</div>
<div className="h-128">Custom height (32rem)</div>
```

#### Border Radius

Use the radius tokens:

```tsx
<div className="rounded-sm">Small radius (0.25rem)</div>
<div className="rounded-md">Medium radius (0.375rem)</div>
<div className="rounded-lg">Large radius (0.5rem)</div>
<div className="rounded-xl">Extra large radius (0.75rem)</div>
<div className="rounded-2xl">2XL radius (1rem)</div>
<div className="rounded-3xl">3XL radius (1.5rem)</div>
```

#### Shadows

Use the shadow tokens:

```tsx
<div className="shadow-sm">Small shadow</div>
<div className="shadow">Default shadow</div>
<div className="shadow-md">Medium shadow</div>
<div className="shadow-lg">Large shadow</div>
<div className="shadow-xl">Extra large shadow</div>
<div className="shadow-2xl">2XL shadow</div>
<div className="shadow-inner">Inner shadow</div>
```

#### Typography

Use the font size tokens:

```tsx
<p className="text-xxs">Extra extra small text (0.625rem)</p>
```

### Shadcn/ui Integration

This project uses shadcn/ui components with custom Dota2-themed colors. The shadcn CSS variables are mapped to our design tokens in `ui/src/index.css`:

```css
:root {
    --primary: hsl(0, 84%, 50%);      /* Rich Red (Dire) */
    --secondary: hsl(38, 92%, 50%);   /* Golden Yellow (Radiant) */
    --background: hsl(0, 0%, 4%);      /* Deep Dark */
    --foreground: hsl(0, 0%, 98%);     /* Light text */
    /* ... etc */
}
```

To add shadcn components:
```bash
npx shadcn@latest add button
npx shadcn@latest add card
npx shadcn@latest add input
```

Components will automatically use your Dota2 color palette.

### Extending the Design System

#### Adding New Colors

Edit `ui/src/index.css` in the `@theme` block:

```css
@theme {
    /* Brand Colors */
    --color-brand-50: hsl(210, 100%, 97%);
    --color-brand-100: hsl(210, 100%, 94%);
    /* ... full scale */
    --color-brand-950: hsl(210, 100%, 12%);
}
```

Then add shadcn-compatible variables in the `:root` block:

```css
:root {
    --brand: hsl(210, 100%, 50%);
    --brand-foreground: hsl(0, 0%, 98%);
}
```

#### Adding New Components

**For shadcn/ui components:**
```bash
npx shadcn@latest add <component-name>
```

**For custom components:**
1. Create the component file in `ui/src/components/ui/`
2. Follow the existing patterns (props interface, variant/size styles)
3. Use design tokens for all colors, spacing, and typography
4. Export from `ui/src/components/ui/index.ts` (if needed)
5. Document usage in this file

#### Adding New Utilities

Add to `ui/src/lib/utils.ts`:

```typescript
export function formatDate(date: Date): string {
    // implementation
}
```

## React/Frontend Guidelines

- Use TanStack Router for client-side routing with type-safe routes
- Use TanStack Query for server state management and data fetching
- Use TanStack Table for complex data tables with sorting, filtering, pagination
- Use Tailwind CSS for utility-first styling
- Use shadcn/ui for accessible, composable UI components
- Use refs for form state that shouldn't trigger rerenders (e.g., filtersRef)
- Update URL history without triggering rerender using router.navigate with replace: true
- Component structure: Sidebar (filters) + Main Content (results) layout
- Keep filter updates local (no API calls) until search button is clicked
- Use TypeScript for type safety across the entire frontend
- Separate concerns: API layer, hooks, components, routes
- Use proper TypeScript types for all API responses and request parameters

### Component Patterns

- **Filters Sidebar**: Left panel with filter inputs and search button
- **Results Table**: Right panel taking remaining space with proper minimums
- **Filters Type**: Define TypeScript type for filter parameters
- **filtersRef**: Use useRef to hold filter state without causing rerenders
- **Search Action**: On search button click, read from filtersRef, make API request, update URL

### URL State Management

- Update URL params without triggering component rerender
- Use router.navigate() with replace: true for seamless updates
- URL params should reflect current filter state for shareable links

### Frontend Best Practices

1. **Use Design Tokens**: Always use the defined color, spacing, and typography tokens instead of hardcoded values
   - Use `bg-primary-500` instead of `bg-red-500`
   - Use `text-foreground-muted` instead of `text-gray-500`
   - Use `border-border` instead of `border-gray-700`
2. **Compose Components**: Build complex UIs by composing smaller, reusable components
3. **Accessibility**: All components include proper ARIA attributes and keyboard navigation
4. **Responsive Design**: Use Tailwind's responsive prefixes for mobile-first design
5. **Consistency**: Follow the established patterns when creating new components
6. **Dark Theme**: This is a dark-first design system. Always test components against dark backgrounds

## Go Guidelines

- Use descriptive naming, code should be self-descriptive
- Keep functions small and focused
- Separate related functionalities into packages
- Define types in dedicated types.go files
- DONT Use pointer types for nullable DB fields (*int64, *string), instead, use default 0 value to determine non-existence.
  - **Exception**: Captain fields (radiant_captain, dire_captain) use `*int64` pointer types to allow NULL values when no captain is assigned
- Separate OpenDota types (prefix OD) from DB models
- Use struct tags for JSON and Bun ORM mapping
- Always pass context.Context for cancellation
- Check errors early, return immediately
- Use zerolog for logging
- Minimal comments, let code explain itself

## Local Development

- API: cd api && go run main.go
- Scraper: cd scraper && go run main.go
- Migrations: cd database && make up
- Frontend: cd ui && npm run dev

Note that to run code locally, flow is typically to do set -a; source ./.env.local; and then run code that requires env vars set (like running api or scraper).

## Environment Variables

- ENVIRONMENT - local or prod
- LOCAL_DB_URL - Local DB connection string
- DB_URL_PARAM_NAME - SSM parameter name for prod

## CI/CD

Push to main triggers build, deploys API and Scraper Lambdas, builds React, uploads to S3, invalidates CloudFront.

## Lambda Config

- API: API Gateway trigger, 256 MB, 30 sec timeout
- Scraper: EventBridge (1/min), 512 MB, 5 min timeout

## Security

RDS in public subnet for cost optimization and ease of development (not having to set up ec2 bastion hosts, and NAT GW costs ~32bucks/mo)

## Performance

- Database: Array columns with GIN indexes, denormalized data, batch inserts (3000 per transaction).
- API: CloudFront caching, response compression, connection pooling via pgx
- Scraper: EventBridge scheduling, batch processing, retry logic (3 attempts)

## Model Response General Guidelines

- Never attempt to run commands yourself of any sort, nor go, git, make, or any command. if a command needs to be ran, say "now run `{COMMAND TO RUN}` to do `{WHATEVER WE NEED TO DO}`".

## Open Questions

- User authentication flow with Cognito
- Real-time match updates
- Player statistics aggregation
- Hero win rate analytics
- Rate limiting implementation
- Monitoring and alerting
