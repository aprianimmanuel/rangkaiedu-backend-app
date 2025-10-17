package repositories

import (
	"database/sql"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
)

// MaterialRepository defines the interface for material data operations
type MaterialRepository interface {
	FindByID(id string) (*models.Material, error)
	FindByClassID(classID string) ([]*models.Material, error)
	FindBySubjectID(subjectID string) ([]*models.Material, error)
	FindByTeacherID(teacherID string) ([]*models.Material, error)
	Create(material *models.Material) error
	Update(material *models.Material) error
	Delete(id string) error
}

// materialRepository implements MaterialRepository interface
type materialRepository struct {
	db *sql.DB
}

// NewMaterialRepository creates a new material repository
func NewMaterialRepository(db *sql.DB) MaterialRepository {
	return &materialRepository{
		db: db,
	}
}

// FindByID finds a material by ID
func (r *materialRepository) FindByID(id string) (*models.Material, error) {
	query := `
		SELECT id, class_id, subject_id, teacher_id, title, description, 
		       file_name, file_path, file_type, file_size, visibility, created_at, updated_at
		FROM materials WHERE id = $1
	`
	
	var material models.Material
	err := r.db.QueryRow(query, id).Scan(
		&material.ID, &material.ClassID, &material.SubjectID, &material.TeacherID,
		&material.Title, &material.Description, &material.FileName, &material.FilePath,
		&material.FileType, &material.FileSize, &material.Visibility, &material.CreatedAt, &material.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &material, nil
}

// FindByClassID finds all materials by class ID
func (r *materialRepository) FindByClassID(classID string) ([]*models.Material, error) {
	query := `
		SELECT id, class_id, subject_id, teacher_id, title, description, 
		       file_name, file_path, file_type, file_size, visibility, created_at, updated_at
		FROM materials WHERE class_id = $1
	`
	
	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var materials []*models.Material
	for rows.Next() {
		var material models.Material
		err := rows.Scan(
			&material.ID, &material.ClassID, &material.SubjectID, &material.TeacherID,
			&material.Title, &material.Description, &material.FileName, &material.FilePath,
			&material.FileType, &material.FileSize, &material.Visibility, &material.CreatedAt, &material.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		materials = append(materials, &material)
	}
	
	return materials, nil
}

// FindBySubjectID finds all materials by subject ID
func (r *materialRepository) FindBySubjectID(subjectID string) ([]*models.Material, error) {
	query := `
		SELECT id, class_id, subject_id, teacher_id, title, description, 
		       file_name, file_path, file_type, file_size, visibility, created_at, updated_at
		FROM materials WHERE subject_id = $1
	`
	
	rows, err := r.db.Query(query, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var materials []*models.Material
	for rows.Next() {
		var material models.Material
		err := rows.Scan(
			&material.ID, &material.ClassID, &material.SubjectID, &material.TeacherID,
			&material.Title, &material.Description, &material.FileName, &material.FilePath,
			&material.FileType, &material.FileSize, &material.Visibility, &material.CreatedAt, &material.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		materials = append(materials, &material)
	}
	
	return materials, nil
}

// FindByTeacherID finds all materials by teacher ID
func (r *materialRepository) FindByTeacherID(teacherID string) ([]*models.Material, error) {
	query := `
		SELECT id, class_id, subject_id, teacher_id, title, description, 
		       file_name, file_path, file_type, file_size, visibility, created_at, updated_at
		FROM materials WHERE teacher_id = $1
	`
	
	rows, err := r.db.Query(query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var materials []*models.Material
	for rows.Next() {
		var material models.Material
		err := rows.Scan(
			&material.ID, &material.ClassID, &material.SubjectID, &material.TeacherID,
			&material.Title, &material.Description, &material.FileName, &material.FilePath,
			&material.FileType, &material.FileSize, &material.Visibility, &material.CreatedAt, &material.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		materials = append(materials, &material)
	}
	
	return materials, nil
}

// Create creates a new material
func (r *materialRepository) Create(material *models.Material) error {
	query := `
		INSERT INTO materials (class_id, subject_id, teacher_id, title, description, 
		                       file_name, file_path, file_type, file_size, visibility)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRow(query, material.ClassID, material.SubjectID, material.TeacherID,
		material.Title, material.Description, material.FileName, material.FilePath,
		material.FileType, material.FileSize, material.Visibility).Scan(
		&material.ID, &material.CreatedAt, &material.UpdatedAt,
	)
	
	return err
}

// Update updates an existing material
func (r *materialRepository) Update(material *models.Material) error {
	query := `
		UPDATE materials 
		SET class_id = $1, subject_id = $2, teacher_id = $3, title = $4, description = $5, 
		    file_name = $6, file_path = $7, file_type = $8, file_size = $9, visibility = $10, 
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $11
	`
	
	_, err := r.db.Exec(query, material.ClassID, material.SubjectID, material.TeacherID,
		material.Title, material.Description, material.FileName, material.FilePath,
		material.FileType, material.FileSize, material.Visibility, material.ID)
	
	return err
}

// Delete deletes a material by ID
func (r *materialRepository) Delete(id string) error {
	query := `DELETE FROM materials WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}