import { StrictMode } from "react"
import ReactDOM from "react-dom/client"
import { RouterProvider, createRouter } from "@tanstack/react-router"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import "./index.css"
import { routeTree } from "./routeTree.gen"

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: (failureCount, error) => {
                // Retry on server errors (5xx) but not on client errors (4xx)
                if (error instanceof Error) {
                    const msg = error.message
                    // Don't retry on client errors (4xx)
                    if (msg.includes("400") || msg.includes("401") || msg.includes("403") || msg.includes("404")) {
                        return false
                    }
                }
                return failureCount < 3 // Retry up to 3 times for server errors (5xx)
            },
            retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 30000), // Exponential backoff
        },
    },
})

const router = createRouter({
    routeTree,
    context: {
        queryClient,
    },
})

declare module "@tanstack/react-router" {
    interface Register {
        router: typeof router
    }
}

const rootElement = document.getElementById("root")!
if (!rootElement.innerHTML) {
    const root = ReactDOM.createRoot(rootElement)
    root.render(
        <StrictMode>
            <QueryClientProvider client={queryClient}>
                <RouterProvider router={router} />
            </QueryClientProvider>
        </StrictMode>,
    )
}
