-- Global search and, later, paginated list search (docs/SECURITY_AND_PERFORMANCE_PLAYBOOK.md
-- section 4.2) filter with ILIKE '%term%', which a plain b-tree index cannot
-- serve. pg_trgm's GIN indexes support that pattern without a full scan,
-- which matters once these tables hold thousands of rows.
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_student_profiles_full_name_trgm ON student_profiles USING GIN (full_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_student_profiles_index_number_trgm ON student_profiles USING GIN (index_number gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_teacher_profiles_full_name_trgm ON teacher_profiles USING GIN (full_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_teacher_profiles_employee_number_trgm ON teacher_profiles USING GIN (employee_number gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_guardians_full_name_trgm ON guardians USING GIN (full_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_guardians_phone_trgm ON guardians USING GIN (phone gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_non_academic_staff_full_name_trgm ON non_academic_staff USING GIN (full_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_non_academic_staff_employee_number_trgm ON non_academic_staff USING GIN (employee_number gin_trgm_ops);
