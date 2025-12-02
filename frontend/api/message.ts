import { axiosClient } from './axios'
import { type APIError } from './errors'


export interface Message {
    content: string
    created_at: string
    group_id: string
    id: string
    sender_id: string
    type: string
}

export type MessageCursor = {
    before: string
    has_before: boolean
}

export type MessageResult = {
    cursor: MessageCursor
    messages: Message[]
}
export type PaginationParams = {
    groupId: string
    limit: number
    before: string
}

export const getMessages = async (pagination: PaginationParams): Promise<MessageResult> => {
    try {
        const response = await axiosClient.get(`/group/${pagination.groupId}/message?before=${pagination.before}&limit=${pagination.limit}`)
        return response.data
    } catch (error: any) {
        if (!error.response) {
            throw {
                success: false,
                error: 'Network error',
                errorDetails: {}
            } satisfies APIError
        }
        const errorData = error.response.data
        throw {
            success: false,
            error: String(errorData.reason) || String(errorData.error),
            errorDetails: errorData
        }
    }
}