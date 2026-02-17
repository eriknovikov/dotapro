import { useRef, useEffect, useState, type KeyboardEvent } from "react"
import type { Filters } from "../api/api"
import { useNavigate } from "@tanstack/react-router"
import { Button, Input } from "./ui/index"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "./ui/index"

interface FiltersSidebarProps {
    initialFilters?: Filters
    isLoading: boolean
    isMobileOpen?: boolean
    onMobileClose?: () => void
}

export function FiltersSidebar({ initialFilters, isLoading, isMobileOpen, onMobileClose }: FiltersSidebarProps) {
    const navigate = useNavigate()

    const leagueInputRef = useRef<HTMLInputElement>(null)
    const teamInputRef = useRef<HTMLInputElement>(null)
    const sortSelectRef = useRef<HTMLSelectElement>(null)

    useEffect(() => {
        if (initialFilters) {
            if (leagueInputRef.current && initialFilters.league !== undefined) {
                leagueInputRef.current.value = String(initialFilters.league)
            }
            if (teamInputRef.current && initialFilters.team !== undefined) {
                teamInputRef.current.value = String(initialFilters.team)
            }
            if (sortSelectRef.current && initialFilters.sort !== undefined) {
                sortSelectRef.current.value = initialFilters.sort
            }
        }
    }, [initialFilters])

    const handleSearchClick = () => {
        const filters: Filters = {}

        if (leagueInputRef.current?.value) {
            filters.league = Number(leagueInputRef.current.value)
        }
        if (teamInputRef.current?.value) {
            filters.team = Number(teamInputRef.current.value)
        }
        if (sortSelectRef.current?.value) {
            filters.sort = sortSelectRef.current.value
        }
        // Reset cursor when filters change
        filters.c = undefined

        navigate({
            to: ".",
            search: filters,
        })
    }

    const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") {
            e.preventDefault()
            handleSearchClick()
        }
    }

    const handleClear = () => {
        if (leagueInputRef.current) leagueInputRef.current.value = ""
        if (teamInputRef.current) teamInputRef.current.value = ""
        if (sortSelectRef.current) sortSelectRef.current.value = ""
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
                <Button
                    variant="ghost"
                    size="icon"
                    onClick={onMobileClose}
                    aria-label="Close filters"
                >
                    <svg
                        className="h-6 w-6"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
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
            <h2 className="hidden md:block text-lg font-semibold text-foreground mb-6">Filters</h2>

            <div className="space-y-6">
                {/* League ID Filter */}
                <div>
                    <label htmlFor="league-id" className="block text-sm font-medium text-foreground mb-2">
                        League ID
                    </label>
                    <Input
                        ref={leagueInputRef}
                        id="league-id"
                        type="number"
                        placeholder="Enter league ID"
                        onKeyDown={handleKeyDown}
                        aria-label="League ID filter"
                    />
                </div>

                {/* Team ID Filter */}
                <div>
                    <label htmlFor="team-id" className="block text-sm font-medium text-foreground mb-2">
                        Team ID
                    </label>
                    <Input
                        ref={teamInputRef}
                        id="team-id"
                        type="number"
                        placeholder="Enter team ID"
                        onKeyDown={handleKeyDown}
                        aria-label="Team ID filter"
                    />
                </div>

                {/* Sort By Filter */}
                <div>
                    <label htmlFor="sort-by" className="block text-sm font-medium text-foreground mb-2">
                        Sort By
                    </label>
                    <Select
                        onValueChange={(value: string) => {
                            if (sortSelectRef.current) {
                                sortSelectRef.current.value = value
                            }
                        }}
                    >
                        <SelectTrigger id="sort-by" aria-label="Sort by filter">
                            <SelectValue placeholder="Default" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="newest">Newest</SelectItem>
                            <SelectItem value="oldest">Oldest</SelectItem>
                        </SelectContent>
                    </Select>
                </div>

                {/* Action Buttons */}
                <div className="space-y-3">
                    <Button
                        onClick={handleSearchClick}
                        disabled={isLoading}
                        className="w-full"
                        aria-label="Search series"
                    >
                        {isLoading ? "Searching..." : "Search"}
                    </Button>

                    <Button
                        onClick={handleClear}
                        disabled={isLoading}
                        variant="secondary"
                        className="w-full"
                        aria-label="Clear filters"
                    >
                        Clear Filters
                    </Button>
                </div>
            </div>
        </aside>
    )
}
