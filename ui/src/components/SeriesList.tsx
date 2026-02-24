import type { Series, Pagination } from "../api/api"
import { SeriesCard } from "./SeriesCard"
import { EmptyState } from "./EmptyState"
import { ErrorState } from "./ErrorState"
import { Spinner } from "./Spinner"
import { SeriesCardSkeleton } from "./Skeleton"
import { Button } from "./ui"
import { useNavigate, useRouter } from "@tanstack/react-router"
import { useSearch } from "@tanstack/react-router"

interface SeriesListProps {
    series: Series[]
    isLoading?: boolean
    error?: Error | null
    pagination?: Pagination
}

export function SeriesList({ series, isLoading, error, pagination }: SeriesListProps) {
    const navigate = useNavigate()
    const router = useRouter()
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
        router.history.back()
    }

    if (error) {
        return <ErrorState error={error} title="Error loading series" />
    }

    if (isLoading) {
        return (
            <div className="w-full">
                {/* Loading header with spinner */}
                <div className="flex items-center justify-center h-16 mb-6">
                    <div className="flex items-center space-x-2">
                        <Spinner size="lg" />
                        <span className="text-foreground-muted">Loading series...</span>
                    </div>
                </div>

                {/* Skeleton cards matching the grid layout */}
                <div className="grid grid-cols-[repeat(auto-fit,minmax(450px,1fr))] gap-6">
                    {Array.from({ length: 15 }).map((_, i) => (
                        <SeriesCardSkeleton key={i} />
                    ))}
                </div>
            </div>
        )
    }

    if (!series || series.length === 0) {
        return (
            <EmptyState
                icon="🤕"
                title="No series found"
                description="Try adjusting your filters to find what you're looking for."
            />
        )
    }

    return (
        <div className="w-full">
            <div className="grid grid-cols-[repeat(auto-fit,minmax(450px,1fr))] gap-6">
                {series.map(s => (
                    <SeriesCard key={s.series_id} series={s} />
                ))}
            </div>

            <div className="flex items-center justify-center gap-4 mt-6">
                <Button
                    onClick={handlePrevious}
                    variant="cool-outline"
                    disabled={isLoading || !search.c}
                >
                    Previous
                </Button>
                <Button
                    onClick={handleLoadMore}
                    variant="cool-outline"
                    disabled={isLoading || !pagination?.nc || !pagination?.has_more}
                >
                    Next
                </Button>
            </div>
        </div>
    )
}
