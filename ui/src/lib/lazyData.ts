/**
 * Lazy-loaded data utilities for optimizing bundle size
 * Large JSON files (heroes, items, neutrals) are loaded on-demand
 */

import type { Hero, Item, PopularData } from "./dotautils"

// ============================================================================
// Lazy Data Loading
// ============================================================================

let heroesData: Record<string, Hero> | null = null
let itemsData: Record<string, Item> | null = null
let neutralsData: Record<string, Item> | null = null
let popularData: PopularData | null = null
let itemsByName: Record<string, Item> | null = null

/**
 * Lazy load heroes data
 */
export async function loadHeroes(): Promise<Record<string, Hero>> {
    if (heroesData) {
        return heroesData
    }

    const heroes = await import("../assets/static_data/heroes.json")
    heroesData = heroes as Record<string, Hero>
    return heroesData
}

/**
 * Lazy load items data
 */
export async function loadItems(): Promise<Record<string, Item>> {
    if (itemsData) {
        return itemsData
    }

    const items = await import("../assets/static_data/items.json")
    itemsData = items as Record<string, Item>
    return itemsData
}

/**
 * Lazy load neutrals data
 */
export async function loadNeutrals(): Promise<Record<string, Item>> {
    if (neutralsData) {
        return neutralsData
    }

    const neutrals = await import("../assets/static_data/neutrals.json")
    neutralsData = neutrals as Record<string, Item>
    return neutralsData
}

/**
 * Lazy load popular data
 */
export async function loadPopular(): Promise<PopularData> {
    if (popularData) {
        return popularData
    }

    const popular = await import("../assets/static_data/popular.json")
    popularData = popular as PopularData
    return popularData
}

/**
 * Get items by name (lazy loaded)
 */
export async function loadItemsByName(): Promise<Record<string, Item>> {
    if (itemsByName) {
        return itemsByName
    }

    const [items, neutrals] = await Promise.all([loadItems(), loadNeutrals()])

    itemsByName = {}
    Object.values(items).forEach(item => {
        itemsByName[item.name] = item
    })
    Object.values(neutrals).forEach(item => {
        itemsByName[item.name] = item
    })

    return itemsByName
}

// ============================================================================
// Synchronous Accessors (for backward compatibility)
// ============================================================================

/**
 * Get hero by ID (synchronous - loads data if not already loaded)
 * Note: This will cause the heroes data to be bundled synchronously
 * For optimal bundle size, use loadHeroes() and await the result
 */
export function getHeroByIdSync(heroId: number): Hero | undefined {
    if (!heroesData) {
        // Fallback to synchronous import for backward compatibility
        // This will bundle the data synchronously
        const heroes = require("../assets/static_data/heroes.json")
        heroesData = heroes as Record<string, Hero>
    }
    return heroesData[heroId.toString()]
}

/**
 * Get item by ID (synchronous - loads data if not already loaded)
 * Note: This will cause the items data to be bundled synchronously
 * For optimal bundle size, use loadItems() and await the result
 */
export function getItemByIdSync(itemId: number): Item | undefined {
    if (!itemsData) {
        const items = require("../assets/static_data/items.json")
        itemsData = items as Record<string, Item>
    }
    return itemsData[itemId.toString()]
}

/**
 * Get neutral by ID (synchronous - loads data if not already loaded)
 * Note: This will cause the neutrals data to be bundled synchronously
 * For optimal bundle size, use loadNeutrals() and await the result
 */
export function getNeutralByIdSync(itemId: number): Item | undefined {
    if (!neutralsData) {
        const neutrals = require("../assets/static_data/neutrals.json")
        neutralsData = neutrals as Record<string, Item>
    }
    return neutralsData[itemId.toString()]
}

/**
 * Get popular data (synchronous - loads data if not already loaded)
 * Note: This will cause the popular data to be bundled synchronously
 * For optimal bundle size, use loadPopular() and await the result
 */
export function getPopularDataSync(): PopularData {
    if (!popularData) {
        const popular = require("../assets/static_data/popular.json")
        popularData = popular as PopularData
    }
    return popularData
}
