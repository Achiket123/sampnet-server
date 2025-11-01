package attendence

import (
	"log"
	"net/http"
	"time"

	"server/database"
	"server/database/models"

	"github.com/gin-gonic/gin"
)

func PostAttendance(c *gin.Context) {
	var attendance models.Attendance
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attendance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Attendance recorded successfully", "attendance": attendance})
}

func UpdateAttendance(c *gin.Context) {
	var attendance models.Attendance
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existingAttendance models.Attendance
	if err := database.DB.Model(&existingAttendance).Where("user_id = ? AND date = ?", attendance.UserID, time.Now().Format("2006-01-02")).First(&existingAttendance).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "PLEASE CHECK IN FIRST"})
		return
	}

	existingAttendance.CheckOutTime = attendance.CheckOutTime
	existingAttendance.CheckOutPhoto = attendance.CheckOutPhoto

	log.Println(existingAttendance)
	if err := database.DB.Save(&existingAttendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attendance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Attendance updated successfully", "attendance": attendance})
}
