import { useEffect } from "react"
import { useLocation } from "@tanstack/react-router"

interface SEOProps {
    title?: string
    description?: string
    image?: string
    type?: "website" | "article"
    noIndex?: boolean
    canonicalUrl?: string
}

const DEFAULT_TITLE = "dotapro.com - Dota 2 professional analytics"
const DEFAULT_DESCRIPTION =
    "Comprehensive Dota 2 match data aggregation system. View series and matches statistics from the professional scene."
const DEFAULT_IMAGE = "/og-image.webp"
const SITE_URL = "https://dotapro.com"

/**
 * SEO component for managing document title and meta tags dynamically
 * This component updates the document head with proper SEO metadata
 */
export function SEO({
    title,
    description,
    image = DEFAULT_IMAGE,
    type = "website",
    noIndex = false,
    canonicalUrl,
}: SEOProps) {
    const location = useLocation()

    useEffect(() => {
        // Build full URL
        const fullUrl = canonicalUrl || `${SITE_URL}${location.pathname}`

        // Set document title
        document.title = title ? `${title} | dotapro.com` : DEFAULT_TITLE

        // Update or create meta tags
        const updateMetaTag = (name: string, content: string, property = false) => {
            const selector = property ? `meta[property="${name}"]` : `meta[name="${name}"]`
            let meta = document.querySelector(selector) as HTMLMetaElement | null

            if (!meta) {
                meta = document.createElement("meta")
                if (property) {
                    meta.setAttribute("property", name)
                } else {
                    meta.setAttribute("name", name)
                }
                document.head.appendChild(meta)
            }

            meta.setAttribute("content", content)
        }

        // Primary Meta Tags
        updateMetaTag("description", description || DEFAULT_DESCRIPTION)
        updateMetaTag("keywords", "Dota 2, matches, statistics, teams, leagues, esports, gaming")

        // Open Graph / Facebook
        updateMetaTag("og:type", type, true)
        updateMetaTag("og:url", fullUrl, true)
        updateMetaTag("og:title", title || "Dotapro - Dota 2 Match Data & Statistics", true)
        updateMetaTag("og:description", description || DEFAULT_DESCRIPTION, true)
        updateMetaTag("og:image", image, true)
        updateMetaTag("og:image:width", "1200", true)
        updateMetaTag("og:image:height", "630", true)
        updateMetaTag("og:image:alt", "Dotapro Logo", true)

        // Twitter
        updateMetaTag("twitter:card", "summary_large_image")
        updateMetaTag("twitter:url", fullUrl, true)
        updateMetaTag("twitter:title", title || "Dotapro - Dota 2 Match Data & Statistics", true)
        updateMetaTag("twitter:description", description || DEFAULT_DESCRIPTION, true)
        updateMetaTag("twitter:image", image, true)

        // Canonical URL
        let canonical = document.querySelector('link[rel="canonical"]') as HTMLLinkElement | null
        if (!canonical) {
            canonical = document.createElement("link")
            canonical.setAttribute("rel", "canonical")
            document.head.appendChild(canonical)
        }
        canonical.setAttribute("href", fullUrl)

        // No Index
        if (noIndex) {
            let robots = document.querySelector('meta[name="robots"]') as HTMLMetaElement | null
            if (!robots) {
                robots = document.createElement("meta")
                robots.setAttribute("name", "robots")
                document.head.appendChild(robots)
            }
            robots.setAttribute("content", "noindex, nofollow")
        }

        // Cleanup function to remove meta tags on unmount
        return () => {
            // Note: We don't remove meta tags on unmount as they should persist
            // This is intentional for SEO purposes
        }
    }, [title, description, image, type, noIndex, canonicalUrl, location.pathname])

    return null // This component doesn't render anything
}

/**
 * Hook to programmatically update SEO metadata
 */
export function useSEO() {
    const location = useLocation()

    const updateSEO = (props: SEOProps) => {
        // This is a no-op in the hook, the SEO component handles the updates
        // This hook is provided for type safety and future extensibility
    }

    return {
        updateSEO,
        currentPath: location.pathname,
        currentUrl: `${SITE_URL}${location.pathname}`,
    }
}
