import { getSeries, type Filters, type GetSeriesResponse } from "@/api"
import { Button, SEO, SeriesFilters, SeriesList } from "@/components"
import { PAGINATION_LIMITS } from "@/constants"
import { useQuery } from "@tanstack/react-query"
import { createFileRoute, useSearch } from "@tanstack/react-router"
import { Funnel } from "lucide-react"
import { useEffect, useRef, useState } from "react"

export const Route = createFileRoute("/series/")({
    component: Series,
    validateSearch: (search: Record<string, unknown>): Filters => {
        const limit = search.limit !== undefined ? Number(search.limit) : undefined
        return {
            league: search.league !== undefined ? Number(search.league) : undefined,
            team: search.team !== undefined ? Number(search.team) : undefined,
            sort: typeof search.sort === "string" ? search.sort : undefined,
            limit:
                limit !== undefined && PAGINATION_LIMITS.includes(limit as (typeof PAGINATION_LIMITS)[number])
                    ? limit
                    : undefined,
            c: search.c !== undefined ? Number(search.c) : undefined,
        }
    },
})

function Series() {
    const search = useSearch({ strict: false })
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

    const { data, isLoading, error } = useQuery({
        queryKey: ["series", search],
        queryFn: async ({ signal }): Promise<GetSeriesResponse> => {
            return getSeries(search, signal)
        },
        retry: 2,
    })

    return (
        <>
            <SEO
                title="Series"
                description="Browse professional Dota 2 series and matches. Filter by league, team, and more."
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
            <SeriesFilters
                filters={search}
                isMobileOpen={isMobileFiltersOpen}
                onMobileClose={() => setIsMobileFiltersOpen(false)}
            />

            {/* Main Content - Results */}
            <main
                className={`px-2 py-6 sm:px-0 lg:ml-72 ${isMobileFiltersOpen ? "overflow-hidden" : "overflow-y-auto"}`}
            >
                <SeriesList
                    series={data?.series || []}
                    isLoading={isLoading}
                    error={error as Error | null}
                    pagination={data?.pagination}
                    limit={search.limit}
                />
            </main>
        </>
    )
}
