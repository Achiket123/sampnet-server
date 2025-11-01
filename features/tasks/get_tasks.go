package tasks

import (
	"net/http"
	"server/database"
	"server/database/models"
	"strconv"

	"github.com/gin-gonic/gin"
)


func GetTaskByID(c *gin.Context) {
	var task models.Task

	if err := database.DB.Preload("AssignedUser").Preload("AssignedByUser").Preload("Organisation").First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}


// GetTeamTasks retrieves tasks assigned to the current user's team.
// It returns a JSON response with the list of tasks or an error message.
func GetTeamTasks(c *gin.Context) {
	var tasks []models.Task

	if err := database.DB.Preload("AssignedUser").Preload("AssignedByUser").Preload("Organisation").Where("assigned_to = ?", c.GetUint("user_id")).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// GetProjectTasks fetches all tasks associated with a specific project.
// It returns a JSON response with the list of tasks or an error message.
func GetProjectTasks(c *gin.Context) {
	var tasks []models.Task
	page, _ := strconv.Atoi(c.Param("page"))
	if err := database.DB.Preload("AssignedUser").Preload("AssignedByUser").Preload("Organisation").Where("project_id = ?", c.GetUint("project_id")).Find(&tasks).Limit(20).Offset((page - 1) * 20).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}	

// GetPersonalTasks retrieves tasks specifically assigned to the current user.
// It returns a JSON response with the list of personal tasks or an error message.
func GetPersonalTasks(c *gin.Context) {
	var tasks []models.Task
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	offset := (page - 1) * pageSize

	query := database.DB.Preload("AssignedUser").Preload("AssignedByUser").Preload("Organisation").Where("assigned_to = ?", c.GetUint("user_id"))

	var totalCount int64
	if err := query.Model(&models.Task{}).Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count tasks"})
		return
	}

	if err := query.Limit(pageSize).Offset(offset).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks":       tasks,
		"total_count": totalCount,
		"page":        page,
		"page_size":   pageSize,
	})
}

// GetOrganisationTasks fetches all tasks associated with the user's organization.
// It returns a JSON response with the list of organization tasks or an error message.
func GetOrganisationTasks(c *gin.Context) {
	var tasks []models.Task
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	orgId := c.Param("organisation_id");
	offset := (page - 1) * pageSize
	// Validate organisation ID
	if orgId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organisation ID is required"})
		return
	}
	query := database.DB.Preload("AssignedUser").Preload("AssignedByUser").Preload("Organisation").Where("organisation_id = ?",orgId)

	var totalCount int64
	if err := query.Model(&models.Task{}).Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count tasks"})
		return
	}

	if err := query.Limit(pageSize).Offset(offset).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks":       tasks,
		"total_count": totalCount,
		"page":        page,
		"page_size":   pageSize,
	})
}

func GetTaskByTitle(c *gin.Context) {
	var task []models.Task

	if err := database.DB.Preload("AssignedUser").Preload("AssignedByUser").Preload("Organisation").Where("title = ?", c.Param("title")).Find(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}
