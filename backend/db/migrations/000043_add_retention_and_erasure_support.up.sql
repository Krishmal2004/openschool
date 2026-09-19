-- Retention (docs/SECURITY_AND_PERFORMANCE_PLAYBOOK.md S11) needs to know
-- *when* a student left, not just that they did, so a nightly job can find
-- records past the retention window. Set once, when the status actually
-- transitions to 'left' (see UpdateStudentEnrollmentStatus).
ALTER TABLE student_profiles ADD COLUMN left_at TIMESTAMPTZ;

-- erased_at marks a profile the admin "erase person" flow (or the nightly
-- retention purge) has already anonymised, so it's never picked up twice
-- and the admin UI can show it was erased rather than merely inactive.
ALTER TABLE student_profiles ADD COLUMN erased_at TIMESTAMPTZ;
