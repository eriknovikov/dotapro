import { useEffect } from "react"
import type { Filters } from "../api/api"

/**
 * Updates the URL query parameters without triggering a rerender.
 * Uses history.pushState to modify the URL directly.
 */
export function useUpdateUrlParams(params: Filters) {
    useEffect(() => {
        const url = new URL(window.location.href)

        // Clear existing params
        url.search = ""

        // Add params that are defined
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
        if (params.page !== undefined) {
            url.searchParams.set("page", String(params.page))
        }

        // Update URL without triggering navigation
        window.history.pushState({}, "", url.toString())
    }, [params])
}
