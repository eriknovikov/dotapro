import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/matches")({
    component: Matches,
})

function Matches() {
    return (
        <main className="flex-1 p-6 min-h-[calc(100vh-4rem)]">
            <div className="max-w-7xl mx-auto">
                <h1 className="text-3xl font-bold text-foreground mb-6">matches</h1>
                <p className="text-foreground-muted">Matches page coming soon...</p>
            </div>
        </main>
    )
}
