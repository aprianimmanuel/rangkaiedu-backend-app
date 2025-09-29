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