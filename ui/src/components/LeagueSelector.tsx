import { useSelect } from "downshift"
import popularData from "../assets/static_data.json"
import { cn } from "../lib/utils"

interface LeagueSelectorProps {
    onSelect: (leagueId: number | undefined) => void
    initialValue?: number
    className?: string
    id?: string
    "aria-label"?: string
}

const leagues = popularData.popular_leagues

const anyLeagueItem = { id: -1, name: "Any league" }

export function LeagueSelector({
    onSelect,
    initialValue,
    className,
    id,
    "aria-label": ariaLabel,
}: LeagueSelectorProps) {
    const initialItem = initialValue ? leagues.find(l => l.id === initialValue) || anyLeagueItem : anyLeagueItem

    const { isOpen, selectedItem, getToggleButtonProps, getMenuProps, highlightedIndex, getItemProps } = useSelect({
        items: [anyLeagueItem, ...leagues],
        initialSelectedItem: initialItem,
        onSelectedItemChange: ({ selectedItem }) => {
            onSelect(selectedItem?.id === -1 ? undefined : selectedItem?.id)
        },
    })

    return (
        <div className={cn("relative w-full", className)}>
            <button
                type="button"
                {...getToggleButtonProps({ id, "aria-label": ariaLabel })}
                className={cn(
                    "flex h-10 w-full items-center justify-between rounded-lg border px-3 py-2 text-sm",
                    "bg-background border-border text-foreground",
                    "focus-visible:outline-none",
                    "transition-all duration-200 ease-in-out",
                )}
            >
                <span className={cn(!selectedItem && "text-foreground-muted")}>
                    {selectedItem?.name || "Select league..."}
                </span>
                <svg
                    className={cn(
                        "h-4 w-4 text-foreground-muted transition-transform duration-300",
                        isOpen && "rotate-180",
                    )}
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
            </button>

            <ul
                {...getMenuProps()}
                className={cn(
                    "absolute z-50 mt-1.5 max-h-72 w-full min-w-75 overflow-auto rounded-lg border border-border-accent bg-background-card shadow-xl",
                    !isOpen && "hidden",
                )}
            >
                {isOpen &&
                    [anyLeagueItem, ...leagues]
                        .filter(league => league.id !== selectedItem?.id)
                        .map(league => {
                            const originalIndex = [anyLeagueItem, ...leagues].findIndex(l => l.id === league.id)
                            return (
                                <li
                                    key={league.id}
                                    {...getItemProps({ item: league, index: originalIndex })}
                                    className={cn(
                                        "cursor-pointer px-4 py-3 text-sm",
                                        highlightedIndex === originalIndex && "bg-red-800/20",
                                    )}
                                >
                                    {league.name}
                                </li>
                            )
                        })}
            </ul>
        </div>
    )
}

interface League {
    id: number
    name: string
}
