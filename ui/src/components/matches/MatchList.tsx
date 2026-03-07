import { useMatches } from "@/hooks/useMatches"
import { MatchCard } from "./MatchCard"
import { Skeleton } from "@/components/ui"
import { EmptyState } from "@/components/EmptyState"
import { ErrorState } from "@/components/ErrorState"
import type { MatchFilters } from "@/types"

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
                    <div key={i} className="bg-background-card/80 rounded-xl shadow-xl p-4 min-w-[270px]">
                        <Skeleton className="h-6 w-3/4 mb-2" />
                        <Skeleton className="h-4 w-1/2 mb-4" />
                        <Skeleton className="h-4 w-full mb-1" />
                        <Skeleton className="h-4 w-2/3" />
                    </div>
                ))}
            </div>
        </div>
    )
}