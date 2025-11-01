package attendence

import (
	"net/http"
	"server/database"
	"server/database/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAttendance(c *gin.Context) {
	userId := c.Param("id")
	offset := c.Param("offset")

	if userId == "" || offset == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}
	var attendance []models.Attendance
	database.DB.Preload("User").Where("user_id = ?", userId).Find(&attendance).Offset(offsetInt)
	c.JSON(http.StatusOK, gin.H{"message": "Attendance fetched successfully", "attendance": attendance})
}

func GetAttendanceByOrganisation(c *gin.Context) {
	var attendance []models.Attendance
	organisationId := c.Param("id")
	offset := c.Param("offset")

	if organisationId == "" || offset == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organisation ID is required"})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}
	database.DB.Preload("User").Where("organisation_id = ?", organisationId).Find(&attendance).Offset(offsetInt)
	c.JSON(http.StatusOK, gin.H{"message": "Attendance fetched successfully", "attendance": attendance})
}

func GetAttendanceByUser(c *gin.Context) {
	var attendance []models.Attendance
	userId := c.Param("id")
	offset := c.Param("offset")
	if userId == "" || offset == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	database.DB.Preload("User").Where("user_id = ?", userId).Find(&attendance).Offset(offsetInt)
	c.JSON(http.StatusOK, gin.H{"message": "Attendance fetched successfully", "attendance": attendance})
}

func GetAttendanceByDateAndUser(c *gin.Context) {
	var attendance models.Attendance
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	if err := database.DB.Preload("User").Where("user_id = ? AND date = ?", userId, time.Now().Format("2006-01-02")	).Find(&attendance).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance fetched successfully", "attendance": attendance})
}
