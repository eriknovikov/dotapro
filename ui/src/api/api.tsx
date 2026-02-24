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

export type PlayerData = {
    name: string | null
    facet: number
    kills: number
    level: number
    deaths: number
    denies: number
    item_0: number
    item_1: number
    item_2: number
    item_3: number
    item_4: number
    item_5: number
    assists: number
    hero_id: number
    last_hits: number
    net_worth: number
    player_id: number
    backpack_0: number
    backpack_1: number
    backpack_2: number
    is_radiant: boolean
    xp_per_min: number
    player_slot: number
    gold_per_min: number
    item_neutral: number
}

export type SeriesMatchDetail = {
    match_id: number
    duration: number
    radiant_win: boolean
    picks_bans: unknown
    players_data: PlayerData[]
    radiant_gold_adv: number[]
    radiant_xp_adv: number[]
    radiant_score: number
    dire_score: number
    radiant_captain: number | null
    dire_captain: number | null
}

export type SeriesDetail = {
    series_id: number
    start_time: string
    team_a: TeamInfo
    team_b: TeamInfo
    league: LeagueInfo
    team_a_score: number
    team_b_score: number
    matches: SeriesMatchDetail[]
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
    const url = buildUrl("/series", params)
    return fetchApi(url, signal, "Failed to fetch series")
}

export async function getSeriesById(id: number, signal: AbortSignal): Promise<SeriesDetail> {
    const url = new URL(`${API_BASE_URL}/series/${id}`)
    return fetchApi(url, signal, "Failed to fetch series details")
}

export async function searchLeagues(query: string, signal: AbortSignal): Promise<LeagueSearchResult[]> {
    const url = buildUrl("/filtersmetadata/leagues", { q: query })
    return fetchApi(url, signal, "Failed to search leagues")
}

export async function searchTeams(query: string, signal: AbortSignal): Promise<TeamSearchResult[]> {
    const url = buildUrl("/filtersmetadata/teams", { q: query })
    return fetchApi(url, signal, "Failed to search teams")
}
