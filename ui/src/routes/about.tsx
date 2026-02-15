import { createFileRoute } from "@tanstack/react-router"
import { Card, CardContent } from "../components/ui/index"

export const Route = createFileRoute("/about")({
    component: About,
})

function About() {
    return (
        <div className="p-6">
            <Card>
                <CardContent>
                    <h1 className="text-2xl font-bold text-foreground mb-4">About</h1>
                    <p className="text-foreground-muted">
                        Dota 2 match data aggregation system. Scrapes from OpenDota API, stores in PostgreSQL, provides REST API.
                    </p>
                </CardContent>
            </Card>
        </div>
    )
}
