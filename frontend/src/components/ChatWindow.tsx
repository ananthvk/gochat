import { Header } from "./GroupHeader";
import { MessageInput } from "./MessageInput";
import { MessagesList } from "./MessagesList";

export function ChatWindow() {
    return <div style={{ backgroundColor: 'cyan' }} className="col-span-8 md:col-span-7 flex flex-col">
        <Header />
        <MessagesList />
        <MessageInput/>
    </div>
}
