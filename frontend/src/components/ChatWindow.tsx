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
    const [liveMessages, setLiveMessages] = useState<Message[]>([])
    const [liveMessageCounter, setLiveMessageCounter] = useState(0)
    const createMessage = useCreateMessage()
    const [forceScrollToEnd, setForceScrollToEnd] = useState(false)
    useEffect(() => {
        setLiveMessages([])
    }, [selectedGroupId])

    useEffect(() => {
        const handleWSMessage = (e: Event) => {
            const msg = (e as any).detail
            // TODO: If the message received was sent by this user, use it as an ack, or for updating chat state when used cross device
            
            if (msg.type === "text_message" && msg.payload?.group_id === selectedGroupId) {
                const payload = msg.payload
                if(payload.sender_id === currentUserId) {
                    return
                }
                // Check if it already exists, to avoid duplicates
                setLiveMessages(prev => {
                    const exists = prev.some(m => m.id === payload.id)
                    if (exists) return prev
                    return [...prev, payload]
                })
            }
        }

        window.addEventListener("ws-message", handleWSMessage)
        return () => window.removeEventListener("ws-message", handleWSMessage)
    }, [selectedGroupId])

    const onSendCurrentMessage = (message: string) => {
        const tempMessageId = `local-${selectedGroupId}-${liveMessageCounter}`
        setLiveMessages(prev => [...prev, {
            content: message,
            group_id: selectedGroupId,
            type: "text",
            created_at: (new Date()).toISOString(),
            id: tempMessageId,
            sender_id: currentUserId,
            status: "pending"
        }])
        setLiveMessageCounter(prev => prev + 1)

        // Optimistic send
        console.log("Sending message")
        setForceScrollToEnd(true)
        createMessage.mutate({ groupId: selectedGroupId, content: message, messageType: "text" },
            {
                onSuccess: (data) => {
                    // TODO: Inefficient, use record type to modify message params quickly
                    setLiveMessages(prev => prev.map(msg => {
                        if (msg.id === tempMessageId) {
                            return { ...msg, status: "sent", id: data.id, created_at: data.created_at }
                        }
                        return msg
                    }))
                    setForceScrollToEnd(false)
                },
                onError: () => {
                    setLiveMessages(prev => prev.map(msg => {
                        if (msg.id === tempMessageId) {
                            return { ...msg, status: "error" }
                        }
                        return msg
                    }))
                    setForceScrollToEnd(false)
                }
            }
        )
    }
    // If no group is selected, display a blank screen
    if (selectedGroupId === "") {
        return <div className={`col-span-0 md:col-span-7 flex flex-col items-center justify-center bg-blue-100`}>
            <div className="text-xl font-semibold">
                Select a group to view chats
            </div>
        </div>
    }
    return <div className={`col-span-0 md:col-span-7 flex flex-col h-screen`}>
        <Header />
        <MessagesList liveMessages={liveMessages} forceScrollToEnd={forceScrollToEnd} />
        <MessageInput onSubmit={onSendCurrentMessage} />
    </div>
}