import { useChatStore } from "../store";
import { Header } from "./GroupHeader";
import { MessageInput } from "./MessageInput";
import { MessagesList } from "./MessagesList";

export function ChatWindow() {
    const isGroupSelected = useChatStore((state) => state.selectedGroupId)
    // If no group is selected, display a blank screen
    if (isGroupSelected === "") {
        return <div className="col-span-8 md:col-span-7 flex flex-col items-center justify-center bg-blue-100">
            <div className="text-xl font-semibold">
                Select a group to view chats
            </div>
        </div>
    }
    return <div className="col-span-8 md:col-span-7 flex flex-col h-screen">
        <Header />
        <MessagesList />
        <MessageInput />
    </div>
}
