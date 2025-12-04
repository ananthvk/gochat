import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { loginUser, type LoginDetails, type LoginResult, getMe, type MeResult } from "../../api/auth"
import { type APIError } from '../../api/errors'
import { useChatStore } from "../store"
import { useEffect } from "react"
import { queryClient } from "../../api/query-client"

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
    const setCurrentUserId = useChatStore(state => state.setCurrentUserId)
    const query = useQuery<MeResult, APIError, MeResult>({
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
            setCurrentUserId(query.data.user ? query.data.user.id : "")
            setLoading(false)
        }
    }, [query.isSuccess])

    useEffect(() => {
        if (query.isError) {
            // Remove the token only if it's an authentication error
            // Otherwise the token may get removed due to network errors
            if (query.error.errorDetails && query.error.errorDetails.error == "not_authenticated") {
                setLoggedIn(false)
                setCurrentUserId("")
                localStorage.removeItem("session_token")
                queryClient.clear();
            }
            setCurrentUserId("")
            setLoading(false)
        }
    }, [query.isError])

    useEffect(() => {
        if (!token) {
            setLoggedIn(false)
            setCurrentUserId("")
            setLoading(false)
        }
    }, [])
    return query
}

export const useLogout = () => {
    const queryClient = useQueryClient();
    const setLoggedIn = useChatStore((state) => state.setIsLoggedIn);
    const setCurrentUserId = useChatStore((state) => state.setCurrentUserId);

    const logout = () => {
        localStorage.removeItem("session_token");
        setLoggedIn(false);
        setCurrentUserId("");
        queryClient.clear();
    };

    return logout;
};