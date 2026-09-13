package identity

import (
	"slices"

	"github.com/openschool-org/openschool/internal/models"
)

var appRolePriority = []string{models.RoleAdmin, models.RoleTeacher, models.RoleStudent, models.RoleParent}

// ResolveAppRole returns the highest-priority OpenSchool role carried by a token.
func ResolveAppRole(tokenRoles []string) string {
	for _, candidate := range appRolePriority {
		if slices.Contains(tokenRoles, candidate) {
			return candidate
		}
	}
	return ""
}
