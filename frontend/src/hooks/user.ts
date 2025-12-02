import { useMutation, useQuery } from "@tanstack/react-query"
import { loginUser, type LoginDetails, type LoginResult, getMe, type MeResult, type MeError } from "../../api/auth"
import { useChatStore } from "../store"
import { useEffect } from "react"

export const useLogin = () => {
    const mutation = useMutation<LoginResult, unknown, LoginDetails>({
        mutationFn: loginUser
    })
    return mutation
}

export const useAuthBootstrap = () => {
    const token = localStorage.getItem("session_token")
    const setLoggedIn = useChatStore(state => state.setIsLoggedIn)
    const setLoading = useChatStore(state => state.setAuthLoading)
    const query = useQuery<MeResult, MeError, MeResult>({
        queryKey: ["auth", "me"],
        queryFn: getMe,
        retry: false,
        staleTime: Infinity,
        enabled: !!token,
    })

    useEffect(() => {
        if (query.isSuccess) {
            // The token is valid
            setLoggedIn(true)
            setLoading(false)
        }
    }, [query.isSuccess])

    useEffect(() => {
        if (query.isError) {
            // Remove the token only if it's an authentication error
            // Otherwise the token may get removed due to network errors
            if (query.error.errorDetails && query.error.errorDetails.error == "not_authenticated") {
                setLoggedIn(false)
                localStorage.removeItem("session_token")
            }
            setLoading(false)
        }
    }, [query.isError])

    useEffect(() => {
        if (!token) {
            setLoggedIn(false)
            setLoading(false)
        }
    }, [])
    return query
}