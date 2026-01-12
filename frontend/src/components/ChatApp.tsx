// This is the main component that contains the web app.

import { ChatWindow } from "./ChatWindow";
import { Sidebar } from "./Sidebar";

// If the user is authenticated, they can access the chat application
export function ChatApp() {
    // TODO: Implement sidebar hide /unhide
    return <div className="gap-x-1 min-h-screen grid grid-cols-10">
        <Sidebar />
        <ChatWindow />
    </div>
}