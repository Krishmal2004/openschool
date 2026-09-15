package academics

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type streamStoreStub struct {
	streamRows int64
	groupRows  int64
	stream     Stream
	group      StreamGroup
	streamErr  error
	groupErr   error
}

func (s *streamStoreStub) createStream(context.Context, string) (Stream, error) { return Stream{}, nil }
func (s *streamStoreStub) getStream(context.Context, uuid.UUID) (Stream, error) {
	return s.stream, s.streamErr
}
func (s *streamStoreStub) listStreams(context.Context) ([]Stream, error) { return nil, nil }
func (s *streamStoreStub) updateStream(context.Context, uuid.UUID, string) (Stream, error) {
	return Stream{}, nil
}
func (s *streamStoreStub) deleteStream(context.Context, uuid.UUID) (int64, error) {
	return s.streamRows, nil
}
func (s *streamStoreStub) createGroup(context.Context, uuid.UUID, string) (StreamGroup, error) {
	return StreamGroup{}, nil
}
func (s *streamStoreStub) getGroup(context.Context, uuid.UUID) (StreamGroup, error) {
	return s.group, s.groupErr
}
func (s *streamStoreStub) listGroups(context.Context, uuid.UUID) ([]StreamGroup, error) {
	return nil, nil
}
func (s *streamStoreStub) updateGroup(context.Context, uuid.UUID, string) (StreamGroup, error) {
	return StreamGroup{}, nil
}
func (s *streamStoreStub) deleteGroup(context.Context, uuid.UUID) (int64, error) {
	return s.groupRows, nil
}

func TestStreamDeletionClassification(t *testing.T) {
	id := uuid.New()
	tests := []struct {
		name  string
		store *streamStoreStub
		group bool
		want  error
	}{
		{name: "stream deleted", store: &streamStoreStub{streamRows: 1}},
		{name: "stream missing", store: &streamStoreStub{streamErr: errors.New("missing")}, want: errStreamNotFound},
		{name: "stream in use", store: &streamStoreStub{stream: Stream{ID: id}}, want: errStreamInUse},
		{name: "group deleted", store: &streamStoreStub{groupRows: 1}, group: true},
		{name: "group missing", store: &streamStoreStub{groupErr: errors.New("missing")}, group: true, want: errStreamGroupNotFound},
		{name: "group in use", store: &streamStoreStub{group: StreamGroup{ID: id}}, group: true, want: errStreamGroupInUse},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &streamService{streams: test.store}
			var err error
			if test.group {
				err = service.deleteGroup(context.Background(), id)
			} else {
				err = service.deleteStream(context.Background(), id)
			}
			if !errors.Is(err, test.want) {
				t.Fatalf("delete error = %v, want %v", err, test.want)
			}
		})
	}
}
