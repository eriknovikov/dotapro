import { SEO } from "@/components"
import { createFileRoute } from "@tanstack/react-router"

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
            <div className="mx-auto min-h-[calc(100vh-4rem)] max-w-4xl py-6">
                <div className="space-y-8">
                    {/* What's dotapro */}
                    <div className="p-10 shadow-md shadow-gray-900">
                        <h2 className="text-foreground mb-4 text-2xl font-bold">What's dotapro?</h2>
                        <p className="text-foreground-muted mb-4 leading-relaxed">
                            It's an open source, completely free-to-use platform about professional{" "}
                            <span className="text-nowrap">Dota 2</span>. No random pub games here – just the good stuff
                            from the pro scene.
                        </p>
                        <p className="text-foreground-muted leading-relaxed">
                            It's by no means a replacement or better alternative of traditional heavy-weight analyitics
                            sites like Dotabuff, Stratz, or OpenDota; it's just that dotapro's focus is completely
                            different. Other platforms just dump a mountain of data on you and say "good luck figuring
                            it out", while dotapro focuses on the analytics that actually matter to Dota players
                            (according to real players), without the noise of public matches data.
                        </p>
                    </div>

                    {/* Who is dotapro for */}
                    <div className="p-10 shadow-md shadow-gray-900">
                        <h2 className="text-foreground mb-4 text-2xl font-bold">Who is dotapro for?</h2>
                        <p className="text-foreground-muted mb-4 leading-relaxed">
                            Pretty much anyone who likes watching pro Dota. But specifically:
                        </p>
                        <ul className="text-foreground-muted ml-6 list-disc space-y-2">
                            <li>
                                <em>Casual fans</em> – Just want to see what your favorite teams are up to without
                                drowning in stats
                            </li>
                            <li>
                                <em>Casters & analysts</em> – Quick access to match history and team performance for
                                prep work
                            </li>
                            <li>
                                <em>Pro players & coaches</em> – Scout opponents, analyze strategies, see what's meta in
                                the pro scene
                            </li>
                            <li>
                                <em>Dota nerds</em> – People who just love looking at match data (I see you)
                            </li>
                        </ul>
                    </div>

                    {/* Credits */}
                    <div className="p-10 shadow-md shadow-gray-900">
                        <h2 className="text-foreground mb-4 text-2xl font-bold">Credits</h2>
                        <p className="text-foreground-muted leading-relaxed">
                            Gotta give credit where it's due – dotapro is built on top of the amazing OpenDota API.
                            Those folks have done incredible work making <span className="text-nowrap">Dota 2</span>{" "}
                            data accessible to everyone, and this project wouldn't exist without them. Big thanks to the
                            OpenDota team! 🙌
                        </p>
                    </div>

                    {/* About the Developer */}
                    <div className="p-10 shadow-md shadow-gray-900">
                        <h2 className="text-foreground mb-4 text-2xl font-bold">About the Developer</h2>
                        <p className="text-foreground-muted mb-4 leading-relaxed">
                            Hi there! I'm Erik, the creator of dotapro.
                        </p>
                        <p className="text-foreground-muted leading-relaxed">
                            As a <span className="text-nowrap">Dota 2</span> fan and player myself, it was quite
                            difficult to follow the pro-scene with traditional platforms, and this also seemed to be an
                            issue for Pro Players I knew, so I decided to build 'dotapro'. The tech stack is pretty
                            modern (Go, React, PostgreSQL, AWS), and it has been super exciting to work on something I
                            love.
                        </p>
                    </div>

                    {/* Contributing */}
                    <div className="p-10 shadow-md shadow-gray-900">
                        <h2 className="text-foreground mb-4 text-2xl font-bold">Contributing</h2>
                        <p className="text-foreground-muted mb-4 leading-relaxed">
                            dotapro.org is open source, and I would appreciate your help greatly! You could:
                        </p>

                        <ul className="text-foreground-muted mb-4 ml-6 list-disc space-y-2">
                            <li>Submit pull requests for bug fixes or new features</li>
                            <li>Open issues for bugs you find or features you'd like to see</li>
                            <li>Share feedback and suggestions (I'm always open to ideas!)</li>
                        </ul>
                    </div>

                    {/* Currently working on */}
                    <div id="features-to-add" className="p-10 shadow-md shadow-gray-900">
                        <h2 className="text-foreground mb-4 text-2xl font-bold">Features on the way</h2>

                        <ul className="text-foreground-muted mb-4 ml-6 list-disc space-y-2">
                            <li>
                                A "no spoiler mode", so that match and series results are turned off with a toggle in
                                the filters
                            </li>
                            <li>Allowing for multiple languages instead of english only</li>
                            <li>Adding an interactive draft mode for both the series and matches</li>
                        </ul>
                        <p className="text-foreground-muted text-sm">
                            The features are listed in the order of intended implementation. If you think you would want
                            a feature more than the other, or want something totally different to be added to dotapro,
                            let us know!
                        </p>
                    </div>
                </div>
            </div>
        </>
    )
}
