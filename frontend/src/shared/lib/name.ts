// Splits a stored full_name into given/family; lossy for multi-word given names, which the joined string cannot recover.
export function splitFullName(fullName: string): { given_name: string; family_name: string } {
  const [given, ...rest] = fullName.trim().split(/\s+/);
  return { given_name: given ?? "", family_name: rest.join(" ") };
}

// Up to two initials from a full name, for avatar placeholders.
export function getInitials(fullName: string): string {
  return fullName
    .split(" ")
    .map((p) => p[0])
    .filter(Boolean)
    .slice(0, 2)
    .join("")
    .toUpperCase();
}
