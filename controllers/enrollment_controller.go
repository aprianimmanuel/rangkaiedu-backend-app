package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/aprianimmanuel/rangkaiedu-backend/middleware"
	"github.com/aprianimmanuel/rangkaiedu-backend/models"
)

// StudentEnrollmentRequest represents the request body for adding a student to a class
type StudentEnrollmentRequest struct {
	StudentID string `json:"student_id" binding:"required,uuid"`
}

// StudentEnrollmentStatusRequest represents the request body for updating enrollment status
type StudentEnrollmentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive graduated"`
}

// AddStudentToClass handles adding a student to a class
func AddStudentToClass(c *gin.Context) {
	// Get current user
	user, err := GetCurrentUser(c)
	if err != nil {
		SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Check if user has required role (teacher or admin)
	if user.Role != string(middleware.RoleTeacher) && user.Role != string(middleware.RoleAdmin) {
		SendErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	// Get class ID from URL parameter
	classID := c.Param("class_id")
	if classID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Class ID is required")
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(classID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid class ID format")
		return
	}

	// Bind request data
	var req StudentEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// Validate student ID format
	if _, err := uuid.Parse(req.StudentID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid student ID format")
		return
	}

	// Get database connection
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// For teachers, verify they are assigned to the class
	if user.Role == string(middleware.RoleTeacher) {
		var exists bool
		err = db.QueryRowContext(ctx,
			`SELECT EXISTS(
				SELECT 1 FROM class_teachers ct
				JOIN teachers t ON ct.teacher_id = t.id
				WHERE ct.class_id = ? AND t.user_id = ?
			)`, classID, user.ID).Scan(&exists)
		if err != nil {
			log.Printf("Failed to check teacher assignment: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher assignment")
			return
		}

		if !exists {
			SendErrorResponse(c, http.StatusForbidden, "Teachers can only manage students in classes they teach")
			return
		}
	}

	// Check if class exists
	var classExists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM classes WHERE id = ?)", classID).Scan(&classExists)
	if err != nil {
		log.Printf("Failed to check class existence: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to check class")
		return
	}

	if !classExists {
		SendErrorResponse(c, http.StatusNotFound, "Class not found")
		return
	}

	// Check if student exists
	var studentExists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM students WHERE id = ?)", req.StudentID).Scan(&studentExists)
	if err != nil {
		log.Printf("Failed to check student existence: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to check student")
		return
	}

	if !studentExists {
		SendErrorResponse(c, http.StatusNotFound, "Student not found")
		return
	}

	// Check if student is already enrolled in the class
	var enrollmentExists bool
	err = db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM student_enrollments WHERE class_id = ? AND student_id = ?)",
		classID, req.StudentID).Scan(&enrollmentExists)
	if err != nil {
		log.Printf("Failed to check enrollment: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to check enrollment")
		return
	}

	if enrollmentExists {
		SendErrorResponse(c, http.StatusConflict, "Student is already enrolled in this class")
		return
	}

	// Create the enrollment
	enrollmentID := uuid.New().String()
	now := time.Now()

	_, err = db.ExecContext(ctx,
		`INSERT INTO student_enrollments (id, class_id, student_id, status, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		enrollmentID, classID, req.StudentID, "active", now)
	if err != nil {
		log.Printf("Failed to create enrollment: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to enroll student")
		return
	}

	// Return the created enrollment
	enrollment := models.StudentEnrollment{
		ID:        enrollmentID,
		ClassID:   classID,
		StudentID: req.StudentID,
		Status:    "active",
		CreatedAt: now,
	}

	c.JSON(http.StatusCreated, enrollment)
}

// GetClassRoster handles retrieving all students enrolled in a class
func GetClassRoster(c *gin.Context) {
	// Get current user
	user, err := GetCurrentUser(c)
	if err != nil {
		SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Check if user has required role (teacher or admin)
	if user.Role != string(middleware.RoleTeacher) && user.Role != string(middleware.RoleAdmin) {
		SendErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	// Get class ID from URL parameter
	classID := c.Param("class_id")
	if classID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Class ID is required")
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(classID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid class ID format")
		return
	}

	// Get database connection
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// For teachers, verify they are assigned to the class
	if user.Role == string(middleware.RoleTeacher) {
		var exists bool
		err = db.QueryRowContext(ctx,
			`SELECT EXISTS(
				SELECT 1 FROM class_teachers ct
				JOIN teachers t ON ct.teacher_id = t.id
				WHERE ct.class_id = ? AND t.user_id = ?
			)`, classID, user.ID).Scan(&exists)
		if err != nil {
			log.Printf("Failed to check teacher assignment: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher assignment")
			return
		}

		if !exists {
			SendErrorResponse(c, http.StatusForbidden, "Teachers can only view rosters for classes they teach")
			return
		}
	}

	// Check if class exists
	var classExists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM classes WHERE id = ?)", classID).Scan(&classExists)
	if err != nil {
		log.Printf("Failed to check class existence: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to check class")
		return
	}

	if !classExists {
		SendErrorResponse(c, http.StatusNotFound, "Class not found")
		return
	}

	// Get class roster
	query := `
		SELECT se.id, se.class_id, se.student_id, se.status, se.created_at, u.name, u.email
		FROM student_enrollments se
		JOIN students s ON se.student_id = s.id
		JOIN users u ON s.user_id = u.id
		WHERE se.class_id = ?
		ORDER BY u.name`

	rows, err := db.QueryContext(ctx, query, classID)
	if err != nil {
		log.Printf("Failed to query class roster: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve class roster")
		return
	}
	defer rows.Close()

	// Process results
	type rosterItem struct {
		ID          string    `json:"id"`
		ClassID     string    `json:"class_id"`
		StudentID   string    `json:"student_id"`
		StudentName string    `json:"student_name"`
		StudentEmail string   `json:"student_email"`
		Status      string    `json:"status"`
		CreatedAt   time.Time `json:"created_at"`
	}

	var roster []rosterItem
	for rows.Next() {
		var item rosterItem
		err := rows.Scan(&item.ID, &item.ClassID, &item.StudentID, &item.Status, &item.CreatedAt, &item.StudentName, &item.StudentEmail)
		if err != nil {
			log.Printf("Failed to scan roster item: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to process class roster")
			return
		}
		roster = append(roster, item)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating roster: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve class roster")
		return
	}

	c.JSON(http.StatusOK, roster)
}

// RemoveStudentFromClass handles removing a student from a class
func RemoveStudentFromClass(c *gin.Context) {
	// Get current user
	user, err := GetCurrentUser(c)
	if err != nil {
		SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Check if user has required role (teacher or admin)
	if user.Role != string(middleware.RoleTeacher) && user.Role != string(middleware.RoleAdmin) {
		SendErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	// Get class ID and student ID from URL parameters
	classID := c.Param("class_id")
	studentID := c.Param("student_id")

	if classID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Class ID is required")
		return
	}

	if studentID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Student ID is required")
		return
	}

	// Validate UUID formats
	if _, err := uuid.Parse(classID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid class ID format")
		return
	}

	if _, err := uuid.Parse(studentID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid student ID format")
		return
	}

	// Get database connection
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// For teachers, verify they are assigned to the class
	if user.Role == string(middleware.RoleTeacher) {
		var exists bool
		err = db.QueryRowContext(ctx,
			`SELECT EXISTS(
				SELECT 1 FROM class_teachers ct
				JOIN teachers t ON ct.teacher_id = t.id
				WHERE ct.class_id = ? AND t.user_id = ?
			)`, classID, user.ID).Scan(&exists)
		if err != nil {
			log.Printf("Failed to check teacher assignment: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher assignment")
			return
		}

		if !exists {
			SendErrorResponse(c, http.StatusForbidden, "Teachers can only manage students in classes they teach")
			return
		}
	}

	// Check if enrollment exists
	var enrollmentExists bool
	err = db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM student_enrollments WHERE class_id = ? AND student_id = ?)",
		classID, studentID).Scan(&enrollmentExists)
	if err != nil {
		log.Printf("Failed to check enrollment: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to check enrollment")
		return
	}

	if !enrollmentExists {
		SendErrorResponse(c, http.StatusNotFound, "Student is not enrolled in this class")
		return
	}

	// Remove the enrollment
	_, err = db.ExecContext(ctx, "DELETE FROM student_enrollments WHERE class_id = ? AND student_id = ?", classID, studentID)
	if err != nil {
		log.Printf("Failed to remove enrollment: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to remove student from class")
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateStudentEnrollmentStatus handles updating a student's enrollment status
func UpdateStudentEnrollmentStatus(c *gin.Context) {
	// Get current user
	user, err := GetCurrentUser(c)
	if err != nil {
		SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Check if user has required role (teacher or admin)
	if user.Role != string(middleware.RoleTeacher) && user.Role != string(middleware.RoleAdmin) {
		SendErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	// Get class ID and student ID from URL parameters
	classID := c.Param("class_id")
	studentID := c.Param("student_id")

	if classID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Class ID is required")
		return
	}

	if studentID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Student ID is required")
		return
	}

	// Validate UUID formats
	if _, err := uuid.Parse(classID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid class ID format")
		return
	}

	if _, err := uuid.Parse(studentID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid student ID format")
		return
	}

	// Bind request data
	var req StudentEnrollmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// Get database connection
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// For teachers, verify they are assigned to the class
	if user.Role == string(middleware.RoleTeacher) {
		var exists bool
		err = db.QueryRowContext(ctx,
			`SELECT EXISTS(
				SELECT 1 FROM class_teachers ct
				JOIN teachers t ON ct.teacher_id = t.id
				WHERE ct.class_id = ? AND t.user_id = ?
			)`, classID, user.ID).Scan(&exists)
		if err != nil {
			log.Printf("Failed to check teacher assignment: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher assignment")
			return
		}

		if !exists {
			SendErrorResponse(c, http.StatusForbidden, "Teachers can only manage students in classes they teach")
			return
		}
	}

	// Check if enrollment exists
	var enrollmentExists bool
	err = db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM student_enrollments WHERE class_id = ? AND student_id = ?)",
		classID, studentID).Scan(&enrollmentExists)
	if err != nil {
		log.Printf("Failed to check enrollment: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to check enrollment")
		return
	}

	if !enrollmentExists {
		SendErrorResponse(c, http.StatusNotFound, "Student is not enrolled in this class")
		return
	}

	// Update the enrollment status
	now := time.Now()
	_, err = db.ExecContext(ctx,
		`UPDATE student_enrollments
		 SET status = ?, created_at = ?  -- Using created_at to store the last updated time
		 WHERE class_id = ? AND student_id = ?`,
		req.Status, now, classID, studentID)
	if err != nil {
		log.Printf("Failed to update enrollment: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to update enrollment status")
		return
	}

	// Return the updated enrollment
	enrollment := models.StudentEnrollment{
		ID:        "", // ID not returned in this endpoint
		ClassID:   classID,
		StudentID: studentID,
		Status:    req.Status,
		CreatedAt: now,
	}

	c.JSON(http.StatusOK, enrollment)
}