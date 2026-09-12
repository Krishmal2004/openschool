package validation

import (
	"errors"
	"regexp"
)

// ErrInvalidPhone is returned when a phone number doesn't match the Sri Lankan format.
var ErrInvalidPhone = errors.New("phone number must be a valid Sri Lankan number (e.g. 0771234567, 0112345678, or +94771234567)")

// sriLankanPhone matches a leading +94 or 0 followed by 9 digits, covering both mobile numbers and landlines with a 2-digit area code.
var sriLankanPhone = regexp.MustCompile(`^(?:\+94|0)\d{9}$`)

// IsValidSriLankanPhone reports whether s is empty (optional fields) or a validly formatted Sri Lankan phone number.
func IsValidSriLankanPhone(s string) bool {
	if s == "" {
		return true
	}
	return sriLankanPhone.MatchString(s)
}
