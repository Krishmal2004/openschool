interface Props {
  label: string;
  bg: string;
  border: string;
  color: string;
}

// Colours come from the attendance status maps, so they stay inline.
export default function StatusTag({ label, bg, border, color }: Props) {
  return (
    <span className="os-status-tag" style={{ borderColor: border, background: bg, color }}>
      {label}
    </span>
  );
}
