package models

import (
	"time"
)

// Class represents a class in a school
type Class struct {
	ID           string    `json:"id" db:"id"`
	SchoolID     string    `json:"school_id" db:"school_id"`
	Name         string    `json:"name" db:"name"`
	GradeLevel   int       `json:"grade_level" db:"grade_level"`
	AcademicYear string    `json:"academic_year" db:"academic_year"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}