import { useCombobox } from "downshift"
import { useRef, useState, useEffect, useCallback } from "react"
import popularData from "../assets/static_data.json"
import { cn } from "../lib/utils"
import { Spinner } from "./ui/spinner"

interface Team {
    id: number
    name: string
}

interface TeamSelectorProps {
    onSelect: (teamId: number | undefined) => void
    initialValue?: number
    placeholder?: string
    className?: string
    id?: string
    "aria-label"?: string
}

export function TeamSelector({ onSelect, initialValue, placeholder = "Search teams...", className, id, "aria-label": ariaLabel }: TeamSelectorProps) {
    const [inputItems, setInputItems] = useState<Team[]>(popularData.popular_teams)
    const [isLoading, setIsLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const abortControllerRef = useRef<AbortController | null>(null)

    // Find initial team name if initialValue is provided
    const initialTeam = initialValue
        ? popularData.popular_teams.find((t) => t.id === initialValue) || { id: initialValue, name: String(initialValue) }
        : null

    const {
        isOpen,
        getMenuProps,
        getInputProps,
        highlightedIndex,
        getItemProps,
        selectedItem,
        setInputValue,
    } = useCombobox({
        items: inputItems,
        initialInputValue: initialTeam?.name || "",
        initialSelectedItem: initialTeam || null,
        onInputValueChange: async ({ inputValue }) => {
            if (!inputValue || inputValue.trim() === "") {
                setInputItems(popularData.popular_teams)
                setError(null)
                return
            }

            // Cancel previous request
            if (abortControllerRef.current) {
                abortControllerRef.current.abort()
            }

            setIsLoading(true)
            setError(null)
            const controller = new AbortController()
            abortControllerRef.current = controller

            try {
                const res = await fetch(`http://localhost:8080/filtersmetadata/teams?q=${encodeURIComponent(inputValue)}`, {
                    signal: controller.signal,
                })
                if (!res.ok) {
                    throw new Error(`Failed to fetch: ${res.status} ${res.statusText}`)
                }
                const data = await res.json()
                setInputItems(data)
            } catch (err) {
                if (err instanceof Error && err.name !== "AbortError") {
                    setError(err.message)
                }
            } finally {
                setIsLoading(false)
            }
        },
        onSelectedItemChange: ({ selectedItem }) => {
            onSelect(selectedItem?.id)
        },
        itemToString: (item) => item?.name || "",
    })

    // Cleanup on unmount
    useEffect(() => {
        return () => {
            if (abortControllerRef.current) {
                abortControllerRef.current.abort()
            }
        }
    }, [])

    const handleClear = useCallback(() => {
        setInputValue("")
        setInputItems(popularData.popular_teams)
        onSelect(undefined)
    }, [setInputValue, onSelect])

    return (
        <div className={cn("relative w-full", className)}>
            <div className="relative">
                <input
                    {...getInputProps({ id, "aria-label": ariaLabel })}
                    placeholder={placeholder}
                    className={cn(
                        "flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-base shadow-sm transition-colors",
                        "placeholder:text-muted-foreground",
                        "focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring",
                        "disabled:cursor-not-allowed disabled:opacity-50",
                        "md:text-sm",
                        "pr-8"
                    )}
                />
                {selectedItem && (
                    <button
                        type="button"
                        onClick={handleClear}
                        className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                        aria-label="Clear selection"
                    >
                        <svg
                            className="h-4 w-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </button>
                )}
            </div>

            <ul
                {...getMenuProps()}
                className={cn(
                    "absolute z-50 mt-1 max-h-60 w-full overflow-auto rounded-md border bg-popover text-popover-foreground shadow-md",
                    "data-[state=open]:animate-in data-[state=closed]:animate-out",
                    "data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0",
                    "data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95",
                    !isOpen && "hidden"
                )}
            >
                {isOpen && (
                    <>
                        {isLoading && (
                            <li className="flex items-center justify-center py-4 px-3 text-sm text-muted-foreground">
                                <Spinner size="sm" className="mr-2" />
                                Loading...
                            </li>
                        )}

                        {error && (
                            <li className="py-4 px-3 text-sm text-destructive">
                                {error}
                            </li>
                        )}

                        {!isLoading && !error && inputItems.length === 0 && (
                            <li className="py-4 px-3 text-sm text-muted-foreground">
                                No teams found
                            </li>
                        )}

                        {!isLoading && !error && inputItems.map((item, index) => (
                            <li
                                {...getItemProps({ item, index })}
                                key={`${item.id}-${index}`}
                                className={cn(
                                    "relative flex cursor-pointer select-none items-center rounded-sm px-3 py-2 text-sm outline-none",
                                    highlightedIndex === index && "bg-accent text-accent-foreground",
                                    selectedItem?.id === item.id && "bg-accent/50"
                                )}
                            >
                                <span className="truncate">{item.name}</span>
                            </li>
                        ))}
                    </>
                )}
            </ul>
        </div>
    )
}
