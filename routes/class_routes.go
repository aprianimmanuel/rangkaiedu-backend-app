package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/aprianimmanuel/backend-app/controllers"
	"github.com/aprianimmanuel/backend-app/middleware"
)

// SetupClassRoutes sets up the class management routes
func SetupClassRoutes(r *gin.Engine) {
	classes := r.Group("/api/classes")
	{
		// Apply authentication middleware to all class routes
		classes.Use(middleware.AuthRequired())
		
		// Apply role-based access control
		classes.Use(middleware.RolesRequired(middleware.RoleTeacher, middleware.RoleAdmin))

		// CRUD operations for classes
		classes.POST("", controllers.CreateClass)
		classes.GET("", controllers.GetAllClasses)
		classes.GET("/:id", controllers.GetClassByID)
		classes.PUT("/:id", controllers.UpdateClass)
		classes.DELETE("/:id", middleware.RoleRequired(middleware.RoleAdmin), controllers.DeleteClass)

		// Class roster management
		classes.POST("/:class_id/students", controllers.AddStudentToClass)
		classes.GET("/:class_id/students", controllers.GetClassRoster)
		classes.DELETE("/:class_id/students/:student_id", controllers.RemoveStudentFromClass)
		classes.PUT("/:class_id/students/:student_id", controllers.UpdateStudentEnrollmentStatus)
	}
}