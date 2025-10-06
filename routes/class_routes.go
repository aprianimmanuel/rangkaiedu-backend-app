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

		// Class roster management (must come before single class routes to avoid conflicts)
		classes.POST("/:class_id/students", controllers.AddStudentToClass)
		classes.GET("/:class_id/students", controllers.GetClassRoster)
		classes.DELETE("/:class_id/students/:student_id", controllers.RemoveStudentFromClass)
		classes.PUT("/:class_id/students/:student_id", controllers.UpdateStudentEnrollmentStatus)

		// CRUD operations for classes
		classes.POST("", controllers.CreateClass)
		classes.GET("", controllers.GetAllClasses)
		classes.GET("/by-id/:classId", controllers.GetClassByID)
		classes.PUT("/by-id/:classId", controllers.UpdateClass)
		classes.DELETE("/by-id/:classId", middleware.RoleRequired(middleware.RoleAdmin), controllers.DeleteClass)
	}
}