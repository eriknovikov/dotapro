# Dotapro API Documentation

## Overview

The Dotapro API provides RESTful endpoints for accessing professional Dota 2 match data, including matches, series, teams, players, and leagues. The API is built with Go and uses cursor-based pagination for efficient data retrieval.

**Base URL**: `https://api.dotapro.com` (production) or `http://localhost:8080` (local)

## Authentication

Currently, the API does not require authentication. This may change in future versions.

## Rate Limiting

Currently, there are no rate limits enforced. This may change in future versions.

## Response Format

All API responses are in JSON format.

### Success Response

```json
{
  "data": { ... },
  "pagination": {
    "nc": 12345678,
    "has_more": true
  }
}
```

### Error Response

```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 400 | Bad Request - Invalid parameters |
| 404 | Not Found |
| 500 | Internal Server Error |
| 504 | Gateway Timeout - Request took too long |

## Pagination

The API uses cursor-based pagination for efficient data retrieval. Use the `c` (cursor) parameter to fetch the next page of results.

### Pagination Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `limit` | integer | Number of results per page (default: 20, max: 100) |
| `c` | integer | Cursor for the next page (from `pagination.nc` in previous response) |

### Pagination Response

| Field | Type | Description |
|-------|------|-------------|
| `nc` | integer | Next cursor value for pagination |
| `has_more` | boolean | Whether more results are available |

---

## Endpoints

### Matches

#### Get Matches

Retrieve a list of matches with optional filtering.

**Endpoint**: `GET /matches`

**Query Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `league` | integer | No | Filter by league ID |
| `team` | integer | No | Filter by team ID |
| `player` | integer | No | Filter by player ID |
| `hero` | integer | No | Filter by hero ID |
| `sort` | string | No | Sort order: `newest` (default) or `oldest` |
| `limit` | integer | No | Results per page (default: 20, max: 100) |
| `c` | integer | No | Cursor for pagination |

**Example Request**:

```bash
curl "http://localhost:8080/matches?league=123&limit=10&sort=newest"
```

**Example Response**:

```json
{
  "matches": [
    {
      "match_id": 7890123456,
      "start_time": "2024-02-28T12:00:00Z",
      "duration": 2345,
      "radiant_win": true,
      "radiant_team": {
        "id": 123456,
        "name": "Team Spirit",
        "tag": "TS",
        "logo_url": "https://example.com/logo.png"
      },
      "dire_team": {
        "id": 789012,
        "name": "Team Liquid",
        "tag": "TL",
        "logo_url": "https://example.com/logo2.png"
      },
      "league": {
        "id": 123,
        "name": "The International 2024",
        "tier": "premium"
      },
      "series_id": 456789,
      "radiant_heroes": [1, 2, 3, 4, 5],
      "dire_heroes": [6, 7, 8, 9, 10]
    }
  ],
  "pagination": {
    "nc": 7890123455,
    "has_more": true
  }
}
```

#### Get Match Details

Retrieve detailed information about a specific match.

**Endpoint**: `GET /matches/{id}`

**Path Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Match ID |

**Example Request**:

```bash
curl "http://localhost:8080/matches/7890123456"
```

**Example Response**:

```json
{
  "match_id": 7890123456,
  "start_time": "2024-02-28T12:00:00Z",
  "duration": 2345,
  "radiant_win": true,
  "patch": "7.36",
  "version": 1,
  "radiant_team": {
    "id": 123456,
    "name": "Team Spirit",
    "tag": "TS",
    "logo_url": "https://example.com/logo.png",
    "captain": 987654
  },
  "dire_team": {
    "id": 789012,
    "name": "Team Liquid",
    "tag": "TL",
    "logo_url": "https://example.com/logo2.png",
    "captain": 654321
  },
  "league": {
    "id": 123,
    "name": "The International 2024",
    "tier": "premium"
  },
  "series_id": 456789,
  "picks_bans": [
    {
      "is_pick": true,
      "hero_id": 1,
      "team": 0,
      "order": 1
    }
  ],
  "players_data": [
    {
      "player_id": 123456,
      "hero_id": 1,
      "kills": 10,
      "deaths": 3,
      "assists": 15,
      "gold_per_min": 650,
      "xp_per_min": 720,
      "hero_damage": 25000,
      "tower_damage": 5000,
      "healing": 3000,
      "level": 25,
      "net_worth": 18000
    }
  ],
  "radiant_gold_adv": [0, 100, 200, 150, ...],
  "radiant_xp_adv": [0, 50, 100, 75, ...],
  "radiant_heroes": [1, 2, 3, 4, 5],
  "dire_heroes": [6, 7, 8, 9, 10],
  "radiant_players": [123456, 234567, 345678, 456789, 567890],
  "dire_players": [678901, 789012, 890123, 901234, 012345]
}
```

---

### Series

#### Get Series

Retrieve a list of series with optional filtering.

**Endpoint**: `GET /series`

**Query Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `league` | integer | No | Filter by league ID |
| `team` | integer | No | Filter by team ID |
| `sort` | string | No | Sort order: `newest` (default) or `oldest` |
| `limit` | integer | No | Results per page (default: 20, max: 100) |
| `c` | integer | No | Cursor for pagination |

**Example Request**:

```bash
curl "http://localhost:8080/series?league=123&limit=10"
```

**Example Response**:

```json
{
  "series": [
    {
      "series_id": 456789,
      "start_time": "2024-02-28T12:00:00Z",
      "team_a": {
        "id": 123456,
        "name": "Team Spirit",
        "tag": "TS",
        "logo_url": "https://example.com/logo.png",
        "score": 2
      },
      "team_b": {
        "id": 789012,
        "name": "Team Liquid",
        "tag": "TL",
        "logo_url": "https://example.com/logo2.png",
        "score": 1
      },
      "league": {
        "id": 123,
        "name": "The International 2024",
        "tier": "premium"
      },
      "team_a_score": 2,
      "team_b_score": 1
    }
  ],
  "pagination": {
    "nc": 456788,
    "has_more": true
  }
}
```

#### Get Series Details

Retrieve detailed information about a specific series, including all matches.

**Endpoint**: `GET /series/{id}`

**Path Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Series ID |

**Example Request**:

```bash
curl "http://localhost:8080/series/456789"
```

**Example Response**:

```json
{
  "series_id": 456789,
  "start_time": "2024-02-28T12:00:00Z",
  "team_a": {
    "id": 123456,
    "name": "Team Spirit",
    "tag": "TS",
    "logo_url": "https://example.com/logo.png"
  },
  "team_b": {
    "id": 789012,
    "name": "Team Liquid",
    "tag": "TL",
    "logo_url": "https://example.com/logo2.png"
  },
  "league": {
    "id": 123,
    "name": "The International 2024",
    "tier": "premium"
  },
  "team_a_score": 2,
  "team_b_score": 1,
  "matches": [
    {
      "match_id": 7890123456,
      "duration": 2345,
      "radiant_win": true,
      "picks_bans": [...],
      "players_data": [...],
      "radiant_gold_adv": [0, 100, 200, ...],
      "radiant_xp_adv": [0, 50, 100, ...],
      "radiant_score": 25,
      "dire_score": 18,
      "radiant_captain": 987654,
      "dire_captain": 654321
    }
  ]
}
```

---

### Filters Metadata

#### Get Available Filters

Retrieve available filter options for leagues, teams, and players.

**Endpoint**: `GET /filters`

**Example Request**:

```bash
curl "http://localhost:8080/filters"
```

**Example Response**:

```json
{
  "leagues": [
    {
      "league_id": 123,
      "name": "The International 2024"
    }
  ],
  "teams": [
    {
      "team_id": 123456,
      "name": "Team Spirit",
      "logo_url": "https://example.com/logo.png"
    }
  ]
}
```

---

## Data Models

### MatchSummary

| Field | Type | Description |
|-------|------|-------------|
| `match_id` | integer | Unique match identifier |
| `start_time` | datetime | Match start time (ISO 8601) |
| `duration` | integer | Match duration in seconds |
| `radiant_win` | boolean | Whether Radiant team won |
| `radiant_team` | TeamInfo | Radiant team information |
| `dire_team` | TeamInfo | Dire team information |
| `league` | LeagueInfo | League information |
| `series_id` | integer | Series identifier |
| `radiant_heroes` | array[int] | Radiant hero IDs |
| `dire_heroes` | array[int] | Dire hero IDs |

### MatchDetail

| Field | Type | Description |
|-------|------|-------------|
| `match_id` | integer | Unique match identifier |
| `start_time` | datetime | Match start time (ISO 8601) |
| `duration` | integer | Match duration in seconds |
| `radiant_win` | boolean | Whether Radiant team won |
| `patch` | string | Dota 2 patch version |
| `version` | integer | Data version |
| `radiant_team` | TeamInfo | Radiant team information |
| `dire_team` | TeamInfo | Dire team information |
| `league` | LeagueInfo | League information |
| `series_id` | integer | Series identifier |
| `picks_bans` | array | Pick/ban data |
| `players_data` | array | Player statistics |
| `radiant_gold_adv` | array[int] | Radiant gold advantage over time |
| `radiant_xp_adv` | array[int] | Radiant XP advantage over time |
| `radiant_heroes` | array[int] | Radiant hero IDs |
| `dire_heroes` | array[int] | Dire hero IDs |
| `radiant_players` | array[int] | Radiant player IDs |
| `dire_players` | array[int] | Dire player IDs |

### TeamInfo

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Unique team identifier |
| `name` | string | Team name |
| `tag` | string | Team tag (optional) |
| `logo_url` | string | Team logo URL (optional) |
| `score` | integer | Team score in series (optional) |
| `captain` | integer | Captain player ID (optional) |

### LeagueInfo

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Unique league identifier |
| `name` | string | League name |
| `tier` | string | League tier (e.g., "premium", "professional") |
| `logo_url` | string | League logo URL (optional) |

### SeriesSummary

| Field | Type | Description |
|-------|------|-------------|
| `series_id` | integer | Unique series identifier |
| `start_time` | datetime | Series start time (ISO 8601) |
| `team_a` | TeamInfo | Team A information |
| `team_b` | TeamInfo | Team B information |
| `league` | LeagueInfo | League information |
| `team_a_score` | integer | Team A score |
| `team_b_score` | integer | Team B score |

### SeriesDetail

| Field | Type | Description |
|-------|------|-------------|
| `series_id` | integer | Unique series identifier |
| `start_time` | datetime | Series start time (ISO 8601) |
| `team_a` | TeamInfo | Team A information |
| `team_b` | TeamInfo | Team B information |
| `league` | LeagueInfo | League information |
| `team_a_score` | integer | Team A score |
| `team_b_score` | integer | Team B score |
| `matches` | array[SeriesMatchDetail] | Matches in the series |

### SeriesMatchDetail

| Field | Type | Description |
|-------|------|-------------|
| `match_id` | integer | Unique match identifier |
| `duration` | integer | Match duration in seconds |
| `radiant_win` | boolean | Whether Radiant team won |
| `picks_bans` | array | Pick/ban data |
| `players_data` | array | Player statistics |
| `radiant_gold_adv` | array[int] | Radiant gold advantage over time |
| `radiant_xp_adv` | array[int] | Radiant XP advantage over time |
| `radiant_score` | integer | Radiant team kills |
| `dire_score` | integer | Dire team kills |
| `radiant_captain` | integer | Radiant captain player ID (optional) |
| `dire_captain` | integer | Dire captain player ID (optional) |

---

## Error Handling

### Error Response Format

```json
{
  "error": "Error message description"
}
```

### Common Errors

| Error | Description |
|-------|-------------|
| `invalid parameter: league_id` | Invalid league ID parameter |
| `invalid parameter: team_id` | Invalid team ID parameter |
| `invalid parameter: player_id` | Invalid player ID parameter |
| `invalid parameter: hero_id` | Invalid hero ID parameter |
| `invalid parameter: limit` | Invalid limit parameter |
| `invalid parameter: cursor` | Invalid cursor parameter |
| `context deadline exceeded` | Request took too long to process |
| `internal server error` | Unexpected server error |

---

## SDKs and Libraries

### JavaScript/TypeScript

```typescript
import { createClient } from '@dotapro/api-client';

const client = createClient({
  baseURL: 'https://api.dotapro.com'
});

// Get matches
const matches = await client.matches.list({
  league: 123,
  limit: 10
});

// Get match details
const match = await client.matches.get(7890123456);
```

### Go

```go
package main

import (
    "context"
    "github.com/dotapro/go-client"
)

func main() {
    client := dotapro.NewClient("https://api.dotapro.com")
    
    // Get matches
    matches, err := client.Matches.List(context.Background(), &dotapro.MatchFilter{
        LeagueID: dotapro.Int64(123),
        Limit:    10,
    })
    
    // Get match details
    match, err := client.Matches.Get(context.Background(), 7890123456)
}
```

---

## Changelog

### Version 1.0.0 (2024-02-28)

- Initial API release
- Matches endpoints
- Series endpoints
- Filters metadata endpoint
- Cursor-based pagination

---

## Support

For API support, please:
- Check the [GitHub Issues](https://github.com/yourusername/dotapro/issues)
- Start a [Discussion](https://github.com/yourusername/dotapro/discussions)
- Contact the maintainers
