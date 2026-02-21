import { cn } from "@/lib/utils"

interface SpinnerProps extends React.HTMLAttributes<HTMLDivElement> {
    size?: "sm" | "md" | "lg" | "xl"
}

const sizeClasses = {
    sm: "w-4 h-4 border-2",
    md: "w-6 h-6 border-2",
    lg: "w-8 h-8 border-3",
    xl: "w-12 h-12 border-4",
}

export function Spinner({ size = "md", className, ...props }: SpinnerProps) {
    return (
        <div
            className={cn(
                "rounded-full border-solid border-transparent spinner-gradient",
                sizeClasses[size],
                className,
            )}
            {...props}
        />
    )
}
