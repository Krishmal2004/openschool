package identity

import (
	"testing"

	"github.com/openschool-org/openschool/internal/models"
)

func TestResolveAppRole(t *testing.T) {
	tests := []struct {
		name  string
		roles []string
		want  string
	}{
		{name: "no application role", roles: []string{"offline_access"}, want: ""},
		{name: "teacher", roles: []string{models.RoleTeacher}, want: models.RoleTeacher},
		{name: "priority", roles: []string{models.RoleStudent, models.RoleAdmin}, want: models.RoleAdmin},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ResolveAppRole(test.roles); got != test.want {
				t.Fatalf("ResolveAppRole() = %q, want %q", got, test.want)
			}
		})
	}
}
