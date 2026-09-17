import { DatePicker, DatePickerInput } from "@carbon/react";
import type { useCreateSession } from "@/features/attendance/queries/useAttendance";
import FormModal from "@/shared/ui/FormModal";
import { toYmd } from "@/shared/lib/date";

interface Props {
  open: boolean;
  sessionDate: string;
  onSessionDateChange: (date: string) => void;
  createSession: ReturnType<typeof useCreateSession>;
  onClose: () => void;
  onCreate: () => void;
}

export default function NewSessionModal({
  open,
  sessionDate,
  onSessionDateChange,
  createSession,
  onClose,
  onCreate,
}: Props) {
  return (
    <FormModal
      open={open}
      title="New attendance session"
      onClose={onClose}
      onSubmit={onCreate}
      isPending={createSession.isPending}
      submitDisabled={!sessionDate}
      submitLabel="Create"
      pendingLabel="Creating…"
      isError={createSession.isError}
      error={createSession.error}
      errorFallback="A session may already exist for this class on this date."
    >
      <p className="os-text-md os-c-secondary os-mb-4">
        One session per class per day. Creating it takes you straight to marking attendance.
      </p>
      <DatePicker
        datePickerType="single"
        dateFormat="Y-m-d"
        value={sessionDate}
        onChange={(dates) => onSessionDateChange(toYmd(dates[0]))}
      >
        <DatePickerInput id="session-date" labelText="Date" placeholder="YYYY-MM-DD" />
      </DatePicker>
    </FormModal>
  );
}
