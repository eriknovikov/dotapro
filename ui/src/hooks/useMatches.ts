import { useQuery } from "@tanstack/react-query"
import { getMatches, getMatchById } from "@/api/api"
import type { MatchFilters } from "@/types"

export function useMatches(filters: MatchFilters) {
    return useQuery({
        queryKey: ["matches", filters],
        queryFn: ({ signal }) => getMatches(filters, signal),
        staleTime: 5 * 60 * 1000, // 5 minutes
        retry: 2,
    })
}

export function useMatch(id: number) {
    return useQuery({
        queryKey: ["match", id],
        queryFn: ({ signal }) => getMatchById(id, signal),
        staleTime: 10 * 60 * 1000, // 10 minutes
    })
}