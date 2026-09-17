import { useMyAttendance } from "@/features/students/queries/useStudentSelf";
import AttendanceHistoryTable from "@/features/attendance/components/AttendanceHistoryTable";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import EmptyState from "@/shared/ui/EmptyState";
import SectionCard from "@/shared/ui/SectionCard";

export default function StudentAttendance() {
  const { data: records, isLoading } = useMyAttendance();
  if (isLoading) return <LoadingSpinner />;
  return (
    <div className="os-p-8">
      <SectionCard title="Attendance History" flush>
        {records?.length ? (
          <AttendanceHistoryTable rows={records} />
        ) : (
          <EmptyState title="No attendance recorded yet" description="Records will show up here once a class session is marked." />
        )}
      </SectionCard>
    </div>
  );
}
