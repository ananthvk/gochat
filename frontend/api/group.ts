import { axiosClient } from './axios'

export interface Group {
    created_at: string
    description: string
    id: string
    name: string
    owner_id: string
}

export type GroupResult = {
    groups: Group[]
}

export type GroupMember = {
    joined_at: string
    role: string
    usr_id: string
    username: string
    name: string
}
export type GroupMemberResult = {
    members: GroupMember[]
}

export const getGroups = async (): Promise<GroupResult> => {
    try {
        const response = await axiosClient.get("/group")
        return response.data
    } catch (error: any) {
        if (!error.response) {
            throw {
                success: false,
                error: 'Network error',
                errorDetails: {}
            }
        }
        const errorData = error.response.data
        throw {
            success: false,
            error: String(errorData.reason) || String(errorData.error),
            errorDetails: errorData
        }
    }
}

export const getGroupMembers = async (groupId: string): Promise<GroupMemberResult> => {
    try {
        const response = await axiosClient.get(`/group/${groupId}/member`)
        return response.data
    } catch (error: any) {
        if (!error.response) {
            throw {
                success: false,
                error: 'Network error',
                errorDetails: {}
            }
        }
        const errorData = error.response.data
        throw {
            success: false,
            error: String(errorData.reason) || String(errorData.error),
            errorDetails: errorData
        }
    }
}