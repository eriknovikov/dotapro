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

export async function getSeriesMOCK(params: Filters, signal: AbortSignal): Promise<GetSeriesResponse> {
    const seriesResponseMock = await import("../../../res.json")
    return seriesResponseMock as GetSeriesResponse
}

export async function getSeries(params: Filters, signal: AbortSignal): Promise<GetSeriesResponse> {
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
