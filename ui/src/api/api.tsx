export type Filters = {
    league?: number
    team?: number
    sort?: string
    limit?: number
    c?: number
}

// @deprecated Use Filters instead
export type Params = Filters

export type TeamInfo = {
    id: number
    name: string
    tag?: string
    logo_url?: string
}

export type LeagueInfo = {
    id: number
    name: string
    tier: string
}

export type LeagueSearchResult = {
    league_id: number
    name: string
}

export type Series = {
    series_id: number
    start_time: string
    team_a: TeamInfo
    team_b: TeamInfo
    league: LeagueInfo
    team_a_score: number
    team_b_score: number
}

export type Pagination = {
    nc?: number
    pc?: number
    has_more: boolean
}

export type GetSeriesResponse = {
    series: Series[]
    pagination: Pagination
}

export async function getSeries(params: Filters, signal: AbortSignal): Promise<GetSeriesResponse> {
    await new Promise(res => setTimeout(res, 2000))
    const url = new URL("http://localhost:8080/series")
    if (params.league !== undefined) {
        url.searchParams.set("league", String(params.league))
    }
    if (params.team !== undefined) {
        url.searchParams.set("team", String(params.team))
    }
    if (params.sort !== undefined) {
        url.searchParams.set("sort", params.sort)
    }
    if (params.limit !== undefined) {
        url.searchParams.set("limit", String(params.limit))
    }
    if (params.c !== undefined) {
        url.searchParams.set("c", String(params.c))
    }
    const res = await fetch(url.toString(), { signal })
    if (!res.ok) {
        throw new Error(`Failed to fetch series: ${res.status} ${res.statusText}`)
    }
    const data = await res.json()!
    return data
}

export async function searchLeagues(query: string, signal: AbortSignal): Promise<LeagueSearchResult[]> {
    const url = new URL("http://localhost:8080/filtersmetadata/leagues")
    url.searchParams.set("q", query)
    const res = await fetch(url.toString(), { signal })
    if (!res.ok) {
        throw new Error(`Failed to search leagues: ${res.status} ${res.statusText}`)
    }
    const data = await res.json()!
    return data
}

export type TeamSearchResult = {
    team_id: number
    name: string
    logo_url?: string
}

export async function searchTeams(query: string, signal: AbortSignal): Promise<TeamSearchResult[]> {
    const url = new URL("http://localhost:8080/filtersmetadata/teams")
    url.searchParams.set("q", query)
    const res = await fetch(url.toString(), { signal })
    if (!res.ok) {
        throw new Error(`Failed to search teams: ${res.status} ${res.statusText}`)
    }
    const data = await res.json()!
    return data
}
