import { useNavigate, useRouter, useSearch } from "@tanstack/react-router"
import type { Pagination, MatchSummary, MatchFilters } from "@/types"
import { Button, EmptyState, ErrorState, MatchCardSkeleton } from ".."
import { MatchCard } from "./MatchCard"

interface MatchListProps {
    matches?: MatchSummary[]
    isLoading?: boolean
    error?: Error | null
    pagination?: Pagination
    limit?: number
}

export function MatchList({ matches, isLoading, error, pagination, limit }: MatchListProps) {
    const navigate = useNavigate()
    const router = useRouter()
    const search = useSearch({ strict: false }) as MatchFilters
    const skeletonCount = limit || 9

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

    if (isLoading) {
        return <MatchListSkeleton count={skeletonCount} />
    }

    if (error) {
        return <ErrorState error={error} title="Error loading matches" />
    }

    if (!matches || matches.length === 0) {
        return (
            <div className="flex min-h-[50vh] items-center justify-center">
                <EmptyState
                    icon="🤕"
                    title="No matches found"
                    description="Try adjusting your filters to find what you're looking for."
                />
            </div>
        )
    }

    return (
        <div className="w-full px-2 md:px-12">
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-[repeat(auto-fit,minmax(270px,1fr))] sm:gap-6">
                {matches.map((match, index) => (
                    <MatchCard
                        key={match.match_id}
                        match={match}
                        style={
                            {
                                animationDelay: `${Math.min(index, 10) * 50}ms`,
                            } as React.CSSProperties
                        }
                        className="card-entrance"
                    />
                ))}
            </div>

            <div className="mt-4 flex items-center justify-center gap-3 sm:mt-6 sm:gap-4">
                <Button onClick={handlePrevious} variant="cool-outline" size="sm" disabled={isLoading || !search.c}>
                    Previous
                </Button>
                <Button
                    onClick={handleLoadMore}
                    variant="cool-outline"
                    size="sm"
                    disabled={isLoading || !pagination?.nc || !pagination?.has_more}
                >
                    Next
                </Button>
            </div>
        </div>
    )
}

function MatchListSkeleton({ count = 9 }: { count?: number }) {
    return (
        <div className="w-full px-2 md:px-12">
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-[repeat(auto-fit,minmax(270px,1fr))] sm:gap-6">
                {Array.from({ length: count }).map((_, i) => (
                    <MatchCardSkeleton key={i} />
                ))}
            </div>
        </div>
    )
}
