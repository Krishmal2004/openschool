#!/usr/bin/env bash
# OpenSchool Unseed Script: Safely removes seeded demo class data without affecting other database records.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [ -f "$ROOT_DIR/backend/.env" ]; then
  set -a
  source <(grep -v '^#' "$ROOT_DIR/backend/.env" | grep -v '^\s*$')
  set +a
fi

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-openschool}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"

EXEC_CMD=""

if docker ps --format '{{.Names}}' | grep -q "^openschool-postgres$"; then
  EXEC_CMD="docker exec -i openschool-postgres psql -U $DB_USER -d $DB_NAME"
elif command -v psql >/dev/null 2>&1; then
  EXEC_CMD="PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME"
else
  echo "❌ Error: Neither docker container 'openschool-postgres' nor local 'psql' is available."
  exit 1
fi

echo "🧹 Removing demo class seeded data from OpenSchool database ($DB_NAME)..."

$EXEC_CMD << 'EOF'
BEGIN;

-- 1. Remove Term Marks for Demo Term
DELETE FROM term_marks WHERE term_id = '44444444-4444-4000-a000-000000000001';

-- 2. Remove Demo Term
DELETE FROM terms WHERE id = '44444444-4444-4000-a000-000000000001';

-- 3. Remove Student Attendance Records & Sessions for Demo Class
DELETE FROM attendance_records WHERE session_id IN (
  SELECT id FROM attendance_sessions WHERE class_id = '00000000-0000-4000-a000-00000000100a'
);
DELETE FROM attendance_sessions WHERE class_id = '00000000-0000-4000-a000-00000000100a';

-- 4. Remove Staff Attendance Records for Demo Teachers
DELETE FROM staff_attendance_records WHERE teacher_id IN (
  '11111111-2222-4000-a000-000000000001',
  '11111111-2222-4000-a000-000000000002'
);

-- 5. Remove Demo Class Subject Assignments & Student Enrollments
DELETE FROM class_subject_teachers WHERE class_id = '00000000-0000-4000-a000-00000000100a';
DELETE FROM class_students WHERE class_id = '00000000-0000-4000-a000-00000000100a';

-- 6. Remove Demo Class
DELETE FROM classes WHERE id = '00000000-0000-4000-a000-00000000100a';

-- 7. Remove Demo Subjects
DELETE FROM subjects WHERE id IN (
  '33333333-3333-4000-a000-000000000001',
  '33333333-3333-4000-a000-000000000002'
);

-- 8. Remove Demo Student Guardians & Guardian Profiles
DELETE FROM student_guardians WHERE guardian_id IN (
  SELECT id FROM guardians WHERE email LIKE '%@demo.openschool.edu.lk'
);
DELETE FROM guardians WHERE email LIKE '%@demo.openschool.edu.lk';

-- 9. Remove Demo Student Profiles
DELETE FROM student_profiles WHERE index_number LIKE 'DEMO-STU-%';

-- 10. Remove Demo Teacher Profiles
DELETE FROM teacher_profiles WHERE employee_number LIKE 'DEMO-EMP-%';

-- 11. Remove Demo Users (Teachers, Students, Parents)
DELETE FROM users WHERE email LIKE '%@demo.openschool.edu.lk';

COMMIT;
EOF

echo "──────────────────────────────────────────────"
echo "✨ Demo class data successfully removed!"
echo "   All other database records remain untouched."
echo "──────────────────────────────────────────────"
