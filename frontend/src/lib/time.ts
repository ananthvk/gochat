export function formatMessageTime(timestamp: string): string {
    const messageDate = new Date(timestamp)
    const now = new Date()

    const diffMs = now.getTime() - messageDate.getTime()
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

    const isToday =
        messageDate.getDate() === now.getDate() &&
        messageDate.getMonth() === now.getMonth() &&
        messageDate.getFullYear() === now.getFullYear()

    const isYesterday =
        new Date(now.getFullYear(), now.getMonth(), now.getDate() - 1).getTime() ===
        new Date(messageDate.getFullYear(), messageDate.getMonth(), messageDate.getDate()).getTime()

    if (isToday) {
        return messageDate.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    } else if (isYesterday) {
        return 'Yesterday'
    } else if (diffDays < 7) {
        return messageDate.toLocaleDateString([], { weekday: 'long' })
    } else {
        return messageDate.toLocaleDateString()
    }
}
