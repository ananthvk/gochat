import { create } from "zustand"

interface Store {
    isLoggedIn: boolean
    selectedGroupId: string
    authLoading: boolean
    setIsLoggedIn: (isLoggedIn: boolean) => void
    setSelectedGroupId: (groupId: string) => void
    setAuthLoading: (value: boolean) => void
}

export const useChatStore = create<Store>((set) => ({
    isLoggedIn: false,
    authLoading: true,
    selectedGroupId: "",
    setIsLoggedIn: (isLoggedIn: boolean) => set({ isLoggedIn: isLoggedIn }),
    setSelectedGroupId: (groupId: string) => set({ selectedGroupId: groupId }),
    setAuthLoading: (value: boolean) => set({ authLoading: value })
}))