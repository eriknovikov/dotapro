import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
}

const MS_PER_MINUTE = 60 * 1000
const MS_PER_HOUR = 60 * MS_PER_MINUTE
const MS_PER_DAY = 24 * MS_PER_HOUR
const MS_PER_MONTH = 30.44 * MS_PER_DAY
const MS_PER_YEAR = 365.25 * MS_PER_DAY

const UNIT_LABELS = ["year", "month", "day", "hour", "minute"] as const

export function formatRelativeTime(dateString: string): string {
    const date = new Date(dateString)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()

    const years = Math.floor(diffMs / MS_PER_YEAR)
    const months = Math.floor((diffMs % MS_PER_YEAR) / MS_PER_MONTH)
    const days = Math.floor((diffMs % MS_PER_MONTH) / MS_PER_DAY)
    const hours = Math.floor((diffMs % MS_PER_DAY) / MS_PER_HOUR)
    const minutes = Math.floor((diffMs % MS_PER_HOUR) / MS_PER_MINUTE)

    const unitValues = [years, months, days, hours, minutes]
    const firstUnitIndex = unitValues.findIndex(value => value > 0)

    if (firstUnitIndex === -1) {
        return "just now"
    }

    const firstValue = unitValues[firstUnitIndex]
    const firstLabel = UNIT_LABELS[firstUnitIndex]
    const firstPart = `${firstValue} ${firstLabel}${firstValue !== 1 ? "s" : ""}`

    const secondValue = unitValues[firstUnitIndex + 1]
    if (!secondValue || secondValue === 0) {
        return firstPart
    }

    const secondLabel = UNIT_LABELS[firstUnitIndex + 1]
    const secondPart = `${secondValue} ${secondLabel}${secondValue !== 1 ? "s" : ""}`

    return `${firstPart} and ${secondPart}`
}
