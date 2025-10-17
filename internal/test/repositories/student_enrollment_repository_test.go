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

func TestStudentEnrollmentRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewStudentEnrollmentRepository(db)

	// Test successful student enrollment retrieval
	expectedEnrollment := &models.StudentEnrollment{
		ID:        "enrollment123",
		ClassID:   "class123",
		StudentID: "student123",
		Status:    "active",
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "class_id", "student_id", "status", "created_at",
	}).AddRow(
		expectedEnrollment.ID, expectedEnrollment.ClassID, expectedEnrollment.StudentID,
		expectedEnrollment.Status, expectedEnrollment.CreatedAt,
	)

	mock.ExpectQuery("SELECT id, class_id, student_id, status, created_at FROM student_enrollments WHERE id = \\$1").
		WithArgs(expectedEnrollment.ID).
		WillReturnRows(rows)

	enrollment, err := repo.FindByID(expectedEnrollment.ID)
	require.NoError(t, err)
	assert.NotNil(t, enrollment)
	assert.Equal(t, expectedEnrollment.ID, enrollment.ID)
	assert.Equal(t, expectedEnrollment.ClassID, enrollment.ClassID)
	assert.Equal(t, expectedEnrollment.StudentID, enrollment.StudentID)
	assert.Equal(t, expectedEnrollment.Status, enrollment.Status)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentEnrollmentRepository_FindByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewStudentEnrollmentRepository(db)

	enrollmentID := "nonexistent123"

	mock.ExpectQuery("SELECT id, class_id, student_id, status, created_at FROM student_enrollments WHERE id = \\$1").
		WithArgs(enrollmentID).
		WillReturnError(sql.ErrNoRows)

	enrollment, err := repo.FindByID(enrollmentID)
	assert.Error(t, err)
	assert.Nil(t, enrollment)
	assert.Equal(t, sql.ErrNoRows, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentEnrollmentRepository_FindByClassID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewStudentEnrollmentRepository(db)

	// Test successful student enrollments retrieval by class ID
	classID := "class123"
	expectedEnrollments := []*models.StudentEnrollment{
		{
			ID:        "enrollment1",
			ClassID:   classID,
			StudentID: "student1",
			Status:    "active",
			CreatedAt: time.Now(),
		},
		{
			ID:        "enrollment2",
			ClassID:   classID,
			StudentID: "student2",
			Status:    "active",
			CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "class_id", "student_id", "status", "created_at",
	})

	for _, enrollment := range expectedEnrollments {
		rows.AddRow(
			enrollment.ID, enrollment.ClassID, enrollment.StudentID,
			enrollment.Status, enrollment.CreatedAt,
		)
	}

	mock.ExpectQuery("SELECT id, class_id, student_id, status, created_at FROM student_enrollments WHERE class_id = \\$1").
		WithArgs(classID).
		WillReturnRows(rows)

	enrollments, err := repo.FindByClassID(classID)
	require.NoError(t, err)
	assert.NotNil(t, enrollments)
	assert.Len(t, enrollments, 2)
	assert.Equal(t, expectedEnrollments[0].ID, enrollments[0].ID)
	assert.Equal(t, expectedEnrollments[1].ID, enrollments[1].ID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentEnrollmentRepository_FindByClassID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewStudentEnrollmentRepository(db)

	// Test empty student enrollments retrieval by class ID
	classID := "class456"

	rows := sqlmock.NewRows([]string{
		"id", "class_id", "student_id", "status", "created_at",
	})

	mock.ExpectQuery("SELECT id, class_id, student_id, status, created_at FROM student_enrollments WHERE class_id = \\$1").
		WithArgs(classID).
		WillReturnRows(rows)

	enrollments, err := repo.FindByClassID(classID)
	require.NoError(t, err)
	assert.Empty(t, enrollments)
	assert.Len(t, enrollments, 0)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentEnrollmentRepository_FindByStudentID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewStudentEnrollmentRepository(db)

	// Test successful student enrollments retrieval by student ID
	studentID := "student123"
	expectedEnrollments := []*models.StudentEnrollment{
		{
			ID:        "enrollment1",
			ClassID:   "class1",
			StudentID: studentID,
			Status:    "active",
			CreatedAt: time.Now(),
		},
		{
			ID:        "enrollment2",
			ClassID:   "class2",
			StudentID: studentID,
			Status:    "active",
			CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "class_id", "student_id", "status", "created_at",
	})

	for _, enrollment := range expectedEnrollments {
		rows.AddRow(
			enrollment.ID, enrollment.ClassID, enrollment.StudentID,
			enrollment.Status, enrollment.CreatedAt,
		)
	}

	mock.ExpectQuery("SELECT id, class_id, student_id, status, created_at FROM student_enrollments WHERE student_id = \\$1").
		WithArgs(studentID).
		WillReturnRows(rows)

	enrollments, err := repo.FindByStudentID(studentID)
	require.NoError(t, err)
	assert.NotNil(t, enrollments)
	assert.Len(t, enrollments, 2)
	assert.Equal(t, expectedEnrollments[0].ID, enrollments[0].ID)
	assert.Equal(t, expectedEnrollments[1].ID, enrollments[1].ID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentEnrollmentRepository_FindByStudentID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewStudentEnrollmentRepository(db)

	// Test empty student enrollments retrieval by student ID
	studentID := "student456"

	rows := sqlmock.NewRows([]string{
		"id", "class_id", "student_id", "status", "created_at",
	})

	mock.ExpectQuery("SELECT id, class_id, student_id, status, created_at FROM student_enrollments WHERE student_id = \\$1").
		WithArgs(studentID).
		WillReturnRows(rows)

	enrollments, err := repo.FindByStudentID(studentID)
	require.NoError(t, err)
	assert.Empty(t, enrollments)
	assert.Len(t, enrollments, 0)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentEnrollmentRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewStudentEnrollmentRepository(db)

	// Test successful student enrollment creation
	enrollment := &models.StudentEnrollment{
		ClassID:   "class123",
		StudentID: "student123",
		Status:    "active",
	}

	enrollmentID := "enrollment789"
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(enrollmentID, createdAt)

	mock.ExpectQuery("INSERT INTO student_enrollments \\(class_id, student_id, status\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id, created_at").
		WithArgs(enrollment.ClassID, enrollment.StudentID, enrollment.Status).
		WillReturnRows(rows)

	err = repo.Create(enrollment)
	require.NoError(t, err)
	assert.Equal(t, enrollmentID, enrollment.ID)
	assert.Equal(t, createdAt, enrollment.CreatedAt)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentEnrollmentRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewStudentEnrollmentRepository(db)

	// Test successful student enrollment update
	enrollment := &models.StudentEnrollment{
		ID:        "enrollment123",
		ClassID:   "class456",
		StudentID: "student456",
		Status:    "inactive",
	}

	mock.ExpectExec("UPDATE student_enrollments SET class_id = \\$1, student_id = \\$2, status = \\$3, updated_at = CURRENT_TIMESTAMP WHERE id = \\$4").
		WithArgs(enrollment.ClassID, enrollment.StudentID, enrollment.Status, enrollment.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(enrollment)
	require.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentEnrollmentRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewStudentEnrollmentRepository(db)

	// Test successful student enrollment deletion
	enrollmentID := "enrollment123"

	mock.ExpectExec("DELETE FROM student_enrollments WHERE id = \\$1").
		WithArgs(enrollmentID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(enrollmentID)
	require.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}