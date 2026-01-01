import { useChatStore } from "../store"
import { getMessages, type Message, type PaginationParams } from "../../api/message"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCheck, faExclamationTriangle, faHourglass } from "@fortawesome/free-solid-svg-icons";
import InfiniteScroll from 'react-infinite-scroll-component';
import { useInfiniteQuery } from "@tanstack/react-query";
import { Loader } from "./Loader";
import { SystemMessage } from "./SystemMessage";
import { useEffect, useRef } from "react";
import type { GroupMember } from "../../api/group";
import { useGroupMembers } from "../hooks/group";
import { queryClient } from "../../api/query-client";

const defaultMessageFetchLimit = 20;

const queryFn = async ({ pageParam }: { pageParam?: PaginationParams }) => {
    if (!pageParam) {
        pageParam = { before: "", groupId: "", limit: defaultMessageFetchLimit }
    }
    const res = await getMessages(pageParam)
    return res
}


function ChatMessage({ message, memberMap }: { message: Message, memberMap?: Record<string, GroupMember> }) {
    const currentUserId = useChatStore((state) => state.currentUserId)
    // TODO: Check if it's empty
    let senderIsCurrentUser = false
    if (currentUserId === message.sender_id) {
        senderIsCurrentUser = true
    }

    // If the status is not present (i.e. it's from the server, it's sent)
    const isSent = (!message.status) || message.status === "sent"
    const isError = message.status === "error"

    const senderName = memberMap?.[message.sender_id]?.name || message.sender_id;
    const username = memberMap?.[message.sender_id]?.username || message.sender_id;

    // TODO: Later map message sender id to username
    return <div className={`p-3 ${senderIsCurrentUser ? "bg-green-100 self-end" : "bg-slate-100"} mb-2 rounded-lg w-fit lg:max-w-6/12 md:max-w-8/12 sm:max-w-10/12 max-w-11/12 shadow-sm`}>
        {senderIsCurrentUser ? <></> :
            <div className="flex flex-row items-center mb-1">
                <a className="font-semibold text-blue-700 transition duration-100 hover:text-blue-900 hover:font-bold hover:underline hover:cursor-pointer">
                    {senderName}
                </a>
                <a className="font-light text-sm ml-3 text-gray-500 hover:text-gray-800 duration-100 transition">
                    @{username}
                </a>
            </div>
        }
        <p>
            {message.content}
        </p>
        <div className="flex flex-row items-center justify-end font-light text-gray-500 text-sm">
            <p className="mr-2">
                {new Date(message.created_at).toLocaleString(undefined, {
                    month: 'short',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                })}
            </p>
            {senderIsCurrentUser && (
                <FontAwesomeIcon
                    icon={isError ? faExclamationTriangle : (isSent ? faCheck : faHourglass)}
                    fontSize={"0.8em"}
                    className={isError ? "text-red-500" : ""}
                />
            )}
        </div>
    </div>
}

export function MessagesList({ liveMessages, forceScrollToEnd }: { liveMessages: Message[], forceScrollToEnd: boolean }) {
    // Temporary workaround flex-col-reverse, later make the div scroll to the end
    const scrollRef = useRef<HTMLDivElement>(null);
    useEffect(() => {
        const messagesDiv = scrollRef.current;
        if (!messagesDiv) return;

        // Threshold: How close to the bottom must we be to auto-scroll?
        const SCROLL_MESSAGES_DISTANCE = 300;

        const isNearBottom =
            messagesDiv.scrollTop + messagesDiv.clientHeight >= messagesDiv.scrollHeight - SCROLL_MESSAGES_DISTANCE;

        if (isNearBottom || forceScrollToEnd) {
            messagesDiv.scrollTo({
                top: messagesDiv.scrollHeight,
                behavior: 'smooth'
            });
        }
    }, [liveMessages, forceScrollToEnd]);
    return <div id="chatMessages" className="flex-1 p-5 overflow-y-scroll flex flex-col-reverse" ref={scrollRef}>
        <InfiniteList liveMessages={liveMessages} />
    </div>
}

function InfiniteList({ liveMessages }: { liveMessages: Message[] }) {
    const selectedGroupId = useChatStore((state) => state.selectedGroupId)
    const { data: memberMap } = useGroupMembers(selectedGroupId);
    const { data, error, fetchNextPage, hasNextPage, status } = useInfiniteQuery({
        initialPageParam: { before: "", groupId: selectedGroupId, limit: defaultMessageFetchLimit },
        queryKey: ["groups", selectedGroupId, "messages"],
        queryFn: queryFn,
        getNextPageParam: (lastPage, _) => {
            if (!lastPage.cursor.has_before)
                return undefined
            return {
                before: lastPage.cursor.before,
                limit: defaultMessageFetchLimit,
                groupId: selectedGroupId
            }
        },
        refetchOnWindowFocus: false,
        refetchOnMount: false,
        refetchOnReconnect: false,
        enabled: selectedGroupId != "",
    })
    useEffect(() => {
        return () => {
            if (!data)
                return
            queryClient.setQueryData(['groups', selectedGroupId, "messages"], (data: any) => ({
                pages: data.pages.slice(0, 1),
                pageParams: data.pageParams.slice(0, 1),
            }))
        }
    }, [selectedGroupId])
    const historyMessages = data ? data.pages.flatMap(p => p.messages) : []

    if (status === 'pending') {
        return <div className="justify-self-center self-center">
            <Loader />
        </div>
    }
    if (status === 'error') {
        return <p className="text-red-600">Error occured while fetching messages {error.message}</p>
    }

    const totalDataLength = historyMessages.length + liveMessages.length;

    if (totalDataLength === 0) {
        return (
            <div className="self-center justify-end"><SystemMessage text="No messages" /></div>
        )
    }

    return (
        <InfiniteScroll
            dataLength={totalDataLength}
            next={() => fetchNextPage()}
            hasMore={hasNextPage}
            loader={<div className="self-center justify-self-center"><Loader /></div>}
            endMessage={<div className="self-center"><SystemMessage text="No more messages" /></div>}
            scrollableTarget="chatMessages"
            inverse={true}
            className="flex-col-reverse flex"
        >
            <div className="flex flex-col">
                {liveMessages.map(message => (
                    <ChatMessage message={message} key={message.id} memberMap={memberMap} />
                ))}
            </div>
            <div className="flex flex-col-reverse">
                {historyMessages.map(message => (
                    <ChatMessage message={message} key={message.id} memberMap={memberMap} />
                ))}
            </div>
        </InfiniteScroll>
    )
}