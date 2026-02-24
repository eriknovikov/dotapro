// ============================================================================
// Types
// ============================================================================

export interface Hero {
    id: number
    name: string
    displayName: string
}

export interface Item {
    id: number
    name: string
    displayName: string
}

export interface PopularLeague {
    id: number
    name: string
}

export interface PopularTeam {
    id: number
    name: string
    logo_url?: string
}

export interface PopularData {
    popular_leagues: PopularLeague[]
    popular_teams: PopularTeam[]
}

// ============================================================================
// Static Data Imports (bundled at build time - no network calls)
// ============================================================================

import heroes from '../assets/static_data/heroes.json'
import items from '../assets/static_data/items.json'
import popular from '../assets/static_data/popular.json'

// Type assertions for imported JSON
const HEROES_DATA = heroes as Record<string, Hero>
const ITEMS_DATA = items as Record<string, Item>
const POPULAR_DATA = popular as PopularData

// ============================================================================
// Data Access Functions (synchronous - data available immediately)
// ============================================================================

/**
 * Get popular data (returns raw object)
 */
export function getPopularData(): PopularData {
    return POPULAR_DATA
}

// ============================================================================
// Helper Functions
// ============================================================================

export function getHeroById(heroId: number): Hero | undefined {
    return HEROES_DATA[heroId.toString()]
}

export function getItemById(itemId: number): Item | undefined {
    return ITEMS_DATA[itemId.toString()]
}
