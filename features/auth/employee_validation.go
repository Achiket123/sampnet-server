package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/database"
	"server/database/models"
	"server/miscallenous"

	"github.com/gin-gonic/gin"
)

func ValidateEmployee(c *gin.Context) {
	token := c.GetHeader("Authorization")
	var employee models.Employee
	tokenData, err := miscallenous.DecodeJWTToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": err.Error()})
		return
	}
	claimsJSON, _ := json.MarshalIndent(tokenData.Claims, "", "  ")

	err = json.Unmarshal(claimsJSON, &employee)
	log.Println(employee.User.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": err.Error()})
		return
	}
	var employeeData models.Employee
	if err2 := database.DB.Preload("User").Preload("Organisation").Where("user_id = ?", employee.User.ID).Find(&employeeData).Error; err2 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid employee data", "message": err2.Error()})
		return
	}
	fmt.Println(employeeData)
	_token, _err := miscallenous.GenerateJWTToken(employeeData, "employee", employeeData.UserID)
	if _err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.Header("Authorization", _token)
	c.JSON(http.StatusOK, gin.H{"message": "Employee validated successfully"})
}
