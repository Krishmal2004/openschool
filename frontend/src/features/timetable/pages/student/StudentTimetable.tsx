import { useMyClassTimetable } from "@/features/timetable/queries/useTimetables";
import TimetableByDay from "@/features/timetable/components/TimetableByDay";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import EmptyState from "@/shared/ui/EmptyState";
import SectionCard from "@/shared/ui/SectionCard";

export default function StudentTimetable() {
  const { data, isLoading, isError } = useMyClassTimetable();
  if (isLoading) return <LoadingSpinner />;
  return (
    <div className="os-p-8">
      <SectionCard title="Class Timetable">
        {isError || !data ? (
          <EmptyState title="No published timetable yet" description="Your class timetable will appear here once it's published by the admin." />
        ) : (
          <div className="os-grid os-gap-6">
            <TimetableByDay entries={data.entries} getRowId={(e) => e.id} middle={{ key: "teacher", header: "Teacher", render: (e) => e.teacher_name ?? "-" }} />
          </div>
        )}
      </SectionCard>
    </div>
  );
}
