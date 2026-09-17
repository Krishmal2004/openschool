export function capitalize(value: string | null | undefined, fallback = "-"): string {
  if (!value) return fallback;
  return value[0].toUpperCase() + value.slice(1);
}
