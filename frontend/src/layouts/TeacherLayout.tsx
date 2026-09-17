import PortalShell from "@/layouts/PortalShell";
import { teacherNav } from "@/layouts/nav/teacher";
import { useMyPosition } from "@/features/positions/queries/usePositions";
import { POSITION_RANK } from "@/shared/lib/constants/people";

export default function TeacherLayout() {
  const { data: position } = useMyPosition();
  const navGroups = teacherNav({
    isLeadership: !!position && position.rank <= POSITION_RANK.vicePrincipal,
    isSectionHead: position?.rank === POSITION_RANK.sectionHead,
  });
  return <PortalShell navGroups={navGroups} />;
}
