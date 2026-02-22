import { useState } from "react"
import { Link } from "@tanstack/react-router"
import { Button } from "./ui/index"

const navLinks = [
    { to: "/", label: "home" },
    { to: "/series", label: "series" },
    { to: "/about", label: "about" },
    { to: "/guide", label: "guide" },
] as const

function NavLink({ to, label, onClick }: { to: string; label: string; onClick?: () => void }) {
    return (
        <Link
            to={to}
            className="text-foreground-muted hover:text-foreground px-3 py-2 text-sm font-medium transition-colors relative after:content-[''] after:absolute after:bottom-0 after:left-0 after:w-full after:h-0.5 after:bg-primary-500 after:scale-x-0 after:transition-transform after:duration-200 after:origin-center"
            activeProps={{ className: "text-foreground after:scale-x-100" }}
            onClick={onClick}
        >
            {label}
        </Link>
    )
}

export function Navbar() {
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

    return (
        <nav className="sticky top-0 z-50 bg-background/80 backdrop-blur-md border-b border-border/50 shadow-lg shadow-primary-500/5">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    <div className="flex-shrink-0 flex items-center gap-2">
                        <picture>
                            <source srcSet="/logo-48.webp" media="(min-width: 768px)" />
                            <img src="/logo-32x32.webp" alt="" className="h-8 md:h-12 w-auto" />
                        </picture>

                        <Link
                            to="/"
                            className="text-3xl font-teko bg-gradient-to-r from-foreground to-foreground bg-clip-text text-transparent hover:from-[hsl(38,92%,50%)] hover:to-[hsl(0,84%,50%)]"
                        >
                            dotapro.com
                        </Link>
                    </div>

                    <div className="hidden md:block">
                        <div className="ml-10 flex items-baseline space-x-4">
                            {navLinks.map(link => (
                                <NavLink key={link.to} to={link.to} label={link.label} />
                            ))}
                        </div>
                    </div>

                    <div className="-mr-2 flex md:hidden">
                        <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                            aria-label="Toggle menu"
                            aria-expanded={isMobileMenuOpen}
                        >
                            <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                {isMobileMenuOpen ? (
                                    <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth={2}
                                        d="M6 18L18 6M6 6l12 12"
                                    />
                                ) : (
                                    <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth={2}
                                        d="M4 6h16M4 12h16M4 18h16"
                                    />
                                )}
                            </svg>
                        </Button>
                    </div>
                </div>
            </div>

            {isMobileMenuOpen && (
                <div className="md:hidden bg-background/95 backdrop-blur-md border-t border-border/50">
                    <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3">
                        {navLinks.map(link => (
                            <NavLink
                                key={link.to}
                                to={link.to}
                                label={link.label}
                                onClick={() => setIsMobileMenuOpen(false)}
                            />
                        ))}
                    </div>
                </div>
            )}
        </nav>
    )
}
