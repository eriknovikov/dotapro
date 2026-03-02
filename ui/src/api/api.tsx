import type {
    Filters,
    Params,
    TeamInfo,
    LeagueInfo,
    LeagueSearchResult,
    Series,
    PlayerData,
    SeriesMatchDetail,
    SeriesDetail,
    Pagination,
    GetSeriesResponse,
    TeamSearchResult,
} from "@/types"

// Re-export types for backward compatibility
export type {
    Filters,
    Params,
    TeamInfo,
    LeagueInfo,
    LeagueSearchResult,
    Series,
    PlayerData,
    SeriesMatchDetail,
    SeriesDetail,
    Pagination,
    GetSeriesResponse,
    TeamSearchResult,
}

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8080"

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
