import {
  usePrefectAppointmentsByStudent,
  useStudentActivities,
  useLeadershipRoles,
  useStudentAwards,
  useDisciplinaryRecords,
} from "@/features/portfolio/queries/useStudentPortfolio";
import { useStudentSocietyMemberships } from "@/features/portfolio/queries/useSocieties";
import { PrefectsTable, SocietiesTable, ActivitiesTable, LeadershipTable, AwardsTable, DisciplineTable } from "@/features/portfolio/components/PortfolioTables";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import PortfolioSection from "@/shared/ui/PortfolioSection";

// Read-only rollup of everything the portfolio feature records for one student.
export default function StudentPortfolioSummary({ studentId }: { studentId: string }) {
  const prefects = usePrefectAppointmentsByStudent(studentId);
  const societies = useStudentSocietyMemberships(studentId);
  const activities = useStudentActivities(studentId);
  const leadership = useLeadershipRoles(studentId);
  const awards = useStudentAwards(studentId);
  const discipline = useDisciplinaryRecords(studentId);

  if ([prefects, societies, activities, leadership, awards, discipline].some((q) => q.isLoading)) return <LoadingSpinner />;

  return (
    <div className="os-grid os-gap-6">
      <PortfolioSection title="Prefect Appointments" isEmpty={!prefects.data?.length} emptyMessage="No prefect appointments recorded.">
        <PrefectsTable rows={prefects.data ?? []} />
      </PortfolioSection>
      <PortfolioSection title="Societies & Clubs" isEmpty={!societies.data?.length} emptyMessage="No society memberships recorded.">
        <SocietiesTable rows={societies.data ?? []} />
      </PortfolioSection>
      <PortfolioSection title="Co-curricular Activities" isEmpty={!activities.data?.length} emptyMessage="No activities logged.">
        <ActivitiesTable rows={activities.data ?? []} />
      </PortfolioSection>
      <PortfolioSection title="Leadership & Awards" isEmpty={!leadership.data?.length && !awards.data?.length} emptyMessage="No leadership roles or awards recorded.">
        {!!leadership.data?.length && <LeadershipTable rows={leadership.data} className="os-mb-4" />}
        {!!awards.data?.length && <AwardsTable rows={awards.data} />}
      </PortfolioSection>
      <PortfolioSection title="Disciplinary Records" isEmpty={!discipline.data?.length} emptyMessage="Disciplinary status is clear.">
        <DisciplineTable rows={discipline.data ?? []} />
      </PortfolioSection>
    </div>
  );
}
