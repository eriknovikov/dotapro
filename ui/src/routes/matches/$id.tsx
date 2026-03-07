import { createFileRoute } from "@tanstack/react-router"
import { useMatch } from "@/hooks/useMatches"
import { PlayersTable } from "@/components/series/PlayersTable"
import { Skeleton } from "@/components/ui"
import { ErrorState } from "@/components/ErrorState"
import { SEO } from "@/components/SEO"
import type { MatchDetail, SeriesMatchDetail } from "@/types"

export const Route = createFileRoute("/matches/$id")({
    component: MatchDetails,
})

function MatchDetails() {
    const { id } = Route.useParams()
    const matchId = parseInt(id)
    const { data: match, isLoading, error } = useMatch(matchId)

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
            <main className="min-h-[calc(100vh-4rem)] flex-1 py-6">
                <div className="mx-auto max-w-7xl">
                    <MatchHeader match={match} />
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
            </main>
        </>
    )
}

function MatchHeader({ match }: { match: MatchDetail }) {
    return (
        <div className="bg-background-accent rounded-lg p-6 mb-6">
            <div className="flex flex-col items-center gap-4 md:flex-row md:items-start md:justify-between">
                <div className="flex items-center gap-4">
                    {match.radiant_team.logo_url && (
                        <img
                            src={match.radiant_team.logo_url}
                            alt={`${match.radiant_team.name} logo`}
                            className="h-12 w-auto"
                        />
                    )}
                    <div>
                        <h1 className="text-2xl font-bold">{match.radiant_team.name}</h1>
                        <p className="text-foreground-muted">{match.radiant_team.tag}</p>
                    </div>
                </div>
                
                <div className="text-center">
                    <div className="text-3xl font-bold">
                        {match.radiant_win ? (
                            <span className="text-success-200">VICTORY</span>
                        ) : (
                            <span className="text-error-300">DEFEAT</span>
                        )}
                    </div>
                    <p className="text-foreground-muted">{Math.floor(match.duration / 60)} minutes</p>
                </div>
                
                <div className="flex items-center gap-4">
                    <div className="text-right">
                        <h1 className="text-2xl font-bold">{match.dire_team.name}</h1>
                        <p className="text-foreground-muted">{match.dire_team.tag}</p>
                    </div>
                    {match.dire_team.logo_url && (
                        <img
                            src={match.dire_team.logo_url}
                            alt={`${match.dire_team.name} logo`}
                            className="h-12 w-auto"
                        />
                    )}
                </div>
            </div>
            
            <div className="mt-4 text-center text-foreground-muted">
                <p>{match.league.name}</p>
                <p>{new Date(match.start_time).toLocaleString()}</p>
                <p>Patch {match.patch}</p>
            </div>
        </div>
    )
}

function MatchDetailsSkeleton() {
    return (
        <main className="min-h-[calc(100vh-4rem)] flex-1 py-6">
            <div className="mx-auto max-w-7xl">
                <div className="bg-background-accent rounded-lg p-6 mb-6">
                    <Skeleton className="h-8 w-3/4 mb-4" />
                    <Skeleton className="h-4 w-1/2" />
                </div>
                <Skeleton className="h-64 w-full" />
            </div>
        </main>
    )
}