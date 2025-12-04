import './App.css'
import { ChatApp } from './components/ChatApp'
import { useChatStore } from './store'
import { LoginPage } from './components/LoginPage'
import { useAuthBootstrap } from './hooks/user'
import { Loader } from './components/Loader'
function App() {
    useAuthBootstrap();
    const isLoggedIn = useChatStore((state) => state.isLoggedIn)
    const authLoading = useChatStore((state) => state.authLoading)
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