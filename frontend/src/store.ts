import { create } from "zustand"
import type { LastMessage } from "../api/group"



interface Store {
    isLoggedIn: boolean
    currentUserId: string
    selectedGroupId: string
    authLoading: boolean
    lastMessages: Map<string, LastMessage>
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
    lastMessages: new Map(),
    setIsLoggedIn: (isLoggedIn: boolean) => set({ isLoggedIn: isLoggedIn }),
    setSelectedGroupId: (groupId: string) => set({ selectedGroupId: groupId }),
    setAuthLoading: (value: boolean) => set({ authLoading: value }),
    setCurrentUserId: (value: string) => set({ currentUserId: value }),
}))