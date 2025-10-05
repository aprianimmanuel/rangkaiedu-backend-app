package models

import (
	"time"
)

// StudentEnrollment represents a student's enrollment in a class
type StudentEnrollment struct {
	ID        string    `json:"id" db:"id"`
	ClassID   string    `json:"class_id" db:"class_id"`
	StudentID string    `json:"student_id" db:"student_id"`
	Status    string    `json:"status" db:"status"` // active, inactive, graduated
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}