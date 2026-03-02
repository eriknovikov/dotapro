import { useState, useEffect } from "react"
import { Link } from "@tanstack/react-router"
import { Button } from "./ui/index"
import { Github, DiscordIcon } from "./Icons"

const navLinks = [
    { to: "/", label: "home", ariaLabel: "Home" },
    { to: "/series", label: "series", ariaLabel: "Series" },
    { to: "/about", label: "about", ariaLabel: "About" },
    { to: "/guide", label: "guide", ariaLabel: "Guide" },
] as const

function NavLink({ to, label, onClick, ariaLabel }: { to: string; label: string; onClick?: () => void; ariaLabel?: string }) {
    return (
        <Link
            to={to}
            className="text-foreground-muted hover:text-foreground px-2 sm:px-3 py-2 text-xs sm:text-sm font-medium transition-colors relative after:content-[''] after:absolute after:bottom-0 after:left-0 after:w-full after:h-0.5 after:bg-primary-500 after:scale-x-0 after:transition-transform after:duration-200 after:origin-center"
            activeProps={{ className: "text-foreground after:scale-x-100" }}
            onClick={onClick}
            aria-label={ariaLabel || label}
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
            <nav className="fixed top-0 left-0 right-0 z-50 bg-background/80 backdrop-blur-md border-b border-border/50 shadow-lg shadow-primary-500/5" role="navigation" aria-label="Main navigation">
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
                            <div className="flex items-center space-x-2 sm:space-x-4" role="menubar">
                                {navLinks.map(link => (
                                    <NavLink key={link.to} to={link.to} label={link.label} ariaLabel={link.ariaLabel} />
                                ))}
                                <a
                                    href="https://github.com/E-nkv/dotapro"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-foreground-muted hover:text-foreground transition-all duration-200 hover:scale-110 px-2 sm:px-3 py-2 flex items-center"
                                    aria-label="GitHub"
                                >
                                    <Github className="h-5 w-5" />
                                </a>
                                <a
                                    href="https://discord.gg/h6sVtge8"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-foreground-muted hover:text-foreground transition-all duration-200 hover:scale-110 px-2 sm:px-3 py-2 flex items-center"
                                    aria-label="Discord"
                                >
                                    <DiscordIcon />
                                </a>
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
                <div className="fixed inset-0 z-50 md:hidden" role="dialog" aria-modal="true" aria-label="Mobile navigation menu">
                    {/* Backdrop */}
                    <div
                        className="absolute inset-0 bg-background/80 backdrop-blur-sm"
                        onClick={() => setIsMobileMenuOpen(false)}
                        aria-hidden="true"
                    />
                    
                    {/* Menu Panel */}
                    <div className="absolute right-0 top-0 bottom-0 w-72 bg-background border-l border-border shadow-xl slide-in-right" role="menu">
                        <div className="p-6 h-full flex flex-col">
                            {/* Close button */}
                            <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => setIsMobileMenuOpen(false)}
                                className="absolute top-4 right-4 h-10 w-10"
                                aria-label="Close menu"
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
                            <div className="mt-16 space-y-1 flex-1" role="menu">
                                {navLinks.map(link => (
                                    <Link
                                        key={link.to}
                                        to={link.to}
                                        onClick={() => setIsMobileMenuOpen(false)}
                                        className="block px-4 py-4 text-lg font-medium rounded-lg hover:bg-background-accent transition-colors"
                                        activeProps={{ className: "text-primary-500 bg-background-accent" }}
                                        role="menuitem"
                                        aria-label={link.ariaLabel}
                                    >
                                        {link.label}
                                    </Link>
                                ))}
                            </div>

                            {/* GitHub link */}
                            <div className="pt-6 border-t border-border">
                                <a
                                    href="https://github.com/E-nkv/dotapro"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="flex items-center gap-3 px-4 py-3 text-foreground-muted hover:text-foreground hover:bg-background-accent rounded-lg transition-colors"
                                    aria-label="View source code on GitHub"
                                >
                                    <Github className="h-5 w-5" />
                                    <span className="font-medium">GitHub</span>
                                </a>
                                <a
                                    href="https://discord.gg/h6sVtge8"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="flex items-center gap-3 px-4 py-3 text-foreground-muted hover:text-foreground hover:bg-background-accent rounded-lg transition-colors"
                                    aria-label="Join Discord server"
                                >
                                    <DiscordIcon />
                                    <span className="font-medium">Discord</span>
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </>
    )
}
