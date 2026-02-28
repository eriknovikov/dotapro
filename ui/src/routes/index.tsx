import { createFileRoute, Link } from "@tanstack/react-router"
import { Button } from "@/components/ui"
import { SEO } from "@/components/SEO"

export const Route = createFileRoute("/")({
    component: Index,
})

function Index() {
    return (
        <>
            <SEO />
            <main className="flex-1 min-h-0 sm:min-h-[calc(100vh-4rem)] overflow-x-hidden">
                <div className="max-w-7xl mx-auto py-8 sm:py-4 md:py-8 lg:py-12 px-3 sm:px-4 md:px-6 lg:px-8">
                    {/* Hero Section */}
                    <div className="flex justify-center mb-8 sm:mb-6 md:mb-12 lg:mb-16">
                        <div className="flex flex-col items-start w-full max-w-2xl px-1 sm:px-2">
                            <h1 className="text-4xl sm:text-4xl md:text-4xl lg:text-5xl xl:text-7xl font-bold text-foreground mb-2 sm:mb-4 md:mb-6">
                                Dota 2 analytics that actually help
                            </h1>
                            <p className="text-lg sm:text-lg md:text-lg lg:text-xl text-foreground-muted mb-3 sm:mb-6 md:mb-10">
                                Level up your game. Enjoy and learn from the pro-scene. Completely for free.
                            </p>

                            {/* CTA Buttons */}
                            <div className="flex flex-row gap-2 sm:gap-3 md:gap-4 w-full sm:w-auto">
                                <Button
                                    asChild
                                    variant="primary"
                                    size="sm"
                                    className="text-sm sm:text-sm md:text-sm lg:text-base px-6 sm:px-6 md:px-6 lg:px-8 flex-1 sm:flex-none"
                                >
                                    <Link to="/series">View Series</Link>
                                </Button>
                                <Button
                                    asChild
                                    variant="outline"
                                    size="sm"
                                    className="text-sm sm:text-sm md:text-sm lg:text-base px-6 sm:px-6 md:px-6 lg:px-8 flex-1 sm:flex-none"
                                >
                                    <Link to="/about">About dotapro</Link>
                                </Button>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Hero Image with Linear-style 3D effect - Full width outside container */}
                <div className="flex items-center justify-center py-1 sm:py-2 md:py-8 lg:py-12 perspective-[2000px] sm:perspective:[2000px] w-full overflow-hidden px-2">
                    <div className="relative transition-shadow duration-300 ease-in-out transform-[rotateX(-18deg)_rotateY(7deg)_rotateZ(-0.75deg)] sm:transform-[rotateX(-25deg)_rotateY(10deg)_rotateZ(-1deg)] hover:shadow-[0_0_0_1px_rgba(0,0,0,.06),0_16px_24px_rgba(0,0,0,.08),0_40px_60px_rgba(0,0,0,.06),0_80px_100px_rgba(0,0,0,.06),0_160px_200px_rgba(0,0,0,.04),0_400px_320px_rgba(0,0,0,.04)] w-full sm:w-auto">
                        {/* Ambient Glow (Light Bleed) */}
                        <div
                            className="absolute -inset-2 sm:-inset-4 bg-primary-500/10 blur-3xl opacity-30"
                            aria-hidden="true"
                        />

                        {/* Main Image Container */}
                        <div className="relative rounded-xl bg-background-card shadow-[0_0_0_1px_rgba(0,0,0,.04),0_4px_6px_rgba(0,0,0,.02),0_10px_15px_rgba(0,0,0,.015),0_20px_25px_rgba(0,0,0,.015),0_40px_50px_rgba(0,0,0,.01),0_100px_80px_rgba(0,0,0,.01)] overflow-hidden">
                            <img
                                src="/hero.webp"
                                alt="Series Screenshot"
                                className="block w-full sm:max-w-[80vw] md:max-w-[90vw] lg:max-w-[80vw] h-auto"
                            />
                        </div>
                    </div>
                </div>
            </main>
        </>
    )
}
