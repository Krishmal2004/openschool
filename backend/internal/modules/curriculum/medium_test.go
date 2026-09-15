package curriculum

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

type mediumStoreStub struct {
	rows   int64
	exists bool
}

func (s *mediumStoreStub) createMedium(context.Context, string) (Medium, error) {
	return Medium{Name: "Sinhala"}, nil
}
func (s *mediumStoreStub) listMediums(context.Context) ([]Medium, error) {
	return []Medium{{Name: "Sinhala"}}, nil
}
func (s *mediumStoreStub) updateMedium(context.Context, uuid.UUID, string) (Medium, error) {
	return Medium{Name: "English"}, nil
}
func (s *mediumStoreStub) deleteMedium(context.Context, uuid.UUID) (int64, error) { return s.rows, nil }
func (s *mediumStoreStub) mediumExists(context.Context, uuid.UUID) (bool, error) {
	return s.exists, nil
}

func TestMediumDeleteDistinguishesMissingAndInUse(t *testing.T) {
	for name, test := range map[string]struct {
		store *mediumStoreStub
		want  error
	}{
		"missing": {store: &mediumStoreStub{}, want: errMediumNotFound},
		"in use":  {store: &mediumStoreStub{exists: true}, want: errMediumInUse},
	} {
		t.Run(name, func(t *testing.T) {
			err := (&mediumService{store: test.store}).delete(context.Background(), uuid.New())
			if err != test.want {
				t.Fatalf("delete() error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestMediumDeleteSucceedsWhenRowIsDeleted(t *testing.T) {
	if err := (&mediumService{store: &mediumStoreStub{rows: 1}}).delete(context.Background(), uuid.New()); err != nil {
		t.Fatal(err)
	}
}
