const API_BASE_URL = "http://localhost:8080"

export type Filters = {
    league?: number
    team?: number
    sort?: string
    limit?: number
    c?: number
}

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
    has_more: boolean
}

export type GetSeriesResponse = {
    series: Series[]
    pagination: Pagination
}

export type TeamSearchResult = {
    team_id: number
    name: string
    logo_url?: string
}

function buildUrl(path: string, params: Record<string, string | number | undefined>): URL {
    const url = new URL(`${API_BASE_URL}${path}`)
    Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
            url.searchParams.set(key, String(value))
        }
    })
    return url
}

async function fetchApi<T>(url: URL, signal: AbortSignal, errorMessage: string): Promise<T> {
    const res = await fetch(url.toString(), { signal })
    if (!res.ok) {
        throw new Error(`${errorMessage}: ${res.status} ${res.statusText}`)
    }
    return res.json()
}

export async function getSeries(params: Filters, signal: AbortSignal): Promise<GetSeriesResponse> {
    await new Promise(res => setTimeout(res, 2000))
    const url = buildUrl("/series", params)
    return fetchApi(url, signal, "Failed to fetch series")
}

export async function searchLeagues(query: string, signal: AbortSignal): Promise<LeagueSearchResult[]> {
    const url = buildUrl("/filtersmetadata/leagues", { q: query })
    return fetchApi(url, signal, "Failed to search leagues")
}

export async function searchTeams(query: string, signal: AbortSignal): Promise<TeamSearchResult[]> {
    const url = buildUrl("/filtersmetadata/teams", { q: query })
    return fetchApi(url, signal, "Failed to search teams")
}
