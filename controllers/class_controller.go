package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/aprianimmanuel/backend-app/middleware"
	"github.com/aprianimmanuel/backend-app/models"
)

// ClassRequest represents the request body for creating/updating a class
type ClassRequest struct {
	SchoolID     string `json:"school_id" binding:"required,uuid"`
	Name         string `json:"name" binding:"required,max=100"`
	GradeLevel   int    `json:"grade_level" binding:"required,min=1,max=12"`
	AcademicYear string `json:"academic_year" binding:"required"`
}

// CreateClass handles the creation of a new class
func CreateClass(c *gin.Context) {
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
	var req ClassRequest
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
		err = db.QueryRow(context.Background(), 
			"SELECT school_id FROM teachers WHERE user_id = $1", user.ID).Scan(&teacherSchoolID)
		if err != nil {
			log.Printf("Failed to get teacher school: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher school")
			return
		}

		if teacherSchoolID != req.SchoolID {
			SendErrorResponse(c, http.StatusForbidden, "Teachers can only create classes in their school")
			return
		}
	}

	// Create the class in the database
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Generate a new UUID for the class
	classID := uuid.New().String()
	now := time.Now()

	// Insert the class
	_, err = db.Exec(ctx,
		`INSERT INTO classes (id, school_id, name, grade_level, academic_year, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		classID, req.SchoolID, req.Name, req.GradeLevel, req.AcademicYear, now, now)
	if err != nil {
		log.Printf("Failed to create class: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to create class")
		return
	}

	// Return the created class
	class := models.Class{
		ID:           classID,
		SchoolID:     req.SchoolID,
		Name:         req.Name,
		GradeLevel:   req.GradeLevel,
		AcademicYear: req.AcademicYear,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	c.JSON(http.StatusCreated, class)
}

// GetAllClasses handles retrieving all classes with optional filtering
func GetAllClasses(c *gin.Context) {
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
	gradeLevel := c.Query("grade_level")
	academicYear := c.Query("academic_year")

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
	query := "SELECT id, school_id, name, grade_level, academic_year, created_at, updated_at FROM classes WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// For teachers, only show classes from their school
	if user.Role == string(middleware.RoleTeacher) {
		var teacherSchoolID string
		err = db.QueryRow(ctx, 
			"SELECT school_id FROM teachers WHERE user_id = $1", user.ID).Scan(&teacherSchoolID)
		if err != nil {
			log.Printf("Failed to get teacher school: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher school")
			return
		}

		query += " AND school_id = $" + string(rune(argIndex+'0'))
		args = append(args, teacherSchoolID)
		argIndex++
	} else if schoolID != "" {
		// For admins, filter by school_id if provided
		query += " AND school_id = $" + string(rune(argIndex+'0'))
		args = append(args, schoolID)
		argIndex++
	}

	// Add other filters
	if gradeLevel != "" {
		query += " AND grade_level = $" + string(rune(argIndex+'0'))
		args = append(args, gradeLevel)
		argIndex++
	}

	if academicYear != "" {
		query += " AND academic_year = $" + string(rune(argIndex+'0'))
		args = append(args, academicYear)
		argIndex++
	}

	// Execute query
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to query classes: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve classes")
		return
	}
	defer rows.Close()

	// Process results
	var classes []models.Class
	for rows.Next() {
		var class models.Class
		err := rows.Scan(&class.ID, &class.SchoolID, &class.Name, &class.GradeLevel, &class.AcademicYear, &class.CreatedAt, &class.UpdatedAt)
		if err != nil {
			log.Printf("Failed to scan class: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to process classes")
			return
		}
		classes = append(classes, class)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating classes: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve classes")
		return
	}

	c.JSON(http.StatusOK, classes)
}

// GetClassByID handles retrieving a specific class by ID
func GetClassByID(c *gin.Context) {
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
	classID := c.Param("id")
	if classID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Class ID is required")
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(classID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid class ID format")
		return
	}

	// Get class from database
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var class models.Class
	err = db.QueryRow(ctx,
		"SELECT id, school_id, name, grade_level, academic_year, created_at, updated_at FROM classes WHERE id = $1",
		classID).Scan(&class.ID, &class.SchoolID, &class.Name, &class.GradeLevel, &class.AcademicYear, &class.CreatedAt, &class.UpdatedAt)
	
	if err == pgx.ErrNoRows {
		SendErrorResponse(c, http.StatusNotFound, "Class not found")
		return
	}
	
	if err != nil {
		log.Printf("Failed to get class: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve class")
		return
	}

	// For teachers, verify they belong to the same school as the class
	if user.Role == string(middleware.RoleTeacher) {
		var teacherSchoolID string
		err = db.QueryRow(ctx, 
			"SELECT school_id FROM teachers WHERE user_id = $1", user.ID).Scan(&teacherSchoolID)
		if err != nil {
			log.Printf("Failed to get teacher school: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher school")
			return
		}

		if teacherSchoolID != class.SchoolID {
			SendErrorResponse(c, http.StatusForbidden, "Teachers can only access classes in their school")
			return
		}
	}

	c.JSON(http.StatusOK, class)
}

// UpdateClass handles updating an existing class
func UpdateClass(c *gin.Context) {
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
	classID := c.Param("id")
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
	var req ClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// Get class from database
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingClass models.Class
	err = db.QueryRow(ctx,
		"SELECT id, school_id, name, grade_level, academic_year, created_at, updated_at FROM classes WHERE id = $1",
		classID).Scan(&existingClass.ID, &existingClass.SchoolID, &existingClass.Name, &existingClass.GradeLevel, &existingClass.AcademicYear, &existingClass.CreatedAt, &existingClass.UpdatedAt)
	
	if err == pgx.ErrNoRows {
		SendErrorResponse(c, http.StatusNotFound, "Class not found")
		return
	}
	
	if err != nil {
		log.Printf("Failed to get class: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve class")
		return
	}

	// For teachers, verify they belong to the same school as the class
	if user.Role == string(middleware.RoleTeacher) {
		var teacherSchoolID string
		err = db.QueryRow(ctx, 
			"SELECT school_id FROM teachers WHERE user_id = $1", user.ID).Scan(&teacherSchoolID)
		if err != nil {
			log.Printf("Failed to get teacher school: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher school")
			return
		}

		if teacherSchoolID != existingClass.SchoolID {
			SendErrorResponse(c, http.StatusForbidden, "Teachers can only update classes in their school")
			return
		}

		// Teachers cannot change the school_id
		if req.SchoolID != existingClass.SchoolID {
			SendErrorResponse(c, http.StatusForbidden, "Teachers cannot change the school of a class")
			return
		}
	}

	// Update the class
	now := time.Now()
	_, err = db.Exec(ctx,
		`UPDATE classes 
		 SET school_id = $1, name = $2, grade_level = $3, academic_year = $4, updated_at = $5
		 WHERE id = $6`,
		req.SchoolID, req.Name, req.GradeLevel, req.AcademicYear, now, classID)
	if err != nil {
		log.Printf("Failed to update class: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to update class")
		return
	}

	// Return the updated class
	updatedClass := models.Class{
		ID:           classID,
		SchoolID:     req.SchoolID,
		Name:         req.Name,
		GradeLevel:   req.GradeLevel,
		AcademicYear: req.AcademicYear,
		CreatedAt:    existingClass.CreatedAt,
		UpdatedAt:    now,
	}

	c.JSON(http.StatusOK, updatedClass)
}

// DeleteClass handles deleting a class
func DeleteClass(c *gin.Context) {
	// Get current user
	user, err := GetCurrentUser(c)
	if err != nil {
		SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Check if user has required role (admin only)
	if user.Role != string(middleware.RoleAdmin) {
		SendErrorResponse(c, http.StatusForbidden, "Only administrators can delete classes")
		return
	}

	// Get class ID from URL parameter
	classID := c.Param("id")
	if classID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Class ID is required")
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(classID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid class ID format")
		return
	}

	// Delete class from database
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if class exists
	var exists bool
	err = db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM classes WHERE id = $1)", classID).Scan(&exists)
	if err != nil {
		log.Printf("Failed to check class existence: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to check class")
		return
	}

	if !exists {
		SendErrorResponse(c, http.StatusNotFound, "Class not found")
		return
	}

	// Delete the class (will cascade to related tables)
	_, err = db.Exec(ctx, "DELETE FROM classes WHERE id = $1", classID)
	if err != nil {
		log.Printf("Failed to delete class: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete class")
		return
	}

	c.Status(http.StatusNoContent)
}