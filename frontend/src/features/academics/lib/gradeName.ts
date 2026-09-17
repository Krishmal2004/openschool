function gradeNumber(name: string): number | null {
  const m = name.match(/\d+/);
  return m ? Number(m[0]) : null;
}

// Warns when a grade name falls outside the school's configured range.
export function gradeNameWarning(name: string, from: number | null, to: number | null): string | null {
  const n = gradeNumber(name);
  if (n === null) return "No number in this name";
  if (from !== null && n < from) return `Below your grade range (${from}-${to})`;
  if (to !== null && n > to) return `Above your grade range (${from}-${to})`;
  return null;
}
