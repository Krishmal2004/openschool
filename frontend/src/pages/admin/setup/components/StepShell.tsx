export default function StepShell({
  icon: Icon,
  title,
  subtitle,
  children,
}: {
  icon: React.ComponentType<{ size?: number; style?: React.CSSProperties }>;
  title: string;
  subtitle: string;
  children: React.ReactNode;
}) {
  return (
    <div>
      <div className="os-wizard-step-header">
        <div
          style={{
            width: "2.5rem",
            height: "2.5rem",
            borderRadius: "50%",
            background: "var(--os-accent-light)",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            flexShrink: 0,
          }}
        >
          <Icon size={20} style={{ fill: "var(--os-accent)" }} />
        </div>
        <div>
          <h2 style={{ margin: 0, fontSize: "1.125rem", fontWeight: 600, color: "var(--os-text-primary)" }}>{title}</h2>
          <p style={{ margin: 0, fontSize: "0.8125rem", color: "var(--os-text-secondary)" }}>{subtitle}</p>
        </div>
      </div>
      {children}
    </div>
  );
}
