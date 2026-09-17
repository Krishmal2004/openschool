import { useMyStudentProfile } from "@/features/students/queries/useStudentSelf";
import { useCurrentAcademicYear } from "@/features/school/queries/useAcademicYears";
import { useStudentEnrollments } from "@/features/students/queries/useEnrollments";
import EnrollmentsTable from "@/features/students/components/EnrollmentsTable";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import EmptyState from "@/shared/ui/EmptyState";
import ErrorMessage from "@/shared/ui/ErrorMessage";
import SectionCard from "@/shared/ui/SectionCard";

export default function StudentEnrollment() {
  const profile = useMyStudentProfile();
  const year = useCurrentAcademicYear();
  const enrollments = useStudentEnrollments(profile.data?.id ?? "", year.data?.id ?? "");

  if (profile.isLoading || year.isLoading || enrollments.isLoading) return <LoadingSpinner />;
  if (profile.isError || year.isError || enrollments.isError) {
    return (
      <div className="os-p-8">
        <ErrorMessage message="Failed to load subject enrollments" onRetry={() => { profile.refetch(); enrollments.refetch(); }} />
      </div>
    );
  }
  return (
    <div className="os-p-8">
      <SectionCard title="Enrolled Subjects" flush>
        {enrollments.data?.length ? (
          <EnrollmentsTable rows={enrollments.data} />
        ) : (
          <EmptyState title="No enrolled subjects" description="Your enrolled subjects will show up here once configured for the current academic year." />
        )}
      </SectionCard>
    </div>
  );
}
