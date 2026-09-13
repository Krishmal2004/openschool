import { useMemo, useState } from "react";
import { Link } from "react-router";
import { Search, Add, Edit, TrashCan } from "@carbon/icons-react";
import {
  Button,
  IconButton,
  Tag,
  TextInput,
  NumberInput,
  ComposedModal,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Pagination,
} from "@carbon/react";
import {
  useSubjects,
  useUpdateSubject,
  useDeleteSubject,
} from "../../../../queries/useSubjects";
import type { Subject } from "../../../../services/subject";
import { usePagination } from "../../../../hooks/usePagination";
import TableSkeleton from "../../../../components/common/TableSkeleton";
import ErrorMessage from "../../../../components/common/ErrorMessage";
import EmptyState from "../../../../components/common/EmptyState";
import ConfirmDeleteModal from "../../../../components/common/ConfirmDeleteModal";
import MutationErrorNotification from "../../../../components/common/MutationErrorNotification";

const SUBJECT_TABLE_HEADERS = ["Code", "Subject", "Type", "Actions"];

export default function SubjectsPanel() {
  const { data: subjects, isLoading, isError, refetch } = useSubjects();
  const updateSubject = useUpdateSubject();
  const deleteSubject = useDeleteSubject();

  const [query, setQuery] = useState("");
  const [editing, setEditing] = useState<Subject | null>(null);
  const [form, setForm] = useState({ name: "", code: "", type: "", max_marks: 100 });
  const [touched, setTouched] = useState<{ name?: boolean; code?: boolean }>({});
  const [toDelete, setToDelete] = useState<Subject | null>(null);

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    return (subjects ?? []).filter(
      (s) =>
        s.name.toLowerCase().includes(q) ||
        s.code.toLowerCase().includes(q) ||
        (s.type ?? "").toLowerCase().includes(q),
    );
  }, [subjects, query]);

  const { page, pageSize, pageItems, totalItems, onChange } = usePagination(filtered, 10);

  const openEdit = (s: Subject) => {
    updateSubject.reset();
    setForm({ name: s.name, code: s.code, type: s.type ?? "", max_marks: s.max_marks });
    setTouched({});
    setEditing(s);
  };

  const handleUpdate = () => {
    setTouched({ name: true, code: true });
    if (!editing || !isValid) return;
    updateSubject.mutate(
      {
        id: editing.id,
        data: {
          name: form.name.trim(),
          code: form.code.trim(),
          type: form.type.trim(),
          max_marks: form.max_marks,
        },
      },
      { onSuccess: () => setEditing(null) },
    );
  };

  const handleDelete = () => {
    if (!toDelete) return;
    deleteSubject.mutate(toDelete.id, { onSettled: () => setToDelete(null) });
  };

  const isValid = form.name.trim() && form.code.trim();
  const searching = query.trim().length > 0;

  return (
    <div>
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          gap: "1rem",
          flexWrap: "wrap",
          margin: "1rem 0",
        }}
      >
        <p style={{ margin: 0, fontSize: "0.8125rem", color: "var(--os-text-secondary)" }}>
          The school&apos;s subject catalogue. Offer a subject to students by
          adding it to a selection group under the Curriculum tab.
        </p>
        <Button renderIcon={Add} kind="primary" size="md" as={Link} to="/subjects/new">
          Add Subject
        </Button>
      </div>

      <div className="os-section">
        <div className="os-toolbar">
          <div className="os-search">
            <Search size={16} className="os-search__icon" />
            <input
              className="os-search__input"
              placeholder="Search by name, code or type…"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
          </div>
        </div>

        <MutationErrorNotification
          isError={deleteSubject.isError}
          error={deleteSubject.error}
          title="Could not delete subject"
          fallback="The subject may be in use by a class or curriculum group."
          onClose={() => deleteSubject.reset()}
          style={{ margin: "0 1.5rem 1rem" }}
        />

        {isLoading ? (
          <TableSkeleton headers={SUBJECT_TABLE_HEADERS} />
        ) : isError ? (
          <ErrorMessage message="Failed to load subjects" onRetry={refetch} />
        ) : filtered.length === 0 ? (
          <EmptyState
            title={searching ? "No subjects found" : "No subjects yet"}
            description={
              searching
                ? `No subject matches "${query}".`
                : "Add the subjects your school teaches to get started."
            }
            action={
              searching ? undefined : (
                <Button renderIcon={Add} kind="primary" as={Link} to="/subjects/new">
                  Add Subject
                </Button>
              )
            }
          />
        ) : (
          <>
            <div className="os-section__header">
              <h2 className="os-section__title">All Subjects</h2>
              <span style={{ fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>
                {filtered.length} records
              </span>
            </div>

            <table className="os-table">
              <thead>
                <tr>
                  <th>Code</th>
                  <th>Subject</th>
                  <th>Type</th>
                  <th>Max Marks</th>
                  <th style={{ width: "6rem", textAlign: "right" }}>Actions</th>
                </tr>
              </thead>
              <tbody>
                {pageItems.map((s) => (
                  <tr key={s.id}>
                    <td>
                      <span className="os-table__mono">{s.code}</span>
                    </td>
                    <td>{s.name}</td>
                    <td>
                      {s.type ? (
                        <Tag type="blue" size="sm">
                          {s.type}
                        </Tag>
                      ) : (
                        <span className="os-table__muted">-</span>
                      )}
                    </td>
                    <td>{s.max_marks}</td>
                    <td>
                      <div
                        style={{
                          display: "flex",
                          gap: "0.25rem",
                          justifyContent: "flex-end",
                        }}
                      >
                        <IconButton
                          label="Edit"
                          kind="ghost"
                          size="sm"
                          onClick={() => openEdit(s)}
                        >
                          <Edit />
                        </IconButton>
                        <IconButton
                          label="Delete"
                          kind="ghost"
                          size="sm"
                          onClick={() => setToDelete(s)}
                        >
                          <TrashCan />
                        </IconButton>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            <Pagination
              totalItems={totalItems}
              page={page}
              pageSize={pageSize}
              pageSizes={[10, 20, 50]}
              onChange={onChange}
            />
          </>
        )}
      </div>

      <ComposedModal open={!!editing} size="sm" onClose={() => setEditing(null)}>
        <ModalHeader title="Edit subject" />
        <ModalBody>
          <MutationErrorNotification
            isError={updateSubject.isError}
            error={updateSubject.error}
            fallback="Failed to update subject"
            style={{ marginBottom: "1rem" }}
          />
          <div style={{ display: "grid", gap: "1rem" }}>
            <TextInput
              id="edit-subject-name"
              labelText="Subject name"
              value={form.name}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
              onBlur={() => setTouched((t) => ({ ...t, name: true }))}
              invalid={!!touched.name && !form.name.trim()}
              invalidText="Subject name is required."
            />
            <TextInput
              id="edit-subject-code"
              labelText="Subject code"
              value={form.code}
              onChange={(e) => setForm((f) => ({ ...f, code: e.target.value }))}
              onBlur={() => setTouched((t) => ({ ...t, code: true }))}
              invalid={!!touched.code && !form.code.trim()}
              invalidText="Subject code is required."
            />
            <TextInput
              id="edit-subject-type"
              labelText="Type (optional)"
              helperText="A descriptive label only, e.g. core, language, aesthetic"
              value={form.type}
              onChange={(e) => setForm((f) => ({ ...f, type: e.target.value }))}
            />
            <NumberInput
              id="edit-subject-max-marks"
              label="Max Marks"
              min={1}
              max={1000}
              value={form.max_marks}
              onChange={(_e, { value }) => setForm((f) => ({ ...f, max_marks: Number(value ?? 100) }))}
              helperText="The maximum marks a student can get for this subject."
            />
          </div>
        </ModalBody>
        <ModalFooter>
          <Button kind="secondary" onClick={() => setEditing(null)}>
            Cancel
          </Button>
          <Button
            kind="primary"
            onClick={handleUpdate}
            disabled={!isValid || updateSubject.isPending}
          >
            {updateSubject.isPending ? "Saving…" : "Save"}
          </Button>
        </ModalFooter>
      </ComposedModal>

      <ConfirmDeleteModal
        open={!!toDelete}
        title="Delete subject"
        description={
          <>
            Delete <strong>{toDelete?.name}</strong> ({toDelete?.code})? This
            cannot be undone, and is blocked if the subject is used by a class or
            curriculum group.
          </>
        }
        isPending={deleteSubject.isPending}
        onClose={() => setToDelete(null)}
        onConfirm={handleDelete}
      />
    </div>
  );
}
