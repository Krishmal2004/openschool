export const GRADE_MIN = 1;
export const GRADE_MAX = 13;

export function gradeRangeErrors(from: number | "", to: number | "") {
  const fromOut = from !== "" && (from < GRADE_MIN || from > GRADE_MAX);
  const toOut = to !== "" && (to < GRADE_MIN || to > GRADE_MAX);
  const inverted = from !== "" && to !== "" && Number(to) < Number(from);
  return { fromOut, toOut, invalid: fromOut || toOut || inverted };
}
