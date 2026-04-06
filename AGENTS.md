# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Stack
- **UI**: React 19 + TypeScript, Vite, Tailwind CSS v4, TanStack Router/Query
- **API**: Go 1.25 (Lambda), Bun + pgx for PostgreSQL
- **Scraper**: Go 1.25 (Lambda), AWS SDK v2, cleanenv, zerolog
- **Database**: PostgreSQL with golang-migrate

## Commands

### Root
- `make build` - Build both Lambda functions (API + Scraper) for Linux ARM64
- `make deploy` - Deploy both to AWS Lambda
- `make run-api` / `make run-scraper` - Run locally (requires `.env.local`)
- `make run-ui` - Run UI dev server (requires `ui/.env.local`)
- `make deploy-ui` - Build UI, upload to S3 with cache headers, invalidate CloudFront

### UI (`ui/`)
- `pnpm dev` / `pnpm build` / `pnpm lint` / `pnpm format`
- Tests not configured yet

### API (`api/`)
- `make test` - Run all Go tests
- `make tidy` - Clean up Go modules
- `make run` - Run locally

### Database (`database/`)
- `make migrate-up` / `make migrate-down` - Run migrations (requires `LOCAL_DB_URL`)
- `make migrate-new name=xxx` - Create migration pair
- `make migrate-force v=1` - Fix dirty state

## Code Style

### Go
- Linter: golangci-lint with errcheck, goconst, govet, revive, staticcheck, unused
- Staticcheck excludes ST1000, ST1005 (no requirement for package comments or doc for exported vars)
- Logs use zerolog

### React/TypeScript
- 4-space indent, 120 char line width
- Prettier handles formatting, ESLint handles quality
- Tailwind CSS v4 (uses `@tailwindcss/vite` plugin, not PostCSS)

## Project Architecture

- **API/Scraper**: Separate Lambda functions deployed via AWS, config via AWS SSM Parameter Store
- **UI**: SPA served from S3 + CloudFront, routing via TanStack Router (file-based)
- **Scraper**: Reads from OpenDota API, writes to PostgreSQL via Bun ORM

## Go Patterns (Non-Obvious)
- Use `*int64` only for nullable DB fields (captains), otherwise use zero values
- Prefix OpenDota types with `OD` to distinguish them
- Always pass `context.Context`, check errors early
- zerolog for logging, minimal comments

## React Patterns (Non-Obvious)
- Use `filtersRef` for non-rerendering filter state
- Use `router.navigate({ replace: true })` for URL updates without rerender

## Database (Non-Obvious)
- `after_match_insert` trigger auto-creates series and updates scores
- Connection pooling via pgx
- Batch inserts use 800 per transaction for performance

## Scraper (Non-Obvious)
- `MAX_BATCHES` default 50 to stay below 60 rpm OpenDota ratelimit
- Retry logic: 3 attempts with HTTP connection pooling

## Gotchas
- Go builds target `linux/arm64` (not your native OS/arch)
- `api/` and `scraper/` are separate Go modules - run `go mod tidy` from each directory
- Database migrations are forward-only by design (no rollback migrations in this project)
- `make deploy-ui` separates index.html (no cache) from assets (1-year cache) and triggers CloudFront invalidation
- RDS in public subnet for cost efficiency (no NAT GW ~$32/mo)
