-- +goose Up
-- Migration: Seed initial development data for Rangkai Edu
-- Description: Inserts sample data for testing relationships and core entities
-- Date: 2025-09-28

-- Insert sample users
INSERT INTO users (id, name, email, phone, role, created_at) VALUES
(gen_random_uuid(), 'Admin User', 'admin@rangkaiedu.com', '+6281234567890', 'admin', CURRENT_TIMESTAMP),
(gen_random_uuid(), 'Teacher One', 'teacher1@schoola.com', '+6281234567891', 'teacher', CURRENT_TIMESTAMP),
(gen_random_uuid(), 'Teacher Two', 'teacher2@schoolb.com', '+6281234567892', 'teacher', CURRENT_TIMESTAMP),
(gen_random_uuid(), 'Student A', 'studenta@schoola.com', '+6281234567893', 'student', CURRENT_TIMESTAMP),
(gen_random_uuid(), 'Student B', 'studentb@schoola.com', '+6281234567894', 'student', CURRENT_TIMESTAMP),
(gen_random_uuid(), 'Student C', 'studentc@schoolb.com', '+6281234567895', 'student', CURRENT_TIMESTAMP),
(gen_random_uuid(), 'Parent One', 'parent1@schoola.com', '+6281234567896', 'parent', CURRENT_TIMESTAMP),
(gen_random_uuid(), 'Parent Two', 'parent2@schoolb.com', '+6281234567897', 'parent', CURRENT_TIMESTAMP);

-- Insert sample schools
INSERT INTO schools (id, name, address, phone, email, created_at) VALUES
(gen_random_uuid(), 'School A', 'Jakarta, Indonesia', '+62211234567', 'info@schoola.com', CURRENT_TIMESTAMP),
(gen_random_uuid(), 'School B', 'Bandung, Indonesia', '+62229876543', 'info@schoolb.com', CURRENT_TIMESTAMP);

-- Assume school IDs: Use subqueries or hardcode if known, but for migration, insert and reference by values if needed. For simplicity, insert sequentially and use lastval or separate vars, but in plain SQL, insert and then query.
-- Better: Insert and capture IDs via RETURNING, but for migration script, use sequential inserts and reference by order.

-- To make it work, we'll insert schools first and use their IDs in subsequent inserts. But since UUIDs are random, use INSERT ... RETURNING but in script, assume separate statements.

-- For practicality in migration, insert all without FK first where possible, but since FKs, insert in order and use temp vars if PL/pgSQL, but keep simple: Insert parents/teachers/students after users, using subqueries for user_id.

-- Insert teachers (link to users with role 'teacher')
INSERT INTO teachers (id, user_id, school_id, employee_id, qualification, hire_date, created_at)
SELECT gen_random_uuid(), u.id, s1.id, 'EMP001', 'Bachelor of Education', '2023-01-01', CURRENT_TIMESTAMP
FROM users u, schools s1
WHERE u.role = 'teacher' AND u.name = 'Teacher One'
AND s1.name = 'School A'
UNION ALL
SELECT gen_random_uuid(), u.id, s2.id, 'EMP002', 'Master of Science', '2023-02-01', CURRENT_TIMESTAMP
FROM users u, schools s2
WHERE u.role = 'teacher' AND u.name = 'Teacher Two'
AND s2.name = 'School B';

-- Insert parents (link to users with role 'parent')
INSERT INTO parents (id, user_id, occupation, created_at)
SELECT gen_random_uuid(), u.id, 'Engineer', CURRENT_TIMESTAMP
FROM users u WHERE u.role = 'parent' AND u.name = 'Parent One'
UNION ALL
SELECT gen_random_uuid(), u.id, 'Doctor', CURRENT_TIMESTAMP
FROM users u WHERE u.role = 'parent' AND u.name = 'Parent Two';

-- Insert students (link to users with role 'student', assign schools, classes later)
-- First insert classes and subjects

-- Insert classes
INSERT INTO classes (id, school_id, name, grade_level, academic_year, created_at)
SELECT gen_random_uuid(), s1.id, 'Class 1A', 1, '2024-2025', CURRENT_TIMESTAMP FROM schools s1 WHERE s1.name = 'School A'
UNION ALL
SELECT gen_random_uuid(), s1.id, 'Class 2A', 2, '2024-2025', CURRENT_TIMESTAMP FROM schools s1 WHERE s1.name = 'School A'
UNION ALL
SELECT gen_random_uuid(), s2.id, 'Class 1B', 1, '2024-2025', CURRENT_TIMESTAMP FROM schools s2 WHERE s2.name = 'School B'
UNION ALL
SELECT gen_random_uuid(), s2.id, 'Class 2B', 2, '2024-2025', CURRENT_TIMESTAMP FROM schools s2 WHERE s2.name = 'School B';

-- Insert subjects
INSERT INTO subjects (id, school_id, name, code, description, created_at)
SELECT gen_random_uuid(), s1.id, 'Mathematics', 'MATH', 'Basic Math', CURRENT_TIMESTAMP FROM schools s1 WHERE s1.name = 'School A'
UNION ALL
SELECT gen_random_uuid(), s1.id, 'Science', 'SCI', 'General Science', CURRENT_TIMESTAMP FROM schools s1 WHERE s1.name = 'School A'
UNION ALL
SELECT gen_random_uuid(), s1.id, 'English', 'ENG', 'Language Arts', CURRENT_TIMESTAMP FROM schools s1 WHERE s1.name = 'School A'
UNION ALL
SELECT gen_random_uuid(), s2.id, 'Mathematics', 'MATH', 'Basic Math', CURRENT_TIMESTAMP FROM schools s2 WHERE s2.name = 'School B'
UNION ALL
SELECT gen_random_uuid(), s2.id, 'History', 'HIST', 'World History', CURRENT_TIMESTAMP FROM schools s2 WHERE s2.name = 'School B';

-- Now insert students
INSERT INTO students (id, user_id, school_id, student_id, date_of_birth, grade_level, class_id, parent_id, enrollment_date, created_at)
SELECT gen_random_uuid(), u.id, s1.id, 'STU001', '2015-05-10', 1, c1.id, NULL, CURRENT_DATE, CURRENT_TIMESTAMP
FROM users u, schools s1, classes c1
WHERE u.role = 'student' AND u.name = 'Student A'
AND s1.name = 'School A'
AND c1.name = 'Class 1A'
UNION ALL
SELECT gen_random_uuid(), u.id, s1.id, 'STU002', '2014-08-15', 2, c2.id, NULL, CURRENT_DATE, CURRENT_TIMESTAMP
FROM users u, schools s1, classes c2
WHERE u.role = 'student' AND u.name = 'Student B'
AND s1.name = 'School A'
AND c2.name = 'Class 2A'
UNION ALL
SELECT gen_random_uuid(), u.id, s2.id, 'STU003', '2015-03-20', 1, c3.id, NULL, CURRENT_DATE, CURRENT_TIMESTAMP
FROM users u, schools s2, classes c3
WHERE u.role = 'student' AND u.name = 'Student C'
AND s2.name = 'School B'
AND c3.name = 'Class 1B';

-- Insert class_teachers (assign teachers to classes/subjects)
INSERT INTO class_teachers (id, class_id, teacher_id, subject_id, academic_year, created_at)
SELECT gen_random_uuid(), c1.id, t1.id, sub1.id, '2024-2025', CURRENT_TIMESTAMP
FROM classes c1, teachers t1, subjects sub1
WHERE c1.name = 'Class 1A' AND t1.employee_id = 'EMP001' AND sub1.code = 'MATH'
UNION ALL
SELECT gen_random_uuid(), c2.id, t1.id, sub2.id, '2024-2025', CURRENT_TIMESTAMP
FROM classes c2, teachers t1, subjects sub2
WHERE c2.name = 'Class 2A' AND t1.employee_id = 'EMP001' AND sub2.code = 'SCI'
UNION ALL
SELECT gen_random_uuid(), c3.id, t2.id, sub4.id, '2024-2025', CURRENT_TIMESTAMP
FROM classes c3, teachers t2, subjects sub4
WHERE c3.name = 'Class 1B' AND t2.employee_id = 'EMP002' AND sub4.code = 'MATH'
UNION ALL
SELECT gen_random_uuid(), c4.id, t2.id, sub5.id, '2024-2025', CURRENT_TIMESTAMP
FROM classes c4, teachers t2, subjects sub5
WHERE c4.name = 'Class 2B' AND t2.employee_id = 'EMP002' AND sub5.code = 'HIST';

-- Insert student_parents (link students to parents)
INSERT INTO student_parents (id, student_id, parent_id, relationship, created_at)
SELECT gen_random_uuid(), st1.id, p1.id, 'mother', CURRENT_TIMESTAMP
FROM students st1, parents p1
WHERE st1.student_id = 'STU001' AND p1.user_id IN (SELECT id FROM users WHERE name = 'Parent One')
UNION ALL
SELECT gen_random_uuid(), st2.id, p1.id, 'father', CURRENT_TIMESTAMP
FROM students st2, parents p1
WHERE st2.student_id = 'STU002' AND p1.user_id IN (SELECT id FROM users WHERE name = 'Parent One')
UNION ALL
SELECT gen_random_uuid(), st3.id, p2.id, 'guardian', CURRENT_TIMESTAMP
FROM students st3, parents p2
WHERE st3.student_id = 'STU003' AND p2.user_id IN (SELECT id FROM users WHERE name = 'Parent Two');

-- +goose Down
-- Delete seed data in reverse order to avoid FK violations

DELETE FROM student_parents;
DELETE FROM class_teachers;
DELETE FROM students;
DELETE FROM parents;
DELETE FROM teachers;
DELETE FROM subjects;
DELETE FROM classes;
DELETE FROM schools;
DELETE FROM users WHERE role != 'admin';  -- Keep admin if needed, but delete all for clean slate
DELETE FROM users WHERE role = 'admin' AND name = 'Admin User';  -- Delete admin too