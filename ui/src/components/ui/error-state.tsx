import { cn } from "@/lib/utils"
import { Button } from "./button"
import { AlertCircle } from "lucide-react"

interface ErrorStateProps extends React.HTMLAttributes<HTMLDivElement> {
    error?: Error | string
    title?: string
    action?: React.ReactNode
    size?: "sm" | "md" | "lg"
}

const sizeClasses = {
    sm: "p-6",
    md: "p-8",
    lg: "p-12",
}

export function ErrorState({
    error,
    title = "Error",
    action,
    size = "md",
    className,
    ...props
}: ErrorStateProps) {
    const errorMessage = typeof error === "string" ? error : error?.message || "An unexpected error occurred"

    return (
        <div
            className={cn(
                "flex flex-col items-center justify-center text-center",
                sizeClasses[size],
                className
            )}
            {...props}
        >
            <div className="mb-4 text-error-500">
                <AlertCircle className="w-12 h-12" />
            </div>
            <h3 className="text-lg font-semibold text-foreground mb-2">{title}</h3>
            <p className="text-foreground-muted max-w-sm mb-4">{errorMessage}</p>
            {action}
        </div>
    )
}
