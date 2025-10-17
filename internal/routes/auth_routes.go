package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/handlers"
)

// SetupAuthRoutes sets up the authentication routes.
func SetupAuthRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {
	auth := r.Group("/api/auth")
	{
		// Legacy routes (kept for backward compatibility)
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/send-otp", authHandler.SendOTP)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
		
		// WhatsApp OTP routes
		auth.POST("/whatsapp/login", authHandler.LoginWithWhatsApp)
		auth.POST("/whatsapp/send-otp", authHandler.SendOTP)
		auth.POST("/whatsapp/verify-otp", authHandler.VerifyOTP)
		
		// New unified authentication routes
		// These would be implemented in other handlers
		// auth.POST("/google", authHandler.GoogleAuth)
		// auth.POST("/facebook", authHandler.FacebookAuth)
		// auth.POST("/whatsapp", authHandler.WhatsAppAuth)
		// auth.POST("/whatsapp/verify", authHandler.WhatsAppVerify)
		// auth.POST("/email", authHandler.EmailAuth)
		// auth.POST("/email/verify", authHandler.EmailVerify)
		
		// MFA routes
		// These would be implemented in MFA handlers
		// auth.POST("/mfa/setup", authHandler.MFASetup)
		// auth.POST("/mfa/verify", authHandler.MFAVerify)
	}
}