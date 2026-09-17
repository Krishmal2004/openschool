import { Button, TextInput, MultiSelect, ComposedModal, ModalHeader, ModalBody, ModalFooter } from "@carbon/react";
import type { Grade } from "@/features/academics/api/grade";
import type { GradeSectionForm } from "@/features/timetable/constants";
import MutationErrorNotification from "@/shared/ui/MutationErrorNotification";

export default function SectionFormModal({
  open,
  isEdit,
  form,
  onFormChange,
  grades,
  selectedGradeIds,
  onSelectedGradeIdsChange,
  onClose,
  onSubmit,
  isPending,
  isError,
  error,
}: {
  open: boolean;
  isEdit: boolean;
  form: GradeSectionForm;
  onFormChange: (form: GradeSectionForm) => void;
  grades: Grade[] | undefined;
  selectedGradeIds: string[];
  onSelectedGradeIdsChange: (ids: string[]) => void;
  onClose: () => void;
  onSubmit: () => void;
  isPending: boolean;
  isError: boolean;
  error: unknown;
}) {
  return (
    <ComposedModal open={open} size="md" onClose={onClose}>
      <ModalHeader title={isEdit ? "Edit grade section" : "New grade section"} />
      <ModalBody>
        <MutationErrorNotification
          isError={isError}
          error={error} className="os-mb-4"
        />
        <div className="os-grid os-gap-4">
          <TextInput
            id="section-name"
            labelText="Section name"
            placeholder="e.g. Junior Secondary"
            value={form.name}
            onChange={(e) => onFormChange({ ...form, name: e.target.value })}
          />
          <div className="os-grid os-grid-cols-2 os-gap-4">
            <TextInput
              id="section-interval-start"
              labelText="Interval start"
              type="time"
              value={form.interval_start_time}
              onChange={(e) => onFormChange({ ...form, interval_start_time: e.target.value })}
            />
            <TextInput
              id="section-interval-end"
              labelText="Interval end"
              type="time"
              value={form.interval_end_time}
              onChange={(e) => onFormChange({ ...form, interval_end_time: e.target.value })}
            />
          </div>
          <MultiSelect
            id="section-grades"
            titleText="Grades in this section"
            label="Select grades…"
            items={grades ?? []}
            itemToString={(g) => g?.name ?? ""}
            selectedItems={(grades ?? []).filter((g) => selectedGradeIds.includes(g.id))}
            onChange={({ selectedItems }) => onSelectedGradeIdsChange((selectedItems ?? []).map((g) => g.id))}
          />
        </div>
      </ModalBody>
      <ModalFooter>
        <Button kind="secondary" onClick={onClose}>
          Cancel
        </Button>
        <Button kind="primary" onClick={onSubmit} disabled={!form.name.trim() || isPending}>
          {isPending ? "Saving…" : "Save"}
        </Button>
      </ModalFooter>
    </ComposedModal>
  );
}
