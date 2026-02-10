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
- Frontend: React + Tailwind + Vite
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
- frontend/ - React + Vite app

## App Flow

Users request from React UI via CloudFront to API Gateway, routes to API Lambda, queries RDS, returns JSON. Scraper Lambda runs every minute via EventBridge, queries OpenDota for new matches, processes and inserts to RDS.

## Core Data

See database/ directory for schemas and queries.

## API Endpoints

See api/docs.yaml for full API documentation.

- GET /matches - List matches with filters (league, team, player, hero, sort by newest/oldest)
- GET /matches/:id - Match details (id, league, radiant, dire, etc)
- GET /series - Series list by team, league (id, team1, team2, score, league, date)
- GET /series/:id - Matches in a series

## Go Guidelines

- Use descriptive naming, code should be self-descriptive
- Keep functions small and focused
- Separate related functionalities into packages
- Define types in dedicated types.go files
- DONT Use pointer types for nullable DB fields (*int64, *string), instead, use the default 0 value to determine non-existence.
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
- Frontend: cd frontend && npm run dev

## Environment Variables

- ENVIRONMENT - local or prod
- LOCAL_DB_URL - Local DB connection string
- DB_URL_PARAM_NAME - SSM parameter name for prod

Note that to run code locally, the flow is typically to do set -a; source ./.env.local; and then run the code that requires the env vars set (like running the api or scraper).

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
