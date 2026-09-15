package authz

import "testing"

func TestResolveAppRole(t *testing.T) {
	tests := []struct {
		name  string
		roles []string
		want  string
	}{
		{name: "no application role", roles: []string{"offline_access"}, want: ""},
		{name: "teacher", roles: []string{RoleTeacher}, want: RoleTeacher},
		{name: "priority", roles: []string{RoleStudent, RoleAdmin}, want: RoleAdmin},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ResolveAppRole(test.roles); got != test.want {
				t.Fatalf("ResolveAppRole() = %q, want %q", got, test.want)
			}
		})
	}
}
