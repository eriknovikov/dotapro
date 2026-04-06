// ============================================================================
// Types
// ============================================================================

export interface Hero {
    id: number
    name: string
    displayName: string
    lore: string
}

export interface ItemAttribute {
    key: string
    display?: string
    value: string | number
}

export interface ItemAbility {
    type: string
    title: string
    description: string
}

export interface Item {
    id: number
    name: string
    displayName: string
    cost?: number
    mc?: number | boolean
    hc?: number | boolean
    cd?: number | boolean
    abilities?: ItemAbility[]
    attrib?: ItemAttribute[]
    components?: string[] | null
    tier?: number | null
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

export interface PopularPlayer {
    id: number
    name: string
}

export interface PopularHero {
    id: number
    name: string
}

export interface PopularData {
    popular_leagues: PopularLeague[]
    popular_teams: PopularTeam[]
    popular_players: PopularPlayer[]
    popular_heroes: PopularHero[]
}

// ============================================================================
// Static Data Imports (bundled at build time - no network calls)
// ============================================================================

import heroes from "../assets/static_data/heroes.json"
import items from "../assets/static_data/items.json"
import neutrals from "../assets/static_data/neutrals.json"
import popular from "../assets/static_data/popular.json"

// Type assertions for imported JSON
const HEROES_DATA = heroes as Record<string, Hero>
const ITEMS_DATA = items as Record<string, Item>
const NEUTRALS_DATA = neutrals as Record<string, Item>
const POPULAR_DATA = popular as PopularData

// Create a reverse lookup map for components (name -> Item)
const ITEMS_BY_NAME: Record<string, Item> = {}
Object.values(ITEMS_DATA).forEach(item => {
    ITEMS_BY_NAME[item.name] = item
})
Object.values(NEUTRALS_DATA).forEach(item => {
    ITEMS_BY_NAME[item.name] = item
})

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

export function getNeutralById(itemId: number): Item | undefined {
    return NEUTRALS_DATA[itemId.toString()]
}

export function getItemByName(name: string): Item | undefined {
    return ITEMS_BY_NAME[name]
}

export function getHeroImageUrl(heroId: number): string {
    return `https://cdn.dota2.com/apps/dota2/images/heroes/${getHeroById(heroId)?.name?.replace("npc_dota_hero_", "")}_lg.png`
}

export function formatRelativeTime(dateString: string): string {
    const date = new Date(dateString)
    const now = new Date()
    const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)

    if (diffInSeconds < 60) {
        return `${diffInSeconds} seconds`
    } else if (diffInSeconds < 3600) {
        const minutes = Math.floor(diffInSeconds / 60)
        return `${minutes} minute${minutes > 1 ? "s" : ""}`
    } else if (diffInSeconds < 86400) {
        const hours = Math.floor(diffInSeconds / 3600)
        return `${hours} hour${hours > 1 ? "s" : ""}`
    } else {
        const days = Math.floor(diffInSeconds / 86400)
        return `${days} day${days > 1 ? "s" : ""}`
    }
}
