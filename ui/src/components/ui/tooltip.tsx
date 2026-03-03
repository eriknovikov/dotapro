import * as React from "react"
import { cn } from "@/lib/utils"

interface TooltipProps {
    children: React.ReactNode
    content: React.ReactNode
    side?: "top" | "right" | "bottom" | "left"
    sideOffset?: number
    className?: string
    contentClassName?: string
}

export function Tooltip({
    children,
    content,
    side = "right",
    sideOffset = 4,
    className,
    contentClassName,
}: TooltipProps) {
    const [open, setOpen] = React.useState(false)
    const [position, setPosition] = React.useState<{ top: number; left: number } | null>(null)
    const triggerRef = React.useRef<HTMLDivElement>(null)
    const contentRef = React.useRef<HTMLDivElement>(null)
    const hasHover = React.useRef(false)

    // Detect if device has hover capability
    React.useEffect(() => {
        hasHover.current = window.matchMedia("(hover: hover)").matches
    }, [])

    const calculatePosition = React.useCallback(() => {
        if (!triggerRef.current || !contentRef.current) return

        const triggerRect = triggerRef.current.getBoundingClientRect()
        const contentRect = contentRef.current.getBoundingClientRect()

        let top = 0
        let left = 0

        switch (side) {
            case "top":
                top = triggerRect.top - contentRect.height - sideOffset
                left = triggerRect.left + (triggerRect.width - contentRect.width) / 2
                break
            case "bottom":
                top = triggerRect.bottom + sideOffset
                left = triggerRect.left + (triggerRect.width - contentRect.width) / 2
                break
            case "left":
                top = triggerRect.top + (triggerRect.height - contentRect.height) / 2
                left = triggerRect.left - contentRect.width - sideOffset
                break
            case "right":
            default:
                top = triggerRect.top + (triggerRect.height - contentRect.height) / 2
                left = triggerRect.right + sideOffset
                break
        }

        // Keep within viewport bounds
        const padding = 8
        const viewportWidth = window.innerWidth
        const viewportHeight = window.innerHeight

        if (left < padding) left = padding
        if (left + contentRect.width > viewportWidth - padding) left = viewportWidth - contentRect.width - padding
        if (top < padding) top = padding
        if (top + contentRect.height > viewportHeight - padding) top = viewportHeight - contentRect.height - padding

        setPosition({ top, left })
    }, [side, sideOffset])

    const handleMouseEnter = () => {
        if (!hasHover.current) return
        setOpen(true)
    }

    const handleMouseLeave = (e: React.MouseEvent) => {
        if (!hasHover.current) return
        const relatedTarget = e.relatedTarget as HTMLElement
        if (contentRef.current?.contains(relatedTarget)) return
        setOpen(false)
        setPosition(null)
    }

    const handleContentMouseEnter = () => {
        if (!hasHover.current) return
        setOpen(true)
    }

    const handleContentMouseLeave = (e: React.MouseEvent) => {
        if (!hasHover.current) return
        const relatedTarget = e.relatedTarget as HTMLElement
        if (triggerRef.current?.contains(relatedTarget)) return
        setOpen(false)
        setPosition(null)
    }

    const handleClickOutside = (e: MouseEvent) => {
        if (!triggerRef.current?.contains(e.target as Node) && !contentRef.current?.contains(e.target as Node)) {
            setOpen(false)
            setPosition(null)
        }
    }

    // Calculate position when content is rendered
    React.useEffect(() => {
        if (open && contentRef.current) {
            calculatePosition()
        }
    }, [open, calculatePosition])

    // Click outside handler
    React.useEffect(() => {
        if (open) {
            document.addEventListener("mousedown", handleClickOutside)
        } else {
            document.removeEventListener("mousedown", handleClickOutside)
        }
        return () => {
            document.removeEventListener("mousedown", handleClickOutside)
        }
    }, [open])

    // Recalculate position on scroll/resize
    React.useEffect(() => {
        if (!open || !position) return
        const handleScroll = () => calculatePosition()
        const handleResize = () => calculatePosition()
        window.addEventListener("scroll", handleScroll, true)
        window.addEventListener("resize", handleResize)
        return () => {
            window.removeEventListener("scroll", handleScroll, true)
            window.removeEventListener("resize", handleResize)
        }
    }, [open, position, calculatePosition])

    return (
        <div
            ref={triggerRef}
            className={cn("inline-block", className)}
            onMouseEnter={handleMouseEnter}
            onMouseLeave={handleMouseLeave}
        >
            {children}
            {open && (
                <div
                    ref={contentRef}
                    className={cn(
                        "animate-in fade-in-0 zoom-in-95 fixed z-50",
                        !position && "pointer-events-none opacity-0",
                        contentClassName,
                    )}
                    style={
                        position
                            ? {
                                  top: `${position.top}px`,
                                  left: `${position.left}px`,
                              }
                            : {
                                  top: "0",
                                  left: "0",
                              }
                    }
                    onMouseEnter={handleContentMouseEnter}
                    onMouseLeave={handleContentMouseLeave}
                >
                    {content}
                </div>
            )}
        </div>
    )
}

export function TooltipProvider({ children }: { children: React.ReactNode }) {
    return <>{children}</>
}

export function TooltipTrigger({
    children,
    ...props
}: { children: React.ReactNode } & React.HTMLAttributes<HTMLDivElement>) {
    return <div {...props}>{children}</div>
}

export function TooltipContent({
    children,
    className,
    ...props
}: { children: React.ReactNode; className?: string } & React.HTMLAttributes<HTMLDivElement>) {
    return (
        <div className={className} {...props}>
            {children}
        </div>
    )
}
