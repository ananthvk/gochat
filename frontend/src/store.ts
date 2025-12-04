import { create } from "zustand"

interface Store {
    isLoggedIn: boolean
    currentUserId: string
    selectedGroupId: string
    authLoading: boolean
    setIsLoggedIn: (isLoggedIn: boolean) => void
    setSelectedGroupId: (groupId: string) => void
    setAuthLoading: (value: boolean) => void
    setCurrentUserId: (value: string) => void
}

export const useChatStore = create<Store>((set) => ({
    isLoggedIn: false,
    authLoading: true,
    selectedGroupId: "",
    currentUserId: "",
    setIsLoggedIn: (isLoggedIn: boolean) => set({ isLoggedIn: isLoggedIn }),
    setSelectedGroupId: (groupId: string) => set({ selectedGroupId: groupId }),
    setAuthLoading: (value: boolean) => set({ authLoading: value }),
    setCurrentUserId: (value: string) => set({ currentUserId: value })
}))