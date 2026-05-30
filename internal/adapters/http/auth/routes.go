package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	r.POST("/api/v1/auth/signup", h.SignUp)
	r.POST("/api/v1/auth/signin", h.SignIn)
	
	authGroup := r.Group("/api/v1/auth")
	authGroup.Use(validateToken)
	{
		authGroup.POST("/complete-signin", h.CompleteSignIn)
		authGroup.GET("/validate-employee", h.ValidateEmployee)
	}
}
