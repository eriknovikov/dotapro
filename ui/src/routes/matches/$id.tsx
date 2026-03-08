import { createFileRoute, useRouter } from "@tanstack/react-router"
import { ChevronLeft, Swords } from "lucide-react"
import { useMatch } from "@/hooks/useMatches"
import { PlayersTable } from "@/components/series/PlayersTable"
import { Skeleton } from "@/components/ui"
import { ErrorState } from "@/components/ErrorState"
import { SEO } from "@/components/SEO"
import { formatRelativeTime, formatDuration } from "@/lib"
import type { MatchDetail, SeriesMatchDetail } from "@/types"

export const Route = createFileRoute("/matches/$id")({
    component: MatchDetails,
})

function MatchDetails() {
    const { id } = Route.useParams()
    const matchId = parseInt(id)
    const { data: match, isLoading, error } = useMatch(matchId)
    const router = useRouter()

    if (isLoading) {
        return <MatchDetailsSkeleton />
    }

    if (error) {
        return <ErrorState error="Failed to load match details" />
    }

    if (!match) {
        return <ErrorState error="Match not found" />
    }

    return (
        <>
            <SEO
                title={`dotapro.org | ${match.radiant_team.name} vs ${match.dire_team.name}`}
                description={`Match details for ${match.radiant_team.name} vs ${match.dire_team.name} in ${match.league.name}`}
            />
            <div className="mx-auto max-w-7xl px-2 py-6">
                <MatchHeader match={match} router={router} />
                <PlayersTable
                    match={
                        {
                            match_id: match.match_id,
                            duration: match.duration,
                            radiant_win: match.radiant_win,
                            picks_bans: match.picks_bans,
                            players_data: match.players_data,
                            radiant_gold_adv: match.radiant_gold_adv,
                            radiant_xp_adv: match.radiant_xp_adv,
                            radiant_score: match.radiant_team.score,
                            dire_score: match.dire_team.score,
                            radiant_captain: match.radiant_team.captain,
                            dire_captain: match.dire_team.captain,
                        } as SeriesMatchDetail
                    }
                    radiantTeam={match.radiant_team}
                    direTeam={match.dire_team}
                />
            </div>
        </>
    )
}

function MatchHeader({ match, router }: { match: MatchDetail; router: ReturnType<typeof useRouter> }) {
    const { radiant_team, dire_team, league, start_time, duration } = match

    return (
        <div className="rounded-xl pb-2" role="region" aria-label="Match header">
            {/* Back Button */}
            <button
                className="text-foreground-muted hover:text-foreground mb-4 flex w-fit items-center gap-2"
                onClick={() => router.history.back()}
                aria-label="Back to matches list"
            >
                <ChevronLeft className="h-4 w-4" />
                Back to Matches
            </button>

            {/* Mobile Layout */}
            <div className="flex flex-col gap-4 px-2 md:hidden">
                {/* Radiant Team */}
                <div className="flex min-w-0 items-center justify-center">
                    <div className="flex min-w-0 items-center gap-2">
                        {radiant_team.logo_url && (
                            <img
                                src={radiant_team.logo_url}
                                alt={`${radiant_team.name} logo`}
                                className="h-8 w-auto shrink-0 rounded-lg object-contain"
                            />
                        )}
                        <span className="font-shantell truncate text-center text-base font-bold">{radiant_team.name}</span>
                    </div>
                </div>

                {/* VS Icon */}
                <div className="flex items-center justify-center">
                    <Swords className="text-foreground-muted h-6 w-6" />
                </div>

                {/* Dire Team */}
                <div className="flex min-w-0 items-center justify-center">
                    <div className="flex min-w-0 items-center gap-2">
                        {dire_team.logo_url && (
                            <img
                                src={dire_team.logo_url}
                                alt={`${dire_team.name} logo`}
                                className="h-8 w-auto shrink-0 rounded-lg object-contain"
                            />
                        )}
                        <span className="font-shantell truncate text-center text-base font-bold">{dire_team.name}</span>
                    </div>
                </div>

                {/* League, Time, and Duration (mobile - simplified) */}
                <div className="text-foreground-muted flex flex-col text-sm">
                    <span className="font-medium">{league.name}</span>
                    <span>~{formatRelativeTime(start_time)} ago</span>
                    <span>{formatDuration(duration)}</span>
                </div>
            </div>

            {/* Desktop Layout */}
            <div className="hidden flex-col md:flex">
                {/* Teams and VS Icon */}
                <div className="flex flex-row items-center gap-25">
                    {/* Radiant Team */}
                    <div className="flex items-center gap-1">
                        {radiant_team.logo_url && (
                            <img
                                src={radiant_team.logo_url}
                                alt={`${radiant_team.name} logo`}
                                className="h-16 w-auto rounded-lg object-contain md:h-24"
                            />
                        )}
                        <span className="font-shantell text-xl font-bold md:text-3xl">{radiant_team.name}</span>
                    </div>

                    {/* VS Icon */}
                    <div className="flex items-center">
                        <Swords className="text-foreground-muted h-6 w-6" />
                    </div>

                    {/* Dire Team */}
                    <div className="flex items-center gap-1">
                        <span className="font-shantell text-xl font-bold md:text-3xl">{dire_team.name}</span>
                        {dire_team.logo_url && (
                            <img
                                src={dire_team.logo_url}
                                alt={`${dire_team.name} logo`}
                                className="h-16 w-auto rounded-lg object-contain md:h-24"
                            />
                        )}
                    </div>
                </div>

                {/* League, Time, and Duration */}
                <div className="text-foreground-muted mt-4 text-sm">
                    <span className="font-medium">{league.name}</span>
                    <span className="mx-2">•</span>
                    <span className="bg-secondary/20 inline-block rounded px-2 py-0.5 text-xs">{league.tier}</span>
                    <span className="mx-2">•</span>
                    <span>~{formatRelativeTime(start_time)} ago</span>
                    <span className="mx-2">•</span>
                    <span>{formatDuration(duration)}</span>
                </div>
            </div>
        </div>
    )
}

function MatchDetailsSkeleton() {
    return (
        <div className="mx-auto max-w-7xl px-2 py-6">
            <div className="rounded-xl pb-2">
                <Skeleton className="mb-4 h-6 w-24" />
                <div className="flex flex-col gap-4 px-2 md:hidden">
                    <Skeleton className="h-8 w-3/4" />
                    <Skeleton className="mx-auto h-6 w-6" />
                    <Skeleton className="h-8 w-3/4" />
                    <div className="flex flex-col gap-2">
                        <Skeleton className="h-4 w-1/2" />
                        <Skeleton className="h-4 w-1/3" />
                        <Skeleton className="h-4 w-1/4" />
                    </div>
                </div>
                <div className="hidden flex-col md:flex">
                    <div className="flex flex-row items-center gap-25">
                        <Skeleton className="h-16 w-16" />
                        <Skeleton className="h-8 w-32" />
                        <Skeleton className="mx-auto h-6 w-6" />
                        <Skeleton className="h-8 w-32" />
                        <Skeleton className="h-16 w-16" />
                    </div>
                    <div className="mt-4 flex gap-2">
                        <Skeleton className="h-4 w-1/4" />
                        <Skeleton className="h-4 w-16" />
                        <Skeleton className="h-4 w-1/3" />
                        <Skeleton className="h-4 w-1/4" />
                    </div>
                </div>
            </div>
            <Skeleton className="h-64 w-full" />
        </div>
    )
}