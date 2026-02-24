import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/about")({
    component: About,
})

function About() {
    return (
        <div className="p-6 min-h-[calc(100vh-4rem)] max-w-4xl mx-auto">
            <div className="space-y-8">
                {/* What's dotapro */}
                <div className="shadow-md shadow-gray-900 p-10">
                    <h2 className="text-2xl font-bold text-foreground mb-4">What's dotapro?</h2>
                    <p className="text-foreground-muted leading-relaxed mb-4">
                        It's an open source, completely free-to-use platform about professional{" "}
                        <span className="text-nowrap">Dota 2</span>. No random pub games here – just the good stuff from
                        the pro scene.
                    </p>
                    <p className="text-foreground-muted leading-relaxed">
                        The key reason for building this was that traditional Dota 2 analytics platforms just dump a
                        mountain of data on you and say "good luck figuring it out". dotapro focuses on the analytics
                        that actually matter to Dota players, presented in a clean, modern UI that doesn't make your
                        eyes bleed. Think of it as OpenDota/Stratz/Liquipedia, but on steroids.
                    </p>
                </div>

                {/* Who is dotapro for */}
                <div className="shadow-md shadow-gray-900 p-10">
                    <h2 className="text-2xl font-bold text-foreground mb-4">Who is dotapro for?</h2>
                    <p className="text-foreground-muted leading-relaxed mb-4">
                        Pretty much anyone who likes watching pro Dota. But specifically:
                    </p>
                    <ul className="text-foreground-muted space-y-2 ml-6 list-disc">
                        <li>
                            <em>Casual fans</em> – Just want to see what your favorite teams are up to without drowning
                            in stats
                        </li>
                        <li>
                            <em>Casters & analysts</em> – Quick access to match history and team performance for prep
                            work
                        </li>
                        <li>
                            <em>Pro players & coaches</em> – Scout opponents, analyze strategies, see what's meta in the
                            pro scene
                        </li>
                        <li>
                            <em>Dota nerds</em> – People who just love looking at match data (I see you)
                        </li>
                    </ul>
                </div>

                {/* Disclaimer */}
                <div className="shadow-md shadow-gray-900 p-10">
                    <h2 className="text-2xl font-bold text-foreground mb-4">Disclaimer</h2>
                    <p className="text-foreground-muted leading-relaxed">
                        Gotta give credit where it's due – dotapro is built on top of the amazing OpenDota API. Those
                        folks have done incredible work making <span className="text-nowrap">Dota 2</span> data
                        accessible to everyone, and this project wouldn't exist without them. Big thanks to the OpenDota
                        team! 🙌
                    </p>
                </div>

                {/* About the Developer */}
                <div className="shadow-md shadow-gray-900 p-10">
                    <h2 className="text-2xl font-bold text-foreground mb-4">About the Developer</h2>
                    <p className="text-foreground-muted leading-relaxed mb-4">
                        Hi there! I'm Erik, the creator of dotapro.
                    </p>
                    <p className="text-foreground-muted leading-relaxed">
                        As a <span className="text-nowrap">Dota 2</span> fan and player myself, it was quite difficult
                        to follow the pro-scene with traditional platforms, and this also seemed to be an issue for Pro
                        Players I knew, so I decided to build 'dotapro'. The tech stack is pretty modern (Go, React,
                        PostgreSQL, AWS), and it has been super exciting to work on something I love.
                    </p>
                </div>

                {/* Contributing */}
                <div className="shadow-md shadow-gray-900 p-10">
                    <h2 className="text-2xl font-bold text-foreground mb-4">Contributing</h2>
                    <p className="text-foreground-muted leading-relaxed mb-4">
                        dotapro is open source and I'd love to have you contribute! Whether you're fixing bugs, adding
                        features, or just have ideas for improvements – all help is appreciated.
                    </p>
                    <p className="text-foreground-muted leading-relaxed mb-4">
                        <strong className="text-foreground">Ways to contribute:</strong>
                    </p>
                    <ul className="text-foreground-muted space-y-2 ml-6 list-disc mb-4">
                        <li>Submit pull requests for bug fixes or new features</li>
                        <li>Open issues for bugs you find or features you'd like to see</li>
                        <li>Share feedback and suggestions (I'm always open to ideas!)</li>
                    </ul>
                    <p className="text-foreground-muted leading-relaxed">
                        The project will be on GitHub once I get it properly public-ready. In the meantime, feel free to
                        reach out if you want to contribute or just chat about the project. I'm pretty responsive and
                        happy to help new contributors get started!
                    </p>
                </div>

                {/* Footer */}
                <div className="text-center text-foreground-muted text-sm pt-4">
                    <p>
                        Made with ❤️ for the <span className="text-nowrap">Dota 2</span> community
                    </p>
                </div>
            </div>
        </div>
    )
}
