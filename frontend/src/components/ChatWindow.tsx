import { useEffect, useState } from "react";
import { useChatStore } from "../store";
import { Header } from "./GroupHeader";
import { MessageInput } from "./MessageInput";
import { MessagesList } from "./MessagesList";
import type { Message } from "../../api/message";
import { useCreateMessage } from "../hooks/message";

export function ChatWindow() {
    const selectedGroupId = useChatStore((state) => state.selectedGroupId)
    const currentUserId = useChatStore((state) => state.currentUserId)
    const [currentMessage, setCurrentMessage] = useState("")
    const [liveMessages, setLiveMessages] = useState<Message[]>([])
    const [liveMessageCounter, setLiveMessageCounter] = useState(0)
    const createMessage = useCreateMessage()
    useEffect(() => {
        setLiveMessages([])
    }, [selectedGroupId])

    const onSendCurrentMessage = () => {
        // Don't send if the message is made of only spaces or is empty
        if (currentMessage.trim().length == 0)
            return
        const tempMessageId = `local-${selectedGroupId}-${liveMessageCounter}`
        setLiveMessages(prev => [...prev, {
            content: currentMessage,
            group_id: selectedGroupId,
            type: "text",
            created_at: (new Date()).toISOString(),
            id: tempMessageId,
            sender_id: currentUserId,
            status: "pending"
        }])
        setLiveMessageCounter(prev => prev + 1)
        setCurrentMessage("")
        // Optimistic send
        // TODO: Inform the user if it fails
        console.log("Sending message")
        createMessage.mutate({ groupId: selectedGroupId, content: currentMessage, messageType: "text" },
            {
                onSuccess: (data) => {
                    // TODO: Inefficient, use record type to modify message params quickly
                    setLiveMessages(prev => prev.map(msg => {
                        if (msg.id === tempMessageId) {
                            return { ...msg, status: "sent", id: data.id, created_at: data.created_at }
                        }
                        return msg
                    }))
                },
                onError: () => {
                    setLiveMessages(prev => prev.map(msg => {
                        if (msg.id === tempMessageId) {
                            return { ...msg, status: "error" }
                        }
                        return msg
                    }))
                }
            }
        )
    }
    // If no group is selected, display a blank screen
    if (selectedGroupId === "") {
        return <div className="col-span-8 md:col-span-7 flex flex-col items-center justify-center bg-blue-100">
            <div className="text-xl font-semibold">
                Select a group to view chats
            </div>
        </div>
    }
    return <div className="col-span-8 md:col-span-7 flex flex-col h-screen">
        <Header />
        <MessagesList liveMessages={liveMessages} />
        <MessageInput message={currentMessage} onMessageChange={(e) => setCurrentMessage(e.target.value)} onSubmit={onSendCurrentMessage} />
    </div>
}