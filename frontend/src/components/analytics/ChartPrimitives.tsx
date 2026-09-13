// Dependency-free chart primitives (hand-rolled SVG/CSS) shared by the dashboard and Analytics page.
import type { GrowthPoint } from "../../services/dashboardAnalytics";
import SectionHeader from "../common/SectionHeader";
import { ACCENT } from "./chartColors";

export function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="os-section">
      <SectionHeader title={title} />
      <div className="os-section__body">{children}</div>
    </div>
  );
}

export function StatTile({ label, value, color }: { label: string; value: string | number; color: string }) {
  return (
    <div style={{ background: "var(--os-layer)", border: "1px solid var(--os-border-subtle)", borderTop: `3px solid ${color}`, padding: "1.25rem 1.5rem" }}>
      <p style={{ margin: "0 0 0.5rem", fontSize: "0.6875rem", fontWeight: 600, textTransform: "uppercase", letterSpacing: "0.08em", color: "var(--os-text-secondary)" }}>
        {label}
      </p>
      <p style={{ margin: 0, fontSize: "2.25rem", fontWeight: 300, color: "var(--os-text-primary)", lineHeight: 1 }}>{value}</p>
    </div>
  );
}

export function BarList({
  rows,
  color = ACCENT,
  formatValue,
}: {
  rows: { label: string; value: number }[];
  color?: string;
  formatValue?: (v: number) => string;
}) {
  const max = Math.max(1, ...rows.map((r) => r.value));
  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "0.625rem" }}>
      {rows.map((r) => (
        <div key={r.label} title={`${r.label}: ${formatValue ? formatValue(r.value) : r.value}`}>
          <div style={{ display: "flex", justifyContent: "space-between", fontSize: "0.75rem", marginBottom: "0.2rem" }}>
            <span style={{ color: "var(--os-text-secondary)" }}>{r.label}</span>
            <span style={{ fontWeight: 600, color: "var(--os-text-primary)" }}>{formatValue ? formatValue(r.value) : r.value}</span>
          </div>
          <div style={{ height: "6px", background: "var(--os-border-subtle)", borderRadius: "3px" }}>
            <div style={{ width: `${(r.value / max) * 100}%`, height: "100%", background: color, borderRadius: "3px" }} />
          </div>
        </div>
      ))}
      {rows.length === 0 && <p style={{ fontSize: "0.8125rem", color: "var(--os-text-tertiary)", margin: 0 }}>No data yet.</p>}
    </div>
  );
}

export function Sparkline({
  points,
  color = ACCENT,
  min,
  max,
  formatValue,
  emptyMessage = "No data yet.",
}: {
  points: { label: string; value: number }[];
  color?: string;
  min?: number;
  max?: number;
  formatValue?: (v: number) => string;
  emptyMessage?: string;
}) {
  if (points.length === 0) return <p style={{ fontSize: "0.8125rem", color: "var(--os-text-tertiary)", margin: 0 }}>{emptyMessage}</p>;
  const w = 100;
  const h = 32;
  const values = points.map((p) => p.value);
  const lo = min ?? Math.min(...values);
  const hi = max ?? Math.max(...values);
  const range = hi - lo || 1;
  const stepX = points.length > 1 ? w / (points.length - 1) : 0;
  const coords = points.map((p, i) => ({ x: i * stepX, y: h - ((p.value - lo) / range) * h }));
  const path = coords.map((c, i) => `${i === 0 ? "M" : "L"}${c.x.toFixed(1)},${c.y.toFixed(1)}`).join(" ");
  const fmt = formatValue ?? ((v: number) => `${v}`);

  return (
    <svg viewBox={`0 0 ${w} ${h}`} preserveAspectRatio="none" style={{ width: "100%", height: "64px" }}>
      <path d={path} style={{ fill: "none", stroke: color }} strokeWidth={2} vectorEffect="non-scaling-stroke" />
      {coords.map((c, i) => (
        <circle key={points[i].label} cx={c.x} cy={c.y} r={1.6} style={{ fill: color }}>
          <title>{`${points[i].label}: ${fmt(points[i].value)}`}</title>
        </circle>
      ))}
    </svg>
  );
}

export function TrendSummary({ title, points, color }: { title: string; points: GrowthPoint[]; color: string }) {
  const latest = points.length > 0 ? points[points.length - 1] : null;
  const prior = points.length > 1 ? points[points.length - 2] : null;
  const delta = latest && prior ? latest.count - prior.count : null;
  const deltaColor = !delta ? "var(--os-text-secondary)" : delta > 0 ? "var(--os-success)" : "var(--os-danger)";

  return (
    <Section title={title}>
      {latest ? (
        <>
          <div style={{ display: "flex", alignItems: "baseline", gap: "0.625rem", marginBottom: "0.75rem", flexWrap: "wrap" }}>
            <span style={{ fontSize: "2rem", fontWeight: 300, color: "var(--os-text-primary)", lineHeight: 1 }}>{latest.count}</span>
            <span style={{ fontSize: "0.75rem", color: "var(--os-text-tertiary)" }}>as of {latest.label}</span>
            {delta !== null && (
              <span style={{ fontSize: "0.75rem", fontWeight: 600, color: deltaColor, marginLeft: "auto" }}>
                {delta > 0 ? "+" : ""}
                {delta} vs {prior!.label}
              </span>
            )}
          </div>
          <Sparkline
            points={points.map((p) => ({ label: p.label, value: p.count }))}
            color={color}
            formatValue={(v) => `${v}`}
          />
        </>
      ) : (
        <p style={{ fontSize: "0.8125rem", color: "var(--os-text-tertiary)", margin: 0 }}>No data yet.</p>
      )}
    </Section>
  );
}

export function DonutChart({
  slices,
  size = 140,
  thickness = 18,
}: {
  slices: { label: string; value: number; color: string }[];
  size?: number;
  thickness?: number;
}) {
  const total = slices.reduce((sum, s) => sum + s.value, 0);
  const radius = (size - thickness) / 2;
  const circumference = 2 * Math.PI * radius;
  let cumulative = 0;

  return (
    <div style={{ display: "flex", alignItems: "center", gap: "1.5rem", flexWrap: "wrap" }}>
      <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} style={{ flexShrink: 0 }}>
        <g transform={`rotate(-90 ${size / 2} ${size / 2})`}>
          {total === 0 ? (
            <circle cx={size / 2} cy={size / 2} r={radius} style={{ fill: "none", stroke: "var(--os-border-subtle)" }} strokeWidth={thickness} />
          ) : (
            slices.map((s) => {
              const fraction = s.value / total;
              const dash = fraction * circumference;
              const offset = -(cumulative / total) * circumference;
              cumulative += s.value;
              return (
                <circle
                  key={s.label}
                  cx={size / 2}
                  cy={size / 2}
                  r={radius}
                  style={{ fill: "none", stroke: s.color }}
                  strokeWidth={thickness}
                  strokeDasharray={`${dash} ${circumference - dash}`}
                  strokeDashoffset={offset}
                >
                  <title>{`${s.label}: ${s.value} (${Math.round(fraction * 100)}%)`}</title>
                </circle>
              );
            })
          )}
        </g>
        <text
          x="50%"
          y="50%"
          textAnchor="middle"
          dominantBaseline="central"
          style={{ fontSize: "1.1rem", fontWeight: 600, fill: "var(--os-text-primary)" }}
        >
          {total}
        </text>
      </svg>
      <div style={{ display: "flex", flexDirection: "column", gap: "0.375rem", flex: 1, minWidth: "8rem" }}>
        {slices.map((s) => (
          <div key={s.label} style={{ display: "flex", alignItems: "center", gap: "0.5rem", fontSize: "0.75rem" }}>
            <span style={{ width: "0.625rem", height: "0.625rem", borderRadius: "2px", background: s.color, flexShrink: 0 }} />
            <span style={{ color: "var(--os-text-secondary)" }}>{s.label}</span>
            <span style={{ fontWeight: 600, color: "var(--os-text-primary)", marginLeft: "auto" }}>
              {s.value}
              {total > 0 && ` (${Math.round((s.value / total) * 100)}%)`}
            </span>
          </div>
        ))}
        {slices.length === 0 && <p style={{ fontSize: "0.8125rem", color: "var(--os-text-tertiary)", margin: 0 }}>No data yet.</p>}
      </div>
    </div>
  );
}
