package repositories

import (
	"database/sql"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
)

// StudentEnrollmentRepository defines the interface for student enrollment data operations
type StudentEnrollmentRepository interface {
	FindByID(id string) (*models.StudentEnrollment, error)
	FindByClassID(classID string) ([]*models.StudentEnrollment, error)
	FindByStudentID(studentID string) ([]*models.StudentEnrollment, error)
	Create(enrollment *models.StudentEnrollment) error
	Update(enrollment *models.StudentEnrollment) error
	Delete(id string) error
}

// studentEnrollmentRepository implements StudentEnrollmentRepository interface
type studentEnrollmentRepository struct {
	db *sql.DB
}

// NewStudentEnrollmentRepository creates a new student enrollment repository
func NewStudentEnrollmentRepository(db *sql.DB) StudentEnrollmentRepository {
	return &studentEnrollmentRepository{
		db: db,
	}
}

// FindByID finds a student enrollment by ID
func (r *studentEnrollmentRepository) FindByID(id string) (*models.StudentEnrollment, error) {
	query := `
		SELECT id, class_id, student_id, status, created_at
		FROM student_enrollments WHERE id = $1
	`
	
	var enrollment models.StudentEnrollment
	err := r.db.QueryRow(query, id).Scan(
		&enrollment.ID, &enrollment.ClassID, &enrollment.StudentID, 
		&enrollment.Status, &enrollment.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &enrollment, nil
}

// FindByClassID finds all student enrollments by class ID
func (r *studentEnrollmentRepository) FindByClassID(classID string) ([]*models.StudentEnrollment, error) {
	query := `
		SELECT id, class_id, student_id, status, created_at
		FROM student_enrollments WHERE class_id = $1
	`
	
	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var enrollments []*models.StudentEnrollment
	for rows.Next() {
		var enrollment models.StudentEnrollment
		err := rows.Scan(
			&enrollment.ID, &enrollment.ClassID, &enrollment.StudentID, 
			&enrollment.Status, &enrollment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, &enrollment)
	}
	
	return enrollments, nil
}

// FindByStudentID finds all student enrollments by student ID
func (r *studentEnrollmentRepository) FindByStudentID(studentID string) ([]*models.StudentEnrollment, error) {
	query := `
		SELECT id, class_id, student_id, status, created_at
		FROM student_enrollments WHERE student_id = $1
	`
	
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var enrollments []*models.StudentEnrollment
	for rows.Next() {
		var enrollment models.StudentEnrollment
		err := rows.Scan(
			&enrollment.ID, &enrollment.ClassID, &enrollment.StudentID, 
			&enrollment.Status, &enrollment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, &enrollment)
	}
	
	return enrollments, nil
}

// Create creates a new student enrollment
func (r *studentEnrollmentRepository) Create(enrollment *models.StudentEnrollment) error {
	query := `
		INSERT INTO student_enrollments (class_id, student_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	
	err := r.db.QueryRow(query, enrollment.ClassID, enrollment.StudentID, enrollment.Status).Scan(
		&enrollment.ID, &enrollment.CreatedAt,
	)
	
	return err
}

// Update updates an existing student enrollment
func (r *studentEnrollmentRepository) Update(enrollment *models.StudentEnrollment) error {
	query := `
		UPDATE student_enrollments 
		SET class_id = $1, student_id = $2, status = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`
	
	_, err := r.db.Exec(query, enrollment.ClassID, enrollment.StudentID, enrollment.Status, enrollment.ID)
	
	return err
}

// Delete deletes a student enrollment by ID
func (r *studentEnrollmentRepository) Delete(id string) error {
	query := `DELETE FROM student_enrollments WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}