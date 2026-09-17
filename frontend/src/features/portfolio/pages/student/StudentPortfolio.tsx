import { useMyStudentProfile } from "@/features/students/queries/useStudentSelf";
import StudentPortfolioSummary from "@/features/portfolio/components/StudentPortfolioSummary";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import ErrorMessage from "@/shared/ui/ErrorMessage";

export default function StudentPortfolio() {
  const { data: profile, isLoading, isError, refetch } = useMyStudentProfile();
  if (isLoading) return <LoadingSpinner />;
  if (isError || !profile) {
    return (
      <div className="os-p-8">
        <ErrorMessage message="Failed to load profile" onRetry={refetch} />
      </div>
    );
  }
  return (
    <div className="os-p-8">
      <div className="os-section os-mb-6">
        <div className="os-section__header">
          <h2 className="os-section__title">Student Portfolio</h2>
        </div>
      </div>
      <StudentPortfolioSummary studentId={profile.id} />
    </div>
  );
}
