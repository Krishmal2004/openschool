import type { Dispatch, SetStateAction } from "react";
import { TextInput, Select, SelectItem, NumberInput } from "@carbon/react";
import type { useDuplicateLevel } from "@/features/curriculum/queries/useCurriculum";
import type { useGrades } from "@/features/academics/queries/useGrades";
import type { Level } from "@/features/curriculum/api/curriculum";
import FormModal from "@/shared/ui/FormModal";
import type { LevelForm } from "@/features/curriculum/constants";

interface Props {
  toDuplicate: Level | null;
  dupForm: LevelForm;
  setDupForm: Dispatch<SetStateAction<LevelForm>>;
  dupLabelTouched: boolean;
  setDupLabelTouched: Dispatch<SetStateAction<boolean>>;
  grades: ReturnType<typeof useGrades>["data"];
  duplicateLevel: ReturnType<typeof useDuplicateLevel>;
  onClose: () => void;
  onDuplicate: () => void;
}

export default function DuplicateLevelModal({
  toDuplicate,
  dupForm,
  setDupForm,
  dupLabelTouched,
  setDupLabelTouched,
  grades,
  duplicateLevel,
  onClose,
  onDuplicate,
}: Props) {
  return (
    <FormModal
      open={!!toDuplicate}
      title={`Duplicate ${toDuplicate?.label ?? ""}`}
      size="md"
      onClose={onClose}
      onSubmit={onDuplicate}
      isPending={duplicateLevel.isPending}
      submitDisabled={!dupForm.label.trim()}
      submitLabel="Duplicate"
      pendingLabel="Duplicating…"
      isError={duplicateLevel.isError}
      error={duplicateLevel.error}
      errorFallback="Failed to duplicate level"
    >
      <p className="os-text-md os-c-secondary os-mb-4">
        Every selection group and its subjects are copied to the new level. Editing one afterwards does not affect
        the other.
      </p>
      <div className="os-grid os-gap-5">
        <TextInput
          id="dup-label"
          labelText="New label"
          helperText="Must differ from every existing level."
          value={dupForm.label}
          onChange={(e) => setDupForm((f) => ({ ...f, label: e.target.value }))}
          onBlur={() => setDupLabelTouched(true)}
          invalid={dupLabelTouched && !dupForm.label.trim()}
          invalidText="A label is required."
        />
        <div className="os-grid os-grid-cols-2-1 os-gap-4">
          <Select
            id="dup-grade"
            labelText="Grade (optional)"
            value={dupForm.grade_id}
            onChange={(e) => setDupForm((f) => ({ ...f, grade_id: e.target.value }))}
          >
            <SelectItem value="" text="No grade" />
            {grades?.map((g) => (
              <SelectItem key={g.id} value={g.id} text={g.name} />
            ))}
          </Select>
          <NumberInput
            id="dup-sort"
            label="Sort order"
            min={0}
            value={dupForm.sort_order}
            onChange={(_e, { value }) => setDupForm((f) => ({ ...f, sort_order: Number(value) || 0 }))}
          />
        </div>
      </div>
    </FormModal>
  );
}
