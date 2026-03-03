import { useLocation } from "@tanstack/react-router"

const SITE_URL = "https://dotapro.org"

/**
 * Hook to programmatically update SEO metadata
 */
export function useSEO() {
    const location = useLocation()

    const updateSEO = () => {
        // This is a no-op in the hook, SEO component handles the updates
        // This hook is provided for type safety and future extensibility
    }

    return {
        updateSEO,
        currentPath: location.pathname,
        currentUrl: `${SITE_URL}${location.pathname}`,
    }
}
