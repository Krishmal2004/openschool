import { POSITION_RANK } from "@/shared/lib/constants/people";
import { useQueries } from "@tanstack/react-query";
import { useMyClasses } from "@/features/teachers/queries/useTeachers";
import { useDailySessions, classSessionsOptions } from "@/features/attendance/queries/useAttendance";
import { studentsByClassOptions } from "@/features/students/queries/useStudents";
import { useCurrentAcademicYear } from "@/features/school/queries/useAcademicYears";
import { useTerms } from "@/features/school/queries/useTerms";
import { useMyPosition, useMyLeadershipOverview } from "@/features/positions/queries/usePositions";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import ErrorMessage from "@/shared/ui/ErrorMessage";
import WelcomeBanner from "@/features/teachers/components/dashboard/WelcomeBanner";
import TodaysClasses from "@/features/teachers/components/dashboard/TodaysClasses";
import RecentSessions from "@/features/teachers/components/dashboard/RecentSessions";
import QuickActions from "@/features/teachers/components/dashboard/QuickActions";
import TodaySummary from "@/features/teachers/components/dashboard/TodaySummary";
import MyClassesPanel from "@/features/teachers/components/dashboard/MyClassesPanel";
import LeadershipPanel from "@/features/teachers/components/dashboard/LeadershipPanel";
import LeadershipOverviewPanel from "@/features/teachers/components/dashboard/LeadershipOverviewPanel";
import TimetableReviewPanel from "@/features/teachers/components/dashboard/TimetableReviewPanel";
import { todayISODate } from "@/shared/lib/date";

// Section Head and above get the extra Leadership panel.


export default function TeacherDashboard() {
  const { teacher: profile, classes: myClasses, isLoading: profileLoading, isError: profileError, refetch } = useMyClasses();
  const { data: currentYear } = useCurrentAcademicYear();
  const { data: terms } = useTerms(currentYear?.id);
  const { data: dailySessions } = useDailySessions(todayISODate());
  const { data: positionSummary } = useMyPosition();
  const showLeadershipPanel = !!positionSummary && positionSummary.rank <= POSITION_RANK.sectionHead;
  const { data: leadershipOverview } = useMyLeadershipOverview(showLeadershipPanel);
  const isSectionHead = positionSummary?.rank === POSITION_RANK.sectionHead;

  const classIds = myClasses.map((c) => c.class_id);

  const studentQueries = useQueries({ queries: classIds.map(studentsByClassOptions) });
  const sessionQueries = useQueries({ queries: classIds.map(classSessionsOptions) });

  const studentCountByClass = new Map(classIds.map((id, i) => [id, studentQueries[i]?.data?.length ?? 0]));
  const totalStudents = [...studentCountByClass.values()].reduce((sum, n) => sum + n, 0);

  const todaySessionByClass = new Map(
    (dailySessions ?? []).filter((s) => classIds.includes(s.class_id)).map((s) => [s.class_id, s]),
  );
  const markedCount = [...todaySessionByClass.values()].filter((s) => s.marked_count > 0).length;
  const pendingCount = Math.max(myClasses.length - markedCount, 0);

  const recentSessions = classIds
    .flatMap((id, i) => {
      const cls = myClasses.find((c) => c.class_id === id);
      return (sessionQueries[i]?.data ?? []).map((s) => ({ session: s, className: cls?.class_name ?? "" }));
    })
    .sort((a, b) => b.session.date.localeCompare(a.session.date))
    .slice(0, 6);

  const currentTerm = terms?.find((t) => t.is_current);

  if (profileLoading) return <LoadingSpinner />;
  if (profileError || !profile) {
    return (
      <div className="os-p-8">
        <ErrorMessage message="Failed to load your teacher profile" onRetry={refetch} />
      </div>
    );
  }

  const subjectSummary = [...new Set(myClasses.flatMap((c) => c.subjects))].join(", ") || "No subjects assigned yet";
  const rankLabel = positionSummary?.rank_label ?? "Teacher";

  return (
    <div className="os-page">
      <WelcomeBanner
        profile={profile}
        subjectSummary={subjectSummary}
        currentYearLabel={currentYear?.label}
        currentTermName={currentTerm?.name}
        pendingCount={pendingCount}
        rankLabel={rankLabel}
      />

      <div className="os-grid os-grid-cols-2-1 os-gap-6 os-items-grid-start">
        <div>
          <TodaysClasses
            loading={profileLoading}
            myClasses={myClasses}
            studentCountByClass={studentCountByClass}
            todaySessionByClass={todaySessionByClass}
          />
          <RecentSessions sessions={recentSessions} />
        </div>

        <div>
          {showLeadershipPanel && leadershipOverview && <LeadershipOverviewPanel overview={leadershipOverview} />}
          {showLeadershipPanel && positionSummary && <LeadershipPanel summary={positionSummary} />}
          {isSectionHead && currentYear && <TimetableReviewPanel academicYearId={currentYear.id} />}
          <QuickActions rankLabel={rankLabel} notifyWholeSchool={positionSummary?.notify_whole_school ?? false} />
          <TodaySummary
            markedCount={markedCount}
            pendingCount={pendingCount}
            myClassCount={myClasses.length}
            totalStudents={totalStudents}
          />
          <MyClassesPanel
            myClasses={myClasses}
            studentCountByClass={studentCountByClass}
            todaySessionByClass={todaySessionByClass}
          />
        </div>
      </div>
    </div>
  );
}
