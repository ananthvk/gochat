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