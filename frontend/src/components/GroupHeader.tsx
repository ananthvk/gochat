import { faUserGroup } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useChatStore } from "../store";
import { useGroupMembers, useGroups } from "../hooks/group";
import { Loader } from "./Loader";
import { useEffect } from "react";

export function Header() {
    const { data: groups, isLoading, isError } = useGroups()
    const selectedGroupId = useChatStore((state) => state.selectedGroupId)
    const setGroupId = useChatStore((state) => state.setSelectedGroupId)
    const members = useGroupMembers(selectedGroupId)
    useEffect(() => {
        if (groups && selectedGroupId && !(selectedGroupId in groups)) {
            setGroupId("")
        }
    }, [groups, selectedGroupId, setGroupId])

    if (selectedGroupId === "") {
        console.log("selected group is empty")
        return <div>
            An error occured
        </div>
    }
    if (isLoading) {
        return <Loader />
    }
    if (isError) {
        return <div className="text-red-600">Error occured while fetching group</div>
    }
    if (!groups)
        return null
    if (!selectedGroupId || !(selectedGroupId in groups)) {
        return <div></div>
    }
    return <div className="p-4 flex flex-row items-center from-blue-500 to-blue-600 bg-linear-to-l">
        <button className="text-white p-3 rounded-full items-start bg-gray-400 mr-3">
            <FontAwesomeIcon icon={faUserGroup} fontSize={"1.5em"} />
        </button>
        <div>
            <p className="text-xl font-bold text-white">
                {groups[selectedGroupId].name}
            </p>
            <div className="flex flex-row text-xs text-gray-50 whitespace-nowrap overflow-hidden">
                {members.data ? Object.keys(members.data).slice(0, 5).map((memberId, index, array) =>
                    <p key={memberId} className="truncate">
                        {members.data[memberId].name}{index < array.length - 1 ? ', ' : ''}
                    </p>
                ) : <></>}
                {members.data && Object.keys(members.data).length > 5 && <p>, ...</p>}
            </div>
        </div>
    </div>
}
