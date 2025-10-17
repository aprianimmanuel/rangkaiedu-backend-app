
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

func TestSubjectRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewSubjectRepository(db)

	// Test successful subject retrieval
	expectedSubject := &models.Subject{
		ID:          "subject123",
		SchoolID:    "school123",
		Name:        "Mathematics",
		Code:        "MATH101",
		Description: "Basic Mathematics",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "school_id", "name", "code", "description", "created_at", "updated_at",
	}).AddRow(
		expectedSubject.ID, expectedSubject.SchoolID, expectedSubject.Name, expectedSubject.Code,
		expectedSubject.Description, expectedSubject.CreatedAt, expectedSubject.UpdatedAt,
	)

	mock.ExpectQuery("SELECT id, school_id, name, code, description, created_at, updated_at FROM subjects WHERE id = \\$1").
		WithArgs(expectedSubject.ID).
		WillReturnRows(rows)

	subject, err := repo.FindByID(expectedSubject.ID)
	require.NoError(t, err)
	assert.NotNil(t, subject)
	assert.Equal(t, expectedSubject.ID, subject.ID)
	assert.Equal(t, expectedSubject.Name, subject.Name)
	assert.Equal(t, expectedSubject.Code, subject.Code)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectRepository_FindByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewSubjectRepository(db)

	subjectID := "nonexistent123"

	mock.ExpectQuery("SELECT id, school_id, name, code, description, created_at, updated_at FROM subjects WHERE id = \\$1").
		WithArgs(subjectID).
		WillReturnError(sql.ErrNoRows)

	subject, err := repo.FindByID(subjectID)
	assert.Error(t, err)
	assert.Nil(t, subject)
	assert.Equal(t, sql.ErrNoRows, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectRepository_FindBySchoolID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewSubjectRepository(db)

	// Test successful subjects retrieval by school ID
	schoolID := "school123"
	expectedSubjects := []*models.Subject{
		{
			ID:          "subject1",
			SchoolID:    schoolID,
			Name:        "Mathematics",
			Code:        "MATH101",
			Description: "Basic Mathematics",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "subject2",
			SchoolID:    schoolID,
			Name:        "Physics",
			Code:        "PHYS101",
		

			Description: "Basic Physics",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "school_id", "name", "code", "description", "created_at", "updated_at",
	})

	for _, subject := range expectedSubjects {
		rows.AddRow(
			subject.ID, subject.SchoolID, subject.Name, subject.Code,
			subject.Description, subject.CreatedAt, subject.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT id, school_id, name, code, description, created_at, updated_at FROM subjects WHERE school_id = \\$1").
		WithArgs(schoolID).
		WillReturnRows(rows)

	subjects, err := repo.FindBySchoolID(schoolID)
	require.NoError(t, err)
	assert.NotNil(t, subjects)
	assert.Len(t, subjects, 2)
	assert.Equal(t, expectedSubjects[0].ID, subjects[0].ID)
	assert.Equal(t, expectedSubjects[1].ID, subjects[1].ID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectRepository_FindBySchoolID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewSubjectRepository(db)

	// Test empty subjects retrieval by school ID
	schoolID := "school456"

	rows := sqlmock.NewRows([]string{
		"id", "school_id", "name", "code", "description", "created_at", "updated_at",
	})

	mock.ExpectQuery("SELECT id, school_id, name, code, description, created_at, updated_at FROM subjects WHERE school_id = \\$1").
		WithArgs(schoolID).
		WillReturnRows(rows)

	subjects, err := repo.FindBySchoolID(schoolID)
	require.NoError(t, err)
	assert.Empty(t, subjects)
	assert.Len(t, subjects, 0)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewSubjectRepository(db)

	// Test successful subject creation
	subject := &models.Subject{
		SchoolID:    "school123",
		Name:        "Chemistry",
		Code:        "CHEM101",
		Description: "Basic Chemistry",
	}

	subjectID := "subject789"
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(subjectID, createdAt, updatedAt)

	mock.ExpectQuery("INSERT INTO subjects \\(school_id, name, code, description\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id, created_at, updated_at").
		WithArgs(subject.SchoolID, subject.Name, subject.Code, subject.Description).
		WillReturnRows(rows)

	err = repo.Create(subject)
	require.NoError(t, err)
	assert.Equal(t, subjectID, subject.ID)
	assert.Equal(t, createdAt, subject.CreatedAt)
	assert.Equal(t, updatedAt, subject.UpdatedAt)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewSubjectRepository(db)

	// Test successful subject update
	subject := &models.Subject{
		ID:          "subject123",
		SchoolID:    "school456",
		Name:        "Updated Mathematics",
		Code:        "MATH102",
		Description: "Advanced Mathematics",
	}

	mock.ExpectExec("UPDATE subjects SET school_id = \\$1, name = \\$2, code = \\$3, description = \\$4, updated_at = CURRENT_TIMESTAMP WHERE id = \\$5").
		WithArgs(subject.SchoolID, subject.Name, subject.Code, subject.Description, subject.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(subject)
	require.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewSubjectRepository(db)

	// Test successful subject deletion
	subjectID := "subject123"

	mock.ExpectExec("DELETE FROM subjects WHERE id = \\$1").
		WithArgs(subjectID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(subjectID)
	require.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}