import type { Filters } from "../api/api"
import { useNavigate } from "@tanstack/react-router"
import { Funnel } from "lucide-react"
import { Button } from "./ui"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "./ui/index"
import { LeagueSelector } from "./LeagueSelector"
import { TeamSelector } from "./TeamSelector"

interface FiltersSidebarProps {
    filters: Filters
    isMobileOpen?: boolean
    onMobileClose?: () => void
}

export function FiltersSidebar({ filters, isMobileOpen, onMobileClose }: FiltersSidebarProps) {
    const navigate = useNavigate()

    const handleLeagueChange = (leagueId: number | undefined) => {
        navigate({
            to: ".",
            search: { ...filters, league: leagueId, c: undefined },
        })
    }

    const handleTeamChange = (teamId: number | undefined) => {
        navigate({
            to: ".",
            search: { ...filters, team: teamId, c: undefined },
        })
    }

    const handleSortChange = (sort: string) => {
        navigate({
            to: ".",
            search: { ...filters, sort: sort === "newest" ? undefined : sort, c: undefined },
        })
    }

    const handleLimitChange = (limit: string) => {
        navigate({
            to: ".",
            search: { ...filters, limit: limit === "20" ? undefined : Number(limit), c: undefined },
        })
    }

    const handleClear = () => {
        navigate({
            to: ".",
            search: { league: undefined, team: undefined, sort: undefined, limit: undefined, c: undefined },
        })
    }

    return (
        <aside
            className={`
                fixed top-16 left-0 z-40
                w-72 h-[calc(100vh-4rem)]
                bg-background border-r border-border
                shrink-0
                transform transition-transform duration-300 ease-in-out
                ${isMobileOpen ? "translate-x-0" : "-translate-x-full md:translate-x-0"}
            `}
        >
            <div className="h-full overflow-y-auto p-6">
                {/* Mobile close button */}
                <div className="flex items-center justify-between mb-6 md:hidden">
                    <h2 className="text-lg font-semibold text-foreground">Filters</h2>
                    <Button variant="ghost" size="icon" onClick={onMobileClose} aria-label="Close filters">
                        <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </Button>
                </div>

                {/* Desktop title */}
                <div className="hidden md:flex items-center gap-2 mb-6">
                    <Funnel className="h-5 w-5 text-foreground" />
                    <h2 className="text-lg font-semibold text-foreground">Filters</h2>
                </div>

                <div className="space-y-12">
                    {/* League Filter */}
                    <div>
                        <label htmlFor="league-select" className="block text-sm font-medium text-foreground mb-2">
                            League
                        </label>
                        <LeagueSelector
                            id="league-select"
                            onSelect={handleLeagueChange}
                            initialValue={filters.league}
                            aria-label="League filter"
                            inputClassName="text-sm"
                        />
                    </div>

                    {/* Team Filter */}
                    <div>
                        <label htmlFor="team-select" className="block text-sm font-medium text-foreground mb-2">
                            Team
                        </label>
                        <TeamSelector
                            id="team-select"
                            onSelect={handleTeamChange}
                            initialValue={filters.team}
                            aria-label="Team filter"
                            inputClassName="text-sm"
                        />
                    </div>

                    {/* Sort By Filter */}
                    <div>
                        <label htmlFor="sort-by" className="block text-sm font-medium text-foreground mb-2">
                            Sort by
                        </label>
                        <Select value={filters.sort || "newest"} onValueChange={handleSortChange}>
                            <SelectTrigger id="sort-by" aria-label="Sort by filter" className="text-sm">
                                <SelectValue placeholder="Default" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="newest" disabled={!filters.sort || filters.sort === "newest"}>
                                    Newest
                                </SelectItem>
                                <SelectItem value="oldest" disabled={filters.sort === "oldest"}>
                                    Oldest
                                </SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    {/* Limit Filter */}
                    <div>
                        <label htmlFor="limit" className="block text-sm font-medium text-foreground mb-2">
                            Series per page
                        </label>
                        <Select value={String(filters.limit || 20)} onValueChange={handleLimitChange}>
                            <SelectTrigger id="limit" aria-label="Results per page filter" className="text-sm">
                                <SelectValue placeholder="20" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="20" disabled={!filters.limit || filters.limit === 20}>
                                    20
                                </SelectItem>
                                <SelectItem value="40" disabled={filters.limit === 40}>
                                    40
                                </SelectItem>
                                <SelectItem value="60" disabled={filters.limit === 60}>
                                    60
                                </SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    {/* Action Buttons */}
                    <div className="flex justify-end">
                        <Button
                            onClick={handleClear}
                            variant="cool-outline"
                            size="sm"
                            className="text-sm"
                            aria-label="Reset defaults"
                        >
                            Reset defaults
                        </Button>
                    </div>
                </div>
            </div>
        </aside>
    )
}
