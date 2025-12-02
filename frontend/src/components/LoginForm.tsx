import { useState } from "react";
import { useLogin } from "../hooks/user";
import { useChatStore } from "../store";

export const LoginForm = () => {
    const [formEmail, setFormEmail] = useState("")
    const [formPassword, setFormPassword] = useState("")
    const loginMutation = useLogin()
    const setIsLoggedIn = useChatStore((state) => state.setIsLoggedIn)

    const handleSubmit = () => {
        loginMutation.mutate({ email: formEmail, password: formPassword }, {
            onError: () => {
                setFormPassword("")
            },
            onSuccess: () => {
                setFormEmail("")
                setFormPassword("")
                setIsLoggedIn(true)
            }
        })
    }

    return <div className="flex flex-col justify-center px-10 py-16 lg:px-12 bg-white rounded-2xl border border-gray-200 shadow-lg w-screen h-screen lg:w-xl lg:h-9/12 hover:shadow-xl transition duration-100">
        <div className="sm:mx-auto sm:w-full sm:max-w-sm">
            <img src="/icons/icon.svg" alt="GoChat" className="mx-auto h-10 w-auto" />
            <h2 className="mt-10 text-center text-2xl/9 font-bold tracking-tight text-gray-900">Sign in to your account</h2>
        </div>

        <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
            <form className="space-y-6" onSubmit={(event) => { event.preventDefault(); handleSubmit(); return false }}>
                <div>
                    <label htmlFor="email" className="block text-sm/6 font-medium text-gray-900">Email address</label>
                    <div className="mt-2">
                        <input value={formEmail} onChange={(event) => { setFormEmail(event.target.value) }} id="email" type="email" name="email" required data-autocomplete="email" placeholder="Enter your email" className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-500 sm:text-sm/6 border border-gray-300" />
                    </div>
                </div>
                <div>
                    <div className="flex items-center justify-between">
                        <label htmlFor="password" className="block text-sm/6 font-medium text-gray-900">Password</label>
                        <div className="text-sm">
                            <a href="#" className="font-semibold text-indigo-600 hover:text-indigo-500">Forgot password?</a>
                        </div>
                    </div>
                    <div className="mt-2">
                        <input value={formPassword} onChange={(event) => { setFormPassword(event.target.value) }} id="password" type="password" name="password" required data-autocomplete="current-password" placeholder="Enter your password" className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-500 sm:text-sm/6 border border-gray-300" />
                    </div>
                </div>
                {loginMutation.isError && (
                    <p className="text-red-600 text-sm mt-2">
                        {(loginMutation.error as any).error}
                    </p>
                )}
                <div>
                    <button type="submit" className="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm/6 font-semibold text-white hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Sign in</button>
                </div>
            </form>
        </div>
    </div>
}