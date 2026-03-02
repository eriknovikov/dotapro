# Dotapro

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5+-3178C6?logo=typescript)

> Professional Dota 2 match data platform. Scrapes from OpenDota, stores in PostgreSQL, serves via API with a React frontend.

## What it does

- Track series with match-by-match breakdowns
- Browse matches with filters (league, team, player, hero)
- Auto-updates every minute
- Beautiful dark-themed UI

## Stack

**Backend**: Go + PostgreSQL + AWS Lambda
**Frontend**: React + TypeScript + Tailwind

## Quick start

```bash
# Clone
git clone https://github.com/E-nkv/dotapro.git
cd dotapro

# Set up env
source .env.local

# Set up database
cd database && make migrate-up

# Run API
cd ../api && go run main.go

# Run scraper (optional, for data)
cd ../scraper && MAX_BATCHES=10 go run main.go

# Run frontend
cd ../ui && pnpm install && pnpm dev
```

Then open `http://localhost:5173`

## Or use Docker

```bash
docker-compose up -d
```

## Project structure

```
api/       - Go backend (REST API)
scraper/   - Go scraper (data ingestion)
database/  - SQL schemas and migrations
ui/        - React frontend
```

## API

| Endpoint       | Method | Description               |
| -------------- | ------ | ------------------------- |
| `/matches`     | GET    | List matches with filters |
| `/matches/:id` | GET    | Get match details         |
| `/series`      | GET    | List series with filters  |
| `/series/:id`  | GET    | Get matches in a series   |

## Deployment

See [ARCHITECTURE.md](ARCHITECTURE.md) for infrastructure details.

- [MANUAL_DEPLOYMENT.md](MANUAL_DEPLOYMENT.md) - First-time AWS setup
- [DEPLOYMENT.md](DEPLOYMENT.md) - Automated CI/CD setup

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT

---

Made with ❤️ for the Dota 2 community
