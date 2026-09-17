import { useMyStudentProfile } from "@/features/students/queries/useStudentSelf";
import { useGuardiansByStudent } from "@/features/guardians/queries/useGuardians";
import GuardiansTable from "@/features/guardians/components/GuardiansTable";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import EmptyState from "@/shared/ui/EmptyState";
import ErrorMessage from "@/shared/ui/ErrorMessage";
import SectionCard from "@/shared/ui/SectionCard";

export default function StudentGuardians() {
  const profile = useMyStudentProfile();
  const guardians = useGuardiansByStudent(profile.data?.id ?? "");

  if (profile.isLoading || guardians.isLoading) return <LoadingSpinner />;
  if (profile.isError || guardians.isError) {
    return (
      <div className="os-p-8">
        <ErrorMessage message="Failed to load guardians" onRetry={() => { profile.refetch(); guardians.refetch(); }} />
      </div>
    );
  }
  return (
    <div className="os-p-8">
      <SectionCard title="My Guardians" flush>
        {guardians.data?.length ? (
          <GuardiansTable rows={guardians.data} />
        ) : (
          <EmptyState title="No guardians linked yet" description="Your linked guardians will appear here once configured by the administration." />
        )}
      </SectionCard>
    </div>
  );
}
