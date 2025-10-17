package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/handlers"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/middleware"
)

// SetupAllRoutes sets up all application routes
func SetupAllRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {
	// Setup auth routes
	SetupAuthRoutes(r, authHandler)
	
	// Setup class routes (commented out until handlers are implemented)
	// SetupClassRoutes(r)
	
	// Setup subject routes (commented out until handlers are implemented)
	// SetupSubjectRoutes(r)
	
	// Setup material routes (commented out until handlers are implemented)
	// SetupMaterialRoutes(r)
	
	// Setup health check routes (commented out until handlers are implemented)
	// SetupHealthRoutes(r)
}

// SetupClassRoutes sets up the class routes
func SetupClassRoutes(r *gin.Engine) {
	classes := r.Group("/api/classes")
	{
		// Apply authentication middleware to all class routes
		classes.Use(middleware.AuthRequired())
		
		// Class roster management (must come before single class routes to avoid conflicts)
		// classes.POST("/:class_id/students", handlers.AddStudentToClass)
		// classes.GET("/:class_id/students", handlers.GetClassRoster)
		// classes.DELETE("/:class_id/students/:student_id", handlers.RemoveStudentFromClass)
		// classes.PUT("/:class_id/students/:student_id", handlers.UpdateStudentEnrollmentStatus)
	
		// CRUD operations for classes
		// classes.POST("", handlers.CreateClass)
		// classes.GET("", handlers.GetAllClasses)
		// classes.GET("/by-id/:classId", handlers.GetClassByID)
		// classes.PUT("/by-id/:classId", handlers.UpdateClass)
		// classes.DELETE("/by-id/:classId", middleware.RoleRequired(middleware.RoleAdmin), handlers.DeleteClass)
	}
}

// SetupSubjectRoutes sets up the subject routes
func SetupSubjectRoutes(r *gin.Engine) {
	subjects := r.Group("/api/subjects")
	{
		// Apply authentication middleware to all subject routes
		subjects.Use(middleware.AuthRequired())
		
		// Apply role-based access control
		subjects.Use(middleware.RolesRequired(middleware.RoleTeacher, middleware.RoleAdmin))

		// CRUD operations for subjects (commented out until handlers are implemented)
		// subjects.POST("", handlers.CreateSubject)
		// subjects.GET("", handlers.GetAllSubjects)
		// subjects.GET("/:id", handlers.GetSubjectByID)
		// subjects.PUT("/:id", handlers.UpdateSubject)
		// subjects.DELETE("/:id", middleware.RoleRequired(middleware.RoleAdmin), handlers.DeleteSubject)
	}
}

// SetupMaterialRoutes sets up the material routes
func SetupMaterialRoutes(r *gin.Engine) {
	materials := r.Group("/api/materials")
	{
		// Apply authentication middleware to all material routes
		materials.Use(middleware.AuthRequired())
		
		// Create material (file upload) - teachers and admins only (commented out until handlers are implemented)
		// materials.POST("", middleware.RolesRequired(middleware.RoleTeacher, middleware.RoleAdmin), handlers.CreateMaterial)

		// Get materials - any authenticated user (commented out until handlers are implemented)
		// materials.GET("", handlers.GetAllMaterials)

		// Get material by ID - any authenticated user (commented out until handlers are implemented)
		// materials.GET("/:id", handlers.GetMaterialByID)

		// Update material - teachers and admins only (commented out until handlers are implemented)
		// materials.PUT("/:id", middleware.RolesRequired(middleware.RoleTeacher, middleware.RoleAdmin), handlers.UpdateMaterial)

		// Delete material - teachers and admins only (commented out until handlers are implemented)
		// materials.DELETE("/:id", middleware.RolesRequired(middleware.RoleTeacher, middleware.RoleAdmin), handlers.DeleteMaterial)

		// Download material - any authenticated user (commented out until handlers are implemented)
		// materials.GET("/:id/download", handlers.DownloadMaterial)
	}
}

// SetupHealthRoutes sets up the health check routes
func SetupHealthRoutes(r *gin.Engine) {
	health := r.Group("/api/health")
	{
		// Basic health check endpoint - fast response for load balancers (commented out until handlers are implemented)
		// health.GET("", handlers.HealthCheckBasic)

		// Standard health check endpoint - comprehensive health status (commented out until handlers are implemented)
		// health.GET("/check", handlers.HealthCheck)

		// Detailed health check endpoint - includes system information (commented out until handlers are implemented)
		// health.GET("/detailed", handlers.HealthCheckDetailed)

		// Database-specific health check endpoint (commented out until handlers are implemented)
		// health.GET("/database", handlers.HealthCheckDatabase)

		// Legacy health check endpoint for backward compatibility
		health.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"message": "pong",
			})
		})
	}
}