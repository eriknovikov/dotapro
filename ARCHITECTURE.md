# Dotapro Architecture

## Overview

Dota 2 professional match data aggregation platform. Scrapes from OpenDota API, stores in PostgreSQL, provides REST API with React frontend for visualization.

## Infrastructure

| Service | Purpose |
|---------|---------|
| AWS RDS (PostgreSQL) | Primary database in public VPC |
| Lambda (API) | REST API server |
| Lambda (Scraper) | Data ingestion from OpenDota |
| S3 + CloudFront | Static frontend hosting |
| API Gateway | API routing |
| EventBridge | Scheduled scraper (1/min) |

## Tech Stack

| Component | Technology |
|-----------|------------|
| API | Go 1.25, Chi router, Bun ORM + pgx, zerolog |
| Database | PostgreSQL 16, golang-migrate |
| Frontend | React, Tailwind v4, Vite, TanStack Router/Query/Table, shadcn/ui, pnpm |

## Directory Structure

```
api/          - REST API Lambda (matches/, series/, config/, db/, errs/, types/, utils/)
scraper/      - Scraper Lambda (main.go, scraper.go, types.go, querybuilder.go, utils.go, config/)
database/     - Schemas, migrations, queries, makefile
ui/           - React app (src/components/, src/routes/, src/api/, src/hooks/, src/lib/)
```

## Data Flow

```
Users → React UI → CloudFront → API Gateway → API Lambda → RDS → JSON
Scraper Lambda (EventBridge 1/min) → OpenDota → RDS
```

## Database Schema

**Tables**: `leagues`, `teams`, `players`, `users`, `series`, `series_match`, `matches`, `matches_metadata`, `user_feeds`, `user_feed_filters`

**Indexes**: GIN on array columns (heroes, players), BTREE on foreign keys and timestamps

**Triggers**: `after_match_insert` auto-creates series and updates scores

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/matches` | List matches (filters: league, team, player, hero, sort, cursor pagination) |
| GET | `/matches/:id` | Match details |
| GET | `/series` | Series list (filters: league, team, sort, cursor pagination) |
| GET | `/series/:id` | Matches in a series |

## Frontend Design System

**Theme**: Dota 2 Dark
- Dire Red: `hsl(0,84%,50%)`
- Radiant Gold: `hsl(38,92%,50%)`
- Deep Dark: `hsl(0,0%,4%)`

**Components**: Button, Input, Select, Card, Badge, Spinner, EmptyState, ErrorState (shadcn/ui + custom)

**Tokens**: Full color scales, custom spacing (18, 88, 128), border radius, shadows, transitions

## Development Guidelines

### React
- TanStack Router for type-safe routing
- TanStack Query for data fetching
- `filtersRef` for non-rerendering filter state
- `router.navigate({ replace: true })` for URL updates without rerender
- TypeScript throughout

### Go
- Descriptive naming, small focused functions
- Separate packages by functionality
- Types in `types.go`, struct tags for JSON/Bun
- `*int64` only for nullable DB fields (captains), otherwise use zero values
- Prefix OpenDota types with `OD`
- Always pass `context.Context`, check errors early
- zerolog for logging, minimal comments

## Local Development

```bash
# Set env vars
set -a; source ./.env.local; set +a

# API
cd api && go run main.go

# Scraper
cd scraper && go run main.go

# Migrations
cd database && make migrate-up

# Frontend
cd ui && npm run dev
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ENVIRON` | `local` or `prod` |
| `LOCAL_DB_URL` | Local DB connection (local only) |
| `DB_URL_PARAM_NAME` | SSM parameter name (prod only) |
| `LOCAL_ADDR` | API address (local only, default `localhost:8080`) |
| `SCRAPING_LIMIT` | Matches per scrape (default 800) |

## CI/CD Pipeline

```
Push to main → build Lambdas → deploy to AWS → build React → upload to S3 → invalidate CloudFront
```

## Lambda Configuration

| Lambda | Trigger | Memory | Timeout |
|--------|---------|--------|---------|
| API | API Gateway | 256 MB | 30 sec |
| Scraper | EventBridge (1/min) | 512 MB | 5 min |

## Performance

- **Database**: Array columns with GIN indexes, denormalized data, batch inserts (800 per transaction)
- **API**: Connection pooling via pgx, 5-10 sec request timeouts
- **Scraper**: Batch processing (800), retry logic (3 attempts), HTTP connection pooling

## Security

RDS in public subnet for cost efficiency (no NAT GW ~$32/mo, no bastion hosts)
