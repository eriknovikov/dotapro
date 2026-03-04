import { createFileRoute } from "@tanstack/react-router"
import { SEO } from "@/components/SEO"

export const Route = createFileRoute("/matches")({
    component: Matches,
})

function Matches() {
    return (
        <>
            <SEO
                title="dotapro.org | Matches"
                description="Browse professional Dota 2 matches. Coming soon with advanced filtering and statistics."
            />
            <main className="min-h-[calc(100vh-4rem)] flex-1 py-6">
                <div className="mx-auto max-w-7xl">
                    <h1 className="text-foreground mb-6 text-3xl font-bold">matches</h1>
                    <p className="text-foreground-muted">Matches page coming soon...</p>
                </div>
            </main>
        </>
    )
}
