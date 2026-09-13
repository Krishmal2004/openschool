// Single source for every color used by the dashboard/Analytics charts — reference these
// instead of a literal hex so the palette stays centralized and theme-editable in one place.
export const ACCENT = "var(--os-accent)";
export const STATUS_COLORS = {
  present: "var(--os-success)",
  late: "var(--os-warning)",
  absent: "var(--os-danger)",
  leave: "var(--os-chart-purple)",
};

// Additional categorical colors, for charts needing more series than the base theme defines.
export const CHART_BLUE = "#0f62fe";
export const CHART_PURPLE = "var(--os-chart-purple)";
export const CHART_PINK = "#d02670";

// Fallback swatches for a variable-length category list (e.g. houses without a set color).
export const CATEGORICAL_PALETTE = [ACCENT, CHART_PURPLE, STATUS_COLORS.present, STATUS_COLORS.late, STATUS_COLORS.absent, CHART_BLUE];

export const GENDER_COLORS: Record<string, string> = {
  Male: ACCENT,
  Female: CHART_PINK,
  Unspecified: "var(--os-text-tertiary)",
};
