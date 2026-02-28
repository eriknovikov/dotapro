import { ChevronLeft, Swords } from "lucide-react"
import { Link, useNavigate } from "@tanstack/react-router"
import type { SeriesDetail } from "@/api/api"
import { formatRelativeTime } from "@/lib"
import { Button } from "@/components/ui"

interface SeriesHeaderProps {
    series: SeriesDetail
}

export function SeriesHeader({ series }: SeriesHeaderProps) {
    const { team_a, team_b, league, team_a_score, team_b_score, start_time } = series

    return (
        <div className="rounded-xl pb-2">
            {/* Back Button */}
            <Link
                className="mb-4 text-foreground-muted hover:text-foreground flex items-center gap-2 w-fit"
                to="/series"
            >
                <ChevronLeft className="h-4 w-4" />
                Back to Series
            </Link>

            {/* Mobile Layout */}
            <div className="flex md:hidden flex-col gap-4 px-2">
                {/* Team A */}
                <div className="flex items-center justify-center min-w-0">
                    <div className="flex min-w-0 gap-2 items-center">
                        {team_a.logo_url && (
                            <img
                                src={team_a.logo_url}
                                alt={team_a.name}
                                className="h-8 w-auto rounded-lg object-contain shrink-0"
                            />
                        )}
                        <span className="text-base font-bold font-shantell truncate text-center">{team_a.name}</span>
                        <span className="font-bold text-sm text-foreground-muted">({team_a_score})</span>
                    </div>
                </div>

                {/* VS Icon */}
                <div className="flex items-center justify-center">
                    <Swords className="h-6 w-6 text-foreground-muted" />
                </div>

                {/* Team B */}
                <div className="flex items-center justify-center min-w-0">
                    <div className="flex min-w-0 gap-2 items-center">
                        {team_b.logo_url && (
                            <img
                                src={team_b.logo_url}
                                alt={team_b.name}
                                className="h-8 w-auto rounded-lg object-contain shrink-0"
                            />
                        )}
                        <span className="text-base font-bold font-shantell truncate text-center">{team_b.name}</span>
                        <span className="font-bold text-sm text-foreground-muted">({team_b_score})</span>
                    </div>
                </div>

                {/* League and Time (mobile - simplified) */}
                <div className="text-sm text-foreground-muted flex flex-col">
                    <span className="font-medium">{league.name}</span>
                    <span>~{formatRelativeTime(start_time)} ago</span>
                </div>
            </div>

            {/* Desktop Layout */}
            <div className="hidden md:flex flex-col">
                {/* Teams and Score */}
                <div className="flex flex-row items-center gap-25">
                    {/* Team A */}
                    <div className="flex items-center gap-1">
                        {team_a.logo_url && (
                            <img
                                src={team_a.logo_url}
                                alt={team_a.name}
                                className="h-16 md:h-24 w-auto rounded-lg object-contain"
                            />
                        )}
                        <span className="text-xl md:text-3xl font-bold font-shantell">{team_a.name}</span>
                    </div>

                    {/* Score */}
                    <div className="text-4xl md:text-5xl font-bold font-teko text-foreground">
                        {team_a_score} - {team_b_score}
                    </div>

                    {/* Team B */}
                    <div className="flex items-center gap-1">
                        <span className="text-xl md:text-3xl font-bold font-shantell">{team_b.name}</span>
                        {team_b.logo_url && (
                            <img
                                src={team_b.logo_url}
                                alt={team_b.name}
                                className="h-16 md:h-24 w-auto rounded-lg object-contain"
                            />
                        )}
                    </div>
                </div>

                {/* League and Time */}
                <div className="mt-4 text-sm text-foreground-muted">
                    <span className="font-medium">{league.name}</span>
                    <span className="mx-2">•</span>
                    <span className="inline-block px-2 py-0.5 rounded text-xs bg-secondary/20">{league.tier}</span>
                    <span className="mx-2">•</span>
                    <span>~{formatRelativeTime(start_time)} ago</span>
                </div>
            </div>
        </div>
    )
}
