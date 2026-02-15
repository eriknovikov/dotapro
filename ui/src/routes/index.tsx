import { createFileRoute, useSearch } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { getSeries, getSeriesMOCK, type Filters, type GetSeriesResponse } from "../api/api"
import { FiltersSidebar } from "../components/FiltersSidebar"
import { SeriesList } from "../components/SeriesList"

export const Route = createFileRoute("/")({
    component: Index,
    validateSearch: (search: Record<string, unknown>): Filters => ({
        league: search.league !== undefined ? Number(search.league) : undefined,
        team: search.team !== undefined ? Number(search.team) : undefined,
        sort: typeof search.sort === "string" ? search.sort : undefined,
    }),
})

function Index() {
    const search = useSearch({ strict: false })

    const { data, isLoading, error } = useQuery({
        queryKey: ["series", search],
        queryFn: async ({ signal }): Promise<GetSeriesResponse> => {
            const filters: Filters = {
                league: search.league !== undefined ? Number(search.league) : undefined,
                team: search.team !== undefined ? Number(search.team) : undefined,
                sort: search.sort as Filters["sort"],
            }
            return getSeriesMOCK(filters, signal)
        },
        retry: 2,
    })

    return (
        <div className="flex min-h-screen bg-background">
            {/* Sidebar - Filters */}
            <FiltersSidebar initialFilters={search} isLoading={isLoading} />

            {/* Main Content - Results */}
            <main className="flex-1 p-6 overflow-auto">
                <div className="max-w-7xl mx-auto">
                    <h1 className="text-2xl font-bold text-foreground mb-6">Series</h1>

                    <SeriesList series={data?.series || []} isLoading={isLoading} error={error as Error | null} />
                </div>
            </main>
        </div>
    )
}
