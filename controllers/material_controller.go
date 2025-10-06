package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/aprianimmanuel/rangkaiedu-backend/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/middleware"
	"github.com/aprianimmanuel/rangkaiedu-backend/models"
	"github.com/aprianimmanuel/rangkaiedu-backend/utils/storage"
)

// MaterialRequest represents the request body for creating/updating a material
type MaterialRequest struct {
	ClassID     string `form:"class_id" json:"class_id" binding:"required,uuid"`
	SubjectID   string `form:"subject_id" json:"subject_id" binding:"required,uuid"`
	Title       string `form:"title" json:"title" binding:"required,max=255"`
	Description string `form:"description" json:"description" binding:"max=1000"`
	Visibility  string `form:"visibility" json:"visibility" binding:"oneof=public private class_only"`
}

// CreateMaterial handles uploading a new teaching material
func CreateMaterial(c *gin.Context) {
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

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(100<<20); err != nil { // 100 MB max memory
		SendErrorResponse(c, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	// Get form values
	classID := c.PostForm("class_id")
	subjectID := c.PostForm("subject_id")
	title := c.PostForm("title")
	description := c.PostForm("description")
	visibility := c.PostForm("visibility")

	// Set default visibility if not provided
	if visibility == "" {
		visibility = "class_only"
	}

	// Validate required fields
	if classID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "class_id is required")
		return
	}

	if subjectID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "subject_id is required")
		return
	}

	if title == "" {
		SendErrorResponse(c, http.StatusBadRequest, "title is required")
		return
	}

	// Validate UUID formats
	if _, err := uuid.Parse(classID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid class ID format")
		return
	}

	if _, err := uuid.Parse(subjectID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid subject ID format")
		return
	}

	// Get the uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "file is required")
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

	// For teachers, verify they teach the class
	if user.Role == string(middleware.RoleTeacher) {
		var exists bool
		err = db.QueryRow(ctx,
			`SELECT EXISTS(
				SELECT 1 FROM class_teachers ct 
				JOIN teachers t ON ct.teacher_id = t.id 
				WHERE ct.class_id = $1 AND t.user_id = $2
			)`, classID, user.ID).Scan(&exists)
		if err != nil {
			log.Printf("Failed to check teacher assignment: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher assignment")
			return
		}

		if !exists {
			SendErrorResponse(c, http.StatusForbidden, "Teachers can only upload materials for classes they teach")
			return
		}
	}

	// Check if class exists
	var classExists bool
	err = db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM classes WHERE id = $1)", classID).Scan(&classExists)
	if err != nil {
		log.Printf("Failed to check class existence: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to check class")
		return
	}

	if !classExists {
		SendErrorResponse(c, http.StatusNotFound, "Class not found")
		return
	}

	// Check if subject exists
	var subjectExists bool
	err = db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM subjects WHERE id = $1)", subjectID).Scan(&subjectExists)
	if err != nil {
		log.Printf("Failed to check subject existence: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to check subject")
		return
	}

	if !subjectExists {
		SendErrorResponse(c, http.StatusNotFound, "Subject not found")
		return
	}

	// Get teacher ID for the material
	var teacherID string
	if user.Role == string(middleware.RoleTeacher) {
		err = db.QueryRow(ctx, "SELECT id FROM teachers WHERE user_id = $1", user.ID).Scan(&teacherID)
		if err != nil {
			log.Printf("Failed to get teacher ID: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to get teacher information")
			return
		}
	} else {
		// For admin, we need to determine which teacher is uploading
		// In a real implementation, this might be passed as a parameter
		// For now, we'll use a placeholder - in practice this should be handled differently
		teacherID = uuid.New().String() // This is a placeholder
	}

	// Initialize storage provider
	cfg := config.Load()
	storageProvider, err := storage.NewStorageProvider(cfg)
	if err != nil {
		log.Printf("Failed to initialize storage provider: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to initialize storage provider")
		return
	}

	// Validate and save the file
	fileMetadata, err := storageProvider.SaveFile(fileHeader, "materials")
	if err != nil {
		log.Printf("Failed to save file: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to save file: "+err.Error())
		return
	}

	// Create the material in the database
	materialID := uuid.New().String()
	now := time.Now()

	_, err = db.Exec(ctx,
		`INSERT INTO materials (id, class_id, subject_id, teacher_id, title, description, file_name, file_path, file_type, file_size, visibility, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		materialID, classID, subjectID, teacherID, title, description, fileMetadata.FileName,
		fileMetadata.FilePath, fileMetadata.FileType, fileMetadata.FileSize,
		visibility, now, now)
	if err != nil {
		log.Printf("Failed to create material: %v", err)
		// Clean up the uploaded file if database insert fails
		storageProvider.DeleteFile(fileMetadata.FilePath)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to create material")
		return
	}

	// Return the created material
	material := models.Material{
		ID:          materialID,
		ClassID:     classID,
		SubjectID:   subjectID,
		TeacherID:   teacherID,
		Title:       title,
		Description: description,
		FileName:    fileMetadata.FileName,
		FilePath:    fileMetadata.FilePath,
		FileType:    fileMetadata.FileType,
		FileSize:    fileMetadata.FileSize,
		Visibility:  visibility,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	c.JSON(http.StatusCreated, material)
}

// GetAllMaterials handles retrieving all materials with optional filtering
func GetAllMaterials(c *gin.Context) {
	// Get current user
	user, err := GetCurrentUser(c)
	if err != nil {
		SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get query parameters
	classID := c.Query("class_id")
	subjectID := c.Query("subject_id")
	teacherID := c.Query("teacher_id")
	visibility := c.Query("visibility")

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

	// Base query - exclude file_path for security
	query := "SELECT id, class_id, subject_id, teacher_id, title, description, file_name, file_type, file_size, visibility, created_at, updated_at FROM materials WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// Add filters based on user role and query parameters
	if classID != "" {
		query += " AND class_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, classID)
		argIndex++
	}

	if subjectID != "" {
		query += " AND subject_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, subjectID)
		argIndex++
	}

	if teacherID != "" {
		query += " AND teacher_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, teacherID)
		argIndex++
	}

	if visibility != "" {
		query += " AND visibility = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, visibility)
		argIndex++
	}

	// For students, only show public materials and materials for their classes
	if user.Role == string(middleware.RoleStudent) {
		// Get student's class IDs
		rows, err := db.Query(ctx, 
			"SELECT class_id FROM student_enrollments WHERE student_id = (SELECT id FROM students WHERE user_id = $1)", user.ID)
		if err != nil {
			log.Printf("Failed to get student classes: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve materials")
			return
		}
		defer rows.Close()

		classIDs := []string{}
		for rows.Next() {
			var classID string
			if err := rows.Scan(&classID); err != nil {
				log.Printf("Failed to scan class ID: %v", err)
				SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve materials")
				return
			}
			classIDs = append(classIDs, classID)
		}

		// Build visibility filter for students
		if len(classIDs) > 0 {
			// Students can see public materials and materials for their classes
			query += " AND (visibility = 'public' OR (visibility = 'class_only' AND class_id = ANY($" + fmt.Sprintf("%d", argIndex) + ")))"
			args = append(args, classIDs)
			argIndex++
		} else {
			// Students with no classes can only see public materials
			query += " AND visibility = 'public'"
		}
	} else if user.Role == string(middleware.RoleTeacher) {
		// Teachers can see public materials and materials for classes they teach
		query += " AND (visibility = 'public' OR teacher_id = (SELECT id FROM teachers WHERE user_id = $" + fmt.Sprintf("%d", argIndex) + "))"
		args = append(args, user.ID)
		argIndex++
	}
	// Admins can see all materials

	// Execute query
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to query materials: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve materials")
		return
	}
	defer rows.Close()

	// Process results
	var materials []models.Material
	for rows.Next() {
		var material models.Material
		err := rows.Scan(&material.ID, &material.ClassID, &material.SubjectID, &material.TeacherID, 
			&material.Title, &material.Description, &material.FileName, &material.FileType, 
			&material.FileSize, &material.Visibility, &material.CreatedAt, &material.UpdatedAt)
		if err != nil {
			log.Printf("Failed to scan material: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to process materials")
			return
		}
		materials = append(materials, material)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating materials: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve materials")
		return
	}

	c.JSON(http.StatusOK, materials)
}

// GetMaterialByID handles retrieving a specific material by ID
func GetMaterialByID(c *gin.Context) {
	// Get current user
	user, err := GetCurrentUser(c)
	if err != nil {
		SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get material ID from URL parameter
	materialID := c.Param("id")
	if materialID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Material ID is required")
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(materialID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid material ID format")
		return
	}

	// Get material from database
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var material models.Material
	err = db.QueryRow(ctx,
		"SELECT id, class_id, subject_id, teacher_id, title, description, file_name, file_path, file_type, file_size, visibility, created_at, updated_at FROM materials WHERE id = $1",
		materialID).Scan(&material.ID, &material.ClassID, &material.SubjectID, &material.TeacherID, 
		&material.Title, &material.Description, &material.FileName, &material.FilePath, 
		&material.FileType, &material.FileSize, &material.Visibility, &material.CreatedAt, &material.UpdatedAt)
	
	if err == pgx.ErrNoRows {
		SendErrorResponse(c, http.StatusNotFound, "Material not found")
		return
	}
	
	if err != nil {
		log.Printf("Failed to get material: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve material")
		return
	}

	// Check access permissions based on user role and material visibility
	canAccess := false
	
	switch user.Role {
	case string(middleware.RoleAdmin):
		// Admins can access all materials
		canAccess = true
	case string(middleware.RoleTeacher):
		// Teachers can access public materials and materials they own
		if material.Visibility == "public" || material.TeacherID == user.ID {
			canAccess = true
		}
	case string(middleware.RoleStudent):
		// Students can access public materials and materials for their classes
		if material.Visibility == "public" {
			canAccess = true
		} else if material.Visibility == "class_only" {
			// Check if student is enrolled in the class
			var enrolled bool
			err = db.QueryRow(ctx, 
				"SELECT EXISTS(SELECT 1 FROM student_enrollments WHERE class_id = $1 AND student_id = (SELECT id FROM students WHERE user_id = $2))", 
				material.ClassID, user.ID).Scan(&enrolled)
			if err != nil {
				log.Printf("Failed to check enrollment: %v", err)
				SendErrorResponse(c, http.StatusInternalServerError, "Failed to check access permissions")
				return
			}
			canAccess = enrolled
		}
	default:
		// Other roles can only access public materials
		canAccess = material.Visibility == "public"
	}

	if !canAccess {
		SendErrorResponse(c, http.StatusForbidden, "Insufficient permissions to access this material")
		return
	}

	// Remove file_path from response for security
	material.FilePath = ""

	c.JSON(http.StatusOK, material)
}

// UpdateMaterial handles updating an existing material
func UpdateMaterial(c *gin.Context) {
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

	// Get material ID from URL parameter
	materialID := c.Param("id")
	if materialID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Material ID is required")
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(materialID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid material ID format")
		return
	}

	// Bind request data
	var req MaterialRequest
	if err := c.ShouldBind(&req); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// Get material from database
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingMaterial models.Material
	err = db.QueryRow(ctx,
		"SELECT id, class_id, subject_id, teacher_id, title, description, file_name, file_path, file_type, file_size, visibility, created_at, updated_at FROM materials WHERE id = $1",
		materialID).Scan(&existingMaterial.ID, &existingMaterial.ClassID, &existingMaterial.SubjectID, &existingMaterial.TeacherID, 
		&existingMaterial.Title, &existingMaterial.Description, &existingMaterial.FileName, &existingMaterial.FilePath, 
		&existingMaterial.FileType, &existingMaterial.FileSize, &existingMaterial.Visibility, &existingMaterial.CreatedAt, &existingMaterial.UpdatedAt)
	
	if err == pgx.ErrNoRows {
		SendErrorResponse(c, http.StatusNotFound, "Material not found")
		return
	}
	
	if err != nil {
		log.Printf("Failed to get material: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve material")
		return
	}

	// For teachers, verify they own the material
	if user.Role == string(middleware.RoleTeacher) {
		var teacherID string
		err = db.QueryRow(ctx, "SELECT id FROM teachers WHERE user_id = $1", user.ID).Scan(&teacherID)
		if err != nil {
			log.Printf("Failed to get teacher ID: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher identity")
			return
		}

		if teacherID != existingMaterial.TeacherID {
			SendErrorResponse(c, http.StatusForbidden, "Teachers can only update their own materials")
			return
		}
	}

	// Update the material
	now := time.Now()
	_, err = db.Exec(ctx,
		`UPDATE materials 
		 SET title = $1, description = $2, visibility = $3, updated_at = $4
		 WHERE id = $5`,
		req.Title, req.Description, req.Visibility, now, materialID)
	if err != nil {
		log.Printf("Failed to update material: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to update material")
		return
	}

	// Return the updated material
	updatedMaterial := models.Material{
		ID:          materialID,
		ClassID:     existingMaterial.ClassID,
		SubjectID:   existingMaterial.SubjectID,
		TeacherID:   existingMaterial.TeacherID,
		Title:       req.Title,
		Description: req.Description,
		FileName:    existingMaterial.FileName,
		FilePath:    existingMaterial.FilePath,
		FileType:    existingMaterial.FileType,
		FileSize:    existingMaterial.FileSize,
		Visibility:  req.Visibility,
		CreatedAt:   existingMaterial.CreatedAt,
		UpdatedAt:   now,
	}

	// Remove file_path from response for security
	updatedMaterial.FilePath = ""

	c.JSON(http.StatusOK, updatedMaterial)
}

// DeleteMaterial handles deleting a material
func DeleteMaterial(c *gin.Context) {
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

	// Get material ID from URL parameter
	materialID := c.Param("id")
	if materialID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Material ID is required")
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(materialID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid material ID format")
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

	// Get material to check ownership and get file path
	var material models.Material
	err = db.QueryRow(ctx,
		"SELECT id, teacher_id, file_path FROM materials WHERE id = $1", materialID).Scan(&material.ID, &material.TeacherID, &material.FilePath)
	
	if err == pgx.ErrNoRows {
		SendErrorResponse(c, http.StatusNotFound, "Material not found")
		return
	}
	
	if err != nil {
		log.Printf("Failed to get material: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve material")
		return
	}

	// For teachers, verify they own the material
	if user.Role == string(middleware.RoleTeacher) {
		var teacherID string
		err = db.QueryRow(ctx, "SELECT id FROM teachers WHERE user_id = $1", user.ID).Scan(&teacherID)
		if err != nil {
			log.Printf("Failed to get teacher ID: %v", err)
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher identity")
			return
		}

		if teacherID != material.TeacherID {
			SendErrorResponse(c, http.StatusForbidden, "Teachers can only delete their own materials")
			return
		}
	}

	// Delete the material from database
	_, err = db.Exec(ctx, "DELETE FROM materials WHERE id = $1", materialID)
	if err != nil {
		log.Printf("Failed to delete material: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete material")
		return
	}

	// Initialize storage provider
	cfg := config.Load()
	storageProvider, err := storage.NewStorageProvider(cfg)
	if err != nil {
		log.Printf("Failed to initialize storage provider: %v", err)
		// Don't return error here as the database record is already deleted
	} else {
		// Delete the file from storage
		if material.FilePath != "" {
			if err := storageProvider.DeleteFile(material.FilePath); err != nil {
				log.Printf("Failed to delete file %s: %v", material.FilePath, err)
				// Don't return error here as the database record is already deleted
			}
		}
	}

	c.Status(http.StatusNoContent)
}

// DownloadMaterial handles downloading a material file
func DownloadMaterial(c *gin.Context) {
	// Get current user
	user, err := GetCurrentUser(c)
	if err != nil {
		SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get material ID from URL parameter
	materialID := c.Param("id")
	if materialID == "" {
		SendErrorResponse(c, http.StatusBadRequest, "Material ID is required")
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(materialID); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid material ID format")
		return
	}

	// Get material from database
	db, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var material models.Material
	err = db.QueryRow(ctx,
		"SELECT id, class_id, teacher_id, title, file_name, file_path, file_type, visibility FROM materials WHERE id = $1",
		materialID).Scan(&material.ID, &material.ClassID, &material.TeacherID, &material.Title, 
		&material.FileName, &material.FilePath, &material.FileType, &material.Visibility)
	
	if err == pgx.ErrNoRows {
		SendErrorResponse(c, http.StatusNotFound, "Material not found")
		return
	}
	
	if err != nil {
		log.Printf("Failed to get material: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve material")
		return
	}

	// Check access permissions based on user role and material visibility
	canAccess := false
	
	switch user.Role {
	case string(middleware.RoleAdmin):
		// Admins can access all materials
		canAccess = true
	case string(middleware.RoleTeacher):
		// Teachers can access public materials and materials they own
		if material.Visibility == "public" {
			canAccess = true
		} else {
			var teacherID string
			err = db.QueryRow(ctx, "SELECT id FROM teachers WHERE user_id = $1", user.ID).Scan(&teacherID)
			if err != nil {
				log.Printf("Failed to get teacher ID: %v", err)
				SendErrorResponse(c, http.StatusInternalServerError, "Failed to verify teacher identity")
				return
			}
			canAccess = (teacherID == material.TeacherID)
		}
	case string(middleware.RoleStudent):
		// Students can access public materials and materials for their classes
		if material.Visibility == "public" {
			canAccess = true
		} else if material.Visibility == "class_only" {
			// Check if student is enrolled in the class
			var enrolled bool
			err = db.QueryRow(ctx, 
				"SELECT EXISTS(SELECT 1 FROM student_enrollments WHERE class_id = $1 AND student_id = (SELECT id FROM students WHERE user_id = $2))", 
				material.ClassID, user.ID).Scan(&enrolled)
			if err != nil {
				log.Printf("Failed to check enrollment: %v", err)
				SendErrorResponse(c, http.StatusInternalServerError, "Failed to check access permissions")
				return
			}
			canAccess = enrolled
		}
	default:
		// Other roles can only access public materials
		canAccess = material.Visibility == "public"
	}

	if !canAccess {
		SendErrorResponse(c, http.StatusForbidden, "Insufficient permissions to access this material")
		return
	}

	// Initialize storage provider
	cfg := config.Load()
	storageProvider, err := storage.NewStorageProvider(cfg)
	if err != nil {
		log.Printf("Failed to initialize storage provider: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to initialize storage provider")
		return
	}

	// Get file URL from storage provider
	fileURL, err := storageProvider.GetFileURL(material.FilePath)
	if err != nil {
		log.Printf("Failed to get file URL: %v", err)
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to get file")
		return
	}

	// For local storage, we can still use c.File
	// For cloud storage, we would redirect to the URL or stream the content
	if cfg.StorageProvider == "local" || cfg.StorageProvider == "" {
		// Check if file exists
		if _, err := os.Stat(material.FilePath); os.IsNotExist(err) {
			SendErrorResponse(c, http.StatusNotFound, "File not found")
			return
		}

		// Set headers for file download
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+material.FileName)
		c.Header("Content-Type", material.FileType)

		// Serve the file
		c.File(material.FilePath)
	} else {
		// For cloud storage, redirect to the file URL
		// In a production environment, you might want to stream the content instead
		c.Redirect(http.StatusFound, fileURL)
	}
}