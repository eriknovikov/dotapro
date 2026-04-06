import { Button, Tooltip } from "../ui"
import { cn, formatRelativeTime, getHeroImageUrl } from "@/lib"
import type { MatchSummary } from "@/types"
import { useNavigate } from "@tanstack/react-router"
import { Eye, Swords, Trophy } from "lucide-react"
import React, { lazy, Suspense } from "react"

const Card = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(({ className, ...props }, ref) => (
    <div
        ref={ref}
        className={cn(
            "bg-background-card/80 text-card-foreground relative min-w-67.5 rounded-xl shadow-xl backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:scale-[1.01] hover:shadow-2xl",
            className,
        )}
        {...props}
    >
        <div className="relative z-10 flex h-full flex-col">{props.children}</div>
    </div>
))
Card.displayName = "Card"

const CardContent = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
    ({ className, ...props }, ref) => <div ref={ref} className={cn("p-6 pt-0", className)} {...props} />,
)
CardContent.displayName = "CardContent"

// Lazy load tooltip components to reduce initial bundle size
const HeroTooltipContent = lazy(() =>
    import("../series/HeroTooltipContent").then(m => ({ default: m.HeroTooltipContent })),
)

// Wrapper components to handle Suspense inside tooltip content
function HeroTooltipWrapper({ heroId }: { heroId: number }) {
    return (
        <Suspense fallback={<div className="h-32 w-80 animate-pulse rounded-md bg-[#1c1d21]" />}>
            <HeroTooltipContent heroId={heroId} />
        </Suspense>
    )
}

export function MatchCard({
    match,
    style,
    className,
}: {
    match: MatchSummary
    style?: React.CSSProperties
    className?: string
}) {
    const navigate = useNavigate()

    // Extract heroes from hero arrays
    const radiantHeroes = match.radiant_heroes.slice(0, 5)
    const direHeroes = match.dire_heroes.slice(0, 5)

    return (
        <Card
            role="article"
            aria-label={`Match: ${match.radiant_team.name} vs ${match.dire_team.name}`}
            style={style}
            className={className}
        >
            <CardContent className="flex flex-1 flex-col items-center p-4">
                {/* Teams and Heroes - Vertical Layout */}
                <div className="mb-3 flex flex-col items-center pb-3">
                    {/* Radiant Team */}
                    <div className="flex w-full items-center gap-3">
                        {match.radiant_team.logo_url && (
                            <img
                                src={match.radiant_team.logo_url}
                                alt={`${match.radiant_team.name} logo`}
                                className="h-8 w-auto max-w-10 shrink-0 rounded transition-all duration-300 select-none hover:brightness-125 hover:saturate-150"
                            />
                        )}
                        <span className="text-foreground font-shantell text-lg font-bold">
                            {match.radiant_team.name}
                        </span>
                        {match.radiant_win && <Trophy className="h-4 w-4 text-yellow-500" />}
                    </div>
                    {/* Radiant Heroes - Aligned below team name */}
                    <div className="mt-2 flex w-full justify-start gap-1">
                        {radiantHeroes.map(heroId => (
                            <Tooltip
                                key={heroId}
                                content={<HeroTooltipWrapper heroId={heroId} />}
                                side="right"
                                contentClassName="p-0 border-none bg-transparent"
                            >
                                <img
                                    src={getHeroImageUrl(heroId)}
                                    alt={`Hero ${heroId}`}
                                    className="h-6 w-auto cursor-pointer"
                                />
                            </Tooltip>
                        ))}
                    </div>

                    {/* VS Icon - Separate line between teams */}
                    <div className="mt-5 mb-4 flex justify-center">
                        <Swords className="text-foreground-muted h-6" />
                    </div>

                    {/* Dire Team */}
                    <div className="flex w-full items-center gap-3">
                        {match.dire_team.logo_url && (
                            <img
                                src={match.dire_team.logo_url}
                                alt={`${match.dire_team.name} logo`}
                                className="h-8 w-auto max-w-10 shrink-0 rounded transition-all duration-300 select-none hover:brightness-125 hover:saturate-150"
                            />
                        )}
                        <span className="text-foreground font-shantell text-lg font-bold">{match.dire_team.name}</span>
                        {!match.radiant_win && <Trophy className="h-4 w-4 text-yellow-500" />}
                    </div>
                    {/* Dire Heroes - Aligned below team name */}
                    <div className="mt-2 flex w-full justify-start gap-1">
                        {direHeroes.map(heroId => (
                            <Tooltip
                                key={heroId}
                                content={<HeroTooltipWrapper heroId={heroId} />}
                                side="right"
                                contentClassName="p-0 border-none bg-transparent"
                            >
                                <img
                                    src={getHeroImageUrl(heroId)}
                                    alt={`Hero ${heroId}`}
                                    className="h-6 w-auto cursor-pointer"
                                />
                            </Tooltip>
                        ))}
                    </div>
                </div>

                {/* Full-width separator line */}
                <hr className="border-border mx-4 my-2 w-full border-t" />

                {/* League and Date */}
                <div className="text-foreground-muted mb-1 flex flex-1 flex-col self-start text-xs">
                    <span>{match.league.name}</span>
                    <span>~{formatRelativeTime(match.start_time)} ago</span>
                </div>
                <div className="mt-auto flex items-center justify-end self-end px-1 sm:px-2">
                    <Button
                        variant="cool-outline"
                        size="sm"
                        className="group/btn px-2 text-xs"
                        onClick={() => navigate({ to: `/matches/${match.match_id}` })}
                        aria-label={`View details for ${match.radiant_team.name} vs ${match.dire_team.name}`}
                    >
                        View match
                        <div className="flex h-4 w-5 items-center justify-center sm:w-6">
                            <Eye className="size-3 transition-all duration-300 group-hover/btn:size-5" />
                        </div>
                    </Button>
                </div>
            </CardContent>
        </Card>
    )
}
