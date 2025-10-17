package models

import (
	"time"
)

// Material represents a teaching material uploaded by a teacher
type Material struct {
	ID          string    `json:"id" db:"id"`
	ClassID     string    `json:"class_id" db:"class_id"`
	SubjectID   string    `json:"subject_id" db:"subject_id"`
	TeacherID   string    `json:"teacher_id" db:"teacher_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	FileName    string    `json:"file_name" db:"file_name"`
	FilePath    string    `json:"file_path" db:"file_path"`
	FileType    string    `json:"file_type" db:"file_type"`
	FileSize    int64     `json:"file_size" db:"file_size"`
	Visibility  string    `json:"visibility" db:"visibility"` // public, private, class_only
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}