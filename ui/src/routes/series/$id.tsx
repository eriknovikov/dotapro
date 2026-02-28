import { createFileRoute, useParams } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { getSeriesById, type SeriesDetail } from "@/api/api"
import { Spinner } from "@/components/Spinner"
import { ErrorState } from "@/components/ErrorState"
import { SeriesHeader, GameTabContent } from "@/components/series"
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui"
import { SEO } from "@/components/SEO"

export const Route = createFileRoute("/series/$id")({
    component: SeriesDetailPage,
})

function SeriesDetailPage() {
    const { id } = useParams({ strict: false })
    const seriesId = Number(id)

    const { data, isLoading, error } = useQuery({
        queryKey: ["series", seriesId],
        queryFn: async ({ signal }): Promise<SeriesDetail> => {
            return getSeriesById(seriesId, signal)
        },
        retry: 2,
    })

    if (isLoading) {
        return (
            <div className="flex items-center justify-center min-h-[calc(100vh-4rem)]">
                <Spinner />
            </div>
        )
    }

    if (error) {
        const errorMessage = (error as Error).message
        const isNotFound = errorMessage.includes("404")
        return (
            <div className="flex items-center justify-center min-h-[calc(100vh-4rem)]">
                <ErrorState
                    title={isNotFound ? "Series Not Found" : "Error Loading Series"}
                    error={isNotFound ? `Series with ID ${seriesId} could not be found.` : errorMessage}
                />
            </div>
        )
    }

    if (!data) {
        return (
            <div className="flex items-center justify-center min-h-[calc(100vh-4rem)]">
                <ErrorState title="Series Not Found" error={`Series with ID ${seriesId} could not be found.`} />
            </div>
        )
    }

    return (
        <>
            <SEO
                title={`${data.team_a.name} vs ${data.team_b.name} - ${data.league.name}`}
                description={`View match details for ${data.team_a.name} vs ${data.team_b.name} in ${data.league.name}. Score: ${data.team_a_score} - ${data.team_b_score}`}
            />
            <div className="py-6 px-2 max-w-7xl mx-auto">
            <SeriesHeader series={data} />

            {data.matches.length > 0 ? (
                <Tabs defaultValue="game-1">
                    <TabsList className="max-w-full overflow-x-auto overflow-y-hidden mb-4">
                        {data.matches.map((_, i) => (
                            <TabsTrigger key={i} value={`game-${i + 1}`} className="px-2 sm:px-4 text-sm sm:text-base">
                                {data.matches.length >= 4 ? (
                                    <>
                                        <span className="hidden sm:inline">Game {i + 1}</span>
                                        <span className="sm:hidden">G{i + 1}</span>
                                    </>
                                ) : (
                                    <>Game {i + 1}</>
                                )}
                            </TabsTrigger>
                        ))}
                    </TabsList>
                    {data.matches.map((match, i) => (
                        <TabsContent key={i} value={`game-${i + 1}`}>
                            <GameTabContent
                                match={match}
                                gameNumber={i + 1}
                                radiantTeam={data.team_a}
                                direTeam={data.team_b}
                            />
                        </TabsContent>
                    ))}
                </Tabs>
            ) : (
                <div className="text-center py-12 text-foreground-muted">No matches available for this series.</div>
            )}
        </div>
        </>
    )
}
