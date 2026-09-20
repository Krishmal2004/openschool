import { describe, it } from "vitest";
import { render } from "@testing-library/react";
import StudentAttendanceRow from "@/features/attendance/components/StudentAttendanceRow";
import { expectNoA11yViolations } from "@/shared/testing/a11y";

const student = { id: "1", full_name: "Jane Doe", index_number: "S001" } as Parameters<typeof StudentAttendanceRow>[0]["student"];

describe("StudentAttendanceRow", () => {
  it("has no accessibility violations, unmarked", async () => {
    const { container } = render(
      <table>
        <tbody>
          <StudentAttendanceRow student={student} idx={0} status={null} note="" readOnly={false} onMark={() => {}} onNoteChange={() => {}} />
        </tbody>
      </table>,
    );
    await expectNoA11yViolations(container);
  });

  it("has no accessibility violations, read-only and marked", async () => {
    const { container } = render(
      <table>
        <tbody>
          <StudentAttendanceRow student={student} idx={0} status="present" note="" readOnly onMark={() => {}} onNoteChange={() => {}} />
        </tbody>
      </table>,
    );
    await expectNoA11yViolations(container);
  });
});
