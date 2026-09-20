import { Link, useParams } from "react-router";
import { Tabs, TabList, Tab, TabPanels } from "@carbon/react";
import { ArrowLeft } from "@carbon/icons-react";
import { useMyChildren } from "@/features/parent/queries/useParent";
import {
  ChildAttendanceTab,
  ChildEnrollmentsTab,
  ChildGuardiansTab,
  ChildMarksTab,
  ChildProgressReportsTab,
  ChildTimetableTab,
} from "@/features/parent/components/ChildAcademicTabs";
import StudentPortfolioSummary from "@/features/portfolio/components/StudentPortfolioSummary";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import EmptyState from "@/shared/ui/EmptyState";
import TabSection from "@/shared/ui/TabSection";
import { getInitials } from "@/shared/lib/name";

export default function ChildDetail() {
  const { id = "" } = useParams();
  const { data: children, isLoading } = useMyChildren();
  const child = children?.find((c) => c.id === id);

  if (isLoading) return <LoadingSpinner />;
  if (!child) {
    return (
      <div className="os-p-8">
        <EmptyState
          title="Child not found"
          description="This student isn't linked to your account."
          action={<Link to="/" className="os-table__link">Back to My Children</Link>}
        />
      </div>
    );
  }

  return (
    <div className="os-bg-layer-hover os-min-h-content">
      <div className="os-profile__banner">
        <div className="os-profile__avatar">{getInitials(child.full_name)}</div>
        <div className="os-flex-1">
          <p className="os-profile__name">{child.full_name}</p>
          <p className="os-profile__meta">
            {child.index_number}
            {child.class_name ? ` · ${child.class_name}` : ""}
            {child.grade_name ? ` · ${child.grade_name}` : ""}
          </p>
        </div>
        <div className="os-profile__actions">
          <Link to="/" className="os-table__link os-flex os-items-center os-gap-1h">
            <ArrowLeft size={16} />
            Back
          </Link>
        </div>
      </div>

      <div className="os-py-6 os-px-8">
        <Tabs>
          <TabList aria-label="Child sections">
            <Tab>Attendance</Tab>
            <Tab>Marks</Tab>
            <Tab>Timetable</Tab>
            <Tab>Enrolments</Tab>
            <Tab>Progress Reports</Tab>
            <Tab>Portfolio Details</Tab>
            <Tab>Guardians</Tab>
          </TabList>
          <TabPanels>
            <TabSection title="Attendance" padded={false}><ChildAttendanceTab studentId={child.id} /></TabSection>
            <TabSection title="Term Marks"><ChildMarksTab studentId={child.id} /></TabSection>
            <TabSection title="Timetable"><ChildTimetableTab studentId={child.id} /></TabSection>
            <TabSection title="Enrolled Subjects"><ChildEnrollmentsTab studentId={child.id} /></TabSection>
            <TabSection title="Narrative Progress Reports"><ChildProgressReportsTab studentId={child.id} /></TabSection>
            <TabSection title="Co-curricular & Portfolio Details"><StudentPortfolioSummary studentId={child.id} /></TabSection>
            <TabSection title="Linked Guardians"><ChildGuardiansTab studentId={child.id} /></TabSection>
          </TabPanels>
        </Tabs>
      </div>
    </div>
  );
}
