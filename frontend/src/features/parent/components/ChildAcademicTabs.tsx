import { useState } from "react";
import { Select, SelectItem } from "@carbon/react";
import { useChildAttendance, useChildMarks } from "@/features/parent/queries/useParent";
import { useCurrentAcademicYear } from "@/features/school/queries/useAcademicYears";
import { useTerms } from "@/features/school/queries/useTerms";
import { useChildTimetable } from "@/features/timetable/queries/useTimetables";
import { useStudentEnrollments } from "@/features/students/queries/useEnrollments";
import { useProgressReports } from "@/features/portfolio/queries/useStudentPortfolio";
import { useGuardiansByStudent } from "@/features/guardians/queries/useGuardians";
import DataGrid from "@/shared/ui/DataGrid";
import AttendanceHistoryTable from "@/features/attendance/components/AttendanceHistoryTable";
import TermMarksTable from "@/features/marks/components/TermMarksTable";
import EnrollmentsTable from "@/features/students/components/EnrollmentsTable";
import ProgressReportsTable from "@/features/portfolio/components/ProgressReportsTable";
import GuardiansTable from "@/features/guardians/components/GuardiansTable";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import EmptyState from "@/shared/ui/EmptyState";
import { WEEKDAYS } from "@/shared/lib/timetable";

type Props = { studentId: string };

export function ChildTimetableTab({ studentId }: Props) {
  const { data, isLoading, isError } = useChildTimetable(studentId);
  if (isLoading) return <LoadingSpinner />;
  if (isError || !data) return <EmptyState title="No published timetable yet" description="This child's class timetable will appear here once it's published." />;

  return (
    <div className="os-grid os-gap-4">
      {WEEKDAYS.map((day) => {
        const entries = data.entries.filter((e) => e.day_of_week === day.value).sort((a, b) => a.period_number - b.period_number);
        if (entries.length === 0) return null;
        return (
          <div key={day.value}>
            <h3 className="os-text-md os-fw-600 os-mt-0 os-mx-0 os-mb-2">{day.label}</h3>
            <DataGrid
              rows={entries}
              getRowId={(e) => e.id}
              pagination={false}
              noHover
              columns={[
                { key: "period", header: "Period", render: (e) => `P${e.period_number}` },
                { key: "subject", header: "Subject", render: (e) => e.subject_name ?? "-" },
                { key: "teacher", header: "Teacher", render: (e) => e.teacher_name ?? "-" },
              ]}
            />
          </div>
        );
      })}
    </div>
  );
}

export function ChildAttendanceTab({ studentId }: Props) {
  const { data: records, isLoading } = useChildAttendance(studentId);
  if (isLoading) return <LoadingSpinner />;
  if (!records?.length) return <EmptyState title="No attendance recorded yet" description="Records will show up here once a class session is marked." />;
  return <AttendanceHistoryTable rows={records} />;
}

export function ChildMarksTab({ studentId }: Props) {
  const { data: currentYear } = useCurrentAcademicYear();
  const { data: terms } = useTerms(currentYear?.id);
  const [termId, setTermId] = useState("");
  const { data: marks, isLoading } = useChildMarks(studentId, termId);

  return (
    <div>
      <Select id="child-marks-term" labelText="Term" value={termId} onChange={(e) => setTermId(e.target.value)} className="os-max-w-20 os-mb-5">
        <SelectItem value="" text="Choose a term…" />
        {terms?.map((t) => <SelectItem key={t.id} value={t.id} text={t.name} />)}
      </Select>
      {!termId ? (
        <EmptyState title="Pick a term" description="Choose a term to see marks recorded for it." />
      ) : isLoading ? (
        <LoadingSpinner />
      ) : marks?.length ? (
        <TermMarksTable rows={marks} />
      ) : (
        <EmptyState title="No marks yet" description="Marks for this term haven't been recorded yet." />
      )}
    </div>
  );
}

export function ChildEnrollmentsTab({ studentId }: Props) {
  const { data: currentYear } = useCurrentAcademicYear();
  const { data: enrollments, isLoading } = useStudentEnrollments(studentId, currentYear?.id ?? "");
  if (isLoading) return <LoadingSpinner />;
  if (!enrollments?.length) return <EmptyState title="No enrolled subjects" description="This child is not enrolled in any subjects for the current year." />;
  return <EnrollmentsTable rows={enrollments} />;
}

export function ChildProgressReportsTab({ studentId }: Props) {
  const { data: reports, isLoading } = useProgressReports(studentId);
  if (isLoading) return <LoadingSpinner />;
  if (!reports?.length) return <EmptyState title="No progress reports yet" description="Narrative progress reports will show up here once posted by teachers." />;
  return <ProgressReportsTable rows={reports} />;
}

export function ChildGuardiansTab({ studentId }: Props) {
  const { data: guardians, isLoading } = useGuardiansByStudent(studentId);
  if (isLoading) return <LoadingSpinner />;
  if (!guardians?.length) return <EmptyState title="No guardians linked" description="No linked guardians found for this child." />;
  return <GuardiansTable rows={guardians} />;
}
