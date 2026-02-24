import { cn } from "@/lib"

interface SkeletonProps extends React.HTMLAttributes<HTMLDivElement> {
    variant?: "default" | "card" | "text" | "button"
}

export function Skeleton({ className, variant = "default", ...props }: SkeletonProps) {
    return (
        <div
            className={cn(
                "animate-pulse rounded-md bg-background-accent",
                variant === "card" && "rounded-xl bg-background-card/80",
                variant === "text" && "h-4 w-3/4",
                variant === "button" && "h-8 w-32 rounded-md",
                className
            )}
            {...props}
        />
    )
}

// SeriesCard skeleton matching the exact layout
export function SeriesCardSkeleton() {
    return (
        <div className="group relative rounded-xl bg-background-card/80 text-card-foreground shadow-xl overflow-hidden backdrop-blur-sm">
            <div className="relative z-10 flex flex-col h-full p-4">
                {/* Teams and Score */}
                <div className="flex items-center justify-between gap-4 border-b border-border mb-3 pb-3">
                    {/* Team A */}
                    <div className="flex items-center min-w-0 flex-1">
                        <div className="flex min-w-0 gap-1 items-center">
                            <Skeleton className="h-7 w-11 rounded shrink-0" variant="default" />
                            <Skeleton className="h-5 w-24" variant="text" />
                        </div>
                    </div>

                    {/* Score */}
                    <div className="flex items-center justify-center shrink-0">
                        <Skeleton className="h-4 w-4" variant="default" />
                        <span className="mx-2 text-gray-400">-</span>
                        <Skeleton className="h-4 w-4" variant="default" />
                    </div>

                    {/* Team B */}
                    <div className="flex items-center justify-end min-w-0 flex-1">
                        <div className="flex min-w-0 gap-1 items-center justify-end">
                            <Skeleton className="h-5 w-24" variant="text" />
                            <Skeleton className="h-7 w-11 rounded shrink-0" variant="default" />
                        </div>
                    </div>
                </div>

                {/* League and Date */}
                <div className="text-xs text-foreground-muted flex flex-col flex-1 gap-2">
                    <Skeleton className="h-3 w-32" variant="default" />
                    <Skeleton className="h-3 w-24" variant="default" />
                </div>

                {/* Button */}
                <div className="flex justify-end items-center mt-auto">
                    <Skeleton className="h-8 w-32 rounded-md" variant="button" />
                </div>
            </div>
        </div>
    )
}
