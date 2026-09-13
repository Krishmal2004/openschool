import { useState } from "react";
import { Link } from "react-router";
import { Locked, UserMultiple, Edit, TrashCan } from "@carbon/icons-react";
import { Tag, SkeletonText, Button } from "@carbon/react";
import { useGuardianStudents, useGuardianNotifications, useDeleteGuardian } from "../../../../queries/useGuardians";
import type { Guardian } from "../../../../services/guardian";
import ConfirmDeleteModal from "../../../../components/common/ConfirmDeleteModal";
import { relationshipLabel } from "../constants";
import EditGuardianModal from "./EditGuardianModal";
import MutationErrorNotification from "../../../../components/common/MutationErrorNotification";
import SectionHeader from "../../../../components/common/SectionHeader";

export default function GuardianDetail({ guardian, onDeleted }: { guardian: Guardian; onDeleted: () => void }) {
  const { data: students, isLoading: studentsLoading } = useGuardianStudents(guardian.id);
  const { data: notifications, isLoading: notificationsLoading } = useGuardianNotifications(guardian.id);
  const deleteGuardian = useDeleteGuardian();

  const [editing, setEditing] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState(false);

  const hasStudents = (students?.length ?? 0) > 0;
  // Blocks delete while the linked-students query is still loading, so a fast double-click can't slip past the "still linked" check.
  const deleteBlocked = studentsLoading || hasStudents;

  return (
    <div className="os-section" style={{ marginTop: 0 }}>
      <SectionHeader
        title={guardian.full_name}
        meta={
          <div style={{ display: "flex", gap: "0.5rem", alignItems: "center" }}>
            <Tag size="sm" type="gray">
              {relationshipLabel(guardian.relationship)}
            </Tag>
            {guardian.user_id ? (
              <Tag size="sm" type="green">
                <Locked size={12} style={{ marginRight: "4px" }} />
                Portal access
              </Tag>
            ) : (
              <Tag size="sm" type="cool-gray">
                No portal login
              </Tag>
            )}
            <Button renderIcon={Edit} kind="ghost" size="sm" onClick={() => setEditing(true)}>
              Edit
            </Button>
            <Button
              renderIcon={TrashCan}
              kind="danger--ghost"
              size="sm"
              onClick={() => {
                deleteGuardian.reset();
                setConfirmDelete(true);
              }}
            >
              Delete
            </Button>
          </div>
        }
      />
      <div className="os-section__body">
        <MutationErrorNotification
          isError={deleteGuardian.isError}
          error={deleteGuardian.error}
          title="Could not delete guardian"
          fallback="This guardian may still be linked to a student."
          onClose={() => deleteGuardian.reset()}
          style={{ marginBottom: "1rem" }}
        />

        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "1rem", marginBottom: "1.5rem" }}>
          <div>
            <p style={{ margin: "0 0 0.15rem", fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>Phone</p>
            <p style={{ margin: 0, fontSize: "0.875rem" }}>{guardian.phone}</p>
          </div>
          <div>
            <p style={{ margin: "0 0 0.15rem", fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>Email</p>
            <p style={{ margin: 0, fontSize: "0.875rem" }}>{guardian.email || "—"}</p>
          </div>
        </div>

        <h3 style={{ fontSize: "0.8125rem", fontWeight: 600, margin: "0 0 0.75rem" }}>Linked students</h3>
        {studentsLoading ? (
          <SkeletonText width="60%" />
        ) : students && students.length > 0 ? (
          <div style={{ display: "flex", gap: "0.5rem", flexWrap: "wrap", marginBottom: "1.5rem" }}>
            {students.map((s) => (
              <Link
                key={s.id}
                to={`/students/${s.id}`}
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "0.375rem",
                  padding: "0.375rem 0.75rem",
                  border: "1px solid var(--os-border-subtle)",
                  fontSize: "0.8125rem",
                  textDecoration: "none",
                  color: "var(--os-text-primary)",
                }}
              >
                <UserMultiple size={14} style={{ fill: "var(--os-accent)" }} />
                {s.full_name}
              </Link>
            ))}
          </div>
        ) : (
          <p style={{ fontSize: "0.8125rem", color: "var(--os-text-tertiary)", marginBottom: "1.5rem" }}>No students linked.</p>
        )}

        <h3 style={{ fontSize: "0.8125rem", fontWeight: 600, margin: "0 0 0.75rem" }}>
          Notification history{notifications && notifications.length > 20 ? " (most recent 20)" : ""}
        </h3>
        {!guardian.user_id ? (
          <p style={{ fontSize: "0.8125rem", color: "var(--os-text-tertiary)" }}>
            This guardian has no portal login, so they haven't received any in-app notifications.
          </p>
        ) : notificationsLoading ? (
          <SkeletonText width="80%" />
        ) : notifications && notifications.length > 0 ? (
          <div>
            {notifications.slice(0, 20).map((n) => (
              <div key={n.recipient_id} style={{ padding: "0.625rem 0", borderBottom: "1px solid var(--os-border-subtle)" }}>
                <div style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
                  <span style={{ fontWeight: 600, fontSize: "0.8125rem" }}>{n.title}</span>
                  <Tag size="sm" type="cool-gray">
                    {n.category}
                  </Tag>
                  <div style={{ flex: 1 }} />
                  <span style={{ fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>
                    {new Date(n.sent_at).toLocaleString()}
                  </span>
                </div>
                <p style={{ margin: "0.25rem 0 0", fontSize: "0.8125rem", color: "var(--os-text-secondary)" }}>{n.message}</p>
              </div>
            ))}
          </div>
        ) : (
          <p style={{ fontSize: "0.8125rem", color: "var(--os-text-tertiary)" }}>No notifications sent yet.</p>
        )}
      </div>

      {editing && <EditGuardianModal guardian={guardian} onClose={() => setEditing(false)} />}

      <ConfirmDeleteModal
        open={confirmDelete}
        title="Delete guardian"
        description={
          studentsLoading ? (
            "Checking linked students…"
          ) : hasStudents ? (
            <>
              <strong>{guardian.full_name}</strong> is still linked to {students?.length} student
              {students && students.length !== 1 ? "s" : ""}. Unlink them from this guardian first
              (from each student's profile), then delete.
            </>
          ) : (
            <>
              Delete <strong>{guardian.full_name}</strong>? This cannot be undone.
            </>
          )
        }
        isPending={deleteGuardian.isPending}
        disabled={deleteBlocked}
        onClose={() => setConfirmDelete(false)}
        onConfirm={() => {
          if (deleteBlocked) return;
          deleteGuardian.mutate(guardian.id, {
            onSuccess: () => {
              setConfirmDelete(false);
              onDeleted();
            },
          });
        }}
      />
    </div>
  );
}
