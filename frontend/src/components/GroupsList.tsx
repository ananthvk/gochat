import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { useChatStore } from "../store"
import { faUser } from "@fortawesome/free-solid-svg-icons"
import { useGroups } from "../hooks/group"
import { Loader } from "./Loader"
import { useEffect } from "react"
import type { Group, GroupResult } from "../../api/group"
import { formatMessageTime } from "../lib/time"
import { queryClient } from "../../api/query-client"

function Group({ group, selected }: { group: Group, selected: boolean }) {
    const setGroup = useChatStore((state) => state.setSelectedGroupId)
    const currentUserId = useChatStore((state) => state.currentUserId)
    return <div className={`${selected ? "bg-gray-200" : "bg-gray-50"} p-3 mt-1 rounded-xl flex flex-row items-center hover:bg-gray-200 transition duration-75`} onClick={
        () => setGroup(group.id)}>
        <button className="text-white p-2 rounded-full items-start bg-gray-400 mr-3">
            <FontAwesomeIcon icon={faUser} fontSize={"1.3em"} />
        </button>
        <div className="flex flex-col flex-1">
            <div className="flex flex-row items-center justify-between w-full">
                <p className="text-xl font-bold">
                    {group.name}
                </p>
                <p className="text-sm text-gray-400">
                    {group.last_message ? formatMessageTime(group.last_message.created_at) : <></>}
                </p>
            </div>
            {group.last_message && group.last_message.type === "text" ?
                <p className="text-base text-gray-600">
                    {group.last_message.sender_id === currentUserId ? "You" : group.last_message.sender_name}: {group.last_message.content}
                </p>
                : <p className="text-base text-gray-500">No messages yet</p>}
        </div>
    </div>
}

export function GroupsList() {
    const selectedGroupId = useChatStore((state) => state.selectedGroupId)
    const setSelectedGroupId = useChatStore((state) => state.setSelectedGroupId)
    const { data: groups, isLoading, isError } = useGroups()

    useEffect(() => {
        if (groups && selectedGroupId && !(selectedGroupId in groups)) {
            setSelectedGroupId("")
        }
    }, [groups, selectedGroupId, setSelectedGroupId])
    useEffect(() => {
        const handleWSMessage = (e: Event) => {
            const msg = (e as any).detail
            if (msg.type === "text_message" && msg.payload?.group_id) {
                queryClient.setQueryData<GroupResult>(
                    ["groups"],
                    (old) => {
                        // TODO: Inefficient, make the api wrapper functions in api/ to return Record instead of lists
                        if (!old) return old
                        let oldGroup: any = null
                        let groups: Group[] = []
                        if (old.groups.length > 0) {
                            groups.push(old.groups[0])
                        }
                        old.groups.forEach((x) => {
                            if (x.id !== msg.payload.group_id) {
                                groups.push(x)
                            } else {
                                oldGroup = x
                            }
                        })
                        if (!oldGroup) {
                            console.log("logic error, group not found in list")
                            return old
                        }
                        // TODO: Fetch the actual sender name
                        const updatedGroup = { ...oldGroup, last_message: { sender_name: "*New message", ...msg.payload } }
                        if (old.groups.length > 0) {
                            groups[0] = updatedGroup
                        } else {
                            groups.push(updatedGroup)
                        }
                        return { groups: groups }
                    }
                )
            }
        }
        window.addEventListener("ws-message", handleWSMessage)
        return () => window.removeEventListener("ws-message", handleWSMessage)
    }, [])

    if (isLoading) {
        return <div className="flex items-center justify-center flex-col h-screen">
            <Loader />
        </div>
    }
    if (isError) {
        return <div className="text-red-600">Unable to fetch groups</div>
    }
    if (!groups)
        return null;
    if (Object.keys(groups).length == 0)
        return <div className="flex flex-row justify-center items-center h-screen font-semibold">You have not joined any groups</div>
    return <div className="flex-1 bg-radial white overflow-y-scroll">
        <div>
            {
                Object.keys(groups).map(grp => <Group group={groups[grp]} selected={selectedGroupId === grp} key={groups[grp].id} />)
            }
        </div>
    </div>
}