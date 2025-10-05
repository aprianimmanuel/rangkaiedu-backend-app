-- +goose Up
-- Migration: Create materials and student_enrollments tables for T1.7
-- Description: Creates tables for managing teaching materials and student class enrollments
-- Date: 2025-10-05

-- Create materials table
CREATE TABLE materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(512) NOT NULL,
    file_type VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL,
    visibility VARCHAR(20) NOT NULL CHECK (visibility IN ('public', 'private', 'class_only')) DEFAULT 'class_only',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for materials updated_at
CREATE TRIGGER update_materials_updated_at 
    BEFORE UPDATE ON materials
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for materials
CREATE INDEX idx_materials_class ON materials(class_id);
CREATE INDEX idx_materials_subject ON materials(subject_id);
CREATE INDEX idx_materials_teacher ON materials(teacher_id);
CREATE INDEX idx_materials_visibility ON materials(visibility);

-- Create student_enrollments table
CREATE TABLE student_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'inactive', 'graduated')) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(class_id, student_id)
);

-- Create indexes for student_enrollments
CREATE INDEX idx_student_enrollments_class ON student_enrollments(class_id);
CREATE INDEX idx_student_enrollments_student ON student_enrollments(student_id);
CREATE INDEX idx_student_enrollments_status ON student_enrollments(status);

-- +goose Down
-- Drop tables and indexes for materials and student_enrollments

DROP TRIGGER IF EXISTS update_materials_updated_at ON materials;

DROP INDEX IF EXISTS idx_materials_class;
DROP INDEX IF EXISTS idx_materials_subject;
DROP INDEX IF EXISTS idx_materials_teacher;
DROP INDEX IF EXISTS idx_materials_visibility;
DROP INDEX IF EXISTS idx_student_enrollments_class;
DROP INDEX IF EXISTS idx_student_enrollments_student;
DROP INDEX IF EXISTS idx_student_enrollments_status;

DROP TABLE IF EXISTS student_enrollments;
DROP TABLE IF EXISTS materials;