import type { Series, Pagination } from "../api/api"
import { SeriesCard } from "./SeriesCard"
import { EmptyState, ErrorState, Spinner, Button } from "./ui/index"
import { useNavigate } from "@tanstack/react-router"
import { useSearch } from "@tanstack/react-router"

interface SeriesListProps {
    series: Series[]
    isLoading?: boolean
    error?: Error | null
    pagination?: Pagination
}

export function SeriesList({ series, isLoading, error, pagination }: SeriesListProps) {
    const navigate = useNavigate()
    const search = useSearch({ strict: false })

    const handleLoadMore = () => {
        if (pagination?.nc) {
            navigate({
                to: ".",
                search: { ...search, c: pagination.nc },
            })
        }
    }

    const handlePrevious = () => {
        if (pagination?.pc) {
            navigate({
                to: ".",
                search: { ...search, c: pagination.pc },
            })
        }
    }

    if (error) {
        return <ErrorState error={error} title="Error loading series" />
    }

    if (isLoading) {
        return (
            <div className="flex items-center justify-center h-64 bg-background-muted border border-border rounded-lg">
                <div className="flex items-center space-x-2">
                    <Spinner />
                    <span className="text-foreground-muted">Loading series...</span>
                </div>
            </div>
        )
    }

    if (!series || series.length === 0) {
        return (
            <EmptyState
                title="No series found"
                description="Try adjusting your filters to find what you're looking for."
            />
        )
    }

    return (
        <div className="w-full">
            <div className="grid grid-cols-[repeat(auto-fit,minmax(420px,1fr))] gap-6">
                {series.map(s => (
                    <SeriesCard key={s.series_id} series={s} />
                ))}
            </div>

            {/* Pagination Controls */}
            <div className="flex items-center justify-center gap-4 mt-6">
                {pagination?.pc && (
                    <Button onClick={handlePrevious} variant="outline" disabled={isLoading}>
                        Previous
                    </Button>
                )}
                {pagination?.has_more && (
                    <Button onClick={handleLoadMore} disabled={isLoading}>
                        Load More
                    </Button>
                )}
            </div>
        </div>
    )
}
