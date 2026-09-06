-- Lets a teacher record "absent for the test" for a student/subject/term
-- instead of a numeric mark. marks stays 0 (NOT NULL, existing CHECK still
-- holds) and is_absent distinguishes it from a genuine zero score, so
-- averages/totals can exclude it explicitly rather than treating an
-- absence as a real zero.
ALTER TABLE term_marks ADD COLUMN is_absent BOOLEAN NOT NULL DEFAULT FALSE;
