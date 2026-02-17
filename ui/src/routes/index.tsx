import { createFileRoute, Link } from "@tanstack/react-router"

export const Route = createFileRoute("/")({
    component: Index,
})

function Index() {
    return (
        <main className="flex-1 p-6 min-h-[calc(100vh-4rem)]">
                <div className="max-w-7xl mx-auto">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 h-[calc(100vh-8rem)]">
                        {/* Matches Card */}
                        <Link
                            to="/matches"
                            className="group relative overflow-hidden rounded-xl border border-border bg-background-card hover:border-primary-500 transition-all duration-300 hover:shadow-lg hover:shadow-primary-500/10"
                        >
                            <div className="absolute inset-0 bg-gradient-to-br from-primary-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
                            <div className="relative h-full flex flex-col items-center justify-center p-12">
                                <div className="mb-6">
                                    <svg
                                        className="h-24 w-24 text-foreground-muted group-hover:text-primary-500 transition-colors duration-300"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                    >
                                        <path
                                            strokeLinecap="round"
                                            strokeLinejoin="round"
                                            strokeWidth={1.5}
                                            d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"
                                        />
                                    </svg>
                                </div>
                                <h2 className="text-4xl font-bold text-foreground mb-2 group-hover:text-primary-500 transition-colors duration-300">
                                    matches
                                </h2>
                                <p className="text-foreground-muted text-lg">Browse and explore Dota 2 matches</p>
                            </div>
                        </Link>

                        {/* Series Card */}
                        <Link
                            to="/series"
                            className="group relative overflow-hidden rounded-xl border border-border bg-background-card hover:border-secondary-500 transition-all duration-300 hover:shadow-lg hover:shadow-secondary-500/10"
                        >
                            <div className="absolute inset-0 bg-gradient-to-br from-secondary-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
                            <div className="relative h-full flex flex-col items-center justify-center p-12">
                                <div className="mb-6">
                                    <svg
                                        className="h-24 w-24 text-foreground-muted group-hover:text-secondary-500 transition-colors duration-300"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                    >
                                        <path
                                            strokeLinecap="round"
                                            strokeLinejoin="round"
                                            strokeWidth={1.5}
                                            d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                                        />
                                    </svg>
                                </div>
                                <h2 className="text-4xl font-bold text-foreground mb-2 group-hover:text-secondary-500 transition-colors duration-300">
                                    series
                                </h2>
                                <p className="text-foreground-muted text-lg">View series and tournament results</p>
                            </div>
                        </Link>
                    </div>
                </div>
            </main>
    )
}
