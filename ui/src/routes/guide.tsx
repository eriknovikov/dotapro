import { createFileRoute } from "@tanstack/react-router"
import { SEO } from "@/components/SEO"

export const Route = createFileRoute("/guide")({
    component: Guide,
})

function Guide() {
    return (
        <>
            <SEO
                title="Guide"
                description="Learn how to get the most out of dotapro's professional Dota 2 analytics."
            />
            <div className="mx-auto max-w-4xl py-6">
                <div className="p-10 shadow-md shadow-gray-900">
                    <h2 className="text-foreground mb-6 text-2xl font-bold">How to get the most out of dotapro</h2>
                    <div className="text-foreground-muted space-y-4 leading-relaxed">
                        <p>If you want to get just the recent series, go to the series pages.</p>
                        <p>
                            In most cases you just want to see series from a specific league (TI15, or Dreamleague
                            Season 28), so use the corresponding filters.
                        </p>
                        <p>
                            Alternatively, you want to see the series from a specific team (Liquid, Spirit, etc), so
                            filter by that.
                        </p>
                        <p>
                            By default, the popular leagues and teams are shown in the corresponding selectors within
                            the filters, but you can type into the selectors to find the league or team you are looking
                            for.
                        </p>
                        <p>
                            It is advised using the platform on desktop, though it has been designed to be responsive in
                            mobile and tablets in case you don't have access to your PC or laptop.
                        </p>
                    </div>
                </div>
            </div>
        </>
    )
}
