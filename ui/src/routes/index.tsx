import { createFileRoute, Link } from "@tanstack/react-router"
import { Button } from "@/components/ui/button"

export const Route = createFileRoute("/")({
    component: Index,
})

function Index() {
    return (
        <main className="flex-1 min-h-[calc(100vh-4rem)]">
            <div className="max-w-7xl mx-auto px-6 py-12">
                {/* Hero Section */}
                <div className="flex justify-center mb-16">
                    <div className="flex flex-col items-start">
                        <h1 className="text-5xl md:text-7xl font-bold text-foreground mb-6">
                            Dota 2 analytics that actually help
                        </h1>
                        <p className="text-xl text-foreground-muted mb-10">
                            Level up your game. Enjoy and learn from the pro-scene. Completely for free.
                        </p>

                        {/* CTA Buttons */}
                        <div className="flex flex-col sm:flex-row gap-4">
                            <Button asChild variant="primary" size="lg" className="text-base px-8">
                                <Link to="/series">View Series</Link>
                            </Button>
                            <Button asChild variant="outline" size="lg" className="text-base px-8">
                                <Link to="/about">Read Guide</Link>
                            </Button>
                        </div>
                    </div>
                </div>
            </div>

            {/* Hero Image with Linear-style 3D effect - Full width outside container */}
            <div className="flex items-center justify-center py-12 perspective:[2000px]">
                <div className="relative transition-all duration-1000 ease-in-out transform-[rotateX(-25deg)_rotateY(10deg)_rotateZ(-1deg)] hover:transform-[rotateX(-15deg)_rotateY(7deg)_rotateZ(0deg)]">
                    {/* Ambient Glow (Light Bleed) */}
                    <div className="absolute -inset-4 bg-primary-500/10 blur-3xl opacity-30" aria-hidden="true" />

                    {/* Main Image Container */}
                    <div className="relative rounded-xl border border-white/10 bg-background-card shadow-[0_0_0_1px_rgba(0,0,0,.04),0_4px_6px_rgba(0,0,0,.02),0_10px_15px_rgba(0,0,0,.015),0_20px_25px_rgba(0,0,0,.015),0_40px_50px_rgba(0,0,0,.01),0_100px_80px_rgba(0,0,0,.01)] overflow-hidden">
                        {/* Top Edge Highlight (Simulates overhead light) */}
                        <div className="absolute inset-0 border-t border-white/20 pointer-events-none" />

                        <img src="/scrsh.png" alt="Series Screenshot" className="block w-full max-w-[80vw] h-auto" />
                    </div>
                </div>
            </div>
        </main>
    )
}
