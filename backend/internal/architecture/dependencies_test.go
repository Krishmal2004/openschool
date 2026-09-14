package architecture

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

const modulePath = "github.com/openschool-org/openschool/"

// sqlcServiceDebt is a ratchet: existing files may be migrated away from sqlc,
// but new service-layer imports are rejected. Remove entries as modules move.
var sqlcServiceDebt = map[string]bool{
	"services/attendance.go": true, "services/audit.go": true, "services/class.go": true,
	"services/curriculum.go": true, "services/curriculum_preset.go": true,
	"services/notifications/notification.go": true, "services/position.go": true,
	"services/prefect.go": true, "services/promotion.go": true, "services/report_export.go": true,
	"services/section_head.go": true, "services/setup.go": true,
	"services/society.go": true, "services/staff_attendance.go": true,
	"services/student_portfolio.go":   true,
	"services/term_mark.go":           true,
	"services/timetable/generate.go":  true,
	"services/timetable/timetable.go": true,
}

func TestDependencyBoundaries(t *testing.T) {
	root := internalRoot(t)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return walkErr
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		imports := fileImports(t, path)

		if strings.HasPrefix(rel, "handlers/") {
			for _, forbidden := range []string{"db/sqlc", "internal/repositories"} {
				if imports[modulePath+forbidden] {
					t.Errorf("%s imports forbidden handler dependency %s", rel, forbidden)
				}
			}
		}

		if strings.HasPrefix(rel, "services/") && imports[modulePath+"db/sqlc"] && !sqlcServiceDebt[rel] {
			t.Errorf("%s adds a new service-layer sqlc dependency", rel)
		}

		if strings.HasPrefix(rel, "modules/") {
			for _, legacy := range []string{"internal/handlers", "internal/services", "internal/repositories"} {
				if imports[modulePath+legacy] {
					t.Errorf("%s imports legacy package %s", rel, legacy)
				}
			}
			if imports[modulePath+"db/sqlc"] && filepath.Base(rel) != "repository.go" {
				t.Errorf("%s imports sqlc outside a module repository adapter", rel)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func internalRoot(t *testing.T) string {
	t.Helper()
	_, current, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate architecture test")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(current), ".."))
}

func fileImports(t *testing.T, path string) map[string]bool {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	imports := make(map[string]bool, len(file.Imports))
	for _, spec := range file.Imports {
		value, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			t.Fatalf("parse import in %s: %v", path, err)
		}
		imports[value] = true
	}
	return imports
}
