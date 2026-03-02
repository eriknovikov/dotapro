import { lazy, Suspense } from "react"
import { Crown, Trophy } from "lucide-react"
import type { SeriesMatchDetail, PlayerData, TeamInfo } from "@/api/api"
import { getHeroImageUrl, getItemImageUrl, cn } from "@/lib"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui"
import { Tooltip } from "@/components/ui/tooltip"

// Lazy load tooltip components to reduce initial bundle size
const HeroTooltipContent = lazy(() => import("./HeroTooltipContent").then(m => ({ default: m.HeroTooltipContent })))
const ItemTooltipContent = lazy(() => import("./ItemTooltipContent").then(m => ({ default: m.ItemTooltipContent })))

interface PlayersTableProps {
    match: SeriesMatchDetail
    radiantTeam: TeamInfo
    direTeam: TeamInfo
}

export function PlayersTable({ match, radiantTeam, direTeam }: PlayersTableProps) {
    const { players_data, radiant_captain, dire_captain, radiant_win, radiant_score, dire_score } = match

    const radiantPlayers = players_data.filter(p => p.is_radiant)
    const direPlayers = players_data.filter(p => !p.is_radiant)

    return (
        <div className="flex flex-col gap-3" role="region" aria-label="Match players statistics">
            <TeamPlayersTable
                players={radiantPlayers}
                team="Radiant"
                captain={radiant_captain}
                teamColor="secondary"
                teamInfo={radiantTeam}
                score={radiant_score}
                isWinner={radiant_win}
            />
            <TeamPlayersTable
                players={direPlayers}
                team="Dire"
                captain={dire_captain}
                teamColor="primary"
                teamInfo={direTeam}
                score={dire_score}
                isWinner={!radiant_win}
            />
        </div>
    )
}

interface TeamPlayersTableProps {
    players: PlayerData[]
    team: "Radiant" | "Dire"
    captain: number | null
    teamColor: "primary" | "secondary"
    teamInfo: TeamInfo
    score: number
    isWinner: boolean
}

const formatNetWorth = (netWorth: number) => {
    if (netWorth >= 1000) {
        return `${(netWorth / 1000).toFixed(1)}K`
    }
    return `${netWorth}K`
}

function TeamPlayersTable({ players, captain, teamInfo, score, isWinner }: TeamPlayersTableProps) {
    return (
        <div
            className="bg-background-accent rounded-lg p-4"
            role="region"
            aria-label={`${teamInfo.name} team statistics`}
        >
            <div className="flex items-center gap-2 mb-3">
                <span className="text-base sm:text-sm font-semibold">{teamInfo.name}</span>
                <span className="text-foreground-muted text-base sm:text-sm" aria-label={`Score: ${score}`}>
                    {score}
                </span>
                {isWinner && (
                    <div
                        className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs text-white font-semibold"
                        aria-label="Winner"
                    >
                        <Trophy className="size-5 text-secondary" aria-hidden="true" />
                    </div>
                )}
            </div>

            <div className="overflow-x-auto">
                <Table className="table-fixed min-w-225">
                    <TableHeader>
                        <TableRow className="border-border hover:bg-transparent">
                            <TableHead className="w-20 px-3 text-xs text-foreground-muted uppercase tracking-wider">
                                Hero
                            </TableHead>
                            <TableHead className="w-48 px-3 text-xs text-foreground-muted uppercase tracking-wider">
                                Player
                            </TableHead>
                            <TableHead className="w-80 px-3 text-xs text-foreground-muted uppercase tracking-wider">
                                Items
                            </TableHead>
                            <TableHead className="w-24 px-3 text-right text-xs text-secondary-500 uppercase tracking-wider">
                                NET
                            </TableHead>
                            <TableHead className="w-32 px-3 text-center text-xs text-foreground-muted uppercase tracking-wider">
                                <span className="text-success-200">K</span>
                                <span className="text-foreground-muted mx-0.5">/</span>
                                <span className="text-error-300">D</span>
                                <span className="text-foreground-muted mx-0.5">/</span>
                                <span className="text-foreground">A</span>
                            </TableHead>
                            <TableHead className="w-20 px-3 text-center text-xs text-foreground-muted uppercase tracking-wider">
                                LH
                            </TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {players.map(player => {
                            const isCaptain = captain === player.player_id
                            return (
                                <TableRow
                                    key={player.player_id}
                                    className="border-border hover:bg-background/50 border-b-2 "
                                >
                                    {/* Hero */}
                                    <TableCell className="px-3">
                                        <Suspense fallback={<img src={getHeroImageUrl(player.hero_id)} className="h-7 w-auto sm:h-8" />}>
                                            <Tooltip
                                                content={<HeroTooltipContent heroId={player.hero_id} />}
                                                side="right"
                                                contentClassName="p-0 border-none bg-transparent"
                                            >
                                                <img
                                                    src={getHeroImageUrl(player.hero_id)}
                                                    className="h-7 w-auto sm:h-8 cursor-pointer"
                                                />
                                            </Tooltip>
                                        </Suspense>
                                    </TableCell>
                                    {/* Player */}
                                    <TableCell className="px-3 py-0">
                                        <div className="flex items-center gap-1">
                                            <span
                                                className={cn("text-md", player.name ?? "")}
                                                title={player.name ?? "unknown"}
                                            >
                                                {player.name ?? "unknown"}
                                            </span>
                                            {isCaptain && (
                                                <span className="inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded text-xs bg-secondary/20 shrink-0 text-secondary">
                                                    <Crown className="w-3 h-3" />
                                                </span>
                                            )}
                                        </div>
                                    </TableCell>
                                    {/* Items */}
                                    <TableCell className="px-3 py-1">
                                        <PlayerItems player={player} />
                                    </TableCell>
                                    {/* NET */}
                                    <TableCell className="px-3 py-0 text-right">
                                        <span className="text-sm text-secondary-500">
                                            {formatNetWorth(player.net_worth)}
                                        </span>
                                    </TableCell>
                                    {/* KDA */}
                                    <TableCell className="px-3 py-0 text-center">
                                        <span className="text-sm">
                                            <span className="text-success-200">{player.kills}</span>
                                            <span className="text-foreground-muted mx-0.5">/</span>
                                            <span className="text-error-300">{player.deaths}</span>
                                            <span className="text-foreground-muted mx-0.5">/</span>
                                            <span className="text-foreground">{player.assists}</span>
                                        </span>
                                    </TableCell>
                                    {/* LH */}
                                    <TableCell className="px-3 py-0 text-center">
                                        <span className="text-sm text-foreground-muted">{player.last_hits}</span>
                                    </TableCell>
                                </TableRow>
                            )
                        })}
                    </TableBody>
                </Table>
            </div>
        </div>
    )
}

// ============================================================================
// Helper Components
// ============================================================================

interface PlayerItemsProps {
    player: PlayerData
}

function PlayerItems({ player }: PlayerItemsProps) {
    return (
        <div className="flex gap-2 items-center">
            {/* Inventory items */}
            <div className="flex flex-col gap-1">
                <div className="flex gap-1">
                    {getInventorySlots(player).map(slot => (
                        <ItemSlot
                            key={slot.key}
                            itemId={slot.itemId}
                            itemUrl={slot.itemUrl}
                            widthClass={slot.widthClass}
                            heightClass={slot.heightClass}
                            emptyBgClass={slot.emptyBgClass}
                        />
                    ))}
                </div>
                {/* Backpack section */}
                <div className="flex gap-1">
                    {getBackpackSlots(player).map(slot => (
                        <ItemSlot
                            key={slot.key}
                            itemId={slot.itemId}
                            itemUrl={slot.itemUrl}
                            widthClass={slot.widthClass}
                            heightClass={slot.heightClass}
                            emptyBgClass={slot.emptyBgClass}
                        />
                    ))}
                </div>
            </div>
            {/* Neutral item */}
            <NeutralItemSlot itemId={player.item_neutral} />
        </div>
    )
}

interface ItemSlotProps {
    key: string
    itemId: number
    itemUrl: string
    widthClass: string
    heightClass: string
    emptyBgClass: string
}

function ItemSlot({ itemId, itemUrl, widthClass, heightClass, emptyBgClass }: ItemSlotProps) {
    const isEmpty = itemId === 0 || itemUrl === ""

    return (
        <Suspense fallback={<div className={`${widthClass} ${heightClass} rounded ${isEmpty ? emptyBgClass : ""}`}>
            {!isEmpty && <img src={itemUrl} className="rounded w-full h-full object-cover" />}
        </div>}>
            <Tooltip
                content={<ItemTooltipContent itemId={itemId} />}
                side="left"
                contentClassName="p-0 border-none bg-transparent"
            >
                <div className={`${widthClass} ${heightClass} rounded ${isEmpty ? emptyBgClass : ""}`}>
                    {!isEmpty && <img src={itemUrl} className="rounded w-full h-full object-cover" />}
                </div>
            </Tooltip>
        </Suspense>
    )
}

function NeutralItemSlot({ itemId }: { itemId: number }) {
    return (
        <Suspense fallback={<div className="w-10 h-7 min-w-10 sm:w-11 sm:h-8 sm:min-w-11 rounded bg-black border border-border shrink-0">
            {itemId > 0 && <img src={getItemImageUrl(itemId)} className="w-full h-full object-cover rounded" />}
        </div>}>
            <Tooltip
                content={<ItemTooltipContent itemId={itemId} />}
                side="right"
                contentClassName="p-0 border-none bg-transparent"
            >
                <div className="w-10 h-7 min-w-10 sm:w-11 sm:h-8 sm:min-w-11 rounded bg-black border border-border shrink-0">
                    {itemId > 0 && <img src={getItemImageUrl(itemId)} className="w-full h-full object-cover rounded" />}
                </div>
            </Tooltip>
        </Suspense>
    )
}

// ============================================================================
// Helper Functions
// ============================================================================

function emptyLast<T extends { isEmpty: boolean }>(a: T, b: T): number {
    if (a.isEmpty && !b.isEmpty) return 1
    if (!a.isEmpty && b.isEmpty) return -1
    return 0
}

function getInventorySlots(player: PlayerData): ItemSlotProps[] {
    return [0, 1, 2, 3, 4, 5]
        .map(slot => ({
            key: `item-${slot}`,
            itemId: player[`item_${slot}` as keyof PlayerData] as number,
            itemUrl: getItemImageUrl(player[`item_${slot}` as keyof PlayerData] as number),
            widthClass: "w-8 sm:w-9",
            heightClass: "h-6 sm:h-7",
            emptyBgClass: "bg-black border border-border",
        }))
        .map(slot => ({ ...slot, isEmpty: slot.itemId === 0 }))
        .sort(emptyLast)
}

function getBackpackSlots(player: PlayerData): ItemSlotProps[] {
    return [0, 1, 2]
        .map(slot => ({
            key: `backpack-${slot}`,
            itemId: player[`backpack_${slot}` as keyof PlayerData] as number,
            itemUrl: getItemImageUrl(player[`backpack_${slot}` as keyof PlayerData] as number),
            widthClass: "w-7 sm:w-8",
            heightClass: "h-5 sm:h-6",
            emptyBgClass: "bg-background border border-border",
        }))
        .map(slot => ({ ...slot, isEmpty: slot.itemId === 0 }))
        .sort(emptyLast)
}
