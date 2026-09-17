import { useMyStudentProfile } from "@/features/students/queries/useStudentSelf";
import { useProgressReports } from "@/features/portfolio/queries/useStudentPortfolio";
import ProgressReportsTable from "@/features/portfolio/components/ProgressReportsTable";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import EmptyState from "@/shared/ui/EmptyState";
import ErrorMessage from "@/shared/ui/ErrorMessage";
import SectionCard from "@/shared/ui/SectionCard";

export default function StudentProgress() {
  const profile = useMyStudentProfile();
  const reports = useProgressReports(profile.data?.id ?? "");

  if (profile.isLoading || reports.isLoading) return <LoadingSpinner />;
  if (profile.isError || reports.isError) {
    return (
      <div className="os-p-8">
        <ErrorMessage message="Failed to load progress reports" onRetry={() => { profile.refetch(); reports.refetch(); }} />
      </div>
    );
  }
  return (
    <div className="os-p-8">
      <SectionCard title="Narrative Progress Reports" flush>
        {reports.data?.length ? (
          <ProgressReportsTable rows={reports.data} />
        ) : (
          <EmptyState title="No progress reports yet" description="Your narrative reports will show up here once created by your teachers." />
        )}
      </SectionCard>
    </div>
  );
}
