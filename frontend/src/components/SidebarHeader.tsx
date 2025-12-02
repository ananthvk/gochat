import { faGear, faUser } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useChatStore } from "../store";
import { logOut } from "../../api/auth";

export function SidebarHeader() {
    const isLoggedIn = useChatStore((state) => state.isLoggedIn)
    const setIsLoggedIn = useChatStore((state) => state.setIsLoggedIn)
    const logout = () => {
        logOut().catch(err => console.log(err))
        setIsLoggedIn(false)
    }
    return <div className="p-5 flex flex-row items-center justify-between  bg-linear-to-r from-blue-50 to-gray-50">
        <div>
            <button className="text-white p-3 rounded-full items-start bg-gray-400 mr-3">
                <FontAwesomeIcon icon={faUser} fontSize={"1.5em"} />
            </button>
            <button className="hover:bg-slate-50 transition duration-300 p-3 rounded-full items-start">
                <FontAwesomeIcon icon={faGear} fontSize={"1.5em"} color="black" />
            </button>
            {isLoggedIn ? <button className="p-3 bg-slate-200 rounded-2xl font-semibold hover:bg-slate-300 transition duration-150" onClick={logout}>Logout</button> : <></>}
        </div>
        <div className="font-extrabold text-4xl from-blue-600 to-blue-900 bg-linear-to-r bg-clip-text text-transparent items-end">
            GoChat
        </div>
    </div>
}