import { useQuery } from "@tanstack/react-query"
import { getGroups, type GroupResult, type Group } from "../../api/group"
import type { APIError } from '../../api/errors'

export const useGroups = () => {
    const query = useQuery<GroupResult, APIError, Record<string, Group>>({
        queryKey: ["groups"],
        queryFn: getGroups,
        staleTime: 15 * 1000, // 15 s
        select: (result) => {
            return result.groups.reduce((acc, g) => {
                acc[g.id] = g
                return acc
            }, {} as Record<string, Group>)
        }
    })
    return query
}