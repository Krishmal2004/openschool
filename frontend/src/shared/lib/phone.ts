// +94 or 0 followed by 9 digits covers Sri Lankan mobiles and landlines; mirrors backend internal/validation/phone.go.
const SRI_LANKAN_PHONE = /^(?:\+94|0)\d{9}$/;

// Empty is valid because most phone fields are optional; callers add their own presence check.
export function isValidSriLankanPhone(value: string): boolean {
  const trimmed = value.trim();
  if (trimmed === "") return true;
  return SRI_LANKAN_PHONE.test(trimmed);
}

export const PHONE_INVALID_TEXT =
  "Enter a valid Sri Lankan number, e.g. 0771234567, 0112345678, or +94771234567.";
