import { Link } from "react-router";
import { ChevronRight } from "@carbon/icons-react";
import { useMyChildren } from "@/features/parent/queries/useParent";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import ErrorMessage from "@/shared/ui/ErrorMessage";
import EmptyState from "@/shared/ui/EmptyState";
import { getInitials } from "@/shared/lib/name";

export default function ParentDashboard() {
  const { data: children, isLoading, isError, refetch } = useMyChildren();

  if (isLoading) return <LoadingSpinner />;
  if (isError) {
    return (
      <div className="os-p-8">
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
              to={`/p/children/${c.id}`} className="os-flex os-items-center os-gap-3h os-p-5 os-bg-layer os-border os-no-underline"
            >
              <div
                className="os-profile__avatar os-w-2t os-h-2t os-text-md"
              >
                {getInitials(c.full_name)}
              </div>
              <div className="os-flex-1 os-min-w-0">
                <p className="os-mt-0 os-mx-0 os-mb-h os-fw-600 os-text-md os-c-primary">
                  {c.full_name}
                </p>
                <p className="os-m-0 os-text-sm os-c-secondary">
                  {c.index_number}
                  {c.class_name ? ` · ${c.class_name}` : ""}
                  {c.grade_name ? ` · ${c.grade_name}` : ""}
                </p>
              </div>
              <ChevronRight size={18} className="os-fill-tertiary os-shrink-0" />
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
