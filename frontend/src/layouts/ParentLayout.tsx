import PortalShell from "@/layouts/PortalShell";
import { PARENT_NAV } from "@/layouts/nav/parent";

export default function ParentLayout() {
  return <PortalShell navGroups={PARENT_NAV} />;
}
