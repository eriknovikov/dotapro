import { Crown, Backpack } from "lucide-react"
import type { SeriesMatchDetail, PlayerData } from "@/api/api"
import { getHeroImageUrl, getItemImageUrl, cn } from "@/lib"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui"

interface PlayersTableProps {
    match: SeriesMatchDetail
}

export function PlayersTable({ match }: PlayersTableProps) {
    const { players_data, radiant_captain, dire_captain } = match

    const radiantPlayers = players_data.filter(p => p.is_radiant)
    const direPlayers = players_data.filter(p => !p.is_radiant)

    return (
        <div className="flex flex-col gap-3">
            <TeamPlayersTable players={radiantPlayers} team="Radiant" captain={radiant_captain} teamColor="secondary" />
            <TeamPlayersTable players={direPlayers} team="Dire" captain={dire_captain} teamColor="primary" />
        </div>
    )
}

interface TeamPlayersTableProps {
    players: PlayerData[]
    team: "Radiant" | "Dire"
    captain: number | null
    teamColor: "primary" | "secondary"
}

function TeamPlayersTable({ players, team, captain, teamColor }: TeamPlayersTableProps) {
    const teamColorClass = teamColor === "secondary" ? "text-secondary" : "text-primary"

    const formatNetWorth = (netWorth: number) => {
        if (netWorth >= 1000) {
            return `${(netWorth / 1000).toFixed(1)}K`
        }
        return `${netWorth}K`
    }

    return (
        <div className="bg-background-accent rounded-lg p-4">
            <div className={cn("text-sm font-semibold mb-3 flex items-center gap-2", teamColorClass)}>
                <span>{team}</span>
                <span className="text-foreground-muted">({players.length})</span>
            </div>

            <Table>
                <TableHeader>
                    <TableRow className="border-border hover:bg-transparent">
                        <TableHead className="w-16 px-2 text-xs text-foreground-muted uppercase tracking-wider">
                            Hero
                        </TableHead>
                        <TableHead className="w-36 px-2 text-xs text-foreground-muted uppercase tracking-wider">
                            Player
                        </TableHead>
                        <TableHead className="w-56 px-2 text-xs text-foreground-muted uppercase tracking-wider">
                            Items
                        </TableHead>
                        <TableHead className="w-20 px-2 text-right text-xs text-secondary-500 uppercase tracking-wider">
                            NET
                        </TableHead>
                        <TableHead className="w-20 px-2 text-center text-xs text-foreground-muted uppercase tracking-wider">
                            <span className="text-success-200">K</span>
                            <span className="text-foreground-muted mx-0.5">/</span>
                            <span className="text-error-300">D</span>
                            <span className="text-foreground-muted mx-0.5">/</span>
                            <span className="text-foreground">A</span>
                        </TableHead>
                        <TableHead className="w-16 px-2 text-center text-xs text-foreground-muted uppercase tracking-wider">
                            LH
                        </TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {players.map(player => {
                        const isCaptain = captain === player.player_id
                        return (
                            <TableRow key={player.player_id} className="border-border hover:bg-background/50">
                                {/* Hero */}
                                <TableCell className="w-16 px-2 py-0">
                                    <img src={getHeroImageUrl(player.hero_id)} className="h-10" />
                                </TableCell>
                                {/* Player */}
                                <TableCell className="w-36 px-2 py-0">
                                    <div className="flex items-center gap-1 min-w-0">
                                        <span className="text-md truncate" title={player.name ?? "unknown"}>
                                            {player.name ?? "unknown"}
                                        </span>
                                        {isCaptain && (
                                            <span className="inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded text-xs bg-secondary/20 text-secondary shrink-0">
                                                <Crown className="w-3 h-3" />
                                            </span>
                                        )}
                                    </div>
                                </TableCell>
                                {/* Items */}
                                <TableCell className="w-56 px-2 py-0">
                                    <div className="flex flex-col gap-1">
                                        {/* Inventory items */}
                                        <div className="flex gap-1">
                                            {[0, 1, 2, 3, 4, 5].map(slot => {
                                                const itemId = player[`item_${slot}` as keyof PlayerData] as number
                                                const itemUrl = getItemImageUrl(itemId)
                                                return (
                                                    <div
                                                        key={slot}
                                                        className={cn(
                                                            "w-6 h-6 rounded",
                                                            itemUrl ? "" : "bg-background border border-border",
                                                        )}
                                                    >
                                                        {itemUrl && <img src={itemUrl} className="w-6 h-6 rounded" />}
                                                    </div>
                                                )
                                            })}
                                            {/* Neutral item */}
                                            {player.item_neutral > 0 && (
                                                <img
                                                    src={getItemImageUrl(player.item_neutral)}
                                                    className="h-4 w-auto rounded"
                                                />
                                            )}
                                        </div>
                                        {/* Backpack section */}
                                        <div className="flex items-center gap-1">
                                            <Backpack className="w-3 h-3 text-foreground-muted" />
                                            {[0, 1, 2].map(slot => {
                                                const itemId = player[`backpack_${slot}` as keyof PlayerData] as number
                                                const itemUrl = getItemImageUrl(itemId)
                                                return (
                                                    <div
                                                        key={`backpack-${slot}`}
                                                        className={cn(
                                                            "w-5 h-5 rounded",
                                                            itemUrl ? "" : "bg-background border border-border",
                                                        )}
                                                    >
                                                        {itemUrl && <img src={itemUrl} className="w-5 h-5 rounded" />}
                                                    </div>
                                                )
                                            })}
                                        </div>
                                    </div>
                                </TableCell>
                                {/* NET */}
                                <TableCell className="w-20 px-2 py-0 text-right">
                                    <span className="text-sm text-secondary-500">
                                        {formatNetWorth(player.net_worth)}
                                    </span>
                                </TableCell>
                                {/* KDA */}
                                <TableCell className="w-20 px-2 py-0 text-center">
                                    <span className="text-sm">
                                        <span className="text-success-200">{player.kills}</span>
                                        <span className="text-foreground-muted mx-0.5">/</span>
                                        <span className="text-error-300">{player.deaths}</span>
                                        <span className="text-foreground-muted mx-0.5">/</span>
                                        <span className="text-foreground">{player.assists}</span>
                                    </span>
                                </TableCell>
                                {/* LH */}
                                <TableCell className="w-16 px-2 py-0 text-center">
                                    <span className="text-sm text-foreground-muted">{player.last_hits}</span>
                                </TableCell>
                            </TableRow>
                        )
                    })}
                </TableBody>
            </Table>
        </div>
    )
}
