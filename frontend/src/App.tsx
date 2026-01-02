import './App.css'
import { ChatApp } from './components/ChatApp'
import { useChatStore } from './store'
import { LoginPage } from './components/LoginPage'
import { useAuthBootstrap } from './hooks/user'
import { Loader } from './components/Loader'
import { initWS, closeWS } from '../api/ws'
import { useEffect } from 'react'
function App() {
    useAuthBootstrap();
    const isLoggedIn = useChatStore((state) => state.isLoggedIn)
    const authLoading = useChatStore((state) => state.authLoading)
    useEffect(() => {
        if (isLoggedIn) {
            initWS()
        }
        if (!isLoggedIn) {
            closeWS()
        }
    }, [isLoggedIn])

    // Only the root app component handles message demultiplexing
    useEffect(() => {
        function handler(e: Event) {
            const msg = (e as CustomEvent).detail
            //if (msg.type === "message") {
            console.log("new message", msg)
            //}
        }
        window.addEventListener("ws-message", handler)
        return () => window.removeEventListener("ws-message", handler)
    }, [])

    if (authLoading) {
        return <div className="text-center font-semibold p-5 h-screen flex items-center justify-center flex-col">
            <div className="mb-3">
                Logging in...
            </div>
            <Loader />
        </div>
    }
    if (isLoggedIn) {
        return <ChatApp />
    } else {
        return <LoginPage />
    }

}
export default App