import { Tag } from "@carbon/react";
import { useMyStudentProfile } from "@/features/students/queries/useStudentSelf";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import ErrorMessage from "@/shared/ui/ErrorMessage";
import { getInitials } from "@/shared/lib/name";

const kvItemStyle = { border: "1px solid var(--os-border-subtle)", background: "var(--os-layer)", padding: "1rem" };

export default function StudentDashboard() {
  const { data: profile, isLoading, isError, refetch } = useMyStudentProfile();

  if (isLoading) return <LoadingSpinner />;
  if (isError || !profile) {
    return (
      <div className="os-p-8">
        <ErrorMessage message="Failed to load your profile" onRetry={refetch} />
      </div>
    );
  }

  return (
    <div className="os-bg-layer-hover os-min-h-content">
      <div className="os-profile__banner">
        <div className="os-profile__avatar">
          {getInitials(profile.full_name)}
        </div>
        <div className="os-flex-1">
          <p className="os-profile__name">{profile.full_name}</p>
          <p className="os-profile__meta">
            {profile.index_number}
            {profile.class_name ? ` · ${profile.class_name}` : ""}
            {profile.grade_name ? ` · ${profile.grade_name}` : ""}
          </p>
        </div>
        {profile.house_name && (
          <div className="os-profile__actions">
            <Tag type="blue" size="sm">
              {profile.house_name} House
            </Tag>
          </div>
        )}
      </div>

      <div className="os-py-6 os-px-8">
        <div className="os-section">
          <div className="os-section__header">
            <h2 className="os-section__title">Student Details</h2>
          </div>
          <div className="os-section__body os-grid os-grid-auto-200 os-gap-4">
            <div className="os-kv-item" style={kvItemStyle}>
              <p className="os-kv-item__label">Full Name</p>
              <p className="os-kv-item__value">{profile.full_name}</p>
            </div>
            <div className="os-kv-item" style={kvItemStyle}>
              <p className="os-kv-item__label">Index Number</p>
              <p className="os-kv-item__value">{profile.index_number}</p>
            </div>
            <div className="os-kv-item" style={kvItemStyle}>
              <p className="os-kv-item__label">Class</p>
              <p className="os-kv-item__value">{profile.class_name || "-"}</p>
            </div>
            <div className="os-kv-item" style={kvItemStyle}>
              <p className="os-kv-item__label">Grade</p>
              <p className="os-kv-item__value">{profile.grade_name || "-"}</p>
            </div>
            <div className="os-kv-item" style={kvItemStyle}>
              <p className="os-kv-item__label">House</p>
              <p className="os-kv-item__value">{profile.house_name || "-"}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
