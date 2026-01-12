import { GroupsList } from "./GroupsList";
import { SidebarHeader } from "./SidebarHeader";

export function Sidebar() {
    return <div className={`flex-col flex col-span-10 md:col-span-3 h-screen`}>
        <SidebarHeader />
        <GroupsList />
    </div>
}