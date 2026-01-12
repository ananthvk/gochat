declare global {
    interface Window {
        ws?: WebSocket
        wsManualClosed?: boolean
    }
}

const RECONNECT_INTERVAL = 20000 // 20s

if (!(import.meta as any).env.VITE_API_BASE_URL) {
    throw new Error("VITE_API_BASE_URL not set")
}

function convertToWebSocket(url: string) {
    if (url.startsWith("http://")) return url.replace("http://", "ws://")
    if (url.startsWith("https://")) return url.replace("https://", "wss://")
    throw new Error("Invalid URL, should start with http:// or https://")
}

export function initWS() {
    if (window.ws) return window.ws

    const baseURL = (import.meta as any).env.VITE_API_BASE_URL
    const token = localStorage.getItem("session_token")
    if (!token) return

    window.wsManualClosed = false

    const ws = new WebSocket(`${convertToWebSocket(baseURL)}/realtime/ws?token=${token}`)

    ws.onopen = () => console.log("WS connected")

    ws.onclose = () => {
        console.log("WS closed")

        window.ws = undefined

        if (!window.wsManualClosed) {
            setTimeout(() => {
                console.log("Reconnecting WS...")
                initWS()
            }, RECONNECT_INTERVAL)
        }
    }

    ws.onerror = (err) => {
        console.error("WS error", err)
        ws.close()
    }

    ws.onmessage = (e) => {
        try {
            const msg = JSON.parse(e.data)
            window.dispatchEvent(new CustomEvent("ws-message", { detail: msg }))
        } catch (err) {
            console.error("Failed to parse WS message", err)
        }
    }

    window.ws = ws
    return ws
}

export function closeWS() {
    window.wsManualClosed = true
    window.ws?.close()
    window.ws = undefined
}
