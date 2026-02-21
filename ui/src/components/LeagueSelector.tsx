import { useQuery } from "@tanstack/react-query"
import { useState, useRef, useEffect, useLayoutEffect, useMemo, useCallback } from "react"
import { searchLeagues } from "../api/api"
import { useDebounce } from "../hooks/useDebounce"
import { cn } from "../lib/utils"
import { Spinner } from "./Spinner"
import popularData from "../assets/static_data.json"

// ============================================================================
// Types
// ============================================================================

interface LeagueSelectorProps {
    onSelect: (leagueId: number | undefined) => void
    initialValue?: number
    className?: string
    id?: string
    "aria-label"?: string
    inputClassName?: string
}

interface League {
    id: number
    name: string
}

// ============================================================================
// Constants
// ============================================================================

const DEBOUNCE_MS = 250
const STALE_TIME_MS = 5000
const NAVIGATION_KEYS = ["ArrowDown", "ArrowUp", "Enter"] as const

// ============================================================================
// Icons
// ============================================================================

const ClearIcon = () => (
    <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
    </svg>
)

const ChevronIcon = ({ isOpen }: { isOpen: boolean }) => (
    <svg
        className={cn("h-4 w-4 text-foreground-muted transition-transform duration-300", isOpen && "rotate-180")}
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        aria-hidden="true"
    >
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
    </svg>
)

// ============================================================================
// Component
// ============================================================================

export function LeagueSelector({
    onSelect,
    initialValue,
    className,
    id,
    "aria-label": ariaLabel,
    inputClassName,
}: LeagueSelectorProps) {
    // ---------------------------------------------------------------------------
    // State
    // ---------------------------------------------------------------------------

    const [isOpen, setIsOpen] = useState(false)
    const [inputValue, setInputValue] = useState("")
    const [highlightedIndex, setHighlightedIndex] = useState(-1)

    // ---------------------------------------------------------------------------
    // Refs
    // ---------------------------------------------------------------------------

    const containerRef = useRef<HTMLDivElement>(null)
    const inputRef = useRef<HTMLDivElement>(null)
    const listRef = useRef<HTMLUListElement>(null)
    const itemRefs = useRef<(HTMLLIElement | null)[]>([])
    const isSelectingRef = useRef(false)
    const isExternalUpdateRef = useRef(false)
    const hasSetInitialValueRef = useRef(false)
    const selectedLeagueIdRef = useRef<number | undefined>(undefined)
    const selectedLeagueNameRef = useRef<string | undefined>(undefined)

    // ---------------------------------------------------------------------------
    // Data fetching
    // ---------------------------------------------------------------------------

    const debouncedQuery = useDebounce(inputValue, DEBOUNCE_MS)

    const {
        data: searchResults = [],
        isLoading,
        isFetching,
    } = useQuery({
        queryKey: ["leagues", "search", debouncedQuery],
        queryFn: async ({ signal }) => {
            if (!debouncedQuery.trim()) return []
            return searchLeagues(debouncedQuery, signal)
        },
        enabled: debouncedQuery.trim().length > 0,
        staleTime: STALE_TIME_MS,
    })

    // ---------------------------------------------------------------------------
    // Computed values
    // ---------------------------------------------------------------------------

    const items = useMemo(() => {
        if (debouncedQuery.trim().length === 0) {
            return popularData.popular_leagues
        }
        return (searchResults || []).map(l => ({ id: l.league_id, name: l.name }))
    }, [debouncedQuery, searchResults])

    const listboxId = `${id}-listbox`
    const hasQuery = debouncedQuery.trim().length > 0
    const showSpinner = isFetching && hasQuery
    const showClearButton = inputValue && !showSpinner
    const showEmptyState = items.length === 0 && !isLoading && !isFetching

    // ---------------------------------------------------------------------------
    // Helper functions
    // ---------------------------------------------------------------------------

    const clearSelectedLeague = useCallback(() => {
        if (selectedLeagueIdRef.current !== undefined) {
            isExternalUpdateRef.current = true
            setInputValue("")
            onSelect(undefined)
            selectedLeagueIdRef.current = undefined
            selectedLeagueNameRef.current = undefined
        }
    }, [onSelect])

    const shouldClearLeagueOnClose = useCallback(() => {
        if (selectedLeagueIdRef.current === undefined) {
            return false
        }
        // Clear if input is empty OR if input differs from the originally selected league name
        return inputValue.trim() === "" || inputValue !== selectedLeagueNameRef.current
    }, [inputValue])

    const closeDropdown = useCallback(() => {
        setIsOpen(false)
        if (shouldClearLeagueOnClose()) {
            clearSelectedLeague()
        }
    }, [shouldClearLeagueOnClose, clearSelectedLeague])

    // ---------------------------------------------------------------------------
    // Effects
    // ---------------------------------------------------------------------------

    // Set initial value when component mounts or initialValue changes
    useEffect(() => {
        if (initialValue !== undefined && !hasSetInitialValueRef.current) {
            const initialItem = items.find(l => l.id === initialValue)
            if (initialItem) {
                isExternalUpdateRef.current = true
                setInputValue(initialItem.name)
                hasSetInitialValueRef.current = true
                selectedLeagueIdRef.current = initialValue
                selectedLeagueNameRef.current = initialItem.name
            }
        } else if (initialValue === undefined && selectedLeagueIdRef.current !== undefined) {
            // Clear the input when initialValue becomes undefined
            isExternalUpdateRef.current = true
            setInputValue("")
            selectedLeagueIdRef.current = undefined
            selectedLeagueNameRef.current = undefined
        }
    }, [initialValue, items])

    // Reset hasSetInitialValueRef when initialValue changes
    useEffect(() => {
        hasSetInitialValueRef.current = false
    }, [initialValue])

    // Close dropdown when clicking outside
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
                closeDropdown()
            }
        }
        document.addEventListener("mousedown", handleClickOutside)
        return () => document.removeEventListener("mousedown", handleClickOutside)
    }, [closeDropdown])

    // Reset highlighted index when items change
    useEffect(() => {
        setHighlightedIndex(-1)
        // Only reset itemRefs if the length changed significantly to avoid disrupting scroll-into-view
        if (itemRefs.current.length !== items.length) {
            itemRefs.current = new Array(items.length).fill(null)
        }
    }, [items])

    // Scroll highlighted item into view
    useEffect(() => {
        if (highlightedIndex >= 0 && itemRefs.current[highlightedIndex]) {
            itemRefs.current[highlightedIndex]?.scrollIntoView({
                block: "nearest",
                behavior: "smooth",
            })
        }
    }, [highlightedIndex])

    // Sync div content with inputValue when it changes externally (e.g., from selection)
    useLayoutEffect(() => {
        if (inputRef.current) {
            // Always clear the div content when inputValue is empty to ensure placeholder shows
            if (inputValue === "") {
                inputRef.current.textContent = ""
                isExternalUpdateRef.current = false
                return
            }
            
            if (isExternalUpdateRef.current) {
                inputRef.current.textContent = inputValue
                // Don't manipulate selection range to avoid unwanted focus
                // The cursor will be placed naturally when user interacts with the input
                isExternalUpdateRef.current = false
            }
        }
    }, [inputValue])

    // ---------------------------------------------------------------------------
    // Event handlers
    // ---------------------------------------------------------------------------

    const handleSelect = useCallback(
        (league: League) => {
            isExternalUpdateRef.current = true
            setInputValue(league.name)
            setIsOpen(false)
            onSelect(league.id)
            selectedLeagueIdRef.current = league.id
            selectedLeagueNameRef.current = league.name
            isSelectingRef.current = true
            inputRef.current?.focus()
            setTimeout(() => {
                isSelectingRef.current = false
            }, 0)
        },
        [onSelect],
    )

    const handleClear = useCallback(() => {
        isExternalUpdateRef.current = true
        setInputValue("")
        onSelect(undefined)
        selectedLeagueIdRef.current = undefined
        selectedLeagueNameRef.current = undefined
        inputRef.current?.focus()
    }, [onSelect])

    const handleInputFocus = useCallback(() => {
        if (!isSelectingRef.current && !isOpen && !isExternalUpdateRef.current) {
            setIsOpen(true)
        }
    }, [isOpen])

    const handleInputChange = useCallback(
        (e: React.FormEvent<HTMLDivElement>) => {
            const text = e.currentTarget.textContent || ""
            setInputValue(text)
            if (!isExternalUpdateRef.current) {
                setIsOpen(true)
            }
        },
        [],
    )

    const handleInputPaste = useCallback((e: React.ClipboardEvent<HTMLDivElement>) => {
        e.preventDefault()
        const text = e.clipboardData.getData("text/plain")
        document.execCommand("insertText", false, text)
    }, [])

    const handleInputBlur = useCallback(() => {
        // Ensure inputValue is in sync with div content on blur
        if (inputRef.current) {
            const text = inputRef.current.textContent || ""
            // Only update state if value actually changed to avoid unnecessary re-renders
            if (text !== inputValue) {
                setInputValue(text)
            }
        }
    }, [inputValue])

    const handleItemMouseEnter = useCallback((index: number) => {
        setHighlightedIndex(index)
    }, [])

    const handleItemClick = useCallback(
        (e: React.MouseEvent, league: League) => {
            e.preventDefault()
            e.stopPropagation()
            handleSelect(league)
        },
        [handleSelect],
    )

    // ---------------------------------------------------------------------------
    // Keyboard navigation
    // ---------------------------------------------------------------------------

    const handleKeyDown = useCallback(
        (e: React.KeyboardEvent) => {
            // Open dropdown on navigation keys when closed
            if (!isOpen) {
                if (NAVIGATION_KEYS.includes(e.key as any)) {
                    // Prevent newline insertion when pressing Enter to open dropdown
                    if (e.key === "Enter") {
                        e.preventDefault()
                    }
                    setIsOpen(true)
                }
                return
            }

            switch (e.key) {
                case "ArrowDown":
                    e.preventDefault()
                    setHighlightedIndex(prev => (prev < items.length - 1 ? prev + 1 : prev))
                    break
                case "ArrowUp":
                    e.preventDefault()
                    setHighlightedIndex(prev => (prev > 0 ? prev - 1 : 0))
                    break
                case "Enter":
                    e.preventDefault()
                    if (highlightedIndex >= 0 && items[highlightedIndex]) {
                        handleSelect(items[highlightedIndex])
                    }
                    break
                case "Escape":
                    closeDropdown()
                    inputRef.current?.focus()
                    break
                case "Tab":
                    closeDropdown()
                    break
            }
        },
        [isOpen, items, highlightedIndex, handleSelect, closeDropdown],
    )

    // ---------------------------------------------------------------------------
    // Render helpers
    // ---------------------------------------------------------------------------

    const renderInput = () => (
        <div className="relative">
            {inputValue === "" && (
                <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-foreground-muted pointer-events-none select-none">
                    Type to search leagues...
                </span>
            )}
            <div
                ref={inputRef}
                id={id}
                contentEditable
                suppressContentEditableWarning
                onInput={handleInputChange}
                onPaste={handleInputPaste}
                onBlur={handleInputBlur}
                onFocus={handleInputFocus}
                onKeyDown={handleKeyDown}
                tabIndex={0}
                aria-label={ariaLabel || "Search leagues"}
                aria-autocomplete="list"
                aria-controls={listboxId}
                aria-expanded={isOpen}
                aria-activedescendant={highlightedIndex >= 0 ? `${listboxId}-option-${highlightedIndex}` : undefined}
                role="combobox"
                className={cn(
                    "min-h-10 w-full rounded-lg border px-3 py-2 text-sm pr-10",
                    "bg-background border-border text-foreground",
                    "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-800/50",
                    "transition-all duration-200 ease-in-out",
                    "whitespace-pre-wrap break-words",
                    inputClassName,
                )}
            />
        </div>
    )

    const renderInputSuffix = () => (
        <div className="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1">
            {showSpinner && <Spinner size="sm" aria-label="Loading leagues" />}
            {showClearButton && (
                <button
                    type="button"
                    onClick={handleClear}
                    className="text-foreground-muted hover:text-foreground transition-colors"
                    aria-label="Clear search"
                >
                    <ClearIcon />
                </button>
            )}
            {!showSpinner && !showClearButton && (
                <button
                    type="button"
                    onClick={() => {
                        setIsOpen(!isOpen)
                        inputRef.current?.focus()
                    }}
                    className="text-foreground-muted hover:text-foreground transition-colors"
                    aria-label={isOpen ? "Close dropdown" : "Open dropdown"}
                    aria-expanded={isOpen}
                >
                    <ChevronIcon isOpen={isOpen} />
                </button>
            )}
        </div>
    )

    const renderEmptyState = () => (
        <li role="option" aria-selected="false" className="px-4 py-3 text-sm text-foreground-muted">
            No leagues found
        </li>
    )

    const renderLeagueItem = (league: League, index: number) => (
        <li
            key={league.id}
            ref={el => {
                itemRefs.current[index] = el
            }}
            id={`${listboxId}-option-${index}`}
            role="option"
            aria-selected={highlightedIndex === index}
            onClick={e => handleItemClick(e, league)}
            onMouseEnter={() => handleItemMouseEnter(index)}
            className={cn("cursor-pointer px-4 py-3 text-sm", highlightedIndex === index && "bg-red-800/20")}
        >
            {league.name}
        </li>
    )

    const renderDropdown = () => (
        <ul
            ref={listRef}
            id={listboxId}
            role="listbox"
            aria-label={ariaLabel || "Leagues"}
            className="absolute z-50 max-h-72 overflow-auto rounded-lg border border-border-accent bg-background-card shadow-xl w-full top-full mt-1"
        >
            {showEmptyState ? renderEmptyState() : items.map(renderLeagueItem)}
        </ul>
    )

    // ---------------------------------------------------------------------------
    // Render
    // ---------------------------------------------------------------------------

    return (
        <div ref={containerRef} className={cn("relative w-full", className)}>
            <div className="relative">
                {renderInput()}
                {renderInputSuffix()}
                {isOpen && renderDropdown()}
            </div>
        </div>
    )
}
