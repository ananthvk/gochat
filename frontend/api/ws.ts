declare global {
    interface Window {
        ws?: WebSocket
    }
}

const RECONNECT_INTERVAL = 15000 // 15s

if (!(import.meta as any).env.VITE_API_BASE_URL) {
    throw new Error("VITE_API_BASE_URL not set")
}

function convertToWebSocket(url: string) {
    // Check if the URL starts with http:// or https://
    if (url.startsWith("http://")) {
        return url.replace("http://", "ws://");
    } else if (url.startsWith("https://")) {
        return url.replace("https://", "wss://");
    } else {
        throw new Error("Invalid URL, should start with http:// or https://")
    }
}


export function initWS() {
    if (window.ws) return window.ws

    const baseURL = (import.meta as any).env.VITE_API_BASE_URL
    const token = localStorage.getItem("session_token");
    if (!token) {
        return
    }

    const ws = new WebSocket(`${convertToWebSocket(baseURL)}/realtime/ws?token=${token}`)

    ws.onopen = () => console.log("WS connected")
    ws.onclose = () => {
        console.log("WS closed")
        // TODO: This only checks once, make it try multiple times
        setTimeout(() => { console.log('Retry'); initWS() }, RECONNECT_INTERVAL)
        window.ws = undefined
    }

    ws.onmessage = (e) => {
        const msg = JSON.parse(e.data)
        window.dispatchEvent(new CustomEvent("ws-message", { detail: msg }))
    }

    window.ws = ws
    return ws
}

export function closeWS() {
    window.ws?.close()
    window.ws = undefined
}