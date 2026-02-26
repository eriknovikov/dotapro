import { ChevronLeft } from "lucide-react"
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
            {/* Teams and Score */}
            <div className="flex flex-col md:flex-row items-center gap-25">
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
            <div className=" mt-4 text-sm text-foreground-muted">
                <span className="font-medium">{league.name}</span>
                <span className="mx-2">•</span>
                <span className="inline-block px-2 py-0.5 rounded text-xs bg-secondary/20">{league.tier}</span>
                <span className="mx-2">•</span>
                <span>~{formatRelativeTime(start_time)} ago</span>
            </div>
        </div>
    )
}
