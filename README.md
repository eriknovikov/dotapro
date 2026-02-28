# Dotapro

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5+-3178C6?logo=typescript)

> A professional Dota 2 match data aggregation and visualization platform.

## Overview

Dotapro is an open-source platform that aggregates professional Dota 2 match data from the OpenDota API and provides a clean, modern interface for analyzing matches, series, player statistics, and team performance. Built with a serverless architecture on AWS, it's designed to be cost-effective, scalable, and easy to deploy.

## Features

- **Series Tracking**: View complete series with match-by-match breakdowns
- **Match Browser** (on the way): Browse professional matches with advanced filtering by league, team, player, and hero
- **Real-time Updates**: Automated and bulk-optimized scraping keeps data fresh (every minute)
- **Responsive Design**: Beautiful dark-themed UI optimized for desktop and mobile
- **Fast Performance**: Optimized queries and caching for quick data retrieval

## Tech Stack

### Backend

- **Language**: Go 1.25+
- **Framework**: Chi router
- **ORM**: Bun with pgx driver
- **Database**: PostgreSQL 16
- **Logging**: zerolog
- **Infrastructure**: AWS Lambda, API Gateway, RDS, S3, CloudFront, EventBridge

### Frontend

- **Framework**: React 18+ with TypeScript
- **Build Tool**: Vite
- **Routing**: TanStack Router
- **Data Fetching**: TanStack Query
- **UI Components**: shadcn/ui
- **Styling**: Tailwind CSS v4
- **Package Manager**: pnpm

## Quick Start

### Prerequisites

- Go 1.25+
- Node.js 20+
- pnpm 9+
- PostgreSQL 16+

### Local Development

1. **Clone the repository**

   ```bash
   git clone https://github.com/E-nkv/dotapro.git
   cd dotapro
   ```

2. **Set up environment variables**

   The `.env.local` file is included in the repository. Review and modify it as needed for your local setup.

3. **Set up the database**

   First, source the environment variables:

   ```bash
   set -a; source .env.local; set +a
   ```

   Then run the migrations:

   ```bash
   cd database
   make migrate-up
   ```

   Initialize the scraper metadata (run this in your preferred DB editor like DBeaver, or via psql):

   ```sql
   INSERT INTO scraper_metadata (id, last_fetched_match_id) VALUES (1, 7944311818);
   ```

4. **Start the API**

   ```bash
   cd ../api
   go run main.go
   ```

5. **Start the scraper** (for data ingestion)

   Set the scraping limit to avoid hitting the OpenDota rate limit:

   ```bash
   cd ../scraper
   SCRAPING_LIMIT=35000 go run main.go
   ```

6. **Start the frontend**

   ```bash
   cd ../ui
   pnpm install
   pnpm dev
   ```

7. **Open your browser**
   Navigate to `http://localhost:5173`

### Using Docker Compose

```bash
docker-compose up -d
```

This will start PostgreSQL, the API, and the UI with hot-reload enabled.

## Project Structure

```
dotapro/
├── api/              # REST API Lambda
│   ├── matches/      # Match endpoints
│   ├── series/       # Series endpoints
│   ├── config/       # Configuration
│   ├── db/           # Database connection
│   ├── errs/         # Error handling
│   ├── types/        # Shared types
│   └── utils/        # Utilities
├── scraper/          # Data scraper Lambda
│   ├── config/       # Configuration
│   └── *.go          # Scraper logic
├── database/         # Database schemas and migrations
│   ├── migrations/   # SQL migrations
│   └── queries/      # SQL queries
├── ui/               # React frontend
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── routes/      # Page routes
│   │   ├── api/         # API client
│   │   ├── hooks/       # Custom hooks
│   │   └── lib/         # Utilities
│   └── public/          # Static assets
├── ARCHITECTURE.md   # Architecture documentation
├── CONTRIBUTING.md   # Contribution guidelines
└── docs/             # Additional documentation
```

## API Documentation

For detailed API documentation, see [docs/api.md](docs/api.md).

### Quick API Reference

| Endpoint       | Method | Description               |
| -------------- | ------ | ------------------------- |
| `/matches`     | GET    | List matches with filters |
| `/matches/:id` | GET    | Get match details         |
| `/series`      | GET    | List series with filters  |
| `/series/:id`  | GET    | Get matches in a series   |

## Deployment

Dotapro is designed for AWS deployment using serverless architecture. See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed infrastructure information.

### Manual Deployment

1. Build the Lambda functions
2. Deploy to AWS Lambda
3. Configure API Gateway
4. Build and deploy the React app to S3
5. Set up CloudFront distribution

### CI/CD

The project includes GitHub Actions workflows for automated deployment on push to main.

## Development

### Running Tests

```bash
# API tests
cd api && go test ./...

# Frontend tests
cd ui && pnpm test
```

### Code Style

- **Go**: Follow standard Go conventions, use `gofmt`
- **TypeScript**: ESLint + Prettier configuration included
- **React**: Functional components with hooks, TypeScript throughout

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [OpenDota](https://www.opendota.com/) for providing the Dota 2 API
- [Valve](https://www.valvesoftware.com/) for Dota 2

## Roadmap

- [ ] User authentication and preferences
- [ ] Custom filters and saved searches
- [ ] Data export functionality
- [ ] Mobile app (React Native)
- [ ] Real-time match updates
- [ ] Advanced analytics and insights
- [ ] Team comparison tools

## Support

- 📖 [Documentation](docs/)
- 🐛 [Issue Tracker](https://github.com/E-nkv/dotapro/issues)
- 💬 [Discussions](https://github.com/E-nkv/dotapro/discussions)

---

Made with ❤️ for the Dota 2 community
