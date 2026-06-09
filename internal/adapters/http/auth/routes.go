package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	r.POST("/api/v1/auth/signup", h.SignUp)
	r.POST("/api/v1/auth/signin", h.SignIn)
	r.GET("/api/v1/auth/verify-email", h.VerifyEmail)
	r.POST("/api/v1/auth/refresh", h.RefreshToken)
	
	authGroup := r.Group("/api/v1/auth")
	authGroup.Use(validateToken)
	{
		authGroup.POST("/complete-signin", h.CompleteSignIn)
		authGroup.GET("/validate-employee", h.ValidateEmployee)
		authGroup.POST("/send-verification", h.SendVerificationEmail)
		authGroup.GET("/me", h.GetMe)
	}
}
