import { useInfiniteQuery, type QueryFunction } from "@tanstack/react-query"
import { useChatStore } from "../store"
import { getMessages, type MessageResult } from "../../api/message"
import type { APIError } from "../../api/errors";
import React from "react";
import { Loader } from "./Loader";

const messageLimit = 100

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

export function MessagesList() {
    return <div className="flex-1 p-5 bg-blue-50 overflow-y-auto">
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
    console.log(data)

    return (status === 'pending' ? (
        <Loader />
    ) : status === 'error' ? (
        <p>Error: {error.error}</p>
    ) : (
        <div>
            {((data as any).pages.map((page: any, i: any) => (
                <React.Fragment key={i}>
                    {page.messages.map((message: any) => (
                        <p key={message.id}>{message.content}</p>
                    ))}
                </React.Fragment>
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
