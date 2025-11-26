import { GroupsList } from "./GroupsList";
import { SidebarHeader } from "./SidebarHeader";

export function Sidebar() {
    return <div style={{ backgroundColor: 'green' }} className="flex flex-col col-span-2 md:col-span-3">
        <SidebarHeader />
        <GroupsList />
    </div>
}