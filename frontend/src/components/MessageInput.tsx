import { faPaperPlane } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useState } from "react";

export function MessageInput({ onSubmit }: {
    onSubmit: (message: string) => void,
}) {
    const [text, setText] = useState("")
    const handleSubmit = () => {
        if (text.trim().length === 0)
            return
        onSubmit(text.trim())
        setText("")
    }
    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key == 'Enter') {
            e.preventDefault()
            handleSubmit()
        }
    }
    return <div className="p-3 flex flex-row items-center gap-3 bg-blue-200">
        <input type="text" placeholder="Type your message here" className="bg-blue-100 flex-1 p-2 rounded-xl outline-0" value={text} onChange={(e) => setText(e.target.value)}
            onKeyDown={handleKeyDown}
        />
        <button className="p-2 rounded-full bg-blue-500 transition text-blue-50 hover:bg-blue-400" onClick={handleSubmit}>
            <FontAwesomeIcon icon={faPaperPlane} className="w-3 h-3" fontSize={"1.1em"} />
        </button>
    </div>
}
