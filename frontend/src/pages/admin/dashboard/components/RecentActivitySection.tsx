import { Link } from "react-router";
import { UserFollow } from "@carbon/icons-react";
import { Tag, SkeletonText } from "@carbon/react";
import EmptyState from "../../../../components/common/EmptyState";
import SectionHeader from "../../../../components/common/SectionHeader";
import { ACCENT } from "../constants";

export type RecentActivityItem = {
  key: string;
  text: string;
  sub: string;
  time: string;
  path: string;
  kind: "student" | "teacher";
};

export default function RecentActivitySection({
  items,
  loading,
}: {
  items: RecentActivityItem[];
  loading: boolean;
}) {
  return (
    <div className="os-section">
      <SectionHeader
        title="Recent Activity"
        meta={
          <Link to="/students" style={{ fontSize: "0.75rem", color: ACCENT, textDecoration: "none" }}>
            View all →
          </Link>
        }
      />
      {loading ? (
        <div style={{ padding: "1.25rem 1.5rem" }}>
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} style={{ marginBottom: "0.75rem" }}>
              <SkeletonText width="60%" />
            </div>
          ))}
        </div>
      ) : items.length === 0 ? (
        <EmptyState
          title="Nothing here yet"
          description="Enrol students and add teachers to see recent activity."
          action={
            <Link to="/students/new" style={{ fontSize: "0.8125rem", color: ACCENT, fontWeight: 500 }}>
              Enrol a student →
            </Link>
          }
        />
      ) : (
        <div>
          {items.map((item) => (
            <Link
              key={item.key}
              to={item.path}
              className="os-list-row"
              style={{ alignItems: "flex-start", textDecoration: "none" }}
            >
              <div style={{ marginTop: "2px", flexShrink: 0 }}>
                <UserFollow size={16} style={{ fill: item.kind === "teacher" ? "var(--os-chart-purple)" : ACCENT }} />
              </div>
              <div style={{ flex: 1, minWidth: 0 }}>
                <p style={{ margin: "0 0 0.15rem", fontSize: "0.875rem", fontWeight: 500, color: "var(--os-text-primary)" }}>
                  {item.text}
                </p>
                {item.sub && <p style={{ margin: 0, fontSize: "0.75rem", color: "var(--os-text-secondary)" }}>{item.sub}</p>}
              </div>
              <div style={{ display: "flex", flexDirection: "column", alignItems: "flex-end", gap: "0.25rem", flexShrink: 0 }}>
                <Tag type={item.kind === "teacher" ? "purple" : "blue"} size="sm">
                  {item.kind === "teacher" ? "Teacher" : "Student"}
                </Tag>
                <span style={{ fontSize: "0.6875rem", color: "var(--os-text-tertiary)" }}>
                  {new Date(item.time).toLocaleDateString("en-LK", { month: "short", day: "numeric" })}
                </span>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
