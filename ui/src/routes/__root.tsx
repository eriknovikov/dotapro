import type { QueryClient } from "@tanstack/react-query"
import { createRootRouteWithContext, Link, Outlet } from "@tanstack/react-router"
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools"
import { cn } from "../lib/utils"

const RootLayout = () => (
    <>
        <nav className="p-2 flex gap-2 bg-background border-b border-border">
            <Link
                to="/"
                className={cn(
                    "px-3 py-1.5 rounded-md text-sm font-medium transition-colors duration-200",
                    "hover:bg-background-accent",
                    "[&.active]:bg-primary-500",
                    "[&.active]:text-foreground"
                )}
            >
                Home
            </Link>
            <Link
                to="/about"
                className={cn(
                    "px-3 py-1.5 rounded-md text-sm font-medium transition-colors duration-200",
                    "hover:bg-background-accent",
                    "[&.active]:bg-primary-500",
                    "[&.active]:text-foreground"
                )}
            >
                About
            </Link>
        </nav>
        <Outlet />
        <TanStackRouterDevtools />
    </>
)

interface RouterCtx {
    queryClient: QueryClient
}
export const Route = createRootRouteWithContext<RouterCtx>()({ component: RootLayout })
