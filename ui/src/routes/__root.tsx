import { Navbar } from "@/components/Navbar"
import type { QueryClient } from "@tanstack/react-query"
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router"
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools"

const RootLayout = () => (
    <div className="min-h-screen bg-background">
        <Navbar />
        <Outlet />
        <TanStackRouterDevtools />
    </div>
)

interface RouterCtx {
    queryClient: QueryClient
}
export const Route = createRootRouteWithContext<RouterCtx>()({ component: RootLayout })
