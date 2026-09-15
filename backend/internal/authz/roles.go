// Package authz owns OpenSchool application roles and role resolution.
package authz

import "slices"

const (
	RoleAdmin   = "admin"
	RoleTeacher = "teacher"
	RoleStudent = "student"
	RoleParent  = "parent"
)

var appRolePriority = []string{RoleAdmin, RoleTeacher, RoleStudent, RoleParent}

// ResolveAppRole returns the highest-priority OpenSchool role carried by a token.
func ResolveAppRole(tokenRoles []string) string {
	for _, candidate := range appRolePriority {
		if slices.Contains(tokenRoles, candidate) {
			return candidate
		}
	}
	return ""
}
