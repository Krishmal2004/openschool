import RemoveIconButton from "../../../../components/common/RemoveIconButton";

export default function RepeatableRow({
  children,
  onRemove,
}: {
  children: React.ReactNode;
  onRemove: () => void;
}) {
  return (
    <div style={{ display: "flex", alignItems: "flex-end", gap: "0.5rem", marginBottom: "0.75rem" }}>
      <div style={{ flex: 1, display: "flex", gap: "0.5rem" }}>{children}</div>
      <RemoveIconButton size="md" onClick={onRemove} />
    </div>
  );
}
