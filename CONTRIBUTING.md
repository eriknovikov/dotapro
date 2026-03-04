# Contributing

Thanks for wanting to contribute! Here's how to get started.

## Quick Start

```bash
# Fork the repo on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/dotapro.git
cd dotapro

# Add the original repo as upstream (for syncing later)
git remote add upstream https://github.com/E-nkv/dotapro.git
```

## Set Up Local Database

The easiest way is with Docker Compose:

```bash
# Start PostgreSQL
docker-compose up -d postgres

# Check it's running
docker ps
```

That's it! The database will be available at `localhost:5432` with:
- User: `postgres`
- Password: `admin`
- Database: `dotapro`

## Set Up the Project

```bash
# Load environment variables
source .env.local

# Run database migrations
cd database && make migrate-up && cd ..

# Install Go dependencies
cd api && go mod download && cd ..

# Install Node dependencies
cd ui && pnpm install && cd ..
```

## Running Everything

```bash
# Terminal 1: Run the API
cd api && go run main.go

# Terminal 2: Run the scraper (optional, for data)
cd scraper && MAX_BATCHES=10 go run main.go

# Terminal 3: Run the frontend
cd ui && pnpm dev
```

Then open `http://localhost:5173`

## Making Changes

1. Create a branch: `git checkout -b feature/your-thing`
2. Make your changes
3. Test them
4. Commit with a clear message
5. Push to your fork: `git push origin feature/your-thing`
6. Open a PR on GitHub

## Commit Messages

Keep it simple and descriptive:
- `feat(api): add hero filter`
- `fix(ui): broken layout on mobile`
- `docs: update contributing guide`

## Testing

- Test your changes manually
- Run tests: `cd api && go test ./...` and `cd ui && pnpm test`
- Don't obsess over 100% coverage - focus on what matters

## Code Style

- **Go**: Follow Effective Go, use `gofmt`
- **React**: Use TypeScript, follow ESLint/Prettier
- Write code that explains itself
- Add comments only when something isn't obvious

## AI is Welcome

Using AI to help write code is fine. Just make sure:
- It works
- It's readable
- You review it carefully

## Pull Requests

1. Fill out the PR template
2. Address feedback promptly
3. Keep it focused - one PR, one thing

---

Be cool, have fun, and thanks again! 🎮
