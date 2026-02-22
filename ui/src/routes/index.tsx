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
                <div className="flex flex-col items-center text-center mb-16">
                    <h1 className="text-5xl md:text-7xl font-bold text-foreground mb-6">
                        Dota Pro Analytics
                    </h1>
                    <p className="text-xl text-foreground-muted max-w-2xl mb-10">
                        Explore professional Dota 2 matches and series with powerful filtering and detailed statistics.
                    </p>
                    
                    {/* CTA Buttons */}
                    <div className="flex flex-col sm:flex-row gap-4">
                        <Button asChild size="lg" className="text-base px-8">
                            <Link to="/series">View Series</Link>
                        </Button>
                        <Button asChild size="lg" variant="outline" className="text-base px-8">
                            <Link to="/about">Read Guide</Link>
                        </Button>
                    </div>
                </div>

                {/* Hero Image with Linear-style 3D effect */}
                <div className="flex items-center justify-center py-12 [perspective:2000px]">
                    <div className="relative transition-all duration-700 ease-out [transform:rotateX(25deg)_rotateY(-10deg)_rotateZ(1deg)] hover:[transform:rotateX(15deg)_rotateY(-5deg)_rotateZ(0deg)]">
                        
                        {/* Ambient Glow (Light Bleed) */}
                        <div className="absolute -inset-4 bg-primary-500/20 blur-3xl opacity-50" aria-hidden="true" />

                        {/* Main Image Container */}
                        <div className="relative rounded-xl border border-white/10 bg-background-card shadow-linear overflow-hidden">
                            
                            {/* Top Edge Highlight (Simulates overhead light) */}
                            <div className="absolute inset-0 border-t border-white/20 pointer-events-none" />
                            
                            <img 
                                src="/scrsh.png" 
                                alt="Series Screenshot" 
                                className="block w-full max-w-5xl h-auto"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </main>
    )
}
