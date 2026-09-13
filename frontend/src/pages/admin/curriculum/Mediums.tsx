import { useState } from "react";
import { Add, Language } from "@carbon/icons-react";
import {
  Button,
  TextInput,
  ComposedModal,
  ModalHeader,
  ModalBody,
  ModalFooter,
} from "@carbon/react";
import {
  useMediums,
  useCreateMedium,
  useUpdateMedium,
  useDeleteMedium,
} from "../../../queries/useCurriculum";
import type { Medium } from "../../../services/curriculum";
import ErrorMessage from "../../../components/common/ErrorMessage";
import EmptyState from "../../../components/common/EmptyState";
import ConfirmDeleteModal from "../../../components/common/ConfirmDeleteModal";
import MutationErrorNotification from "../../../components/common/MutationErrorNotification";
import SectionHeader from "../../../components/common/SectionHeader";
import ListRowSkeleton from "../../../components/common/ListRowSkeleton";

export default function Mediums() {
  const { data: mediums, isLoading, isError, refetch } = useMediums();
  const createMedium = useCreateMedium();
  const updateMedium = useUpdateMedium();
  const deleteMedium = useDeleteMedium();

  const [modal, setModal] = useState<"create" | "edit" | null>(null);
  const [editing, setEditing] = useState<Medium | null>(null);
  const [name, setName] = useState("");
  const [nameTouched, setNameTouched] = useState(false);
  const [toDelete, setToDelete] = useState<Medium | null>(null);

  const openCreate = () => {
    createMedium.reset();
    setName("");
    setNameTouched(false);
    setEditing(null);
    setModal("create");
  };

  const openEdit = (m: Medium) => {
    updateMedium.reset();
    setName(m.name);
    setNameTouched(false);
    setEditing(m);
    setModal("edit");
  };

  const handleSave = () => {
    setNameTouched(true);
    if (!name.trim()) return;
    const data = { name: name.trim() };
    if (modal === "create") {
      createMedium.mutate(data, { onSuccess: () => setModal(null) });
    } else if (editing) {
      updateMedium.mutate(
        { id: editing.id, data },
        { onSuccess: () => setModal(null) },
      );
    }
  };

  const handleDelete = () => {
    if (!toDelete) return;
    deleteMedium.mutate(toDelete.id, { onSettled: () => setToDelete(null) });
  };

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">Mediums</h1>
          <p className="os-page__subtitle">
            Languages of instruction. Used to restrict a subject within a
            selection group to a single medium.
          </p>
        </div>
        <Button renderIcon={Add} kind="primary" size="md" onClick={openCreate}>
          New Medium
        </Button>
      </div>

      <div className="os-section">
        <SectionHeader
          title="Mediums"
          meta={mediums && <span className="os-section__meta">{mediums.length} total</span>}
        />

        {isLoading && (
          <div>
            {Array.from({ length: 3 }).map((_, i) => (
              <ListRowSkeleton key={i} leadingWidth="1.5rem" titleWidth="30%" subtitleWidth={null} trailingWidth={null} />
            ))}
          </div>
        )}
        {isError && (
          <ErrorMessage message="Could not load mediums." onRetry={refetch} />
        )}

        <MutationErrorNotification
          isError={deleteMedium.isError}
          error={deleteMedium.error}
          title="Could not delete medium"
          fallback="The medium may be in use by a group subject or enrollment."
          onClose={() => deleteMedium.reset()}
          style={{ margin: "0 1.5rem 1rem" }}
        />

        {!isLoading && !isError && mediums?.length === 0 && (
          <EmptyState
            title="No mediums"
            description="Add the languages your school teaches in, for example Sinhala, Tamil or English."
            action={
              <Button renderIcon={Add} kind="primary" onClick={openCreate}>
                New Medium
              </Button>
            }
          />
        )}

        {!isLoading && mediums && mediums.length > 0 && (
          <div>
            {mediums.map((m) => (
              <div key={m.id} className="os-list-row">
                <Language size={20} style={{ fill: "var(--os-accent)", flexShrink: 0 }} />
                <span style={{ flex: 1, fontWeight: 600, fontSize: "0.9rem", color: "var(--os-text-primary)" }}>
                  {m.name}
                </span>
                <Button kind="ghost" size="sm" onClick={() => openEdit(m)}>
                  Edit
                </Button>
                <Button
                  kind="danger--ghost"
                  size="sm"
                  onClick={() => setToDelete(m)}
                >
                  Delete
                </Button>
              </div>
            ))}
          </div>
        )}
      </div>

      <ComposedModal open={!!modal} size="sm" onClose={() => setModal(null)}>
        <ModalHeader title={modal === "create" ? "New medium" : "Edit medium"} />
        <ModalBody>
          <MutationErrorNotification
            isError={createMedium.isError || updateMedium.isError}
            error={createMedium.error ?? updateMedium.error}
            fallback="Failed to save medium"
            style={{ marginBottom: "1rem" }}
          />
          <TextInput
            id="medium-name"
            labelText="Name"
            placeholder="e.g. English"
            value={name}
            onChange={(e) => setName(e.target.value)}
            onBlur={() => setNameTouched(true)}
            invalid={nameTouched && !name.trim()}
            invalidText="Medium name is required."
          />
        </ModalBody>
        <ModalFooter>
          <Button kind="secondary" onClick={() => setModal(null)}>
            Cancel
          </Button>
          <Button
            kind="primary"
            onClick={handleSave}
            disabled={
              !name.trim() || createMedium.isPending || updateMedium.isPending
            }
          >
            {createMedium.isPending || updateMedium.isPending
              ? "Saving…"
              : "Save"}
          </Button>
        </ModalFooter>
      </ComposedModal>

      <ConfirmDeleteModal
        open={!!toDelete}
        title="Delete medium"
        description={
          <>
            Delete <strong>{toDelete?.name}</strong>? This cannot be undone.
          </>
        }
        isPending={deleteMedium.isPending}
        onClose={() => setToDelete(null)}
        onConfirm={handleDelete}
      />
    </div>
  );
}
