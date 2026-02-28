# Contributing to Dotapro

Thank you for your interest in contributing to Dotapro! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [AI-Assisted Development](#ai-assisted-development)
- [Submitting Changes](#submitting-changes)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment. Be constructive, welcoming, and respectful to all contributors.

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js 20+
- pnpm 9+
- PostgreSQL 16+
- Git

### Setting Up Your Development Environment

1. **Fork the repository**
   ```bash
   # Fork the repository on GitHub, then clone your fork
   git clone https://github.com/E-nkv/dotapro.git
   cd dotapro
   ```

2. **Add the upstream remote**
   ```bash
   git remote add upstream https://github.com/E-nkv/dotapro.git
   ```

3. **Set up environment variables**

   The `.env.local` file is included in the repository. Review and modify it as needed for your local setup.

4. **Set up the database**

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

5. **Install dependencies**
   ```bash
   # API dependencies
   cd api && go mod download

   # UI dependencies
   cd ../ui && pnpm install
   ```

## Development Workflow

### Branching Strategy

- `main` - Production-ready code
- `develop` - Integration branch for features
- `feature/*` - New features
- `bugfix/*` - Bug fixes
- `hotfix/*` - Critical production fixes

### Creating a Feature Branch

```bash
# Ensure your main branch is up to date
git checkout main
git pull upstream main

# Create a new feature branch
git checkout -b feature/your-feature-name
```

### Making Changes

1. **Write code following the [Coding Standards](#coding-standards)**
2. **Add tests for new functionality** (see [Testing](#testing) for guidelines)
3. **Ensure all tests pass**
4. **Commit your changes with descriptive messages**

### Commit Message Format

Follow conventional commits:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(api): add hero filter to matches endpoint

Implement hero filtering in the matches list endpoint with
proper validation and error handling.

Closes #123
```

```
fix(ui): correct series card layout on mobile

Fix flexbox layout issue that caused series cards to overflow
on screens smaller than 640px.
```

## Coding Standards

### Go (Backend)

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Keep functions small and focused
- Use descriptive variable and function names
- Always handle errors
- **Code should be self-explanatory where it makes sense** - avoid needless comments
- **Add concise and precise comments where it makes sense** - for complex logic, non-obvious decisions, or exported functions
- Use `context.Context` for all operations that may timeout or cancel

```go
// Good - self-explanatory, no comment needed
func (s *MatchService) GetMatch(ctx context.Context, id int64) (*Match, error) {
    match, err := s.db.GetMatch(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get match: %w", err)
    }
    return match, nil
}

// Good - comment explains non-obvious behavior
// Use exponential backoff to avoid overwhelming the API during retries
func (s *Scraper) fetchWithRetry(ctx context.Context, url string) ([]byte, error) {
    // ...
}
```

### TypeScript/React (Frontend)

- Use TypeScript for all new code
- Follow ESLint and Prettier configurations
- Use functional components with hooks
- Prefer composition over inheritance
- Keep components small and focused
- Use proper TypeScript types (avoid `any`)
- **Code should be self-explanatory where it makes sense** - avoid needless comments
- **Add concise and precise comments where it makes sense** - for complex logic, non-obvious decisions, or exported functions

```tsx
// Good - self-explanatory, no comment needed
interface MatchCardProps {
  match: Match;
  onSelect: (id: number) => void;
}

export function MatchCard({ match, onSelect }: MatchCardProps) {
  return (
    <Card onClick={() => onSelect(match.id)}>
      <MatchHeader match={match} />
      <MatchStats match={match} />
    </Card>
  );
}

// Good - comment explains non-obvious behavior
// Using keydown instead of click to support keyboard navigation
export function MatchCard({ match, onSelect }: MatchCardProps) {
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' || e.key === ' ') {
      onSelect(match.id);
    }
  };
  // ...
}
```

### SQL

- Use uppercase for SQL keywords
- Use lowercase for table and column names
- Add comments for complex queries
- Use parameterized queries to prevent SQL injection

```sql
-- Good
SELECT
    m.id,
    m.radiant_team_id,
    m.dire_team_id,
    m.start_time
FROM matches m
WHERE m.league_id = $1
ORDER BY m.start_time DESC
LIMIT $2;
```

## Testing

### Testing Philosophy

We take a pragmatic approach to testing:

- **Add tests for business logic** - core functionality that could break
- **Add tests for API endpoints** - ensure contracts are maintained
- **Add tests for utility functions** - verify edge cases are handled
- **Be flexible with UI components** - visual components that are hard to test don't require extensive test coverage
- **Skip tests for obvious functionality** - simple getters/setters or trivial operations don't need tests

### Go Tests

```bash
# Run all tests
cd api && go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

### Frontend Tests

```bash
# Run all tests
cd ui && pnpm test

# Run tests with coverage
pnpm test -- --coverage

# Run tests in watch mode
pnpm test -- --watch
```

### Test Coverage

- Aim for reasonable coverage on business logic and API endpoints
- Don't obsess over 100% coverage - focus on critical paths
- Write meaningful tests that provide value, not just for the sake of coverage

## AI-Assisted Development

**AI usage for code generation is accepted and encouraged**, as long as the contribution meets these standards:

- **Functional**: The code works as intended and doesn't introduce bugs
- **Readable**: The code is easy to understand and follows project conventions
- **Maintainable**: The code is well-structured and can be easily modified by others
- **Performant**: The code is efficient and doesn't introduce performance regressions

When using AI tools:
- Review all generated code carefully
- Ensure it follows the project's coding standards
- Add tests for AI-generated functionality
- Don't blindly accept AI suggestions - use your judgment

## Submitting Changes

### Pull Request Process

1. **Update your branch**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create a Pull Request**
   - Go to the repository on GitHub
   - Click "New Pull Request"
   - Select your feature branch
   - Fill in the PR template

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests added/updated (where applicable)
- [ ] All tests passing

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Comments added where needed (avoiding needless comments)
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Commit messages follow conventions

## Related Issues
Closes #123
```

### Review Process

1. Automated checks must pass (CI/CD)
2. At least one maintainer approval required
3. Address all review comments
4. Squash commits if requested
5. Merge when approved

## Reporting Issues

### Bug Reports

When reporting a bug, please include:

1. **Description**: Clear description of the bug
2. **Steps to reproduce**: Detailed steps to reproduce the issue
3. **Expected behavior**: What you expected to happen
4. **Actual behavior**: What actually happened
5. **Environment**:
   - OS and version
   - Go/Node versions
   - Browser version (if applicable)
6. **Screenshots**: If applicable
7. **Additional context**: Any other relevant information

### Feature Requests

When requesting a feature, please include:

1. **Description**: Clear description of the feature
2. **Use case**: Why this feature would be useful
3. **Proposed solution**: How you envision the feature working
4. **Alternatives considered**: Other approaches you've considered
5. **Additional context**: Any other relevant information

## Getting Help

- Check existing [documentation](docs/)
- Search [existing issues](https://github.com/E-nkv/dotapro/issues)
- Start a [discussion](https://github.com/E-nkv/dotapro/discussions)
- Join our community chat (if available)

## Recognition

Contributors who make significant contributions will be recognized in the project's contributors list and may be invited to become maintainers.

---

Thank you for contributing to Dotapro! 🎮
