import { useNavigate, useRouter, useSearch } from "@tanstack/react-router"
import type { Pagination, MatchSummary, MatchFilters } from "@/types"
import { EmptyState } from "@/components/EmptyState"
import { ErrorState } from "@/components/ErrorState"
import { Skeleton } from "@/components/ui"
import { Button } from "@/components/ui"
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
                {matches.map(match => (
                    <MatchCard key={match.match_id} match={match} />
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
                    <div key={i} className="bg-background-card/80 min-w-67.5 rounded-xl p-4 shadow-xl">
                        <Skeleton className="mb-2 h-6 w-3/4" />
                        <Skeleton className="mb-4 h-4 w-1/2" />
                        <Skeleton className="mb-1 h-4 w-full" />
                        <Skeleton className="h-4 w-2/3" />
                    </div>
                ))}
            </div>
        </div>
    )
}
