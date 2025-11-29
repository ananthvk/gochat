import { faUserGroup } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useChatStore, useGroupStore } from "../store";

export function Header() {
    const groups = useGroupStore((state) => state.groups)
    const selectedGroupId = useChatStore((state) => state.selectedGroupId)
    if (selectedGroupId === "") {
        console.log("selected group is empty")
        return <div>
            An error occured
        </div>
    }
    return <div className="p-4 flex flex-row items-center from-blue-500 to-blue-600 bg-linear-to-l">
        <button className="text-white p-3 rounded-full items-start bg-gray-400 mr-3">
            <FontAwesomeIcon icon={faUserGroup} fontSize={"1.5em"} />
        </button>
        <p className="text-xl font-bold text-white">
            {groups[selectedGroupId].name}
        </p>
    </div>
}
