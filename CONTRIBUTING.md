# Contributing

Thanks for wanting to contribute!

## Quick start

```bash
# Fork and clone
git clone https://github.com/YOUR_USERNAME/dotapro.git
cd dotapro

# Set up env
source .env.local

# Set up database
cd database && make migrate-up

# Install deps
cd api && go mod download
cd ../ui && pnpm install
```

## Making changes

1. Create a branch: `git checkout -b feature/your-thing`
2. Make your changes
3. Test them
4. Commit with a clear message
5. Push and open a PR

## Commit messages

Keep it simple: `feat(api): add hero filter` or `fix(ui): broken layout on mobile`

## Testing

- Test your changes
- Run tests: `cd api && go test ./...` and `cd ui && pnpm test`
- Don't obsess over 100% coverage - focus on what matters

## Code style

- Go: follow Effective Go, use gofmt
- React: use TypeScript, follow ESLint/Prettier
- Write code that explains itself
- Add comments only when something isn't obvious

## AI is welcome

Using AI to help write code is fine. Just make sure:
- It works
- It's readable
- You review it carefully

## PRs

Fill out the template. Address feedback. That's it.

---

Be cool, have fun, and thanks again! 🎮
