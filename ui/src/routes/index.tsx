import { Button, SEO } from "@/components"
import { createFileRoute, Link } from "@tanstack/react-router"

export const Route = createFileRoute("/")({
    component: Index,
})

function Index() {
    return (
        <>
            <SEO />
            <main className="min-h-0 flex-1 overflow-x-hidden sm:min-h-[calc(100vh-4rem)]">
                <div className="mx-auto max-w-7xl px-3 py-8 sm:px-4 sm:py-4 md:px-6 md:py-8 lg:px-8 lg:py-12">
                    {/* Hero Section */}
                    <div className="mb-8 flex justify-center sm:mb-6 md:mb-12 lg:mb-16">
                        <div className="flex w-full max-w-2xl flex-col items-start px-1 sm:px-2">
                            <h1 className="text-foreground mb-2 text-4xl font-bold sm:mb-4 sm:text-4xl md:mb-6 md:text-4xl lg:text-5xl xl:text-7xl">
                                Dota 2 analytics that actually help
                            </h1>
                            <p className="text-foreground-muted mb-3 text-lg sm:mb-6 sm:text-lg md:mb-10 md:text-lg lg:text-xl">
                                Level up your game. Enjoy and learn from the pro-scene. Completely for free.
                            </p>

                            {/* CTA Buttons */}
                            <div className="flex w-full flex-row items-center gap-2 sm:w-auto sm:gap-3 md:gap-4">
                                <Button
                                    asChild
                                    variant="primary"
                                    size="sm"
                                    className="px-6 text-sm sm:px-6 sm:text-sm md:px-6 md:text-sm lg:px-8 lg:text-base"
                                >
                                    <Link to="/series">View Series</Link>
                                </Button>
                                <Button
                                    asChild
                                    variant="cool-outline"
                                    size="sm"
                                    className="px-6 text-sm sm:px-6 sm:text-sm md:px-6 md:text-sm lg:px-8 lg:text-base"
                                >
                                    <Link to="/matches">View Matches</Link>
                                </Button>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Hero Image with Linear-style 3D effect - Full width outside container */}
                <div className="sm:perspective:[2000px] bg-gradient-mesh relative flex w-full items-center justify-center overflow-hidden px-2 py-1 perspective-[2000px] sm:py-2 md:py-8 lg:py-12">
                    {/* Scanline effect overlay */}
                    <div className="scanline-effect pointer-events-none absolute inset-0 z-10" aria-hidden="true" />

                    <div className="hero-float relative w-full transform-[rotateX(-18deg)_rotateY(7deg)_rotateZ(-0.75deg)] transition-shadow duration-300 ease-in-out hover:shadow-[0_0_0_1px_rgba(0,0,0,.06),0_16px_24px_rgba(0,0,0,.08),0_40px_60px_rgba(0,0,0,.06),0_80px_100px_rgba(0,0,0,.06),0_160px_200px_rgba(0,0,0,.04),0_400px_320px_rgba(0,0,0,.04)] sm:w-auto sm:transform-[rotateX(-25deg)_rotateY(10deg)_rotateZ(-1deg)]">
                        {/* Ambient Glow (Light Bleed) - with Radiant gold tint */}
                        <div
                            className="bg-primary-500/10 absolute -inset-2 opacity-30 blur-3xl sm:-inset-4"
                            aria-hidden="true"
                        />

                        {/* Main Image Container */}
                        <div className="bg-background-card relative overflow-hidden rounded-xl shadow-[0_0_0_1px_rgba(0,0,0,.04),0_4px_6px_rgba(0,0,0,.02),0_10px_15px_rgba(0,0,0,.015),0_20px_25px_rgba(0,0,0,.015),0_40px_50px_rgba(0,0,0,.01),0_100px_80px_rgba(0,0,0,.01)]">
                            <img
                                src="/hero.webp"
                                alt="Series Screenshot"
                                width="2525"
                                height="1222"
                                className="block h-auto w-full sm:max-w-[80vw] md:max-w-[90vw] lg:max-w-[80vw]"
                            />
                        </div>
                    </div>
                </div>
            </main>
        </>
    )
}
