import { useQuery } from "@tanstack/react-query"
import { getGroups, type GroupResult, type Group, getGroupMembers, type GroupMember } from "../../api/group"
import type { APIError } from '../../api/errors'
import { useChatStore } from "../store"

export const useGroups = () => {
    const isLoggedIn = useChatStore((state) => state.isLoggedIn)
    const query = useQuery<GroupResult, APIError, Record<string, Group>>({
        queryKey: ["groups"],
        queryFn: getGroups,
        staleTime: 15 * 1000, // 15 s
        select: (result) => {
            return result.groups.reduce((acc, g) => {
                acc[g.id] = g
                return acc
            }, {} as Record<string, Group>)
        },
        enabled: isLoggedIn
    })
    return query
}

export const useGroupMembers = (groupId: string) => {
    return useQuery({
        queryKey: ["groups", groupId, "members"],
        queryFn: () => getGroupMembers(groupId),
        enabled: !!groupId,
        select: (data) => {
            const memberMap: Record<string, GroupMember> = {};
            data.members.forEach(member => {
                memberMap[member.usr_id] = member;
            });
            return memberMap;
        },
        // 5 min cache
        staleTime: 1000 * 60 * 5,
        refetchOnWindowFocus: false,
    });
};