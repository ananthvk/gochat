export function formatChatTime(timestamp: string) {
    const date = new Date(timestamp);
    const now = new Date();

    const diffMs = now.getTime() - date.getTime();
    const diffSec = Math.floor(diffMs / 1000);
    const diffMin = Math.floor(diffSec / 60);
    const diffHour = Math.floor(diffMin / 60);

    const isToday = now.getDate() === date.getDate() &&
        now.getMonth() === date.getMonth() &&
        now.getFullYear() === date.getFullYear();

    const yesterday = new Date(now);
    yesterday.setDate(now.getDate() - 1);
    const isYesterday = yesterday.getDate() === date.getDate() &&
        yesterday.getMonth() === date.getMonth() &&
        yesterday.getFullYear() === date.getFullYear();

    if (diffSec < 10) return "Just now";
    if (diffMin < 1) return `${diffSec} sec ago`;
    if (diffMin < 60) return `${diffMin} min ago`;
    if (diffHour < 24 && isToday)
        return date.toLocaleTimeString("en-US", {
            hour: "2-digit",
            minute: "2-digit",
        });
    if (isYesterday)
        return `Yesterday ${date.toLocaleTimeString("en-US", {
            hour: "2-digit",
            minute: "2-digit",
        })}`;
    // older than yesterday
    return date.toLocaleDateString("en-US") + " " + date.toLocaleTimeString("en-US", { hour: "2-digit", minute: "2-digit" });
}
