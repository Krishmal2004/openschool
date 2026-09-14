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

// legacyRepositoryDebt is a ratchet for the horizontal repository package.
// Entries are removed as their remaining feature consumers migrate to modules.
var legacyRepositoryDebt = map[string]bool{
	"repositories/auth.go":          true,
	"repositories/job_checks.go":    true,
	"repositories/job_scheduler.go": true,
	"repositories/school.go":        true,
	"repositories/student.go":       true,
	"repositories/teacher.go":       true,
	"repositories/user.go":          true,
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
		if strings.HasPrefix(rel, "repositories/") && !legacyRepositoryDebt[rel] {
			t.Errorf("%s adds a new horizontal repository; add persistence to its owning module", rel)
		}

		if strings.HasPrefix(rel, "handlers/") {
			for _, forbidden := range []string{"db/sqlc", "internal/repositories"} {
				if imports[modulePath+forbidden] {
					t.Errorf("%s imports forbidden handler dependency %s", rel, forbidden)
				}
			}
		}

		if strings.HasPrefix(rel, "services/") && imports[modulePath+"db/sqlc"] {
			t.Errorf("%s imports forbidden service-layer sqlc dependency", rel)
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
