import { createFileRoute, useSearch } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { useState } from "react"
import { getSeries, type Filters, type GetSeriesResponse } from "@/api/api"
import { FiltersSidebar } from "@/components/FiltersSidebar"
import { SeriesList } from "@/components/SeriesList"
import { Button } from "@/components/ui"

export const Route = createFileRoute("/series/")({
    component: Series,
    validateSearch: (search: Record<string, unknown>): Filters => {
        const validLimits = [20, 40, 60]
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
    const [isMobileFiltersOpen, setIsMobileFiltersOpen] = useState(false)

    const { data, isLoading, error } = useQuery({
        queryKey: ["series", search],
        queryFn: async ({ signal }): Promise<GetSeriesResponse> => {
            return getSeries(search, signal)
        },
        retry: 2,
    })

    return (
        <>
            {/* Mobile Filters Toggle Button */}
            <div className="md:hidden fixed bottom-4 right-4 z-50">
                <Button
                    onClick={() => setIsMobileFiltersOpen(true)}
                    size="lg"
                    className="shadow-lg"
                    aria-label="Open filters"
                >
                    <svg className="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"
                        />
                    </svg>
                    Filters
                </Button>
            </div>

            {/* Overlay for mobile filters */}
            {isMobileFiltersOpen && (
                <div
                    className="fixed inset-0 bg-black/50 z-30 md:hidden"
                    onClick={() => setIsMobileFiltersOpen(false)}
                    aria-hidden="true"
                />
            )}

            <div className="flex relative min-h-[calc(100vh-4rem)]">
                {/* Sidebar - Filters */}
                <FiltersSidebar
                    filters={search}
                    isMobileOpen={isMobileFiltersOpen}
                    onMobileClose={() => setIsMobileFiltersOpen(false)}
                />

                {/* Main Content - Results */}
                <main className="flex-1 py-6 ml-72 overflow-y-auto">
                    <SeriesList
                        series={data?.series || []}
                        isLoading={isLoading}
                        error={error as Error | null}
                        pagination={data?.pagination}
                    />
                </main>
            </div>
        </>
    )
}
