import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { useChatStore, useGroupStore } from "../store"
import { faUser } from "@fortawesome/free-solid-svg-icons"

function Group({ groupId, selected }: { groupId: string, selected: boolean }) {
    const setGroup = useChatStore((state) => state.setSelectedGroupId)
    const groups = useGroupStore((state) => state.groups)
    return <div className={`${selected ? "bg-gray-200" : "bg-gray-50"} p-3 mt-1 rounded-xl flex flex-row items-center hover:bg-gray-200 transition duration-75`} onClick={
        () => setGroup(groupId)}>
        <button className="text-white p-2 rounded-full items-start bg-gray-400 mr-3">
            <FontAwesomeIcon icon={faUser} fontSize={"1.3em"} />
        </button>
        <div>
            <p className="text-xl font-bold">
                {groups[groupId].name}
            </p>
            <p className="text-base text-gray-600">
                Last message here
            </p>
        </div>
    </div>
}

export function GroupsList() {
    const selectedGroupId = useChatStore((state) => state.selectedGroupId)
    const groups = useGroupStore((state) => state.groups)
    return <div className="flex-1 bg-radial white overflow-y-scroll">
        <ul>
            {
                Object.keys(groups).map(grp => <Group groupId={grp} key={grp} selected={selectedGroupId === grp} />)
            }
        </ul>
    </div>
}