import { useState } from "react";
import { Select, SelectItem } from "@carbon/react";
import { useMyMarks } from "@/features/students/queries/useStudentSelf";
import { useCurrentAcademicYear } from "@/features/school/queries/useAcademicYears";
import { useTerms } from "@/features/school/queries/useTerms";
import TermMarksTable from "@/features/marks/components/TermMarksTable";
import LoadingSpinner from "@/shared/ui/LoadingSpinner";
import EmptyState from "@/shared/ui/EmptyState";
import SectionCard from "@/shared/ui/SectionCard";

export default function StudentMarks() {
  const { data: currentYear } = useCurrentAcademicYear();
  const { data: terms } = useTerms(currentYear?.id);
  const [termId, setTermId] = useState("");
  const { data: marks, isLoading } = useMyMarks(termId);

  return (
    <div className="os-p-8">
      <SectionCard title="Term Marks">
        <Select id="my-marks-term" labelText="Term" value={termId} onChange={(e) => setTermId(e.target.value)} className="os-max-w-20 os-mb-5">
          <SelectItem value="" text="Choose a term…" />
          {terms?.map((t) => <SelectItem key={t.id} value={t.id} text={t.name} />)}
        </Select>
        {!termId ? (
          <EmptyState title="Pick a term" description="Choose a term to see your marks for it." />
        ) : isLoading ? (
          <LoadingSpinner />
        ) : marks?.length ? (
          <TermMarksTable rows={marks} />
        ) : (
          <EmptyState title="No marks yet" description="Marks for this term haven't been recorded yet." />
        )}
      </SectionCard>
    </div>
  );
}
