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

        const response = await axiosClient.get('/auth/me')

        return {
            success: true,
            user: response.data.user,
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
// This method logs out the user, it always succeeds and never results in an error even if the request fails
// Irrespective of what happens, the token is removed from local storage
export const logOut = async () => {
    const token = localStorage.getItem('session_token')
    if (!token)
        return
    try {
        await axiosClient.post('/auth/logout')
    } catch (error: any) {

    } finally {
        localStorage.removeItem("session_token")
        localStorage.removeItem("session_expiry")
    }
}