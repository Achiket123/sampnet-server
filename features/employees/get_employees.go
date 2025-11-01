package employees

import (
	"net/http"
	"server/database"
	"server/database/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetEmployees(c *gin.Context) {
	var users []models.Employee
	organisation_id, _ := strconv.Atoi(c.Param("organisation_id"))
	if organisation_id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organisation ID is required"})
		return
	}
	database.DB.Preload("User").Preload("Organisation").Where("organisation_id = ?", organisation_id).Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetEmployeeByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.Employee
	err := database.DB.Preload("User").Preload("Organisation").First(&user, id).Error
	if err != nil {	
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UpdateEmployee(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.Employee
	c.BindJSON(&user)
	database.DB.Model(&user).Where("id = ?", id).Updates(user)
	c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully"})
}

func DeleteEmployee(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.Employee
	database.DB.Model(&user).Where("id = ?", id).Update("deleted_at", time.Now())
	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}

func SearchEmployee(c *gin.Context) {
	query := c.Query("query")
	var users []models.Employee
	database.DB.Where("first_name LIKE ? OR last_name LIKE ?", "%"+query+"%", "%"+query+"%").Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}
