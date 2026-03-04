import { DiscordIcon } from "./Icons"

export function Footer() {
    return (
        <footer className="border-border/50 mt-auto w-full border-t">
            <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
                <div className="text-foreground-muted text-center text-sm">
                    <p className="text-nowrap">
                        <a href="/" className="hover:text-foreground transition-colors">
                            dotapro.org
                        </a>
                        {" © 2026"}
                        <a
                            href="https://github.com/E-nkv/dotapro"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-foreground-muted hover:text-foreground ml-4 inline-block align-middle transition-all duration-200 hover:scale-110"
                            aria-label="GitHub"
                        >
                            <img src="/github.svg" alt="GitHub" className="h-5 w-5 brightness-0 invert" />
                        </a>
                        <a
                            href="https://discord.gg/h6sVtge8"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-foreground-muted hover:text-foreground ml-4 inline-block align-middle transition-all duration-200 hover:scale-110"
                            aria-label="Discord"
                        >
                            <DiscordIcon />
                        </a>
                    </p>
                </div>
            </div>
        </footer>
    )
}
