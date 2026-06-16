package invites

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	domain "server/internal/domain/invites"
)

type Handler struct {
	uc domain.UseCase
}

func NewHandler(uc domain.UseCase) *Handler {
	return &Handler{uc: uc}
}

func getUserAuthInfo(c *gin.Context) (uint, uint, string, bool) {
	userIDVal, existsUser := c.Get("userID")
	orgIDVal, existsOrg := c.Get("organisationID")
	roleVal, existsRole := c.Get("role")
	if !existsUser || !existsOrg || !existsRole {
		return 0, 0, "", false
	}
	userID, okUser := userIDVal.(uint)
	orgID, okOrg := orgIDVal.(uint)
	role, okRole := roleVal.(string)
	if !okUser || !okOrg || !okRole {
		return 0, 0, "", false
	}
	return userID, orgID, role, true
}

func (h *Handler) InviteEmployee(c *gin.Context) {
	userID, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if role != "boss" && role != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
		return
	}

	var req struct {
		Email        string `json:"email" binding:"required,email"`
		FirstName    string `json:"first_name" binding:"required"`
		LastName     string `json:"last_name" binding:"required"`
		PhoneNumber  string `json:"phone_number"`
		EmploymentID int    `json:"employment_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invite := &domain.EmployeeInvite{
		Email:           req.Email,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		PhoneNumber:     req.PhoneNumber,
		EmploymentID:    req.EmploymentID,
		OrganisationID:  orgID,
		InvitedByUserID: userID,
	}

	if err := h.uc.InviteEmployee(c.Request.Context(), invite); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent successfully"})
}

func (h *Handler) AcceptInvite(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pair, err := h.uc.AcceptInvite(c.Request.Context(), req.Token, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Authorization", pair.AccessToken)
	c.JSON(http.StatusOK, gin.H{
		"message":       "Invitation accepted successfully",
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
	})
}

func (h *Handler) GetInvites(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if role != "boss" && role != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
		return
	}

	invites, err := h.uc.GetInvitesByOrg(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invites": invites})
}

func (h *Handler) ResendInvite(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if role != "boss" && role != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
		return
	}

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.ResendInvite(c.Request.Context(), req.Email, orgID); err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": errMsg})
			return
		}
		if strings.Contains(errMsg, "cannot resend") {
			c.JSON(http.StatusConflict, gin.H{"error": errMsg})
			return
		}
		if strings.Contains(errMsg, "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation email resent successfully"})
}