import { Link } from "react-router";
import { ChevronRight } from "@carbon/icons-react";
import { useMyChildren } from "../../queries/useParent";
import LoadingSpinner from "../../components/common/LoadingSpinner";
import ErrorMessage from "../../components/common/ErrorMessage";
import EmptyState from "../../components/common/EmptyState";
import { getInitials } from "../../lib/name";

export default function ParentDashboard() {
  const { data: children, isLoading, isError, refetch } = useMyChildren();

  if (isLoading) return <LoadingSpinner />;
  if (isError) {
    return (
      <div style={{ padding: "2rem" }}>
        <ErrorMessage message="Failed to load your children" onRetry={refetch} />
      </div>
    );
  }

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">My Children</h1>
          <p className="os-page__subtitle">
            Select a child to see their attendance, marks, and class details.
          </p>
        </div>
      </div>

      {children && children.length === 0 && (
        <EmptyState
          title="No children linked yet"
          description="Your school administrator hasn't linked any students to your account. Contact them if this looks wrong."
        />
      )}

      {children && children.length > 0 && (
        <div className="os-quick-grid">
          {children.map((c) => (
            <Link
              key={c.id}
              to={`/p/children/${c.id}`}
              style={{
                display: "flex",
                alignItems: "center",
                gap: "0.875rem",
                padding: "1.25rem",
                background: "var(--os-layer)",
                border: "1px solid var(--os-border-subtle)",
                textDecoration: "none",
                transition: "border-color 0.15s ease",
              }}
            >
              <div
                className="os-profile__avatar"
                style={{ width: "2.75rem", height: "2.75rem", fontSize: "0.9rem" }}
              >
                {getInitials(c.full_name)}
              </div>
              <div style={{ flex: 1, minWidth: 0 }}>
                <p style={{ margin: "0 0 0.2rem", fontWeight: 600, fontSize: "0.9375rem", color: "var(--os-text-primary)" }}>
                  {c.full_name}
                </p>
                <p style={{ margin: 0, fontSize: "0.8125rem", color: "var(--os-text-secondary)" }}>
                  {c.index_number}
                  {c.class_name ? ` · ${c.class_name}` : ""}
                  {c.grade_name ? ` · ${c.grade_name}` : ""}
                </p>
              </div>
              <ChevronRight size={18} style={{ fill: "var(--os-text-tertiary)", flexShrink: 0 }} />
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
