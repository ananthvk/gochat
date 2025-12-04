import { useChatStore } from "../store"
import { getMessages, type Message, type PaginationParams } from "../../api/message"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCheck, faExclamationTriangle, faHourglass } from "@fortawesome/free-solid-svg-icons";
import InfiniteScroll from 'react-infinite-scroll-component';
import { formatChatTime } from "../lib/formatChatTime";
import { useInfiniteQuery } from "@tanstack/react-query";
import { Loader } from "./Loader";
import { SystemMessage } from "./SystemMessage";
import { useEffect, useRef } from "react";

const defaultMessageFetchLimit = 20;

const queryFn = async ({ pageParam }: { pageParam?: PaginationParams }) => {
    if (!pageParam) {
        pageParam = { before: "", groupId: "", limit: defaultMessageFetchLimit }
    }
    const res = await getMessages(pageParam)
    return res
}


function ChatMessage(message: Message) {
    const currentUserId = useChatStore((state) => state.currentUserId)
    // TODO: Check if it's empty
    let senderIsCurrentUser = false
    if (currentUserId === message.sender_id) {
        senderIsCurrentUser = true
    }

    // If the status is not present (i.e. it's from the server, it's sent)
    const isSent = (!message.status) || message.status === "sent"
    const isError = message.status === "error"

    // TODO: Later map message sender id to username
    return <div className={`p-3 ${senderIsCurrentUser ? "bg-green-100 self-end" : "bg-slate-100"} mb-2 rounded-lg w-fit lg:max-w-6/12 md:max-w-8/12 sm:max-w-10/12 max-w-11/12 shadow-sm`}>
        {senderIsCurrentUser ? <></> :
            <a className="font-semibold text-slate-900 hover:text-slate-600 transition duration-100">
                {message.sender_id}
            </a>
        }
        <p>
            {message.content}
        </p>
        <div className="flex flex-row items-center justify-end font-light text-gray-500 text-sm">
            <p className="mr-2">
                {formatChatTime(message.created_at)}
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

export function MessagesList({ liveMessages }: { liveMessages: Message[] }) {
    // Temporary workaround flex-col-reverse, later make the div scroll to the end
    const scrollRef = useRef<HTMLDivElement>(null);
    useEffect(() => {
        const messagesDiv = scrollRef.current;
        if (!messagesDiv) return;

        // Threshold: How close to the bottom must we be to auto-scroll?
        const SCROLL_MESSAGES_DISTANCE = 300;

        const isNearBottom =
            messagesDiv.scrollTop + messagesDiv.clientHeight >= messagesDiv.scrollHeight - SCROLL_MESSAGES_DISTANCE;

        if (isNearBottom) {
            messagesDiv.scrollTo({
                top: messagesDiv.scrollHeight,
                behavior: 'smooth'
            });
        }
    }, [liveMessages]);
    return <div id="chatMessages" className="flex-1 p-5 overflow-y-scroll flex flex-col-reverse" ref={scrollRef}>
        <InfiniteList liveMessages={liveMessages} />
    </div>
}

function InfiniteList({ liveMessages }: { liveMessages: Message[] }) {
    const selectedGroupId = useChatStore((state) => state.selectedGroupId)
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
        enabled: selectedGroupId != ""
    })
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
                    <ChatMessage {...message} key={message.id || `live-${message.created_at}`} />
                ))}
            </div>
            <div className="flex flex-col-reverse">
                {historyMessages.map(message => (
                    <ChatMessage {...message} key={message.id} />
                ))}
            </div>
        </InfiniteScroll>
    )
}