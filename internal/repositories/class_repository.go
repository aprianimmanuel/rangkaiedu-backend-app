package repositories

import (
	"database/sql"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
)

// ClassRepository defines the interface for class data operations
type ClassRepository interface {
	FindByID(id string) (*models.Class, error)
	FindBySchoolID(schoolID string) ([]*models.Class, error)
	Create(class *models.Class) error
	Update(class *models.Class) error
	Delete(id string) error
}

// classRepository implements ClassRepository interface
type classRepository struct {
	db *sql.DB
}

// NewClassRepository creates a new class repository
func NewClassRepository(db *sql.DB) ClassRepository {
	return &classRepository{
		db: db,
	}
}

// FindByID finds a class by ID
func (r *classRepository) FindByID(id string) (*models.Class, error) {
	query := `
		SELECT id, school_id, name, grade_level, academic_year, created_at, updated_at
		FROM classes WHERE id = $1
	`
	
	var class models.Class
	err := r.db.QueryRow(query, id).Scan(
		&class.ID, &class.SchoolID, &class.Name, &class.GradeLevel, 
		&class.AcademicYear, &class.CreatedAt, &class.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &class, nil
}

// FindBySchoolID finds all classes by school ID
func (r *classRepository) FindBySchoolID(schoolID string) ([]*models.Class, error) {
	query := `
		SELECT id, school_id, name, grade_level, academic_year, created_at, updated_at
		FROM classes WHERE school_id = $1
	`
	
	rows, err := r.db.Query(query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var classes []*models.Class
	for rows.Next() {
		var class models.Class
		err := rows.Scan(
			&class.ID, &class.SchoolID, &class.Name, &class.GradeLevel, 
			&class.AcademicYear, &class.CreatedAt, &class.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		classes = append(classes, &class)
	}
	
	return classes, nil
}

// Create creates a new class
func (r *classRepository) Create(class *models.Class) error {
	query := `
		INSERT INTO classes (school_id, name, grade_level, academic_year)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRow(query, class.SchoolID, class.Name, class.GradeLevel, class.AcademicYear).Scan(
		&class.ID, &class.CreatedAt, &class.UpdatedAt,
	)
	
	return err
}

// Update updates an existing class
func (r *classRepository) Update(class *models.Class) error {
	query := `
		UPDATE classes 
		SET school_id = $1, name = $2, grade_level = $3, academic_year = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`
	
	_, err := r.db.Exec(query, class.SchoolID, class.Name, class.GradeLevel, class.AcademicYear, class.ID)
	
	return err
}

// Delete deletes a class by ID
func (r *classRepository) Delete(id string) error {
	query := `DELETE FROM classes WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}