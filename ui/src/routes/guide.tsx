import { SEO } from "@/components"
import { createFileRoute, Link } from "@tanstack/react-router"

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
            <div className="mx-auto flex max-w-4xl flex-col gap-10 py-6">
                <div className="p-10 shadow-md shadow-gray-900">
                    <h2 className="text-foreground mb-6 text-2xl font-bold">How to get the most out of dotapro</h2>
                    <div className="text-foreground-muted space-y-4 leading-relaxed">
                        <p>
                            If you want to get just the recent series, go to the <Link to="/series">series page</Link>.
                            Here you can search the series by league and team. You could find data like "the latest
                            series from The International 2025", or "from Team Liquid". Or even mixing them up -
                            "getting the latest series from Team Liquid in The International 2025". Most of the time you
                            will be searching for the series from a specific League like{" "}
                            <Link to="/series" search={{ league: 19435 }} className="text-blue-200">
                                PGL Wallachia 2026 Season 7
                            </Link>
                            , unless you are hunting for a specific team, like{" "}
                            <Link to="/series" search={{ team: 8291895 }} className="text-blue-200">
                                Tundra
                            </Link>
                            .
                        </p>
                        <p className="">
                            Alternatively, you may want to get the recent{" "}
                            <Link to="/matches" className="text-blue-200">
                                matches
                            </Link>{" "}
                            instead. Here, you can search for the matches directly, using the same filters plus
                            searching by hero or player.
                        </p>
                        <p>
                            It is advised using the platform on desktop, though it has been designed to be responsive in
                            mobile and tablets in case you don't have access to your PC or laptop.
                        </p>
                    </div>
                </div>
                <div className="p-10 shadow-md shadow-gray-900">
                    <h2 className="text-foreground mb-4 text-2xl font-bold">Want a new feature?</h2>
                    <p className="text-foreground-muted leading-relaxed">
                        Take a look at the{" "}
                        <Link to="/about" hash="features-to-add" className="text-blue-200">
                            features to add section
                        </Link>
                        .
                    </p>
                </div>
            </div>
        </>
    )
}
