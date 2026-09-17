import { CheckmarkFilled, CloseFilled, Time, Certificate } from "@carbon/icons-react";
import { STATUS_STYLES, type Status } from "@/features/attendance/constants";

const STATUS_ICONS = { present: CheckmarkFilled, absent: CloseFilled, late: Time, excused: Certificate };

export default function StatusButton({
  value,
  selected,
  onClick,
}: {
  value: NonNullable<Status>;
  selected: boolean;
  onClick: () => void;
}) {
  const Icon = STATUS_ICONS[value];
  const className = selected ? `os-status-toggle os-status-toggle--selected os-status-toggle--${value}` : "os-status-toggle";
  return (
    <button className={className} onClick={onClick} aria-pressed={selected}>
      <Icon size={12} />
      {STATUS_STYLES[value].label}
    </button>
  );
}
