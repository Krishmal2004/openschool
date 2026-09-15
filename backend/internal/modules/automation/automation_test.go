package automation

import (
	"context"
	"errors"
	"testing"
)

type testJob struct {
	name string
	run  func(context.Context) (Result, error)
}

func (j testJob) Name() string        { return j.name }
func (j testJob) Schedule() string    { return "0 * * * *" }
func (j testJob) Description() string { return "test job" }
func (j testJob) Run(ctx context.Context) (Result, error) {
	return j.run(ctx)
}

func TestRunChecksAggregatesFindingsAndErrors(t *testing.T) {
	wantErr := errors.New("check failed")
	result, err := runChecks(context.Background(),
		func(context.Context) checkOutcome { return checkOutcome{findings: 2, label: "two findings"} },
		func(context.Context) checkOutcome { return checkOutcome{err: wantErr} },
	)

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected joined error %v, got %v", wantErr, err)
	}
	if result.Findings != 2 || result.Summary != "two findings" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestRunJobRecoveringConvertsPanicToError(t *testing.T) {
	job := testJob{name: "panic-test", run: func(context.Context) (Result, error) {
		panic("boom")
	}}

	result, err := runJobRecovering(context.Background(), job)
	if err == nil {
		t.Fatal("expected panic to become an error")
	}
	if result.Summary == "" {
		t.Fatal("expected failed result to include a summary")
	}
}

func TestNewSchedulerRejectsDuplicateNames(t *testing.T) {
	job := testJob{name: "duplicate", run: func(context.Context) (Result, error) { return Result{}, nil }}
	defer func() {
		if recover() == nil {
			t.Fatal("expected duplicate job names to panic")
		}
	}()
	NewScheduler([]Job{job, job}, nil)
}

func TestAlgorithmHelpers(t *testing.T) {
	if got := medianInt64([]int64{9, 1, 5, 3}); got != 4 {
		t.Fatalf("medianInt64() = %d, want 4", got)
	}
	if got := clamp01(-0.5); got != 0 {
		t.Fatalf("clamp01(-0.5) = %v, want 0", got)
	}
	if got := clamp01(1.5); got != 1 {
		t.Fatalf("clamp01(1.5) = %v, want 1", got)
	}
}
