import { useInfiniteQuery, type QueryFunction } from "@tanstack/react-query"
import { useChatStore } from "../store"
import { getMessages, type Message, type MessageResult } from "../../api/message"
import type { APIError } from "../../api/errors";
import { Loader } from "./Loader";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCheck } from "@fortawesome/free-solid-svg-icons";

const messageLimit = 50

const getMessagesQueryFn: QueryFunction<
    MessageResult,
    [string, string, string],
    string
> = async ({ pageParam, queryKey }) => {
    const [, groupId] = queryKey;

    return getMessages({
        groupId,
        before: pageParam ?? "",
        limit: messageLimit,
    });
};
function formatChatTime(timestamp: string) {
    const date = new Date(timestamp);
    const now = new Date();

    const diffMs = now.getTime() - date.getTime();
    const diffSec = Math.floor(diffMs / 1000);
    const diffMin = Math.floor(diffSec / 60);
    const diffHour = Math.floor(diffMin / 60);

    const isToday =
        now.getDate() === date.getDate() &&
        now.getMonth() === date.getMonth() &&
        now.getFullYear() === date.getFullYear();

    const yesterday = new Date(now);
    yesterday.setDate(now.getDate() - 1);
    const isYesterday =
        yesterday.getDate() === date.getDate() &&
        yesterday.getMonth() === date.getMonth() &&
        yesterday.getFullYear() === date.getFullYear();

    if (diffSec < 10) return "Just now";
    if (diffMin < 1) return `${diffSec} sec ago`;
    if (diffMin < 60) return `${diffMin} min ago`;
    if (diffHour < 24 && isToday)
        return date.toLocaleTimeString("en-US", {
            hour: "2-digit",
            minute: "2-digit",
        });
    if (isYesterday)
        return `Yesterday ${date.toLocaleTimeString("en-US", {
            hour: "2-digit",
            minute: "2-digit",
        })}`;
    // older than yesterday
    return date.toLocaleDateString("en-US") + " " + date.toLocaleTimeString("en-US", { hour: "2-digit", minute: "2-digit" });
}


function ChatMessage(message: Message) {
    const currentUserId = useChatStore((state) => state.currentUserId)
    // TODO: Check if it's empty
    let senderIsCurrentUser = false
    if (currentUserId === message.sender_id) {
        senderIsCurrentUser = true
    }

    // TODO: Later map message sender id to username
    return <div className={`p-3 ${senderIsCurrentUser ? "bg-green-100 self-end" : "bg-slate-100"} align- mb-2 rounded-lg w-fit lg:max-w-6/12 md:max-w-8/12 sm:max-w-10/12 max-w-11/12 shadow-sm`}>
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
            {senderIsCurrentUser ? <FontAwesomeIcon icon={faCheck} fontSize={"0.8em"} /> : <></>}
        </div>
    </div>
}

export function MessagesList() {
    return <div className="flex-1 p-5 overflow-y-auto">
        <InfiniteList />
    </div>
}

function InfiniteList() {
    const selectedGroupId = useChatStore((state) => state.selectedGroupId);

    const { data,
        error,
        fetchNextPage,
        hasNextPage,
        isFetching,
        isFetchingNextPage,
        status,
    } = useInfiniteQuery<
        MessageResult,
        APIError,
        MessageResult,
        [string, string, string],
        string
    >({
        queryKey: ["groups", selectedGroupId, "messages"],
        queryFn: getMessagesQueryFn,
        initialPageParam: "",
        getNextPageParam: (lastPage) => lastPage.cursor.has_before ? lastPage.cursor.before : undefined,
        refetchOnWindowFocus: false,
        refetchOnReconnect: false,
        refetchOnMount: false,
        // Keep messages cached indefinitely
        staleTime: Infinity,
    });

    return (status === 'pending' ? (
        <Loader />
    ) : status === 'error' ? (
        <p>Error: {error.error}</p>
    ) : (
        <div>
            {((data as any).pages.map((page: any, i: any) => (
                <div key={i} className="flex flex-col">
                    {page.messages.map((message: Message) => (
                        <ChatMessage key={message.id} {...message} />
                    ))}
                </div>
            )))}
            <div>{isFetching && !isFetchingNextPage ? 'Fetching...' : null}</div>
            <div>
                <button
                    onClick={() => fetchNextPage()}
                    disabled={!hasNextPage || isFetching}
                >
                    {isFetchingNextPage
                        ? 'Loading more...'
                        : hasNextPage
                            ? 'Load More'
                            : 'Nothing more to load'}
                </button>
            </div>
        </div>
    )
    )
}
