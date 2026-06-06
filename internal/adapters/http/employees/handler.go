package employees

import (
	"net/http"
	authDomain "server/internal/domain/auth"
	domain "server/internal/domain/employees"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h *Handler) AddEmployee(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if role != "boss" && role != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Only owners (boss) or managers can add employees"})
		return
	}

	var req struct {
		EmploymentID int    `json:"employment_id" binding:"required"`
		FirstName    string `json:"first_name" binding:"required"`
		LastName     string `json:"last_name" binding:"required"`
		Email        string `json:"email" binding:"required,email"`
		PhoneNumber  string `json:"phone_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emp := &domain.Employee{
		EmploymentID:   req.EmploymentID,
		OrganisationID: orgID,
		Email:          req.Email,
		User: authDomain.User{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
		},
	}

	if err := h.uc.AddEmployee(c.Request.Context(), emp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add employee", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee created successfully", "employee_id": emp.UserID})
}

func (h *Handler) MakeManager(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if role != "boss" && role != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Only owners (boss) or managers can make a manager"})
		return
	}

	var req struct {
		UserID         uint   `json:"user_id" binding:"required"`
		OrganisationID uint   `json:"organisation_id" binding:"required"`
		Type           string `json:"type" binding:"required"`
		Salary         string `json:"salary" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.OrganisationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Organisation ID mismatch"})
		return
	}

	manager := &domain.Manager{
		UserID:         req.UserID,
		OrganisationID: orgID,
		Type:           req.Type,
		Salary:         req.Salary,
	}

	if err := h.uc.MakeManager(c.Request.Context(), manager); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager created successfully", "manager_id": manager.UserID})
}

func (h *Handler) IsEmployee(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	msg, data, token, err := h.uc.IsEmployeeOrManager(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msg, "data": data, "token": token})
}

func (h *Handler) CreateBoss(c *gin.Context) {
	var boss domain.Boss
	if err := c.ShouldBindJSON(&boss); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.CreateBoss(c.Request.Context(), &boss); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Boss created successfully", "boss": boss})
}

func (h *Handler) GetBoss(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	boss, err := h.uc.GetBoss(c.Request.Context(), userIDVal.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"boss": boss})
}

func (h *Handler) CreateEmployee(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if role != "boss" && role != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Only owners (boss) or managers can create employees"})
		return
	}

	var emp domain.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emp.OrganisationID = orgID

	if err := h.uc.AddEmployee(c.Request.Context(), &emp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee created successfully", "employee": emp})
}

func (h *Handler) GetEmployees(c *gin.Context) {
	userID, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	paramOrgID, err := strconv.Atoi(c.Param("organisation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID"})
		return
	}

	if uint(paramOrgID) != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Organisation ID mismatch"})
		return
	}

	users, err := h.uc.GetEmployees(c.Request.Context(), uint(paramOrgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if role == "employee" {
		for i := range users {
			if users[i].UserID != userID {
				users[i].Salary = ""
				users[i].User.PhoneNumber = ""
				users[i].User.Email = ""
				users[i].EmploymentID = 0
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *Handler) GetEmployeeByID(c *gin.Context) {
	userID, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	if role == "employee" && uint(id) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Employees can only view their own data"})
		return
	}

	user, err := h.uc.GetEmployee(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	if user.OrganisationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Employee belongs to a different organisation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *Handler) UpdateEmployee(c *gin.Context) {
	userID, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	if role == "employee" && uint(id) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Employees can only update their own data"})
		return
	}

	existing, err := h.uc.GetEmployee(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	if existing.OrganisationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Employee belongs to a different organisation"})
		return
	}

	var emp domain.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emp.UserID = uint(id)
	emp.OrganisationID = orgID

	if err := h.uc.UpdateEmployee(c.Request.Context(), &emp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully"})
}

func (h *Handler) DeleteEmployee(c *gin.Context) {
	_, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if role != "boss" && role != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Only owners (boss) or managers can delete employees"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	emp, err := h.uc.GetEmployee(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	if emp.OrganisationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Employee belongs to a different organisation"})
		return
	}

	if err := h.uc.DeleteEmployee(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}

func (h *Handler) SearchEmployee(c *gin.Context) {
	userID, orgID, role, ok := getUserAuthInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	query := c.Query("query")
	users, err := h.uc.SearchEmployees(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var filtered []domain.Employee
	for _, u := range users {
		if u.OrganisationID == orgID {
			if role == "employee" && u.UserID != userID {
				u.Salary = ""
				u.User.PhoneNumber = ""
				u.User.Email = ""
				u.EmploymentID = 0
			}
			filtered = append(filtered, u)
		}
	}

	c.JSON(http.StatusOK, gin.H{"users": filtered})
}
