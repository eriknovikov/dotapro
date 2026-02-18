import type { Filters } from "../api/api"
import { useNavigate } from "@tanstack/react-router"
import { Button } from "./ui/index"
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

    const handleClear = () => {
        navigate({ to: ".", search: { c: undefined } })
    }

    return (
        <aside
            className={`
                fixed top-16 left-0 z-40
                w-64 min-h-[calc(100vh-4rem)]
                bg-background border-r border-border p-6
                shrink-0
                transform transition-transform duration-300 ease-in-out
                ${isMobileOpen ? "translate-x-0" : "-translate-x-full md:translate-x-0"}
            `}
        >
            {/* Mobile close button */}
            <div className="flex items-center justify-between mb-6 md:hidden">
                <h2 className="text-lg font-semibold text-foreground">Filters</h2>
                <Button variant="ghost" size="icon" onClick={onMobileClose} aria-label="Close filters">
                    <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </Button>
            </div>

            {/* Desktop title */}
            <h2 className="hidden md:block text-lg font-semibold text-foreground mb-6">Filters</h2>

            <div className="space-y-6">
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
                    />
                </div>

                {/* Sort By Filter */}
                <div>
                    <label htmlFor="sort-by" className="block text-sm font-medium text-foreground mb-2">
                        Sort By
                    </label>
                    <Select value={filters.sort || "newest"} onValueChange={handleSortChange}>
                        <SelectTrigger id="sort-by" aria-label="Sort by filter">
                            <SelectValue placeholder="Default" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="newest" disabled={!filters.sort || filters.sort === "newest"}>Newest</SelectItem>
                            <SelectItem value="oldest" disabled={filters.sort === "oldest"}>Oldest</SelectItem>
                        </SelectContent>
                    </Select>
                </div>

                {/* Action Buttons */}
                <div className="space-y-3">
                    <Button onClick={handleClear} variant="secondary" className="w-full" aria-label="Clear filters">
                        Clear Filters
                    </Button>
                </div>
            </div>
        </aside>
    )
}
