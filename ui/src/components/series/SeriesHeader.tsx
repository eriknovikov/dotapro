import { ChevronLeft, Swords } from "lucide-react"
import { useRouter } from "@tanstack/react-router"
import type { SeriesDetail } from "@/api/api"
import { formatRelativeTime } from "@/lib"

interface SeriesHeaderProps {
    series: SeriesDetail
}

export function SeriesHeader({ series }: SeriesHeaderProps) {
    const router = useRouter()
    const { team_a, team_b, league, team_a_score, team_b_score, start_time } = series

    return (
        <div className="rounded-xl pb-2" role="region" aria-label="Series header">
            {/* Back Button */}
            <button
                className="text-foreground-muted hover:text-foreground mb-4 flex w-fit items-center gap-2"
                onClick={() => router.history.back()}
                aria-label="Back to series list"
            >
                <ChevronLeft className="h-4 w-4" />
                Back to Series
            </button>

            {/* Mobile Layout */}
            <div className="flex flex-col gap-4 px-2 md:hidden">
                {/* Team A */}
                <div className="flex min-w-0 items-center justify-center">
                    <div className="flex min-w-0 items-center gap-2">
                        {team_a.logo_url && (
                            <img
                                src={team_a.logo_url}
                                alt={`${team_a.name} logo`}
                                className="h-8 w-auto shrink-0 rounded-lg object-contain"
                            />
                        )}
                        <span className="font-shantell truncate text-center text-base font-bold">{team_a.name}</span>
                        <span className="text-foreground-muted text-sm font-bold">({team_a_score})</span>
                    </div>
                </div>

                {/* VS Icon */}
                <div className="flex items-center justify-center">
                    <Swords className="text-foreground-muted h-6 w-6" />
                </div>

                {/* Team B */}
                <div className="flex min-w-0 items-center justify-center">
                    <div className="flex min-w-0 items-center gap-2">
                        {team_b.logo_url && (
                            <img
                                src={team_b.logo_url}
                                alt={`${team_b.name} logo`}
                                className="h-8 w-auto shrink-0 rounded-lg object-contain"
                            />
                        )}
                        <span className="font-shantell truncate text-center text-base font-bold">{team_b.name}</span>
                        <span className="text-foreground-muted text-sm font-bold">({team_b_score})</span>
                    </div>
                </div>

                {/* League and Time (mobile - simplified) */}
                <div className="text-foreground-muted flex flex-col text-sm">
                    <span className="font-medium">{league.name}</span>
                    <span>~{formatRelativeTime(start_time)} ago</span>
                </div>
            </div>

            {/* Desktop Layout */}
            <div className="hidden flex-col md:flex">
                {/* Teams and Score */}
                <div className="flex flex-row items-center gap-25">
                    {/* Team A */}
                    <div className="flex items-center gap-1">
                        {team_a.logo_url && (
                            <img
                                src={team_a.logo_url}
                                alt={`${team_a.name} logo`}
                                className="h-16 w-auto rounded-lg object-contain md:h-24"
                            />
                        )}
                        <span className="font-shantell text-xl font-bold md:text-3xl">{team_a.name}</span>
                    </div>

                    {/* Score */}
                    <div className="font-teko text-foreground text-4xl font-bold md:text-5xl">
                        {team_a_score} - {team_b_score}
                    </div>

                    {/* Team B */}
                    <div className="flex items-center gap-1">
                        <span className="font-shantell text-xl font-bold md:text-3xl">{team_b.name}</span>
                        {team_b.logo_url && (
                            <img
                                src={team_b.logo_url}
                                alt={`${team_b.name} logo`}
                                className="h-16 w-auto rounded-lg object-contain md:h-24"
                            />
                        )}
                    </div>
                </div>

                {/* League and Time */}
                <div className="text-foreground-muted mt-4 text-sm">
                    <span className="font-medium">{league.name}</span>
                    <span className="mx-2">•</span>
                    <span className="bg-secondary/20 inline-block rounded px-2 py-0.5 text-xs">{league.tier}</span>
                    <span className="mx-2">•</span>
                    <span>~{formatRelativeTime(start_time)} ago</span>
                </div>
            </div>
        </div>
    )
}
