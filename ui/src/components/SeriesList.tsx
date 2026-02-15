import type { Series } from "../api/api"
import { Badge, Card, CardContent, EmptyState, ErrorState, Spinner } from "./ui/index"

interface SeriesListProps {
    series: Series[]
    isLoading?: boolean
    error?: Error | null
}

// Helper function to format date
function formatDate(dateString: string): string {
    const date = new Date(dateString)
    return date.toLocaleDateString("en-US", {
        year: "numeric",
        month: "short",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    })
}

function SeriesCard({ series }: { series: Series }) {
    return (
        <Card className="hover:shadow-md transition-shadow duration-200">
            <CardContent className="p-4">
                {/* Teams and Score */}
                <div className="flex items-center justify-between">
                    {/* Team A */}
                    <div className="flex items-center flex-1">
                        {series.team_a.logo_url && (
                            <img src={series.team_a.logo_url} alt={series.team_a.name} className="w-8 h-8 rounded" />
                        )}
                        <div className="flex flex-col">
                            <span className="font-medium text-foreground text-lg">{series.team_a.name}</span>
                        </div>
                    </div>

                    {/* Score or VS*/}
                    <div className="flex items-center">
                        <span className="font-bold text-lg text-primary-600">{series.team_a_score}</span>
                        <span className="mx-2 text-gray-400">-</span>
                        <span className="font-bold text-lg text-primary-600">{series.team_b_score}</span>
                    </div>

                    {/* Team B */}
                    <div className="flex items-center flex-1">
                        <div className="flex flex-col items-end ">
                            <span className="font-medium text-foreground text-lg">{series.team_b.name}</span>
                        </div>
                        {series.team_b.logo_url && (
                            <img src={series.team_b.logo_url} alt={series.team_b.name} className="w-8 h-8 rounded" />
                        )}
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}

export function SeriesList({ series, isLoading, error }: SeriesListProps) {
    if (error) {
        return <ErrorState error={error} title="Error loading series" />
    }

    if (isLoading) {
        return (
            <div className="flex items-center justify-center h-64 bg-background-muted border border-border rounded-lg">
                <div className="flex items-center space-x-2">
                    <Spinner />
                    <span className="text-foreground-muted">Loading series...</span>
                </div>
            </div>
        )
    }

    if (!series || series.length === 0) {
        return (
            <EmptyState
                title="No series found"
                description="Try adjusting your filters to find what you're looking for."
            />
        )
    }

    return (
        <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3  gap-4">
            {series.map(s => (
                <SeriesCard key={s.series_id} series={s} />
            ))}
        </div>
    )
}
