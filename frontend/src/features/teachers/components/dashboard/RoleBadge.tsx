import { Tag } from "@carbon/react";
import type { PositionRankLabel } from "@/features/positions/api/position";

// One colour per rank so the badge reads as a hierarchy; gray for no position.
const RANK_TAG_TYPE: Record<PositionRankLabel, "purple" | "blue" | "teal" | "cyan" | "green" | "gray"> = {
  Principal: "purple",
  "Vice Principal": "blue",
  "Section Head": "teal",
  "Class Teacher": "cyan",
  "Subject Teacher": "green",
  Teacher: "gray",
};

export default function RoleBadge({ rankLabel }: { rankLabel: PositionRankLabel }) {
  return (
    <Tag type={RANK_TAG_TYPE[rankLabel]} size="sm">
      {rankLabel}
    </Tag>
  );
}
