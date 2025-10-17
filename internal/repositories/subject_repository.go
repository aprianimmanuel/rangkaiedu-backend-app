package repositories

import (
	"database/sql"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
)

// SubjectRepository defines the interface for subject data operations
type SubjectRepository interface {
	FindByID(id string) (*models.Subject, error)
	FindBySchoolID(schoolID string) ([]*models.Subject, error)
	Create(subject *models.Subject) error
	Update(subject *models.Subject) error
	Delete(id string) error
}

// subjectRepository implements SubjectRepository interface
type subjectRepository struct {
	db *sql.DB
}

// NewSubjectRepository creates a new subject repository
func NewSubjectRepository(db *sql.DB) SubjectRepository {
	return &subjectRepository{
		db: db,
	}
}

// FindByID finds a subject by ID
func (r *subjectRepository) FindByID(id string) (*models.Subject, error) {
	query := `
		SELECT id, school_id, name, code, description, created_at, updated_at
		FROM subjects WHERE id = $1
	`
	
	var subject models.Subject
	err := r.db.QueryRow(query, id).Scan(
		&subject.ID, &subject.SchoolID, &subject.Name, &subject.Code, 
		&subject.Description, &subject.CreatedAt, &subject.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &subject, nil
}

// FindBySchoolID finds all subjects by school ID
func (r *subjectRepository) FindBySchoolID(schoolID string) ([]*models.Subject, error) {
	query := `
		SELECT id, school_id, name, code, description, created_at, updated_at
		FROM subjects WHERE school_id = $1
	`
	
	rows, err := r.db.Query(query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var subjects []*models.Subject
	for rows.Next() {
		var subject models.Subject
		err := rows.Scan(
			&subject.ID, &subject.SchoolID, &subject.Name, &subject.Code, 
			&subject.Description, &subject.CreatedAt, &subject.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, &subject)
	}
	
	return subjects, nil
}

// Create creates a new subject
func (r *subjectRepository) Create(subject *models.Subject) error {
	query := `
		INSERT INTO subjects (school_id, name, code, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRow(query, subject.SchoolID, subject.Name, subject.Code, subject.Description).Scan(
		&subject.ID, &subject.CreatedAt, &subject.UpdatedAt,
	)
	
	return err
}

// Update updates an existing subject
func (r *subjectRepository) Update(subject *models.Subject) error {
	query := `
		UPDATE subjects 
		SET school_id = $1, name = $2, code = $3, description = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`
	
	_, err := r.db.Exec(query, subject.SchoolID, subject.Name, subject.Code, subject.Description, subject.ID)
	
	return err
}

// Delete deletes a subject by ID
func (r *subjectRepository) Delete(id string) error {
	query := `DELETE FROM subjects WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}