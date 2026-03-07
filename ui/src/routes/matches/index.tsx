import type { MatchFilters } from "@/api/api"
import { FiltersSidebar } from "@/components/FiltersSidebar"
import { MatchList } from "@/components/matches/MatchList"
import { SEO } from "@/components/SEO"
import { Button } from "@/components/ui"
import { PAGINATION_LIMITS } from "@/constants"
import { createFileRoute, useNavigate, useSearch } from "@tanstack/react-router"
import { Funnel } from "lucide-react"
import { useEffect, useRef, useState } from "react"

export const Route = createFileRoute("/matches/")({
    component: Matches,
    validateSearch: (search: Record<string, unknown>): MatchFilters => {
        const limit = search.limit !== undefined ? Number(search.limit) : undefined
        return {
            league: search.league !== undefined ? Number(search.league) : undefined,
            team: search.team !== undefined ? Number(search.team) : undefined,
            player: search.player !== undefined ? Number(search.player) : undefined,
            hero: search.hero !== undefined ? Number(search.hero) : undefined,
            sort: typeof search.sort === "string" ? search.sort : undefined,
            limit: limit !== undefined && PAGINATION_LIMITS.includes(limit as typeof PAGINATION_LIMITS[number]) ? limit : undefined,
            c: search.c !== undefined ? Number(search.c) : undefined,
        }
    },
})

function Matches() {
    const search = useSearch({ strict: false }) as MatchFilters
    const navigate = useNavigate()
    const [isMobileFiltersOpen, setIsMobileFiltersOpen] = useState(false)
    const scrollPositionRef = useRef(0)

    // Lock body scroll when filter modal is open
    useEffect(() => {
        if (isMobileFiltersOpen) {
            // Save current scroll position
            scrollPositionRef.current = window.scrollY
            // Lock scroll
            document.body.style.overflow = "hidden"
        } else {
            // Unlock scroll
            document.body.style.overflow = ""
            // Restore scroll position
            window.scrollTo(0, scrollPositionRef.current)
        }
        return () => {
            document.body.style.overflow = ""
        }
    }, [isMobileFiltersOpen])

    // Close filter modal when navbar hamburger is clicked
    useEffect(() => {
        const handleNavbarMenuOpen = () => {
            if (isMobileFiltersOpen) {
                setIsMobileFiltersOpen(false)
            }
        }

        window.addEventListener("navbar-menu-open", handleNavbarMenuOpen)
        return () => {
            window.removeEventListener("navbar-menu-open", handleNavbarMenuOpen)
        }
    }, [isMobileFiltersOpen])

    return (
        <>
            <SEO
                title="Matches"
                description="Browse professional Dota 2 matches with advanced filtering and statistics."
            />
            {/* Mobile/Tablet Filters Toggle Button */}
            <div className="fixed right-4 bottom-4 z-50 lg:hidden">
                <Button
                    onClick={() => setIsMobileFiltersOpen(true)}
                    size="lg"
                    className="shadow-lg"
                    aria-label="Open filters"
                >
                    <Funnel className="mr-2 h-5 w-5" />
                    Filters
                </Button>
            </div>

            {/* Sidebar - Filters (absolute positioned to not affect footer width) */}
            <FiltersSidebar
                filters={search}
                isMobileOpen={isMobileFiltersOpen}
                onMobileClose={() => setIsMobileFiltersOpen(false)}
            />

            {/* Main Content - Results */}
            <main
                className={`px-2 py-6 sm:px-0 lg:ml-72 ${isMobileFiltersOpen ? "overflow-hidden" : "overflow-y-auto"}`}
            >
                <div className="mx-auto max-w-7xl">
                    <MatchList filters={search} />
                </div>
            </main>
        </>
    )
}
