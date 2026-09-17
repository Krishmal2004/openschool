import { useState } from "react";
import { Checkbox, MultiSelect } from "@carbon/react";
import { useTeachers } from "@/features/teachers/queries/useTeachers";
import { useGrades } from "@/features/academics/queries/useGrades";
import { useAssignPrincipal, useAssignVicePrincipal } from "@/features/positions/queries/usePositions";
import FormModal from "@/shared/ui/FormModal";
import EntityCombobox from "@/shared/ui/EntityCombobox";

export function AssignPrincipalModal({ currentTeacherId, onClose }: { currentTeacherId: string; onClose: () => void }) {
  const { data: teachers } = useTeachers();
  const assign = useAssignPrincipal();
  const [choice, setChoice] = useState(currentTeacherId);
  return (
    <FormModal
      open
      title={currentTeacherId ? "Change Principal" : "Assign Principal"}
      onClose={onClose}
      onSubmit={() => choice && assign.mutate({ teacher_id: choice }, { onSuccess: onClose })}
      isPending={assign.isPending}
      submitLabel="Assign"
      submitDisabled={!choice}
      isError={assign.isError}
      error={assign.error}
      errorFallback="Failed to assign Principal"
    >
      <EntityCombobox id="principal-teacher" labelText="Teacher" items={teachers ?? []} selectedId={choice} onSelect={setChoice} getId={(t) => t.id} itemToString={(t) => t.full_name} placeholder="Search teachers by name…" />
    </FormModal>
  );
}

export function AddVicePrincipalModal({ onClose }: { onClose: () => void }) {
  const { data: teachers } = useTeachers();
  const { data: grades } = useGrades();
  const assign = useAssignVicePrincipal();
  const [teacherId, setTeacherId] = useState("");
  const [wholeSchool, setWholeSchool] = useState(false);
  const [gradeIds, setGradeIds] = useState<string[]>([]);
  return (
    <FormModal
      open
      title="Add Vice Principal"
      onClose={onClose}
      onSubmit={() => teacherId && assign.mutate({ teacher_id: teacherId, notify_whole_school: wholeSchool, grade_ids: wholeSchool ? [] : gradeIds }, { onSuccess: onClose })}
      isPending={assign.isPending}
      submitLabel="Add"
      submitDisabled={!teacherId}
      isError={assign.isError}
      error={assign.error}
      errorFallback="Failed to add Vice Principal"
    >
      <div className="os-grid os-gap-4">
        <EntityCombobox id="vp-teacher" labelText="Teacher" items={teachers ?? []} selectedId={teacherId} onSelect={setTeacherId} getId={(t) => t.id} itemToString={(t) => t.full_name} placeholder="Search teachers by name…" />
        <Checkbox id="vp-whole-school" labelText="Can notify the whole school" checked={wholeSchool} onChange={(_e, { checked }) => setWholeSchool(checked)} />
        {!wholeSchool && (
          <MultiSelect
            id="vp-grades"
            titleText="Grades this Vice Principal can notify"
            label="Select grades…"
            items={grades ?? []}
            itemToString={(g) => g?.name ?? ""}
            selectedItems={(grades ?? []).filter((g) => gradeIds.includes(g.id))}
            onChange={({ selectedItems }) => setGradeIds((selectedItems ?? []).map((g) => g.id))}
          />
        )}
      </div>
    </FormModal>
  );
}
