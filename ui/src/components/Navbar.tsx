import { useState, useEffect } from "react"
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
            className="text-foreground-muted hover:text-foreground px-2 sm:px-3 py-2 text-xs sm:text-sm font-medium transition-colors relative after:content-[''] after:absolute after:bottom-0 after:left-0 after:w-full after:h-0.5 after:bg-primary-500 after:scale-x-0 after:transition-transform after:duration-200 after:origin-center"
            activeProps={{ className: "text-foreground after:scale-x-100" }}
            onClick={onClick}
        >
            {label}
        </Link>
    )
}

export function Navbar() {
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

    // Lock body scroll when menu is open
    useEffect(() => {
        if (isMobileMenuOpen) {
            document.body.style.overflow = 'hidden'
        } else {
            document.body.style.overflow = ''
        }
        return () => {
            document.body.style.overflow = ''
        }
    }, [isMobileMenuOpen])

    // Close menu on escape key
    useEffect(() => {
        const handleEscape = (e: KeyboardEvent) => {
            if (e.key === 'Escape' && isMobileMenuOpen) {
                setIsMobileMenuOpen(false)
            }
        }
        
        document.addEventListener('keydown', handleEscape)
        return () => document.removeEventListener('keydown', handleEscape)
    }, [isMobileMenuOpen])

    return (
        <>
            <nav className="fixed top-0 left-0 right-0 z-50 bg-background/80 backdrop-blur-md border-b border-border/50 shadow-lg shadow-primary-500/5">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex items-center justify-between h-14 sm:h-16">
                        <div className="shrink-0 flex items-center gap-2">
                            <picture>
                                <source srcSet="/logo-48.webp" media="(min-width: 768px)" />
                                <img src="/logo-32x32.webp" alt="" className="h-7 sm:h-8 md:h-10 lg:h-12 w-auto" />
                            </picture>

                            <Link
                                to="/"
                                className="text-2xl sm:text-3xl font-teko bg-linear-to-r from-foreground to-foreground bg-clip-text text-transparent hover:from-[hsl(38,92%,50%)] hover:to-[hsl(0,84%,50%)]"
                            >
                                dotapro.com
                            </Link>
                        </div>

                        <div className="hidden md:block">
                            <div className="flex items-baseline space-x-2 sm:space-x-4">
                                {navLinks.map(link => (
                                    <NavLink key={link.to} to={link.to} label={link.label} />
                                ))}
                            </div>
                        </div>

                        <div className="flex md:hidden">
                            <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => {
                                    // Dispatch event to close filter modal if open
                                    window.dispatchEvent(new CustomEvent('navbar-menu-open'))
                                    setIsMobileMenuOpen(!isMobileMenuOpen)
                                }}
                                aria-label="Toggle menu"
                                aria-expanded={isMobileMenuOpen}
                                className="h-10 w-10 sm:h-12 sm:w-12"
                            >
                                <svg 
                                    className={`h-6 w-6 transition-transform duration-200 ${isMobileMenuOpen ? 'rotate-90' : ''}`} 
                                    fill="none" 
                                    viewBox="0 0 24 24" 
                                    stroke="currentColor"
                                >
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
            </nav>

            {/* Full-screen overlay mobile menu */}
            {isMobileMenuOpen && (
                <div className="fixed inset-0 z-50 md:hidden">
                    {/* Backdrop */}
                    <div 
                        className="absolute inset-0 bg-background/80 backdrop-blur-sm"
                        onClick={() => setIsMobileMenuOpen(false)}
                    />
                    
                    {/* Menu Panel */}
                    <div className="absolute right-0 top-0 bottom-0 w-72 bg-background border-l border-border shadow-xl slide-in-right">
                        <div className="p-6 h-full flex flex-col">
                            {/* Close button */}
                            <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => setIsMobileMenuOpen(false)}
                                className="absolute top-4 right-4 h-10 w-10"
                            >
                                <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth={2}
                                        d="M6 18L18 6M6 6l12 12"
                                    />
                                </svg>
                            </Button>
                            
                            {/* Nav links */}
                            <div className="mt-16 space-y-1 flex-1">
                                {navLinks.map(link => (
                                    <Link
                                        key={link.to}
                                        to={link.to}
                                        onClick={() => setIsMobileMenuOpen(false)}
                                        className="block px-4 py-4 text-lg font-medium rounded-lg hover:bg-background-accent transition-colors"
                                        activeProps={{ className: "text-primary-500 bg-background-accent" }}
                                    >
                                        {link.label}
                                    </Link>
                                ))}
                            </div>

                            {/* GitHub link */}
                            <div className="pt-6 border-t border-border">
                                <a
                                    href="https://GITHUBPROJECT"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="flex items-center gap-3 px-4 py-3 text-foreground-muted hover:text-foreground hover:bg-background-accent rounded-lg transition-colors"
                                >
                                    <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                                        <path fillRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clipRule="evenodd" />
                                    </svg>
                                    <span className="font-medium">GitHub</span>
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </>
    )
}
