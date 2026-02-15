import { cn } from "@/lib/utils"
import { Button } from "./button"

interface EmptyStateProps extends React.HTMLAttributes<HTMLDivElement> {
    title?: string
    description?: string
    action?: React.ReactNode
    size?: "sm" | "md" | "lg"
    icon?: React.ReactNode
}

const sizeClasses = {
    sm: "p-6",
    md: "p-8",
    lg: "p-12",
}

export function EmptyState({
    title = "No data found",
    description,
    action,
    size = "md",
    icon,
    className,
    ...props
}: EmptyStateProps) {
    return (
        <div
            className={cn(
                "flex flex-col items-center justify-center text-center",
                sizeClasses[size],
                className
            )}
            {...props}
        >
            {icon && <div className="mb-4 text-foreground-muted">{icon}</div>}
            <h3 className="text-lg font-semibold text-foreground mb-2">{title}</h3>
            {description && (
                <p className="text-foreground-muted max-w-sm mb-4">{description}</p>
            )}
            {action}
        </div>
    )
}
