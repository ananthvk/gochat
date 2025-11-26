import { faPaperPlane } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

export function MessageInput() {
    return <div style={{ backgroundColor: 'orange' }} className="p-3 flex flex-row items-center gap-3">
        <input type="text" placeholder="Type your message here" className="bg-amber-400 flex-1 p-2 rounded-xl outline-0" />
        <button className="p-2 rounded-full hover:bg-orange-500 transition">
            <FontAwesomeIcon icon={faPaperPlane} className="w-5 h-5"/>
        </button>
    </div>
}
