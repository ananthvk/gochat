import { GroupsList } from "./GroupsList";
import { SidebarHeader } from "./SidebarHeader";

export function Sidebar() {
    return <div className="flex flex-col col-span-2 md:col-span-3 h-screen">
        <SidebarHeader />
        <GroupsList />
    </div>
}