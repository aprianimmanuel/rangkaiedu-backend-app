package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aprianimmanuel/backend-app/controllers"
)

// SetupAuthRoutes sets up the authentication routes.
func SetupAuthRoutes(r *gin.Engine) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.RegisterHandler)
		auth.POST("/login", controllers.LoginHandler)
		auth.POST("/send-otp", controllers.SendOTPHandler)
		auth.POST("/verify-otp", controllers.VerifyOTPHandler)
	}
}