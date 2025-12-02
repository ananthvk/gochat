import Axios, { type InternalAxiosRequestConfig } from "axios"
if (!(import.meta as any).env.VITE_API_BASE_URL) {
    throw new Error("VITE_API_BASE_URL not set")
}
export const axiosClient = Axios.create({
    baseURL: (import.meta as any).env.VITE_API_BASE_URL
})
axiosClient.interceptors.request.use(function (config): InternalAxiosRequestConfig {
    const token = localStorage.getItem("session_token");
    if (token) {
        config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
})