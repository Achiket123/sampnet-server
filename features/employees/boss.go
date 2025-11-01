package employees

import (
	"fmt"
	"net/http"
	"server/database"
	"server/database/models"

	"github.com/gin-gonic/gin"
)

func CreateBoss(c *gin.Context) {
	var boss models.Boss
	if err := c.ShouldBindJSON(&boss); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&boss).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Boss created successfully", "boss": boss})
}

func GetBoss(c *gin.Context) {
	var boss models.Boss
	if err := database.DB.Preload("User").Preload("Organisation").First(&boss, "user_id = ?", c.GetUint("user_id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"boss": boss})

}

func CreateEmployee(c *gin.Context) {
	var newEmployee models.Employee

	if err := c.ShouldBindJSON(&newEmployee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(newEmployee.UserID, newEmployee.OrganisationID, newEmployee.EmploymentID)
	// Check if user exists
	var user models.UserModel
	if err := database.DB.First(&user, newEmployee.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if organisation exists
	var organisation models.Organisation
	if err := database.DB.First(&organisation, newEmployee.OrganisationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organisation not found"})
		return
	}

	if err := database.DB.Create(&newEmployee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
		return
	}

	// Load the relationships
	if err := database.DB.Preload("User").Preload("Organisation").Where("user_id = ?",newEmployee.UserID).Find(&newEmployee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load employee relationships"})
		return
	}
	fmt.Println(newEmployee.UserID, newEmployee.OrganisationID, newEmployee.EmploymentID)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Employee created successfully",
		"employee": newEmployee,
	})
}
