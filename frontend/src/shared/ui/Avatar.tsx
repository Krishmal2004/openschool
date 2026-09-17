import { getInitials } from "@/shared/lib/name";

const SIZE = { sm: "2.25rem", md: "3.25rem" };
const FONT_SIZE = { sm: "0.75rem", md: "1.1rem" };

interface Props {
  name: string;
  size?: "sm" | "md";
}

export default function Avatar({ name, size = "md" }: Props) {
  if (size === "md") {
    return <div className="os-profile__avatar">{getInitials(name)}</div>;
  }
  return (
    <div className="os-bg-accent-dark os-rounded-full os-c-layer os-flex os-items-center os-justify-center os-fw-700 os-shrink-0" style={{ width: SIZE[size], height: SIZE[size], fontSize: FONT_SIZE[size] }}
    >
      {getInitials(name)}
    </div>
  );
}
