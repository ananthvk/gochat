import { useChatStore } from "../store";
import { Header } from "./GroupHeader";
import { MessageInput } from "./MessageInput";
import { MessagesList } from "./MessagesList";

export function ChatWindow() {
    const isGroupSelected = useChatStore((state) => state.selectedGroupId)
    // If no group is selected, display a blank screen
    if (isGroupSelected === "") {
        return <div className="col-span-8 md:col-span-7 flex flex-col items-center justify-center bg-linear-to-tr from-cyan-500 from-0% via-purple-300 via-50% to-blue-100 to-100%">
            <div className="text-xl font-semibold">
                Select a group to view chats
            </div>
        </div>
    }
    return <div className="col-span-8 md:col-span-7 flex flex-col">
        <Header />
        <MessagesList />
        <MessageInput />
    </div>
}
