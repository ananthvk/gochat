import { create } from "zustand"


interface Store {
    isLoggedIn: boolean
    selectedGroupId: string
    authLoading: boolean
    setIsLoggedIn: (isLoggedIn: boolean) => void
    setSelectedGroupId: (groupId: string) => void
    setAuthLoading: (value: boolean) => void
}

interface Group {
    created_at: string
    description: string
    id: string
    name: string
    owner_id: string
}

interface GroupStore {
    groups: Record<string, Group>
}

export const useChatStore = create<Store>((set) => ({
    isLoggedIn: false,
    authLoading: true,
    selectedGroupId: "",
    setIsLoggedIn: (isLoggedIn: boolean) => set({ isLoggedIn: isLoggedIn }),
    setSelectedGroupId: (groupId: string) => set({ selectedGroupId: groupId }),
    setAuthLoading: (value: boolean) => set({ authLoading: value })
}))

const hardcodedGroups: Record<string, Group> = {
    "01KB36WR9MK8QF86Z79Q6FJCY4": {
        "created_at": "2025-11-27T23:18:33.972065+05:30",
        "description": "This is the first group",
        "id": "01KB36WR9MK8QF86Z79Q6FJCY4",
        "name": "One",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "01KB36X0C41M5CAN43KEEFNBDJ": {
        "created_at": "2025-11-27T23:18:42.244732+05:30",
        "description": "This is the second group",
        "id": "01KB36X0C41M5CAN43KEEFNBDJ",
        "name": "Two",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "01KB36X6H14M1GFCDQM4C03WSC": {
        "created_at": "2025-11-27T23:18:48.545519+05:30",
        "description": "This is the third group",
        "id": "01KB36X6H14M1GFCDQM4C03WSC",
        "name": "Three",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "01KB36WR9MK8QF86Z79Q6FJCY8": {
        "created_at": "2025-11-27T23:18:33.972065+05:30",
        "description": "This is the fourth group",
        "id": "01KB36WR9MK8QF86Z79Q6FJCY8",
        "name": "Four",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "01KB36X0C41M5CAN43KEEFNBD9": {
        "created_at": "2025-11-27T23:18:42.244732+05:30",
        "description": "This is the fifth group",
        "id": "01KB36X0C41M5CAN43KEEFNBD9",
        "name": "Five",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "01KB36X6H14M1GFCDQM4C03WS7": {
        "created_at": "2025-11-27T23:18:48.545519+05:30",
        "description": "This is the sixth group",
        "id": "01KB36X6H14M1GFCDQM4C03WS7",
        "name": "Six",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "02KB36X0C41M5CAN43KEEFNBD9": {
        "created_at": "2025-11-27T23:18:42.244732+05:30",
        "description": "This is the seventh group",
        "id": "02KB36X0C41M5CAN43KEEFNBD9",
        "name": "Seven",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "02KB36X6H14M1GFCDQM4C03WS7": {
        "created_at": "2025-11-27T23:18:48.545519+05:30",
        "description": "This is the eighth group",
        "id": "02KB36X6H14M1GFCDQM4C03WS7",
        "name": "Eight",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "03KB36X6H14M1GFCDQM4C03WS8": {
        "created_at": "2025-11-27T23:19:00.123456+05:30",
        "description": "This is the ninth group",
        "id": "03KB36X6H14M1GFCDQM4C03WS8",
        "name": "Nine",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "04KB36X6H14M1GFCDQM4C03WS9": {
        "created_at": "2025-11-27T23:19:10.654321+05:30",
        "description": "This is the tenth group",
        "id": "04KB36X6H14M1GFCDQM4C03WS9",
        "name": "Ten",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "05KB36X6H14M1GFCDQM4C03WS0": {
        "created_at": "2025-11-27T23:19:20.789012+05:30",
        "description": "This is the eleventh group",
        "id": "05KB36X6H14M1GFCDQM4C03WS0",
        "name": "Eleven",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "06KB36X6H14M1GFCDQM4C03WS1": {
        "created_at": "2025-11-27T23:19:30.345678+05:30",
        "description": "This is the twelfth group",
        "id": "06KB36X6H14M1GFCDQM4C03WS1",
        "name": "Twelve",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "07KB36X6H14M1GFCDQM4C03WS2": {
        "created_at": "2025-11-27T23:19:40.901234+05:30",
        "description": "This is the thirteenth group",
        "id": "07KB36X6H14M1GFCDQM4C03WS2",
        "name": "Thirteen",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "08KB36X6H14M1GFCDQM4C03WS3": {
        "created_at": "2025-11-27T23:19:50.567890+05:30",
        "description": "This is the fourteenth group",
        "id": "08KB36X6H14M1GFCDQM4C03WS3",
        "name": "Fourteen",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "09KB36X6H14M1GFCDQM4C03WS4": {
        "created_at": "2025-11-27T23:20:00.123456+05:30",
        "description": "This is the fifteenth group",
        "id": "09KB36X6H14M1GFCDQM4C03WS4",
        "name": "Fifteen",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "10KB36X6H14M1GFCDQM4C03WS5": {
        "created_at": "2025-11-27T23:20:10.789012+05:30",
        "description": "This is the sixteenth group",
        "id": "10KB36X6H14M1GFCDQM4C03WS5",
        "name": "Sixteen",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "11KB36X6H14M1GFCDQM4C03WS6": {
        "created_at": "2025-11-27T23:20:20.345678+05:30",
        "description": "This is the seventeenth group",
        "id": "11KB36X6H14M1GFCDQM4C03WS6",
        "name": "Seventeen",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    },
    "12KB36X6H14M1GFCDQM4C03WS7": {
        "created_at": "2025-11-27T23:20:30.901234+05:30",
        "description": "This is the eighteenth group",
        "id": "12KB36X6H14M1GFCDQM4C03WS7",
        "name": "Eighteen",
        "owner_id": "01KAH3YBVGZYPEQZP3CQNB7Q5M"
    }
}

export const useGroupStore = create<GroupStore>((set) => ({
    groups: { ...hardcodedGroups },
    addGroup: (group: Group) => set((state) => ({ groups: { ...state.groups, group } })),
    setGroups: (newGroups: Record<string, Group>) => set({ groups: newGroups })
}))