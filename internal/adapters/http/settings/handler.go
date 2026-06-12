package settings

import (
	"net/http"
	"server/internal/domain/settings"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc settings.UseCase
}

func NewHandler(uc settings.UseCase) *Handler {
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

func (h *Handler) GetOrgSettings(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	org, err := h.uc.GetOrgSettings(orgID, role)
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organisation": org})
}

func (h *Handler) UpdateOrgSettings(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req settings.OrgSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.UpdateOrgSettings(orgID, role, &req)
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "organisation settings updated successfully"})
}

func (h *Handler) DeleteOrganisation(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		ConfirmationCode string `json:"confirmation_code"`
	}
	// Try parsing body, fallback to query param
	_ = c.ShouldBindJSON(&req)
	if req.ConfirmationCode == "" {
		req.ConfirmationCode = c.Query("confirmation_code")
	}

	if req.ConfirmationCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "confirmation code is required"})
		return
	}

	err := h.uc.DeleteOrganisation(orgID, role, req.ConfirmationCode)
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err == settings.ErrInvalidInput {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid confirmation code"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "organisation deleted successfully"})
}

func (h *Handler) GetPlanInfo(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	plan, err := h.uc.GetPlanInfo(orgID, role)
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plan": plan})
}

func (h *Handler) GetRolePermissions(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	perms, err := h.uc.GetRolePermissions(orgID, role)
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role_permissions": perms})
}

func (h *Handler) UpdateRolePermissions(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		RolePermissions []settings.RolePermissionsSet `json:"role_permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.UpdateRolePermissions(orgID, role, req.RolePermissions)
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role permissions updated successfully"})
}

func (h *Handler) GetAttendancePolicy(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	policy, err := h.uc.GetAttendancePolicy(orgID, role)
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attendance_policy": policy})
}

func (h *Handler) UpdateAttendancePolicy(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req settings.AttendancePolicy
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.UpdateAttendancePolicy(orgID, role, &req)
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attendance policy updated successfully"})
}

func (h *Handler) GetLeavePolicies(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	policies, err := h.uc.GetLeavePolicies(orgID, role)
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leave_policies": policies})
}

func (h *Handler) UpdateLeavePolicies(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		LeavePolicies []settings.LeavePolicyEntry `json:"leave_policies"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.UpdateLeavePolicies(orgID, role, req.LeavePolicies)
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "leave policies updated successfully"})
}

func (h *Handler) SaveLeavePolicy(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req settings.LeavePolicyEntry
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.SaveLeavePolicy(orgID, role, &req)
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) DeleteLeavePolicy(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	policyID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid leave policy id"})
		return
	}

	err = h.uc.DeleteLeavePolicy(orgID, role, uint(policyID))
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "leave policy deleted successfully"})
}

func (h *Handler) GetTaskTypes(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	taskTypes, err := h.uc.GetTaskTypes(orgID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task_types": taskTypes})
}

func (h *Handler) SaveTaskType(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req settings.TaskTypeEntry
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.SaveTaskType(orgID, role, &req)
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task_type": req})
}

func (h *Handler) DeleteTaskType(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	typeID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task type id"})
		return
	}

	err = h.uc.DeleteTaskType(orgID, role, uint(typeID))
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task type deleted successfully"})
}

func (h *Handler) GetUserProfile(c *gin.Context) {
	userID, _, _, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	profile, err := h.uc.GetUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *Handler) UpdateUserProfile(c *gin.Context) {
	userID, _, _, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		City        string `json:"city"`
		Country     string `json:"country"`
		DateOfBirth string `json:"date_of_birth"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dob time.Time
	if req.DateOfBirth != "" {
		var err error
		dob, err = time.Parse(time.RFC3339, req.DateOfBirth)
		if err != nil {
			// Fallback to simple date format if parse fails
			dob, err = time.Parse("2006-01-02", req.DateOfBirth)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date of birth format, use YYYY-MM-DD or RFC3339"})
				return
			}
		}
	}

	profile := settings.UserProfile{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		City:        req.City,
		Country:     req.Country,
		DateOfBirth: dob,
	}

	err := h.uc.UpdateUserProfile(userID, &profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user profile updated successfully"})
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	userID, _, _, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.UpdatePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if err == settings.ErrInvalidInput {
			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect current password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

func (h *Handler) GetNotificationPreferences(c *gin.Context) {
	userID, _, _, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	prefs, err := h.uc.GetNotificationPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification_preferences": prefs})
}

func (h *Handler) UpdateNotificationPreferences(c *gin.Context) {
	userID, _, _, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		NotificationPreferences []settings.NotificationPreferenceEntry `json:"notification_preferences"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.UpdateNotificationPreferences(userID, req.NotificationPreferences)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification preferences updated successfully"})
}

func (h *Handler) ExportOrgData(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	data, err := h.uc.ExportOrgData(orgID, role)
	if err != nil {
		if err == settings.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, data)
}
