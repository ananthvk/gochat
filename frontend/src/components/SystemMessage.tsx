import { faInfoCircle } from "@fortawesome/free-solid-svg-icons"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"

type MessageType = 'info' | 'alert'

export function SystemMessage({ text, messageType = 'info' }: { text: string, messageType?: MessageType }) {

    return <div className={`${messageType === 'info' ? 'bg-blue-200' : 'bg-red-200'} p-2 mt-3 rounded-xl font-semibold flex flex-row items-center shadow-sm mb-3`}>
        <FontAwesomeIcon icon={faInfoCircle} className="mr-2 text-gray-700" />
        <p>
            {text}
        </p>
    </div>
}