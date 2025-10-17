
package repositories

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/repositories"
)

func TestMaterialRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMaterialRepository(db)

	// Test successful material retrieval
	expectedMaterial := &models.Material{
		ID:          "material123",
		ClassID:     "class123",
		SubjectID:   "subject123",
		TeacherID:   "teacher123",
		Title:       "Mathematics Chapter 1",
		Description: "Introduction to Algebra",
		FileName:    "chapter1.pdf",
		FilePath:    "/materials/chapter1.pdf",
		FileType:    "application/pdf",
		FileSize:    1024,
		Visibility:  "public",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "class_id", "subject_id", "teacher_id", "title", "description",
		"file_name", "file_path", "file_type", "file_size", "visibility", "created_at", "updated_at",
	}).AddRow(
		expectedMaterial.ID, expectedMaterial.ClassID, expectedMaterial.SubjectID, expectedMaterial.TeacherID,
		expectedMaterial.Title, expectedMaterial.Description, expectedMaterial.FileName, expectedMaterial.FilePath,
		expectedMaterial.FileType, expectedMaterial.FileSize, expectedMaterial.Visibility,
		expectedMaterial.CreatedAt, expectedMaterial.UpdatedAt,
	)

	mock.ExpectQuery("SELECT id, class_id, subject_id, teacher_id, title, description, file_name, file_path, file_type, file_size, visibility, created_at, updated_at FROM materials WHERE id = \\$1").
		WithArgs(expectedMaterial.ID).
		WillReturnRows(rows)

	material, err := repo.FindByID(expectedMaterial.ID)
	require.NoError(t, err)
	assert.NotNil(t, material)
	assert.Equal(t, expectedMaterial.ID, material.ID)
	assert.Equal(t, expectedMaterial.Title, material.Title)
	assert.Equal(t, expectedMaterial.Description, material.Description)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMaterialRepository_FindByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMaterialRepository(db)

	materialID := "nonexistent123"

	mock.ExpectQuery("SELECT id, class_id, subject_id, teacher_id, title, description, file_name, file_path, file_type, file_size, visibility, created_at, updated_at FROM materials WHERE id = \\$1").
		WithArgs(materialID).
		WillReturnError(sql.ErrNoRows)

	material, err := repo.FindByID(materialID)
	assert.Error(t, err)
	assert.Nil(t, material)
	assert.Equal(t, sql.ErrNoRows, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMaterialRepository_FindByClassID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMaterialRepository(db)

	// Test successful materials retrieval by class ID
	classID := "class123"
	expectedMaterials := []*models.Material{
		{
			ID:          "material1",
			ClassID:     classID,
			SubjectID:   "subject1",
			TeacherID:   "teacher1",
			Title:       "Mathematics Chapter 1",
			Description: "Introduction to Algebra",
			FileName:    "chapter1.pdf",
			FilePath:    "/materials/chapter1.pdf",
			FileType:    "application/pdf",
			FileSize:    1024,
			Visibility:  "public",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "material2",
			ClassID:     classID,
			SubjectID:   "subject1",
			TeacherID:   "teacher1",
			Title:       "Mathematics Chapter 2",
			Description: "Linear Equations",
			FileName:    "chapter2.pdf",
			FilePath:    "/materials/chapter2.pdf",
			FileType:    "application/pdf",
			FileSize:    2048,
			Visibility:  "public",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "class_id", "subject_id", "teacher_id", "title", "description",
		"file_name", "file_path", "file_type", "file_size", "visibility", "created_at", "updated_at",
	})

	for _, material := range expectedMaterials {
		rows.AddRow(
			material.ID, material.ClassID, material.SubjectID, material.TeacherID,
			material.Title, material.Description, material.FileName, material.FilePath,
			material.FileType, material.FileSize, material.Visibility,
			material.CreatedAt, material.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT id, class_id, subject_id, teacher_id, title, description, file_name, file_path, file_type, file_size, visibility, created_at, updated_at FROM materials WHERE class_id = \\$1").
		WithArgs(classID).
		WillReturnRows(rows)

	materials, err := repo.FindByClassID(classID)
	require.NoError(t, err)
	assert.NotNil(t, materials)
	assert.Len(t, materials, 2)
	assert.Equal(t, expectedMaterials[0].ID, materials[0].ID)
	assert.Equal(t, expectedMaterials[1].ID, materials[1].ID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMaterialRepository_FindByClassID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMaterialRepository(db)

	// Test empty materials retrieval by class ID
	classID := "class456"

	rows := sqlmock.NewRows([]string{
		"id", "class_id", "subject_id", "teacher_id", "title", "description",
		"file_name", "file_path", "file_type", "file_size", "visibility", "created_at", "updated_at",
	})

	mock.ExpectQuery("SELECT id, class_id, subject_id, teacher_id, title, description, file_name, file_path, file_type, file_size, visibility, created_at, updated_at FROM materials WHERE class_id = \\$1").
		WithArgs(classID).
		WillReturnRows(rows)

	materials, err := repo.FindByClassID(classID)
	require.NoError(t, err)
	assert.Empty(t, materials)
	assert.Len(t, materials, 0)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMaterialRepository_FindBySubjectID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMaterialRepository(db)

	// Test successful materials retrieval by subject ID
	subjectID := "subject123"
	expectedMaterials := []*models.Material{
		{
			ID:          "material1",
			ClassID:     "class1",
			SubjectID:   subjectID,
			TeacherID:   "teacher1",
			Title:       "Mathematics Chapter 1",
			Description: "Introduction to Algebra",
			FileName:    "chapter1.pdf",
			FilePath:    "/materials/chapter1.pdf",
			FileType:    "application/pdf",
			FileSize:    1024,
			Visibility:  "public",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "material2",
			ClassID:     "class2",
			SubjectID:   subjectID,
			TeacherID:   "teacher2",
			Title:       "Advanced Mathematics",
			Description: "Calculus Introduction",
			FileName:    "calculus.pdf",
			FilePath:    "/materials/calculus.pdf",
			FileType:    "application/pdf",
			FileSize:    4096,
			Visibility:  "public",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "class_id", "subject_id", "teacher_id", "title", "description",
		"file_name", "file_path", "file_type", "file_size", "visibility", "created_at", "updated_at",
	})

	for _, material := range expectedMaterials {
		rows.AddRow(
			material.ID, material.ClassID, material.SubjectID, material.TeacherID,
			material.Title, material.Description, material.FileName, material.FilePath,
			material.FileType, material.FileSize, material.Visibility,
			material.CreatedAt, material.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT id, class_id, subject_id, teacher_id, title, description, file_name, file_path, file_type, file_size, visibility, created_at, updated_at FROM materials WHERE subject_id = \\$1").
		WithArgs(subjectID).
		WillReturnRows(rows)

	materials, err := repo.FindBySubjectID(subjectID)
	require.NoError(t, err)
	assert.NotNil(t, materials)
	assert.Len(t, materials, 2)
	assert.Equal(t, expectedMaterials[0].ID, materials[0].ID)
	assert.Equal(t, expectedMaterials[1].ID, materials[1].ID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMaterialRepository_FindBySubjectID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMaterialRepository(db)

	// Test empty materials retrieval by subject ID
	subjectID := "subject456"

	rows := sqlmock.NewRows([]string{
		"id", "class_id", "subject_id", "teacher_id", "title", "description",
		"file_name", "file_path", "file_type", "file_size", "visibility", "created_at", "updated_at",
	})

	mock.ExpectQuery("SELECT id, class_id, subject_id, teacher_id, title, description, file_name, file_path, file_type, file_size, visibility, created_at, updated_at FROM materials WHERE subject_id = \\$1").
		WithArgs(subjectID).
		WillReturnRows(rows)

	materials, err := repo.FindBySubjectID(subjectID)
	require.NoError(t, err)
	assert.Empty(t, materials)
	assert.Len(t, materials, 0)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMaterialRepository_FindByTeacherID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMaterialRepository(db)

	// Test successful materials retrieval by teacher ID
	teacherID := "teacher123"
	expectedMaterials := []*models.Material{
		{
			ID:          "material1",
			ClassID:     "class1",
			SubjectID:   "subject1",
			TeacherID:   teacherID,
			Title:       "Mathematics Chapter 1",
			Description: "Introduction to Algebra",
			FileName:    "chapter1.pdf",
			FilePath:    "/materials/chapter1.pdf",
			FileType:    "application/pdf",
			FileSize:    1024,
			Visibility:  "public",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "material2",
			ClassID:     "class2",
			SubjectID:   "subject2",
			TeacherID:   teacherID,
			Title:       "Physics Chapter 1",
			Description: "Introduction to Mechanics",
			FileName:    "mechanics.pdf",
			FilePath:    "/materials/mechanics.pdf",
			FileType:    "application/pdf",
			FileSize:    2048,
			Visibility:  "public",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "class_id", "subject_id", "teacher_id", "title", "description",
		"file_name", "file_path", "file_type", "file_size", "visibility", "created_at", "updated_at",
	})

	for _, material := range expectedMaterials {
		rows.AddRow(
			material.ID, material.ClassID, material.SubjectID, material.TeacherID,
			material.Title, material.Description, material.FileName, material.FilePath,
			material.FileType, material.FileSize, material.Visibility,
			material.CreatedAt, material.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT id, class_id, subject_id, teacher_id, title, description, file_name, file_path, file_type, file_size, visibility, created_at, updated_at FROM materials WHERE teacher_id = \\$1").
		WithArgs(teacherID).
		WillReturnRows(rows)

	materials, err := repo.FindByTeacherID(teacherID)
	require.NoError(t, err)
	assert.NotNil(t, materials)
	assert.Len(t, materials, 2)
	assert.Equal(t, expectedMaterials[0].ID, materials[0].ID)
	assert.Equal(t, expectedMaterials[1].ID, materials[1].ID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMaterialRepository_FindByTeacherID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMaterialRepository(db)

	// Test empty materials retrieval by teacher ID
	teacherID := "teacher456"

	rows := sqlmock.NewRows([]string{
		"id", "class_id", "subject_id", "teacher_id", "title", "description",
		"file_name", "file_path", "file_type", "file_size", "visibility", "created_at", "updated_at",
	})

	mock.ExpectQuery("SELECT id, class_id, subject_id, teacher_id, title, description, file_name, file_path, file_type, file_size, visibility, created_at, updated_at FROM materials WHERE teacher_id = \\$1").
		WithArgs(teacherID).
		WillReturnRows(rows)

	materials, err := repo.FindByTeacherID(teacherID)
	require.NoError(t, err)
	assert.Empty(t, materials)
	assert.Len(t, materials, 0)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMaterialRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMaterialRepository(db)

	// Test successful material creation
	material := &models.Material{
		ClassID:     "class123",
		SubjectID:   "subject123",
		TeacherID:   "teacher123",
		Title:       "Chemistry Chapter 1",
		Description: "Introduction to Chemistry",
		FileName:    "chemistry.pdf",
		FilePath:    "/materials/chemistry.pdf",
		FileType:    "application/pdf",
		FileSize:    3072,
		Visibility:  "public",
	}

	materialID := "material789"
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(materialID, createdAt, updatedAt)

	mock.ExpectQuery("INSERT INTO materials \\(class_id, subject_id, teacher_id, title, description, file_name, file_path, file_type, file_size, visibility\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7, \\$8, \\$9, \\$10\\) RETURNING id, created_at, updated_at").
		WithArgs(material.ClassID, material.SubjectID, material.TeacherID, material.Title, material.Description,
			material.FileName, material.FilePath, material.FileType, material.FileSize, material.Visibility).
		WillReturnRows(rows)

	err = repo.Create(material)
	require.NoError(t, err)
	assert.Equal(t, materialID, material.ID)
	assert.Equal(t, createdAt, material.CreatedAt)
	assert.Equal(t, updatedAt, material.UpdatedAt)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMaterialRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMaterialRepository(db)

	// Test successful material update
	material := &models.Material{
		ID:          "material123",
		ClassID:     "class456",
		SubjectID:   "subject456",
		TeacherID:   "teacher456",
		Title:       "Updated Mathematics Chapter 1",
		Description: "Advanced Algebra",
		FileName:    "updated_chapter1.pdf",
		FilePath:    "/materials/updated_chapter1.pdf",
		FileType:    "application/pdf",
		FileSize:    1536,
		Visibility:  "private",
	}

	mock.ExpectExec("UPDATE materials SET class_id = \\$1, subject_id = \\$2, teacher_id = \\$3, title = \\$4, description = \\$5, file_name = \\$6, file_path = \\$7, file_type = \\$8, file_size = \\$9, visibility = \\$10, updated_at = CURRENT_TIMESTAMP WHERE id = \\$11").
		WithArgs(material.ClassID, material.SubjectID, material.TeacherID, material.Title, material.Description,
			material.FileName, material.FilePath, material.FileType, material.FileSize, material.Visibility, material.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(material)
	require.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMaterialRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMaterialRepository(db)

	// Test successful material deletion
	materialID := "material123"

	mock.ExpectExec("DELETE FROM materials WHERE id = \\$1").
		WithArgs(materialID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(materialID)
	require.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}