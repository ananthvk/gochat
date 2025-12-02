import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { useChatStore } from "../store"
import { faUser } from "@fortawesome/free-solid-svg-icons"
import { useGroups } from "../hooks/group"
import { Loader } from "./Loader"
import { useEffect } from "react"

function Group({ groupId, selected, name }: { groupId: string, selected: boolean, name: string }) {
    const setGroup = useChatStore((state) => state.setSelectedGroupId)
    return <div className={`${selected ? "bg-gray-200" : "bg-gray-50"} p-3 mt-1 rounded-xl flex flex-row items-center hover:bg-gray-200 transition duration-75`} onClick={
        () => setGroup(groupId)}>
        <button className="text-white p-2 rounded-full items-start bg-gray-400 mr-3">
            <FontAwesomeIcon icon={faUser} fontSize={"1.3em"} />
        </button>
        <div>
            <p className="text-xl font-bold">
                {name}
            </p>
            <p className="text-base text-gray-600">
                Last message here
            </p>
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
        <ul>
            {
                Object.keys(groups).map(grp => <Group groupId={grp} key={grp} selected={selectedGroupId === grp} name={groups[grp].name} />)
            }
        </ul>
    </div>
}