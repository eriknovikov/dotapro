export { cn } from "./utils"
export { formatRelativeTime } from "./utils"
export { formatDuration } from "./utils"
export { copyToClipboard } from "./utils"
export { getHeroImageUrl } from "./utils"
export { getItemImageUrl } from "./utils"
export { getHeroById } from "./dotautils"
export { getNeutralById } from "./dotautils"
export { getItemById } from "./dotautils"
export { getPopularData } from "./dotautils"
export type { Hero, Item, PopularLeague, PopularTeam, PopularData } from "./dotautils"

// Lazy data utilities for bundle optimization
export {
    loadHeroes,
    loadItems,
    loadNeutrals,
    loadPopular,
    loadItemsByName,
    getHeroByIdSync,
    getItemByIdSync,
    getNeutralByIdSync,
    getPopularDataSync,
} from "./lazyData"
