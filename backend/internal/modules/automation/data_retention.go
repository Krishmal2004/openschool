package automation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/openschool-org/openschool/internal/modules/notifications"
)

// DataRetentionAgentName is this agent's stable job_settings/job_runs identifier.
const DataRetentionAgentName = "data_retention_agent"

// studentRetentionYears bounds how long a left student's personal data is
// kept before the nightly purge anonymises it (S11,
// docs/SECURITY_AND_PERFORMANCE_PLAYBOOK.md section 3): 7 years after
// leaving, matching common school-records retention practice.
const studentRetentionYears = 7

// DataRetentionAgent anonymises left students' personal data once it has
// sat past the retention window, keeping the profile row itself so
// historical marks/attendance stay attributable in aggregate.
type DataRetentionAgent struct {
	checks   *Repository
	notifSvc *notifications.NotificationService
}

// NewDataRetentionAgent constructs a DataRetentionAgent with its dependencies.
func NewDataRetentionAgent(checks *Repository, notifSvc *notifications.NotificationService) *DataRetentionAgent {
	return &DataRetentionAgent{checks: checks, notifSvc: notifSvc}
}

// Name returns this agent's job_settings/job_runs identifier.
func (a *DataRetentionAgent) Name() string { return DataRetentionAgentName }

// Schedule returns this agent's cron expression: daily 03:00.
func (a *DataRetentionAgent) Schedule() string { return "0 3 * * *" }

// Description returns the one-line summary shown on the Automation panel.
func (a *DataRetentionAgent) Description() string {
	return fmt.Sprintf("Anonymises left students' personal data %d years after they left, keeping the profile row so historical marks/attendance stay attributable.", studentRetentionYears)
}

// Run anonymises every student past the retention window and notifies admins of what it did.
func (a *DataRetentionAgent) Run(ctx context.Context) (Result, error) {
	candidates, err := a.checks.ListStudentsPastRetention(ctx, studentRetentionYears)
	if err != nil {
		return Result{}, fmt.Errorf("retention scan: %w", err)
	}
	if len(candidates) == 0 {
		return Result{Summary: "no records past the retention window"}, nil
	}

	var errs []error
	names := make([]string, 0, len(candidates))
	erased := 0
	for _, c := range candidates {
		if err := a.checks.AnonymizeStudentProfile(ctx, c.ID); err != nil {
			errs = append(errs, fmt.Errorf("anonymise %s: %w", c.ID, err))
			continue
		}
		if c.HasUser {
			if err := a.checks.EraseStudentUser(ctx, c.UserID); err != nil {
				// The profile is already anonymised — the PDPA-sensitive part is
				// done — so a failure here is logged, not treated as the whole
				// candidate having failed.
				log.Printf("data-retention: profile %s anonymised but failed to scrub user %s: %v", c.ID, c.UserID, err)
			}
		}
		names = append(names, c.FullName)
		erased++
	}

	summary := fmt.Sprintf("anonymised %d student(s) past the %d-year retention window: %s", erased, studentRetentionYears, strings.Join(names, ", "))
	if err := notifyAdmins(ctx, a.checks, a.notifSvc, "Nightly retention purge ran", summary, "general", SeverityElevated); err != nil {
		errs = append(errs, fmt.Errorf("ran purge but failed to notify admins: %w", err))
	}

	result := Result{Summary: summary, Findings: erased}
	if len(errs) > 0 {
		return result, errors.Join(errs...)
	}
	return result, nil
}
