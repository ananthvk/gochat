import { faGear, faUser } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

export function SidebarHeader() {
    return <div className="p-5 flex flex-row items-center justify-between  bg-linear-to-r from-blue-50 to-gray-50">
        <div>
            <button className="text-white p-3 rounded-full items-start bg-gray-400 mr-3">
                <FontAwesomeIcon icon={faUser} fontSize={"1.5em"} />
            </button>
            <button className="hover:bg-slate-50 transition duration-300 p-3 rounded-full items-start">
                <FontAwesomeIcon icon={faGear} fontSize={"1.5em"} color="black" />
            </button>
        </div>
        <div className="font-extrabold text-4xl from-blue-600 to-blue-900 bg-linear-to-r bg-clip-text text-transparent items-end">
            GoChat
        </div>
    </div>
}