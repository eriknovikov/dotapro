import * as React from "react"
import type { Series } from "@/api/api"
import { cn, formatRelativeTime } from "@/lib"
import { Eye, Swords } from "lucide-react"
import { useNavigate } from "@tanstack/react-router"
import { Button } from "@/components/ui"

const Card = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(({ className, ...props }, ref) => (
    <div
        ref={ref}
        className={cn(
            "group relative rounded-xl bg-background-card/80 text-card-foreground shadow-xl transition-all duration-200 ease-in-out hover:shadow-2xl overflow-hidden backdrop-blur-sm",
            className,
        )}
        {...props}
    >
        <div className="absolute inset-0 pointer-events-none overflow-hidden">
            <div className="absolute inset-0 bg-linear-to-br from-primary-500/20 via-primary-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-200 ease-in-out" />
        </div>

        <div className="absolute inset-0 pointer-events-none">
            <div className="absolute inset-0 bg-linear-to-t from-primary-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-200 ease-in-out" />
        </div>

        <div className="relative z-10 flex flex-col h-full">{props.children}</div>
    </div>
))
Card.displayName = "Card"

const CardContent = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
    ({ className, ...props }, ref) => <div ref={ref} className={cn("p-6 pt-0", className)} {...props} />,
)
CardContent.displayName = "CardContent"

export function SeriesCard({ series }: { series: Series }) {
    const navigate = useNavigate()

    return (
        <Card>
            <CardContent className="p-4 flex flex-col flex-1">
                {/* Teams and Score */}
                <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-1 sm:gap-4 border-b border-border mb-3 pb-3">
                    {/* Team A */}
                    <div className="flex items-center justify-center min-w-0 flex-1 sm:justify-start">
                        <div className="flex min-w-0 gap-1 sm:gap-2 items-center">
                            {series.team_a.logo_url && (
                                <img
                                    src={series.team_a.logo_url}
                                    className="h-6 sm:h-7 w-auto max-w-8 sm:max-w-11 rounded shrink-0 select-none"
                                />
                            )}
                            <span className="font-bold text-foreground text-sm sm:text-md font-shantell truncate text-center sm:text-left">
                                {series.team_a.name}
                            </span>
                            <span className="font-bold text-xs text-foreground-muted sm:hidden">
                                ({series.team_a_score})
                            </span>
                        </div>
                    </div>

                    {/* VS Icon (mobile) / Score (desktop) */}
                    <div className="flex items-center justify-center shrink-0 select-none py-0.5 sm:py-0 sm:px-2">
                        <Swords className="h-4 w-4 text-foreground-muted sm:hidden" />
                        <div className="hidden sm:flex items-center">
                            <span className="font-bold text-xs text-foreground-muted">{series.team_a_score}</span>
                            <span className="mx-2 text-gray-400">{"-"}</span>
                            <span className="font-bold text-xs text-foreground-muted">{series.team_b_score}</span>
                        </div>
                    </div>

                    {/* Team B */}
                    <div className="flex items-center justify-center min-w-0 flex-1 sm:justify-end">
                        <div className="flex min-w-0 gap-1 sm:gap-2 items-center">
                            {series.team_b.logo_url && (
                                <img
                                    src={series.team_b.logo_url}
                                    className="h-6 sm:h-7 w-auto max-w-8 sm:max-w-11 rounded shrink-0 select-none"
                                />
                            )}
                            <span className="font-bold text-foreground text-sm sm:text-md font-shantell truncate text-center sm:text-right">
                                {series.team_b.name}
                            </span>
                            <span className="font-bold text-xs text-foreground-muted sm:hidden">
                                ({series.team_b_score})
                            </span>
                        </div>
                    </div>
                </div>

                {/* League and Date */}
                <div className="text-xs text-foreground-muted flex flex-col flex-1">
                    <span>{series.league.name}</span>
                    <span>~{formatRelativeTime(series.start_time)} ago</span>
                </div>
                <div className="flex justify-end items-center mt-auto px-1 sm:px-2">
                    <Button
                        variant="cool-outline"
                        size="sm"
                        className="px-2 group/btn text-xs"
                        onClick={() => navigate({ to: `/series/${series.series_id}` })}
                    >
                        View series
                        <div className="flex justify-center items-center w-5 sm:w-6 h-4">
                            <Eye className="size-3 transition-all duration-300 group-hover/btn:size-5" />
                        </div>
                    </Button>
                </div>
            </CardContent>
        </Card>
    )
}
