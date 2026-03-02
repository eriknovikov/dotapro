import { Github, DiscordIcon } from "./Icons"

export function Footer() {
    return (
        <footer className="border-t border-border/50 mt-auto">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
                <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
                    <p className="text-center text-foreground-muted text-sm">
                        Made with ❤️ for the <span className="text-nowrap">Dota 2</span> community
                        <span className="mx-2">•</span>
                        <a href="/" className="hover:text-foreground transition-colors">
                            dotapro.com
                        </a>
                        {" © 2026"}
                    </p>
                    <div className="flex items-center gap-6">
                        <a
                            href="https://github.com/E-nkv/dotapro"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-foreground-muted hover:text-foreground transition-all duration-200 hover:scale-110"
                            aria-label="GitHub"
                        >
                            <Github className="h-5 w-5" />
                        </a>
                        <a
                            href="https://discord.gg/h6sVtge8"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-foreground-muted hover:text-foreground transition-all duration-200 hover:scale-110"
                            aria-label="Discord"
                        >
                            <DiscordIcon />
                        </a>
                    </div>
                </div>
            </div>
        </footer>
    )
}
