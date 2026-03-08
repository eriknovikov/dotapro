import type { Filters, MatchFilters } from "@/types"

const MATCHES_FILTERS_KEY = "dotapro_matches_filters"
const SERIES_FILTERS_KEY = "dotapro_series_filters"

/**
 * Save filters to localStorage for the specified type
 */
export function saveFilters(type: "matches" | "series", filters: Filters | MatchFilters): void {
    try {
        const key = type === "matches" ? MATCHES_FILTERS_KEY : SERIES_FILTERS_KEY
        localStorage.setItem(key, JSON.stringify(filters))
    } catch (error) {
        // Silently fail if localStorage is not available
        console.warn("Failed to save filters to localStorage:", error)
    }
}

/**
 * Load filters from localStorage for the specified type
 */
export function loadFilters(type: "matches" | "series"): Filters | MatchFilters | null {
    try {
        const key = type === "matches" ? MATCHES_FILTERS_KEY : SERIES_FILTERS_KEY
        const stored = localStorage.getItem(key)
        if (stored) {
            return JSON.parse(stored)
        }
    } catch (error) {
        console.warn("Failed to load filters from localStorage:", error)
    }
    return null
}

/**
 * Clear filters from localStorage for the specified type
 */
export function clearFilters(type: "matches" | "series"): void {
    try {
        const key = type === "matches" ? MATCHES_FILTERS_KEY : SERIES_FILTERS_KEY
        localStorage.removeItem(key)
    } catch (error) {
        console.warn("Failed to clear filters from localStorage:", error)
    }
}

/**
 * Hook to manage filter persistence
 */
export function useFilterPersistence(type: "matches" | "series") {
    return {
        save: (filters: Filters | MatchFilters) => saveFilters(type, filters),
        load: () => loadFilters(type),
        clear: () => clearFilters(type),
    }
}
