package employees

import (
	"net/http"
	"server/database"
	"server/database/models"
	"server/miscallenous"

	"github.com/gin-gonic/gin"
)

func AddEmployees(c *gin.Context) {
	var employee struct {
		EmploymentID int    `json:"employment_id" binding:"required"`
		FirstName    string `json:"first_name" binding:"required"`
		LastName     string `json:"last_name" binding:"required"`
		Email        string `json:"email" binding:"required,email"`
		PhoneNumber  string `json:"phone_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.UserModel{
		FirstName:   employee.FirstName,
		LastName:    employee.LastName,
		Email:       employee.Email,
		PhoneNumber: employee.PhoneNumber,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user","error":err.Error()})
		return
	}

	newEmployee := models.Employee{
		UserID:       user.ID,
		EmploymentID: employee.EmploymentID,
	}

	if err := database.DB.Create(&newEmployee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee created successfully", "employee_id": newEmployee.UserID})
}

func MakeManager(c *gin.Context) {
	var managerData struct {
		UserID         uint   `json:"user_id" binding:"required"`
		OrganisationID uint   `json:"organisation_id" binding:"required"`
		Type           string `json:"type" binding:"required"`
		Salary         string `json:"salary" binding:"required"`
	}

	if err := c.ShouldBindJSON(&managerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists
	var user models.UserModel
	if err := database.DB.First(&user, managerData.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the organisation exists
	var organisation models.Organisation
	if err := database.DB.First(&organisation, managerData.OrganisationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organisation not found"})
		return
	}

	manager := models.Manager{
		UserID:         managerData.UserID,
		OrganisationID: managerData.OrganisationID,
		Type:           managerData.Type,
		Salary:         managerData.Salary,
	}

	if err := database.DB.Create(&manager).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create manager"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager created successfully", "manager_id": manager.UserID})
}

func IsEmployee(c *gin.Context) {
	var employee models.Employee
	var manager models.Manager

	userID := c.Param("user_id")

	// Check for manager first
	isManager := database.DB.Preload("User").Where("user_id = ?", userID).First(&manager).Error == nil
	// Check for employee
	isEmployee := database.DB.Preload("User").Where("user_id = ?", userID).First(&employee).Error == nil

	if !isManager && !isEmployee {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var token string
	var err error

	if isManager {
		// Return manager token if user is a manager (regardless of employee status)
		token, err = miscallenous.GenerateJWTToken(manager, "manager", manager.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Manager found", "manager": manager, "token": token})
		return
	}

	// At this point user is only an employee
	token, err = miscallenous.GenerateJWTToken(employee, "employee", employee.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Employee found", "employee": employee, "token": token})
}
