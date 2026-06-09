package invites

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	r.POST("/api/v1/employees/invite", validateToken, h.InviteEmployee)
	r.POST("/api/v1/invites/accept", h.AcceptInvite) // public
	r.GET("/api/v1/invites", validateToken, h.GetInvites)
}