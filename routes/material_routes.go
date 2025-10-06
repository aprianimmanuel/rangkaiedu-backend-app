package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/aprianimmanuel/rangkaiedu-backend/controllers"
	"github.com/aprianimmanuel/rangkaiedu-backend/middleware"
)

// SetupMaterialRoutes sets up the material management routes
func SetupMaterialRoutes(r *gin.Engine) {
	materials := r.Group("/api/materials")
	{
		// Apply authentication middleware to all material routes
		materials.Use(middleware.AuthRequired())

		// Create material (file upload) - teachers and admins only
		materials.POST("", middleware.RolesRequired(middleware.RoleTeacher, middleware.RoleAdmin), controllers.CreateMaterial)

		// Get materials - any authenticated user
		materials.GET("", controllers.GetAllMaterials)

		// Get material by ID - any authenticated user
		materials.GET("/:id", controllers.GetMaterialByID)

		// Update material - teachers and admins only
		materials.PUT("/:id", middleware.RolesRequired(middleware.RoleTeacher, middleware.RoleAdmin), controllers.UpdateMaterial)

		// Delete material - teachers and admins only
		materials.DELETE("/:id", middleware.RolesRequired(middleware.RoleTeacher, middleware.RoleAdmin), controllers.DeleteMaterial)

		// Download material - any authenticated user
		materials.GET("/:id/download", controllers.DownloadMaterial)
	}
}