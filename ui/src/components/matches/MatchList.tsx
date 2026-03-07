import { EmptyState } from "@/components/EmptyState"
import { ErrorState } from "@/components/ErrorState"
import { Skeleton } from "@/components/ui"
import { useMatches } from "@/hooks/useMatches"
import type { MatchFilters } from "@/types"
import { MatchCard } from "./MatchCard"

interface MatchListProps {
    filters: MatchFilters
}

export function MatchList({ filters }: MatchListProps) {
    const { data, isLoading, error } = useMatches(filters)

    if (isLoading) {
        return <MatchListSkeleton />
    }

    if (error) {
        return <ErrorState error="Failed to load matches" />
    }

    if (!data?.matches.length) {
        return <EmptyState title="No matches found" />
    }

    return (
        <div className="w-full px-2 md:px-12">
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-[repeat(auto-fit,minmax(270px,1fr))] sm:gap-6">
                {data.matches.map(match => (
                    <MatchCard key={match.match_id} match={match} />
                ))}
            </div>
        </div>
    )
}

function MatchListSkeleton() {
    return (
        <div className="w-full px-2 md:px-12">
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-[repeat(auto-fit,minmax(270px,1fr))] sm:gap-6">
                {Array.from({ length: 9 }).map((_, i) => (
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
