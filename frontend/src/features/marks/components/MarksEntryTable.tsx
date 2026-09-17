import { Checkbox, NumberInput, Tag } from "@carbon/react";
import type { Student } from "@/features/students/api/student";
import type { useMarksDraft } from "@/features/marks/hooks/useMarksDraft";

interface Props {
  students: Student[];
  maxMarks: number;
  draft: ReturnType<typeof useMarksDraft>;
}

// Editable roster; stays a plain table because every row holds form controls.
export default function MarksEntryTable({ students, maxMarks, draft }: Props) {
  return (
    <table className="os-table">
      <thead>
        <tr>
          <th className="os-w-3">#</th>
          <th>Student Name</th>
          <th>Index Number</th>
          <th className="os-w-8">Marks</th>
          <th className="os-w-5">Absent</th>
          <th className="os-w-7">Status</th>
        </tr>
      </thead>
      <tbody>
        {students.map((student, i) => {
          const row = draft.get(student.id);
          return (
            <tr key={student.id}>
              <td className="os-table__muted">{i + 1}</td>
              <td className="os-fw-500">{student.full_name}</td>
              <td className="os-table__mono">{student.index_number}</td>
              <td>
                <NumberInput
                  id={`marks-${student.id}`}
                  label=""
                  min={0}
                  max={maxMarks}
                  value={row.marks}
                  disabled={row.isAbsent}
                  onChange={(_e, { value }) => draft.setMarks(student.id, Number(value))}
                  size="sm"
                  hideSteppers
                />
              </td>
              <td>
                <Checkbox id={`absent-${student.id}`} labelText="AB" checked={row.isAbsent} onChange={(_e, { checked }) => draft.setAbsent(student.id, checked)} />
              </td>
              <td>
                {row.isAbsent && <Tag type="red" size="sm">AB</Tag>}
                {draft.isUnsaved(student.id) && <Tag type="high-contrast" size="sm">Unsaved</Tag>}
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}
