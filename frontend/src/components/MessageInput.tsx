import { faPaperPlane } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import type React from "react";

export function MessageInput({ message, onSubmit, onMessageChange }: {
    message: string,
    onSubmit: () => void,
    onMessageChange: (e: React.ChangeEvent<HTMLInputElement>) => void
}) {
    return <div className="p-3 flex flex-row items-center gap-3 bg-blue-200">
        <input type="text" placeholder="Type your message here" className="bg-blue-100 flex-1 p-2 rounded-xl outline-0" value={message} onChange={onMessageChange}
            onKeyDown={(e) => {
                if (e.key === "Enter") onSubmit()
            }}
        />
        <button className="p-2 rounded-full bg-blue-500 transition text-blue-50 hover:bg-blue-400" onClick={onSubmit}>
            <FontAwesomeIcon icon={faPaperPlane} className="w-3 h-3" fontSize={"1.1em"} />
        </button>
    </div>
}
