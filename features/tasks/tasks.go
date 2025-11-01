// Package tasks provides functionality for managing tasks in the application.
package tasks

import (
	"net/http"
	"server/database"
	"server/database/models"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateTask handles the creation of a new task.
// It binds the JSON request to a Task model, saves it to the database,
// and returns the created task with its related entities.
func CreateTask(c *gin.Context) {
	var task models.Task

	// Bind JSON request to task model
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	// Save task to database
	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "message": err.Error()})
		return
	}

	// Fetch the created task with related entities
	database.DB.Preload("AssignedUser").Preload("AssignedByUser").Preload("Organisation").First(&task)
	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "task": task})
}

// UpdateTask handles the updating of an existing task.
// It binds the JSON request to a Task model, updates it in the database,
// and returns the updated task with its related entities.
func UpdateTask(c *gin.Context) {
	var task models.Task

	// Bind JSON request to task model
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Find existing task
	var existingTask models.Task
	if err := database.DB.First(&existingTask, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Update task in database
	if err := database.DB.Model(&existingTask).Updates(task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Fetch the updated task with related entities
	if err := database.DB.Preload("AssignedUser").Preload("AssignedByUser").
		Preload("Organisation").Preload("Team").Preload("Project").
		First(&task, existingTask.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"task":    task,
	})
}

// DeleteTask handles the deletion of a task by its ID.
// It attempts to delete the task from the database and returns a success message if successful.
func SoftDeleteTask(c *gin.Context) {
	var task models.Task

	// Find the task by ID
	if err := database.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Perform soft delete
	if err := database.DB.Model(&task).Update("deleted_at", time.Now()).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task soft deleted successfully"})
}
