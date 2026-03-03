# API Documentation

**Base URL**: `https://api.dotapro.org` (prod) or `http://localhost:8080` (local)

No auth required. No rate limits (for now).

## Response format

All responses are JSON.

**Success**:

```json
{
  "data": { ... },
  "pagination": {
    "nc": 12345678,
    "has_more": true
  }
}
```

**Error**:

```json
{
  "error": "Error message"
}
```

## Pagination

Use cursor-based pagination. Pass `c` (cursor) from the previous response to get the next page.

| Param   | Type | Description                              |
| ------- | ---- | ---------------------------------------- |
| `limit` | int  | Results per page (default: 20, max: 100) |
| `c`     | int  | Cursor for next page                     |

## Endpoints

### Matches

**GET /matches** - List matches

| Param    | Type   | Description                    |
| -------- | ------ | ------------------------------ |
| `league` | int    | Filter by league ID            |
| `team`   | int    | Filter by team ID              |
| `player` | int    | Filter by player ID            |
| `hero`   | int    | Filter by hero ID              |
| `sort`   | string | `newest` (default) or `oldest` |
| `limit`  | int    | Results per page               |
| `c`      | int    | Cursor for pagination          |

```bash
curl "http://localhost:8080/matches?league=123&limit=10"
```

**GET /matches/:id** - Get match details

```bash
curl "http://localhost:8080/matches/7890123456"
```

### Series

**GET /series** - List series

| Param    | Type   | Description                    |
| -------- | ------ | ------------------------------ |
| `league` | int    | Filter by league ID            |
| `team`   | int    | Filter by team ID              |
| `sort`   | string | `newest` (default) or `oldest` |
| `limit`  | int    | Results per page               |
| `c`      | int    | Cursor for pagination          |

```bash
curl "http://localhost:8080/series?league=123&limit=10"
```

**GET /series/:id** - Get matches in a series

```bash
curl "http://localhost:8080/series/456789"
```

### Filters

**GET /filters** - Get available filter options

```bash
curl "http://localhost:8080/filters"
```

Returns available leagues and teams for filtering.

## HTTP Status Codes

| Code | Description  |
| ---- | ------------ |
| 200  | Success      |
| 400  | Bad request  |
| 404  | Not found    |
| 500  | Server error |
| 504  | Timeout      |

## Common errors

- `invalid parameter: X` - Invalid parameter value
- `context deadline exceeded` - Request took too long
- `internal server error` - Something broke on our end

---

That's it. Simple and fast.
