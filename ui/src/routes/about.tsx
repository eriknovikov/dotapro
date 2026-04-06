import { SEO } from "@/components"
import { Card, CardContent } from "@/components/ui/card"
import { createFileRoute } from "@tanstack/react-router"
import { Heart, Lightbulb, Users } from "lucide-react"

export const Route = createFileRoute("/about")({
    component: About,
})

function About() {
    return (
        <>
            <SEO
                title="About"
                description="Learn about dotapro - an open source, completely free-to-use platform about professional Dota 2 analytics."
            />
            <div className="mx-4 min-h-[calc(100vh-4rem)] max-w-4xl space-y-8 py-6 sm:mx-auto">
                {/* Hero Section */}
                <div className="space-y-4 py-8 text-center">
                    <h1 className="text-foreground text-3xl font-bold tracking-tight sm:text-3xl md:text-4xl lg:text-5xl">
                        About dotapro
                    </h1>
                    <p className="text-muted-foreground mx-auto max-w-2xl text-base sm:text-lg md:text-lg lg:text-xl">
                        An open source platform built for the professional <span className="text-nowrap">Dota 2</span>{" "}
                        community.
                    </p>
                </div>

                {/* What's dotapro */}
                <Card className="border-border/50">
                    <CardContent className="pt-6">
                        <h2 className="text-foreground mb-4 text-xl font-semibold">What is dotapro?</h2>
                        <p className="text-muted-foreground mb-4 leading-relaxed">
                            dotapro is an open source, free-to-use platform focused exclusively on professional{" "}
                            <span className="text-nowrap">Dota 2</span>. While established platforms like Dotabuff,
                            Stratz, and OpenDota offer comprehensive data across all matches, dotapro's focus is
                            different.
                        </p>
                        <p className="text-muted-foreground leading-relaxed">
                            We distill the noise of public match data and deliver the analytics that actually matter to
                            Dota players—from casual fans following their favorite teams to pro players scouting
                            opponents before big matches.
                        </p>
                    </CardContent>
                </Card>

                {/* Who is it for */}
                <Card className="border-border/50">
                    <CardContent className="pt-6">
                        <div className="bg-primary/10 mb-4 flex h-12 w-12 items-center justify-center rounded-lg">
                            <Users className="text-primary h-6 w-6" />
                        </div>
                        <h2 className="text-foreground mb-6 text-xl font-semibold">Who is dotapro for?</h2>
                        <ul className="text-muted-foreground space-y-3">
                            <li>
                                <strong className="text-foreground font-medium">Casual fans</strong> – Stay updated on
                                your favorite teams without drowning in statistics.
                            </li>
                            <li>
                                <strong className="text-foreground font-medium">Casters and analysts</strong> – Quick
                                access to match history and team performance for broadcast prep.
                            </li>
                            <li>
                                <strong className="text-foreground font-medium">Pro players & coaches</strong> – Scout
                                opponents, analyze strategies, track the evolving meta.
                            </li>
                            <li>
                                <strong className="text-foreground font-medium">Dota enthusiasts</strong> – Deep dive
                                into professional match data with clean, focused tooling.
                            </li>
                        </ul>
                    </CardContent>
                </Card>

                {/* Credits */}
                <Card className="border-border/50">
                    <CardContent className="pt-6">
                        <div className="bg-primary/10 mb-4 flex h-12 w-12 items-center justify-center rounded-lg">
                            <Heart className="text-primary h-6 w-6" />
                        </div>
                        <h2 className="text-foreground mb-4 text-xl font-semibold">Powered by OpenDota</h2>
                        <p className="text-muted-foreground leading-relaxed">
                            dotapro wouldn't exist without the incredible work of the{" "}
                            <a
                                href="https://www.opendota.com"
                                target="_blank"
                                rel="noopener noreferrer"
                                className="text-blue-500 hover:underline"
                            >
                                OpenDota
                            </a>{" "}
                            team. Their API makes professional <span className="text-nowrap">Dota 2</span> data
                            accessible to everyone, and we're committed to building on that foundation.
                        </p>
                    </CardContent>
                </Card>

                {/* Roadmap */}
                <Card id="features-to-add" className="border-border/50">
                    <CardContent className="pt-6">
                        <div className="bg-primary/10 mb-4 flex h-12 w-12 items-center justify-center rounded-lg">
                            <Lightbulb className="text-primary h-6 w-6" />
                        </div>
                        <h2 className="text-foreground mb-4 text-xl font-semibold">Coming soon</h2>
                        <ul className="text-muted-foreground mb-4 space-y-3">
                            <li className="flex items-start gap-3">
                                <span className="bg-primary mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full" />
                                <span>
                                    <strong className="text-foreground">No spoiler mode</strong> – Toggle to hide match
                                    and series results when browsing
                                </span>
                            </li>
                            <li className="flex items-start gap-3">
                                <span className="bg-primary mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full" />
                                <span>
                                    <strong className="text-foreground">Multi-language support</strong> – Expanding
                                    beyond English to serve the global community
                                </span>
                            </li>
                            <li className="flex items-start gap-3">
                                <span className="bg-primary mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full" />
                                <span>
                                    <strong className="text-foreground">Interactive draft mode</strong> – Visualize hero
                                    picks and bans for both series and matches
                                </span>
                            </li>
                        </ul>
                        <p className="text-muted-foreground text-sm">
                            Have a feature request?{" "}
                            <a
                                href="https://github.com/nk1e/dotapro/issues"
                                target="_blank"
                                rel="noopener noreferrer"
                                className="text-blue-500 hover:underline"
                            >
                                Open an issue
                            </a>{" "}
                            and help shape the future of dotapro.
                        </p>
                    </CardContent>
                </Card>

                {/* Contributing */}
                <Card className="border-border/50 from-primary/5 bg-gradient-to-br to-transparent">
                    <CardContent className="pt-6">
                        <h2 className="text-foreground mb-4 text-xl font-semibold">Open to contributions</h2>
                        <p className="text-muted-foreground leading-relaxed">
                            dotapro is open source and contributions are welcome. Whether it's fixing bugs, improving
                            the UI, or adding new features—every pull request makes the platform better for everyone.
                        </p>
                    </CardContent>
                </Card>
            </div>
        </>
    )
}
