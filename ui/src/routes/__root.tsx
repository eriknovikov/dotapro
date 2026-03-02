import { Navbar } from "@/components/Navbar"
import { Footer } from "@/components/Footer"
import { ErrorBoundary } from "@/components/ErrorBoundary"
import type { QueryClient } from "@tanstack/react-query"
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router"
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools"

const RootLayout = () => (
    <ErrorBoundary>
        <div className="min-h-screen bg-background flex flex-col overflow-x-hidden">
            <Navbar />
            <div className="pt-14 sm:pt-16 flex-1">
                <Outlet />
            </div>
            <Footer />
            {import.meta.env.DEV && <TanStackRouterDevtools />}
        </div>
    </ErrorBoundary>
)

interface RouterCtx {
    queryClient: QueryClient
}
export const Route = createRootRouteWithContext<RouterCtx>()({ component: RootLayout })
