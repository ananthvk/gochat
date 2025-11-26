import { faUserGroup } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

export function Header() {
    return <div style={{ backgroundColor: 'yellow' }} className="p-4 flex flex-row items-center">
        <button className="hover:bg-slate-50 text-white hover:text-gray-300 transition duration-300 p-3 rounded-full items-start bg-gray-400 mr-3">
            <FontAwesomeIcon icon={faUserGroup} fontSize={"1.5em"} />
        </button>
        <p className="text-xl font-bold">
            Group Name
        </p>
    </div>
}
