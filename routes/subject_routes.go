package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/aprianimmanuel/rangkaiedu-backend/controllers"
	"github.com/aprianimmanuel/rangkaiedu-backend/middleware"
)

// SetupSubjectRoutes sets up the subject management routes
func SetupSubjectRoutes(r *gin.Engine) {
	subjects := r.Group("/api/subjects")
	{
		// Apply authentication middleware to all subject routes
		subjects.Use(middleware.AuthRequired())
		
		// Apply role-based access control
		subjects.Use(middleware.RolesRequired(middleware.RoleTeacher, middleware.RoleAdmin))

		// CRUD operations for subjects
		subjects.POST("", controllers.CreateSubject)
		subjects.GET("", controllers.GetAllSubjects)
		subjects.GET("/:id", controllers.GetSubjectByID)
		subjects.PUT("/:id", controllers.UpdateSubject)
		subjects.DELETE("/:id", middleware.RoleRequired(middleware.RoleAdmin), controllers.DeleteSubject)
	}
}