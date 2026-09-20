#!/usr/bin/env bash
# OpenSchool Seed Script: Populates a demo class with teachers, students, parents, attendance, and term marks.
# All seeded entities use deterministic UUIDs so they can be safely removed by scripts/unseed_demo_class.sh.
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
  echo "   Please start Postgres first using 'make setup' or 'cd backend && docker compose up -d'."
  exit 1
fi

echo "🌱 Seeding demo class data into OpenSchool database ($DB_NAME)..."

$EXEC_CMD << 'EOF'
BEGIN;

-- 1. Ensure School exists (shared single-row school record)
INSERT INTO school (id, name, address, phone, email, created_at, updated_at)
VALUES ('00000000-0000-4000-a000-000000000001', 'OpenSchool Central High', '123 Main Street, Colombo', '0112345678', 'info@openschool.edu.lk', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 2. Demo Academic Year (2026)
INSERT INTO academic_years (id, year, start_date, end_date, is_current, created_at)
VALUES ('00000000-0000-4000-a000-000000002026', '2026', '2026-01-01', '2026-12-31', TRUE, NOW())
ON CONFLICT (year) DO UPDATE SET is_current = TRUE;

-- 3. Demo Grade 10
INSERT INTO grades (id, name, sort_order, created_at)
VALUES ('00000000-0000-4000-a000-000000000010', 'Grade 10', 10, NOW())
ON CONFLICT (name) DO NOTHING;

-- 4. Demo Teachers
-- Form Teacher: Mr. Kavinda Perera
INSERT INTO users (id, email, full_name, role, is_active, created_at, updated_at)
VALUES ('11111111-1111-4000-a000-000000000001', 'kavinda.teacher@demo.openschool.edu.lk', 'Mr. Kavinda Perera', 'teacher', TRUE, NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO teacher_profiles (id, user_id, full_name, employee_number, joined_date, phone, created_at, updated_at)
VALUES ('11111111-2222-4000-a000-000000000001', '11111111-1111-4000-a000-000000000001', 'Mr. Kavinda Perera', 'DEMO-EMP-1001', '2022-01-15', '0771112233', NOW(), NOW())
ON CONFLICT (employee_number) DO NOTHING;

-- Subject Teacher: Ms. Nimali Fernando
INSERT INTO users (id, email, full_name, role, is_active, created_at, updated_at)
VALUES ('11111111-1111-4000-a000-000000000002', 'nimali.teacher@demo.openschool.edu.lk', 'Ms. Nimali Fernando', 'teacher', TRUE, NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO teacher_profiles (id, user_id, full_name, employee_number, joined_date, phone, created_at, updated_at)
VALUES ('11111111-2222-4000-a000-000000000002', '11111111-1111-4000-a000-000000000002', 'Ms. Nimali Fernando', 'DEMO-EMP-1002', '2023-03-01', '0772223344', NOW(), NOW())
ON CONFLICT (employee_number) DO NOTHING;

-- 5. Demo Class "Grade 10-A"
INSERT INTO classes (id, grade_id, academic_year_id, form_teacher_id, name, created_at)
VALUES (
  '00000000-0000-4000-a000-00000000100a',
  '00000000-0000-4000-a000-000000000010',
  '00000000-0000-4000-a000-000000002026',
  '11111111-2222-4000-a000-000000000001',
  'Grade 10-A',
  NOW()
)
ON CONFLICT (grade_id, academic_year_id, name) DO NOTHING;

-- 6. Demo Subjects ("Mathematics", "Science")
INSERT INTO subjects (id, code, name, category, created_at)
VALUES
  ('33333333-3333-4000-a000-000000000001', 'MATH10-DEMO', 'Mathematics (Demo)', 'core', NOW()),
  ('33333333-3333-4000-a000-000000000002', 'SCI10-DEMO', 'Science (Demo)', 'core', NOW())
ON CONFLICT (code) DO NOTHING;

-- Assign Teachers to Class Subjects
INSERT INTO class_subject_teachers (class_id, subject_id, teacher_id)
VALUES
  ('00000000-0000-4000-a000-00000000100a', '33333333-3333-4000-a000-000000000001', '11111111-2222-4000-a000-000000000001'),
  ('00000000-0000-4000-a000-00000000100a', '33333333-3333-4000-a000-000000000002', '11111111-2222-4000-a000-000000000002')
ON CONFLICT (class_id, subject_id) DO NOTHING;

-- 7. Deterministic Demo Students & Guardians
-- Student 1
INSERT INTO users (id, email, full_name, role, is_active)
VALUES ('22222222-1111-4000-a000-000000000001', 'amaya.j@demo.openschool.edu.lk', 'Amaya Jayawardena', 'student', TRUE)
ON CONFLICT (email) DO NOTHING;
INSERT INTO student_profiles (id, user_id, full_name, index_number, phone)
VALUES ('22222222-2222-4000-a000-000000000001', '22222222-1111-4000-a000-000000000001', 'Amaya Jayawardena', 'DEMO-STU-1001', '0710000001')
ON CONFLICT (index_number) DO NOTHING;

INSERT INTO users (id, email, full_name, role, is_active)
VALUES ('33333333-1111-4000-a000-000000000001', 'sunil.j@demo.openschool.edu.lk', 'Sunil Jayawardena', 'parent', TRUE)
ON CONFLICT (email) DO NOTHING;
INSERT INTO guardians (id, user_id, full_name, relationship, phone, email)
VALUES ('33333333-2222-4000-a000-000000000001', '33333333-1111-4000-a000-000000000001', 'Sunil Jayawardena', 'father', '0770000001', 'sunil.j@demo.openschool.edu.lk')
ON CONFLICT (id) DO NOTHING;
INSERT INTO student_guardians (student_id, guardian_id, is_primary_contact)
VALUES ('22222222-2222-4000-a000-000000000001', '33333333-2222-4000-a000-000000000001', TRUE)
ON CONFLICT DO NOTHING;

-- Student 2
INSERT INTO users (id, email, full_name, role, is_active)
VALUES ('22222222-1111-4000-a000-000000000002', 'bhanuka.s@demo.openschool.edu.lk', 'Bhanuka Silva', 'student', TRUE)
ON CONFLICT (email) DO NOTHING;
INSERT INTO student_profiles (id, user_id, full_name, index_number, phone)
VALUES ('22222222-2222-4000-a000-000000000002', '22222222-1111-4000-a000-000000000002', 'Bhanuka Silva', 'DEMO-STU-1002', '0710000002')
ON CONFLICT (index_number) DO NOTHING;

INSERT INTO users (id, email, full_name, role, is_active)
VALUES ('33333333-1111-4000-a000-000000000002', 'kamani.s@demo.openschool.edu.lk', 'Kamani Silva', 'parent', TRUE)
ON CONFLICT (email) DO NOTHING;
INSERT INTO guardians (id, user_id, full_name, relationship, phone, email)
VALUES ('33333333-2222-4000-a000-000000000002', '33333333-1111-4000-a000-000000000002', 'Kamani Silva', 'mother', '0770000002', 'kamani.s@demo.openschool.edu.lk')
ON CONFLICT (id) DO NOTHING;
INSERT INTO student_guardians (student_id, guardian_id, is_primary_contact)
VALUES ('22222222-2222-4000-a000-000000000002', '33333333-2222-4000-a000-000000000002', TRUE)
ON CONFLICT DO NOTHING;

-- Student 3
INSERT INTO users (id, email, full_name, role, is_active)
VALUES ('22222222-1111-4000-a000-000000000003', 'chamindu.p@demo.openschool.edu.lk', 'Chamindu Perera', 'student', TRUE)
ON CONFLICT (email) DO NOTHING;
INSERT INTO student_profiles (id, user_id, full_name, index_number, phone)
VALUES ('22222222-2222-4000-a000-000000000003', '22222222-1111-4000-a000-000000000003', 'Chamindu Perera', 'DEMO-STU-1003', '0710000003')
ON CONFLICT (index_number) DO NOTHING;

INSERT INTO users (id, email, full_name, role, is_active)
VALUES ('33333333-1111-4000-a000-000000000003', 'rohan.p@demo.openschool.edu.lk', 'Rohan Perera', 'parent', TRUE)
ON CONFLICT (email) DO NOTHING;
INSERT INTO guardians (id, user_id, full_name, relationship, phone, email)
VALUES ('33333333-2222-4000-a000-000000000003', '33333333-1111-4000-a000-000000000003', 'Rohan Perera', 'father', '0770000003', 'rohan.p@demo.openschool.edu.lk')
ON CONFLICT (id) DO NOTHING;
INSERT INTO student_guardians (student_id, guardian_id, is_primary_contact)
VALUES ('22222222-2222-4000-a000-000000000003', '33333333-2222-4000-a000-000000000003', TRUE)
ON CONFLICT DO NOTHING;

-- Student 4
INSERT INTO users (id, email, full_name, role, is_active)
VALUES ('22222222-1111-4000-a000-000000000004', 'dilini.w@demo.openschool.edu.lk', 'Dilini Wickramasinghe', 'student', TRUE)
ON CONFLICT (email) DO NOTHING;
INSERT INTO student_profiles (id, user_id, full_name, index_number, phone)
VALUES ('22222222-2222-4000-a000-000000000004', '22222222-1111-4000-a000-000000000004', 'Dilini Wickramasinghe', 'DEMO-STU-1004', '0710000004')
ON CONFLICT (index_number) DO NOTHING;

INSERT INTO users (id, email, full_name, role, is_active)
VALUES ('33333333-1111-4000-a000-000000000004', 'anusha.w@demo.openschool.edu.lk', 'Anusha Wickramasinghe', 'parent', TRUE)
ON CONFLICT (email) DO NOTHING;
INSERT INTO guardians (id, user_id, full_name, relationship, phone, email)
VALUES ('33333333-2222-4000-a000-000000000004', '33333333-1111-4000-a000-000000000004', 'Anusha Wickramasinghe', 'mother', '0770000004', 'anusha.w@demo.openschool.edu.lk')
ON CONFLICT (id) DO NOTHING;
INSERT INTO student_guardians (student_id, guardian_id, is_primary_contact)
VALUES ('22222222-2222-4000-a000-000000000004', '33333333-2222-4000-a000-000000000004', TRUE)
ON CONFLICT DO NOTHING;

-- Student 5
INSERT INTO users (id, email, full_name, role, is_active)
VALUES ('22222222-1111-4000-a000-000000000005', 'eshan.r@demo.openschool.edu.lk', 'Eshan Ratnayake', 'student', TRUE)
ON CONFLICT (email) DO NOTHING;
INSERT INTO student_profiles (id, user_id, full_name, index_number, phone)
VALUES ('22222222-2222-4000-a000-000000000005', '22222222-1111-4000-a000-000000000005', 'Eshan Ratnayake', 'DEMO-STU-1005', '0710000005')
ON CONFLICT (index_number) DO NOTHING;

INSERT INTO users (id, email, full_name, role, is_active)
VALUES ('33333333-1111-4000-a000-000000000005', 'sarath.r@demo.openschool.edu.lk', 'Sarath Ratnayake', 'parent', TRUE)
ON CONFLICT (email) DO NOTHING;
INSERT INTO guardians (id, user_id, full_name, relationship, phone, email)
VALUES ('33333333-2222-4000-a000-000000000005', '33333333-1111-4000-a000-000000000005', 'Sarath Ratnayake', 'father', '0770000005', 'sarath.r@demo.openschool.edu.lk')
ON CONFLICT (id) DO NOTHING;
INSERT INTO student_guardians (student_id, guardian_id, is_primary_contact)
VALUES ('22222222-2222-4000-a000-000000000005', '33333333-2222-4000-a000-000000000005', TRUE)
ON CONFLICT DO NOTHING;

-- Enroll all 5 students into Grade 10-A
INSERT INTO class_students (class_id, student_id)
VALUES
  ('00000000-0000-4000-a000-00000000100a', '22222222-2222-4000-a000-000000000001'),
  ('00000000-0000-4000-a000-00000000100a', '22222222-2222-4000-a000-000000000002'),
  ('00000000-0000-4000-a000-00000000100a', '22222222-2222-4000-a000-000000000003'),
  ('00000000-0000-4000-a000-00000000100a', '22222222-2222-4000-a000-000000000004'),
  ('00000000-0000-4000-a000-00000000100a', '22222222-2222-4000-a000-000000000005')
ON CONFLICT DO NOTHING;

-- 8. Seed 3 Days of Teacher Attendance (staff_attendance_records)
INSERT INTO staff_attendance_records (teacher_id, date, status, note)
VALUES
  ('11111111-2222-4000-a000-000000000001', CURRENT_DATE - INTERVAL '2 days', 'present', 'On time'),
  ('11111111-2222-4000-a000-000000000001', CURRENT_DATE - INTERVAL '1 day',  'present', 'On time'),
  ('11111111-2222-4000-a000-000000000001', CURRENT_DATE,                     'present', 'On time'),
  ('11111111-2222-4000-a000-000000000002', CURRENT_DATE - INTERVAL '2 days', 'present', 'On time'),
  ('11111111-2222-4000-a000-000000000002', CURRENT_DATE - INTERVAL '1 day',  'late',    'Traffic delay 10m'),
  ('11111111-2222-4000-a000-000000000002', CURRENT_DATE,                     'present', 'On time')
ON CONFLICT (teacher_id, date) DO NOTHING;

-- 9. Seed 3 Days of Student Attendance (attendance_sessions & attendance_records)
-- Fixed Session IDs per day:
-- Day -2: 55555555-1111-4000-a000-000000000001
-- Day -1: 55555555-1111-4000-a000-000000000002
-- Day  0: 55555555-1111-4000-a000-000000000003
INSERT INTO attendance_sessions (id, class_id, taken_by, date)
VALUES
  ('55555555-1111-4000-a000-000000000001', '00000000-0000-4000-a000-00000000100a', '11111111-2222-4000-a000-000000000001', CURRENT_DATE - INTERVAL '2 days'),
  ('55555555-1111-4000-a000-000000000002', '00000000-0000-4000-a000-00000000100a', '11111111-2222-4000-a000-000000000001', CURRENT_DATE - INTERVAL '1 day'),
  ('55555555-1111-4000-a000-000000000003', '00000000-0000-4000-a000-00000000100a', '11111111-2222-4000-a000-000000000001', CURRENT_DATE)
ON CONFLICT (class_id, date) DO NOTHING;

INSERT INTO attendance_records (session_id, student_id, status)
VALUES
  ('55555555-1111-4000-a000-000000000001', '22222222-2222-4000-a000-000000000001', 'present'),
  ('55555555-1111-4000-a000-000000000001', '22222222-2222-4000-a000-000000000002', 'present'),
  ('55555555-1111-4000-a000-000000000001', '22222222-2222-4000-a000-000000000003', 'present'),
  ('55555555-1111-4000-a000-000000000001', '22222222-2222-4000-a000-000000000004', 'late'),
  ('55555555-1111-4000-a000-000000000001', '22222222-2222-4000-a000-000000000005', 'present'),

  ('55555555-1111-4000-a000-000000000002', '22222222-2222-4000-a000-000000000001', 'present'),
  ('55555555-1111-4000-a000-000000000002', '22222222-2222-4000-a000-000000000002', 'present'),
  ('55555555-1111-4000-a000-000000000002', '22222222-2222-4000-a000-000000000003', 'absent'),
  ('55555555-1111-4000-a000-000000000002', '22222222-2222-4000-a000-000000000004', 'present'),
  ('55555555-1111-4000-a000-000000000002', '22222222-2222-4000-a000-000000000005', 'present'),

  ('55555555-1111-4000-a000-000000000003', '22222222-2222-4000-a000-000000000001', 'present'),
  ('55555555-1111-4000-a000-000000000003', '22222222-2222-4000-a000-000000000002', 'late'),
  ('55555555-1111-4000-a000-000000000003', '22222222-2222-4000-a000-000000000003', 'present'),
  ('55555555-1111-4000-a000-000000000003', '22222222-2222-4000-a000-000000000004', 'present'),
  ('55555555-1111-4000-a000-000000000003', '22222222-2222-4000-a000-000000000005', 'present')
ON CONFLICT (session_id, student_id) DO NOTHING;

-- 10. Seed Term 1 & Term Marks for Mathematics & Science
INSERT INTO terms (id, academic_year_id, name, start_date, end_date, is_current, sort_order)
VALUES ('44444444-4444-4000-a000-000000000001', '00000000-0000-4000-a000-000000002026', 'Term 1 (Demo)', '2026-01-01', '2026-04-30', TRUE, 1)
ON CONFLICT (academic_year_id, name) DO NOTHING;

INSERT INTO term_marks (student_id, subject_id, term_id, marks, max_marks, entered_by)
VALUES
  ('22222222-2222-4000-a000-000000000001', '33333333-3333-4000-a000-000000000001', '44444444-4444-4000-a000-000000000001', 85, 100, '11111111-1111-4000-a000-000000000001'),
  ('22222222-2222-4000-a000-000000000001', '33333333-3333-4000-a000-000000000002', '44444444-4444-4000-a000-000000000001', 88, 100, '11111111-1111-4000-a000-000000000001'),

  ('22222222-2222-4000-a000-000000000002', '33333333-3333-4000-a000-000000000001', '44444444-4444-4000-a000-000000000001', 92, 100, '11111111-1111-4000-a000-000000000001'),
  ('22222222-2222-4000-a000-000000000002', '33333333-3333-4000-a000-000000000002', '44444444-4444-4000-a000-000000000001', 95, 100, '11111111-1111-4000-a000-000000000001'),

  ('22222222-2222-4000-a000-000000000003', '33333333-3333-4000-a000-000000000001', '44444444-4444-4000-a000-000000000001', 78, 100, '11111111-1111-4000-a000-000000000001'),
  ('22222222-2222-4000-a000-000000000003', '33333333-3333-4000-a000-000000000002', '44444444-4444-4000-a000-000000000001', 72, 100, '11111111-1111-4000-a000-000000000001'),

  ('22222222-2222-4000-a000-000000000004', '33333333-3333-4000-a000-000000000001', '44444444-4444-4000-a000-000000000001', 64, 100, '11111111-1111-4000-a000-000000000001'),
  ('22222222-2222-4000-a000-000000000004', '33333333-3333-4000-a000-000000000002', '44444444-4444-4000-a000-000000000001', 70, 100, '11111111-1111-4000-a000-000000000001'),

  ('22222222-2222-4000-a000-000000000005', '33333333-3333-4000-a000-000000000001', '44444444-4444-4000-a000-000000000001', 90, 100, '11111111-1111-4000-a000-000000000001'),
  ('22222222-2222-4000-a000-000000000005', '33333333-3333-4000-a000-000000000002', '44444444-4444-4000-a000-000000000001', 94, 100, '11111111-1111-4000-a000-000000000001')
ON CONFLICT (student_id, subject_id, term_id) DO NOTHING;

COMMIT;
EOF

echo "──────────────────────────────────────────────"
echo "✅ Demo class 'Grade 10-A' successfully seeded!"
echo "   - Class: Grade 10-A (Academic Year 2026)"
echo "   - Teachers: 2 (Form Teacher: Mr. Kavinda Perera, Subject Teacher: Ms. Nimali Fernando)"
echo "   - Students: 5 enrolled"
echo "   - Parents/Guardians: 5 linked"
echo "   - Attendance: 3 Days (Teachers & Students)"
echo "   - Term Marks: Term 1 marks for Mathematics & Science"
echo "   - Clean Up: Run 'make unseed' or './scripts/unseed_demo_class.sh' to remove demo data."
echo "──────────────────────────────────────────────"
