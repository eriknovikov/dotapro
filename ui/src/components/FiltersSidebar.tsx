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
        <>
            {/* Backdrop overlay for mobile/tablet */}
            {isMobileOpen && (
                <div
                    className="fixed inset-0 bg-black/60 z-40 lg:hidden"
                    onClick={onMobileClose}
                    aria-hidden="true"
                    role="presentation"
                />
            )}

            <aside
                className={`
                    fixed z-50
                    ${isMobileOpen ? "translate-y-0" : "translate-y-full"}
                    lg:fixed lg:top-16 lg:left-0 lg:translate-x-0 lg:translate-y-0
                    lg:w-72 lg:h-[calc(100vh-4rem)]
                    bottom-0 left-0 right-0
                    h-[90vh]
                    bg-background border-t lg:border-r border-border
                    rounded-t-2xl lg:rounded-none
                    shadow-2xl
                    transition-transform duration-300 ease-in-out
                `}
                role="complementary"
                aria-label="Filters"
            >
                <div className="h-full overflow-y-auto p-6">
                    {/* Mobile close button */}
                    <div className="flex items-center justify-between mb-6 lg:hidden">
                        <div className="flex items-center gap-2">
                            <Funnel className="h-5 w-5 text-foreground" />
                            <h2 className="text-lg font-semibold text-foreground">Filters</h2>
                        </div>
                        <Button
                            variant="ghost"
                            size="icon"
                            onClick={onMobileClose}
                            aria-label="Close filters"
                            className="h-12 w-12"
                        >
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
                    <div className="hidden lg:flex items-center gap-2 mb-6">
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
                                <SelectItem value="10" disabled={filters.limit === 10}>
                                    10
                                </SelectItem>
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
                    <div className="flex justify-end gap-1">
                        <Button
                            onClick={handleClear}
                            variant="outline"
                            size="sm"
                            className="text-sm"
                            aria-label="Reset defaults"
                        >
                            Reset defaults
                        </Button>
                        <Button
                            onClick={onMobileClose}
                            variant="primary"
                            size="sm"
                            className="text-sm lg:hidden"
                            aria-label="Close filters"
                        >
                            Done
                        </Button>
                    </div>
                </div>
            </div>
        </aside>
        </>
    )
}
