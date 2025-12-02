import { LoginForm } from "./LoginForm";

export function LoginPage() {
    return <div className="h-screen md:grid md:grid-cols-2">
        <div className="hidden md:justify-center md:items-center md:flex md:flex-col bg-linear-to-br to-blue-950 from-cyan-950">
            <h1 className="text-4xl font-bold text-white mb-4 drop-shadow-lg">
                Welcome to GoChat
            </h1>
            <p className="text-xl font-medium text-white drop-shadow-md">
                Connect, share, and chat with your community
            </p>
        </div>
        <div className="bg-gray-50 h-screen flex items-center justify-center">
            <LoginForm />
        </div>
    </div>
}