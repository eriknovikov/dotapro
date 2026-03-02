import { createFileRoute, useSearch, useNavigate } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { useState, useEffect, useRef } from "react"
import { Funnel } from "lucide-react"
import { getSeries, type Filters, type GetSeriesResponse } from "@/api/api"
import { FiltersSidebar } from "@/components/FiltersSidebar"
import { SeriesList } from "@/components/SeriesList"
import { Button } from "@/components/ui"
import { SEO } from "@/components/SEO"

export const Route = createFileRoute("/series/")({
    component: Series,
    validateSearch: (search: Record<string, unknown>): Filters => {
        const validLimits = [10, 20, 40, 60]
        const limit = search.limit !== undefined ? Number(search.limit) : undefined
        return {
            league: search.league !== undefined ? Number(search.league) : undefined,
            team: search.team !== undefined ? Number(search.team) : undefined,
            sort: typeof search.sort === "string" ? search.sort : undefined,
            limit: limit !== undefined && validLimits.includes(limit) ? limit : undefined,
            c: search.c !== undefined ? Number(search.c) : undefined,
        }
    },
})

function Series() {
    const search = useSearch({ strict: false })
    const navigate = useNavigate()
    const [isMobileFiltersOpen, setIsMobileFiltersOpen] = useState(false)
    const [hasSetMobileDefault, setHasSetMobileDefault] = useState(false)
    const scrollPositionRef = useRef(0)

    // Set default limit to 10 on mobile (for better performance)
    useEffect(() => {
        if (!hasSetMobileDefault && search.limit === undefined && window.innerWidth < 1024) {
            navigate({
                to: ".",
                search: { ...search, limit: 10 },
                replace: true,
            })
            setHasSetMobileDefault(true)
        }
    }, [search, navigate, hasSetMobileDefault])

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
            <div className="lg:hidden fixed bottom-4 right-4 z-50">
                <Button
                    onClick={() => setIsMobileFiltersOpen(true)}
                    size="lg"
                    className="shadow-lg"
                    aria-label="Open filters"
                >
                    <Funnel className="h-5 w-5 mr-2" />
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
                className={`min-h-[calc(100vh-4rem)] py-6 px-2 sm:px-0 lg:ml-72 ${isMobileFiltersOpen ? "overflow-hidden" : "overflow-y-auto"}`}
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
