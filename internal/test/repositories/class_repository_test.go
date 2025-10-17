
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

func TestClassRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewClassRepository(db)

	// Test successful class retrieval
	expectedClass := &models.Class{
		ID:           "class123",
		SchoolID:     "school123",
		Name:         "Mathematics",
		GradeLevel:   10,
		AcademicYear: "2023/2024",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "school_id", "name", "grade_level", "academic_year", "created_at", "updated_at",
	}).AddRow(
		expectedClass.ID, expectedClass.SchoolID, expectedClass.Name, expectedClass.GradeLevel,
		expectedClass.AcademicYear, expectedClass.CreatedAt, expectedClass.UpdatedAt,
	)

	mock.ExpectQuery("SELECT id, school_id, name, grade_level, academic_year, created_at, updated_at FROM classes WHERE id = \\$1").
		WithArgs(expectedClass.ID).
		WillReturnRows(rows)

	class, err := repo.FindByID(expectedClass.ID)
	require.NoError(t, err)
	assert.NotNil(t, class)
	assert.Equal(t, expectedClass.ID, class.ID)
	assert.Equal(t, expectedClass.Name, class.Name)
	assert.Equal(t, expectedClass.GradeLevel, class.GradeLevel)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassRepository_FindByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewClassRepository(db)

	classID := "nonexistent123"

	mock.ExpectQuery("SELECT id, school_id, name, grade_level, academic_year, created_at, updated_at FROM classes WHERE id = \\$1").
		WithArgs(classID).
		WillReturnError(sql.ErrNoRows)

	class, err := repo.FindByID(classID)
	assert.Error(t, err)
	assert.Nil(t, class)
	assert.Equal(t, sql.ErrNoRows, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassRepository_FindBySchoolID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewClassRepository(db)

	// Test successful classes retrieval by school ID
	schoolID := "school123"
	expectedClasses := []*models.Class{
		{
			ID:           "class1",
			SchoolID:     schoolID,
			Name:         "Mathematics",
			GradeLevel:   10,
			AcademicYear: "2023/2024",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           "class2",
			SchoolID:     schoolID,
			Name:         "Physics",
			GradeLevel:   11,
			AcademicYear: "2023/2024",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "school_id", "name", "grade_level", "academic_year", "created_at", "updated_at",
	})

	for _, class := range expectedClasses {
		rows.AddRow(
			class.ID, class.SchoolID, class.Name, class.GradeLevel,
			class.AcademicYear, class.CreatedAt, class.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT id, school_id, name, grade_level, academic_year, created_at, updated_at FROM classes WHERE school_id = \\$1").
		WithArgs(schoolID).
		WillReturnRows(rows)

	classes, err := repo.FindBySchoolID(schoolID)
	require.NoError(t, err)
	assert.NotNil(t, classes)
	assert.Len(t, classes, 2)
	assert.Equal(t, expectedClasses[0].ID, classes[0].ID)
	assert.Equal(t, expectedClasses[1].ID, classes[1].ID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassRepository_FindBySchoolID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewClassRepository(db)

	// Test empty classes retrieval by school ID
	schoolID := "school456"

	rows := sqlmock.NewRows([]string{
		"id", "school_id", "name", "grade_level", "academic_year", "created_at", "updated_at",
	})

	mock.ExpectQuery("SELECT id, school_id, name, grade_level, academic_year, created_at, updated_at FROM classes WHERE school_id = \\$1").
		WithArgs(schoolID).
		WillReturnRows(rows)

	classes, err := repo.FindBySchoolID(schoolID)
	require.NoError(t, err)
	assert.Empty(t, classes) // Check for empty slice instead of non-nil
	assert.Len(t, classes, 0)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewClassRepository(db)

	// Test successful class creation
	class := &models.Class{
		SchoolID:     "school123",
		Name:         "Chemistry",
		GradeLevel:   12,
		AcademicYear: "2023/2024",
	}

	classID := "class789"
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(classID, createdAt, updatedAt)

	mock.ExpectQuery("INSERT INTO classes \\(school_id, name, grade_level, academic_year\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id, created_at, updated_at").
		WithArgs(class.SchoolID, class.Name, class.GradeLevel, class.AcademicYear).
		WillReturnRows(rows)

	err = repo.Create(class)
	require.NoError(t, err)
	assert.Equal(t, classID, class.ID)
	assert.Equal(t, createdAt, class.CreatedAt)
	assert.Equal(t, updatedAt, class.UpdatedAt)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewClassRepository(db)

	// Test successful class update
	class := &models.Class{
		ID:           "class123",
		SchoolID:     "school456",
		Name:         "Updated Mathematics",
		GradeLevel:   11,
		AcademicYear: "2024/2025",
	}

	mock.ExpectExec("UPDATE classes SET school_id = \\$1, name = \\$2, grade_level = \\$3, academic_year = \\$4, updated_at = CURRENT_TIMESTAMP WHERE id = \\$5").
		WithArgs(class.SchoolID, class.Name, class.GradeLevel, class.AcademicYear, class.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(class)
	require.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewClassRepository(db)

	// Test successful class deletion
	classID := "class123"

	mock.ExpectExec("DELETE FROM classes WHERE id = \\$1").
		WithArgs(classID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(classID)
	require.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
