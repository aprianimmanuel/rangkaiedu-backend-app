package controllers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/aprianimmanuel/rangkaiedu-backend/middleware"
	"github.com/aprianimmanuel/rangkaiedu-backend/models"
)

// SubjectRequest represents the request body for creating/updating a subject
type SubjectRequest struct {
	SchoolID    string `json:"school_id" binding:"required,uuid"`
	Name        string `json:"name" binding:"required,max=100"`
	Code        string `json:"code" binding:"required,max=20"`
	Description string `json:"description" binding:"max=500"`
}

// CreateSubject handles the creation of a new subject
func CreateSubject(c *gin.Context) {
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

	// Bind request data
	var req SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// For teachers, verify they belong to the same school
	if user.Role == string(middleware.RoleTeacher) {
		// Check if teacher belongs to the specified school
		db, err := GetDBConnection()
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
			return
		}
		defer db.Close()

		var teacherSchoolID string
		err = db.QueryRowContext(context.Background(),
			"SELECT school_id FROM teachers WHERE user_id = ?", user.ID).Scan(&teacherSchoolID)
		if err != nil {
			log.Printf("Failed to get teacher school: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher school")
			return
		}

		if teacherSchoolID != req.SchoolID {
			SendErrorResponse(c, http.StatusForbidden, "Teachers can only create subjects in their school")
			return
		}
	}

	// Create the subject in the database
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Generate a new UUID for the subject
	subjectID := uuid.New().String()
	now := time.Now()

	// Insert the subject
	_, err = db.ExecContext(ctx,
		`INSERT INTO subjects (id, school_id, name, code, description, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		subjectID, req.SchoolID, req.Name, req.Code, req.Description, now, now)
	if err != nil {
		log.Printf("Failed to create subject: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to create subject")
		return
	}

	// Return the created subject
	subject := models.Subject{
		ID:          subjectID,
		SchoolID:    req.SchoolID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	c.JSON(http.StatusCreated, subject)
}

// GetAllSubjects handles retrieving all subjects with optional filtering
func GetAllSubjects(c *gin.Context) {
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

	// Get query parameters
	schoolID := c.Query("school_id")

	// Build query based on user role and filters
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Base query
	query := "SELECT id, school_id, name, code, description, created_at, updated_at FROM subjects WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// For teachers, only show subjects from their school
	if user.Role == string(middleware.RoleTeacher) {
		var teacherSchoolID string
		err = db.QueryRowContext(ctx,
			"SELECT school_id FROM teachers WHERE user_id = ?", user.ID).Scan(&teacherSchoolID)
		if err != nil {
			log.Printf("Failed to get teacher school: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher school")
			return
		}

		query += " AND school_id = ?"
		args = append(args, teacherSchoolID)
		argIndex++
	} else if schoolID != "" {
		// For admins, filter by school_id if provided
		query += " AND school_id = ?"
		args = append(args, schoolID)
		argIndex++
	}

	// Execute query
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to query subjects: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve subjects")
		return
	}
	defer rows.Close()

	// Process results
	var subjects []models.Subject
	for rows.Next() {
		var subject models.Subject
		err := rows.Scan(&subject.ID, &subject.SchoolID, &subject.Name, &subject.Code, &subject.Description, &subject.CreatedAt, &subject.UpdatedAt)
		if err != nil {
			log.Printf("Failed to scan subject: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to process subjects")
			return
		}
		subjects = append(subjects, subject)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating subjects: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve subjects")
		return
	}

	c.JSON(http.StatusOK, subjects)
}

// GetSubjectByID handles retrieving a specific subject by ID
func GetSubjectByID(c *gin.Context) {
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

	// Get subject ID from URL parameter
	subjectID := c.Param("id")
	if subjectID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Subject ID is required")
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(subjectID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid subject ID format")
		return
	}

	// Get subject from database
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var subject models.Subject
	err = db.QueryRowContext(ctx,
		"SELECT id, school_id, name, code, description, created_at, updated_at FROM subjects WHERE id = ?",
		subjectID).Scan(&subject.ID, &subject.SchoolID, &subject.Name, &subject.Code, &subject.Description, &subject.CreatedAt, &subject.UpdatedAt)
	
	if err == sql.ErrNoRows {
		SendErrorResponse(c, http.StatusNotFound, "Subject not found")
		return
	}
	
	if err != nil {
		log.Printf("Failed to get subject: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve subject")
		return
	}

	// For teachers, verify they belong to the same school as the subject
	if user.Role == string(middleware.RoleTeacher) {
		var teacherSchoolID string
		err = db.QueryRowContext(ctx,
			"SELECT school_id FROM teachers WHERE user_id = ?", user.ID).Scan(&teacherSchoolID)
		if err != nil {
			log.Printf("Failed to get teacher school: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher school")
			return
		}

		if teacherSchoolID != subject.SchoolID {
			SendErrorResponse(c, http.StatusForbidden, "Teachers can only access subjects in their school")
			return
		}
	}

	c.JSON(http.StatusOK, subject)
}

// UpdateSubject handles updating an existing subject
func UpdateSubject(c *gin.Context) {
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

	// Get subject ID from URL parameter
	subjectID := c.Param("id")
	if subjectID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Subject ID is required")
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(subjectID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid subject ID format")
		return
	}

	// Bind request data
	var req SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// Get subject from database
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingSubject models.Subject
	err = db.QueryRowContext(ctx,
		"SELECT id, school_id, name, code, description, created_at, updated_at FROM subjects WHERE id = ?",
		subjectID).Scan(&existingSubject.ID, &existingSubject.SchoolID, &existingSubject.Name, &existingSubject.Code, &existingSubject.Description, &existingSubject.CreatedAt, &existingSubject.UpdatedAt)
	
	if err == sql.ErrNoRows {
		SendErrorResponse(c, http.StatusNotFound, "Subject not found")
		return
	}
	
	if err != nil {
		log.Printf("Failed to get subject: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve subject")
		return
	}

	// For teachers, verify they belong to the same school as the subject
	if user.Role == string(middleware.RoleTeacher) {
		var teacherSchoolID string
		err = db.QueryRowContext(ctx,
			"SELECT school_id FROM teachers WHERE user_id = ?", user.ID).Scan(&teacherSchoolID)
		if err != nil {
			log.Printf("Failed to get teacher school: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher school")
			return
		}

		if teacherSchoolID != existingSubject.SchoolID {
			SendErrorResponse(c, http.StatusForbidden, "Teachers can only update subjects in their school")
			return
		}

		// Teachers cannot change the school_id
		if req.SchoolID != existingSubject.SchoolID {
			SendErrorResponse(c, http.StatusForbidden, "Teachers cannot change the school of a subject")
			return
		}
	}

	// Update the subject
	now := time.Now()
	_, err = db.ExecContext(ctx,
		`UPDATE subjects
		 SET school_id = ?, name = ?, code = ?, description = ?, updated_at = ?
		 WHERE id = ?`,
		req.SchoolID, req.Name, req.Code, req.Description, now, subjectID)
	if err != nil {
		log.Printf("Failed to update subject: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to update subject")
		return
	}

	// Return the updated subject
	updatedSubject := models.Subject{
		ID:          subjectID,
		SchoolID:    req.SchoolID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		CreatedAt:   existingSubject.CreatedAt,
		UpdatedAt:   now,
	}

	c.JSON(http.StatusOK, updatedSubject)
}

// DeleteSubject handles deleting a subject
func DeleteSubject(c *gin.Context) {
	// Get current user
	user, err := GetCurrentUser(c)
	if err != nil {
		SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Check if user has required role (admin only)
	if user.Role != string(middleware.RoleAdmin) {
		SendErrorResponse(c, http.StatusForbidden, "Only administrators can delete subjects")
		return
	}

	// Get subject ID from URL parameter
	subjectID := c.Param("id")
	if subjectID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Subject ID is required")
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(subjectID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid subject ID format")
		return
	}

	// Delete subject from database
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if subject exists
	var exists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM subjects WHERE id = ?)", subjectID).Scan(&exists)
	if err != nil {
		log.Printf("Failed to check subject existence: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to check subject")
		return
	}

	if !exists {
		SendErrorResponse(c, http.StatusNotFound, "Subject not found")
		return
	}

	// Delete the subject (will cascade to related tables)
	_, err = db.ExecContext(ctx, "DELETE FROM subjects WHERE id = ?", subjectID)
	if err != nil {
		log.Printf("Failed to delete subject: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete subject")
		return
	}

	c.Status(http.StatusNoContent)
}