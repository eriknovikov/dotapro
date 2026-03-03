import { Navbar } from "@/components/Navbar"
import { Footer } from "@/components/Footer"
import { ErrorBoundary } from "@/components/ErrorBoundary"
import type { QueryClient } from "@tanstack/react-query"
import { createRootRouteWithContext, Outlet, useLocation } from "@tanstack/react-router"
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools"

const RootLayout = () => {
    const location = useLocation()
    const isSeriesPage = location.pathname === "/series"

    return (
        <ErrorBoundary>
            <div className="bg-background flex min-h-screen flex-col overflow-x-hidden">
                <Navbar />
                <div className="flex-1 pt-14 sm:pt-16">
                    <Outlet />
                </div>
                <div className={isSeriesPage ? "lg:ml-72" : ""}>
                    <Footer />
                </div>
                {import.meta.env.DEV && <TanStackRouterDevtools />}
            </div>
        </ErrorBoundary>
    )
}

interface RouterCtx {
    queryClient: QueryClient
}
export const Route = createRootRouteWithContext<RouterCtx>()({ component: RootLayout })
