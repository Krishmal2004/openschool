import RemoveIconButton from "@/shared/ui/RemoveIconButton";

export default function RepeatableRow({
  children,
  onRemove,
}: {
  children: React.ReactNode;
  onRemove: () => void;
}) {
  return (
    <div className="os-flex os-items-end os-gap-2 os-mb-3">
      <div className="os-flex-1 os-flex os-gap-2">{children}</div>
      <RemoveIconButton size="md" onClick={onRemove} />
    </div>
  );
}
