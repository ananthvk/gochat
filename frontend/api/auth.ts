import { axiosClient } from './axios'
export type LoginDetails = {
    email: string,
    password: string
}

export type LoginResult = {
    success: boolean,
    error: string,
    errorDetails: any
}

export const loginUser = async ({ email, password }: LoginDetails): Promise<LoginResult> => {
    try {
        const response = await axiosClient.post('/auth/login', {
            email,
            password
        })

        const { token, expiry } = response.data.authenticate
        localStorage.setItem('session_expiry', expiry)
        localStorage.setItem('session_token', token)

        return { success: true, error: "", errorDetails: {} }

    } catch (error: any) {
        if (!error.response) {
            throw {
                success: false,
                error: 'Network error - please check your connection',
                errorDetails: {}
            }
        }

        const errorData = error.response.data

        throw {
            success: false,
            error: String(errorData.reason) || String(errorData.error) || 'Authentication failed',
            errorDetails: errorData
        }
    }
}



export type User = {
    id: string,
    email: string,
    name: string,
    username: string,
    activated: boolean,
    created_at: string
}

export type MeResult = {
    success: boolean,
    user?: User,
    error: string,
    errorDetails: any
}
export type MeError = {
    success: boolean,
    error: string,
    errorDetails: any
}

export const getMe = async (): Promise<MeResult> => {
    try {
        const token = localStorage.getItem('session_token')
        if (!token) {
            throw {
                success: false,
                error: 'No authentication token found',
                errorDetails: {}
            }
        }

        const response = await axiosClient.get('/auth/me', {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })

        return {
            success: true,
            user: response.data.user,
            error: "",
            errorDetails: {}
        }

    } catch (error: any) {
        if (!error.response) {
            throw {
                success: false,
                error: 'Network error - please check your connection',
                errorDetails: {}
            }
        }

        const errorData = error.response.data

        throw {
            success: false,
            error: String(errorData.reason) || String(errorData.error) || 'Failed to get user info',
            errorDetails: errorData
        }
    }
}