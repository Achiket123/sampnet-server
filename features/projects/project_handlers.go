package projects

import (
	"net/http"
	"server/database"
	"server/database/models"

	"github.com/gin-gonic/gin"
)

func CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := database.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project created successfully", "project": project})
}

func GetProject(c *gin.Context) {
	var project models.Project
	if err := database.DB.Preload("Team").Preload("CreatedByUser").First(&project, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project fetched successfully", "project": project})
}
func GetProjectsByOrganisation(c *gin.Context) {
	var projects []models.Project
	if err := database.DB.Preload("Team").Preload("CreatedByUser").Where("organisation_id = ?", c.Param("organisation_id")).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Projects fetched successfully", "projects": projects})
}

func UpdateProject(c *gin.Context) {
	var project models.Project
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := database.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully", "project": project})
}

func GetProjectsByTeam(c *gin.Context) {
	var projects []models.Project
	if err := database.DB.Preload("Team").Preload("CreatedByUser").Where("team_id = ?", c.Param("team_id")).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Projects fetched successfully", "projects": projects})

}
func GetProjectsWithLessData(c *gin.Context) {
	var projects []models.Project
	organisationID := c.Param("organisation_id")

	if err := database.DB.Where("organisation_id = ?", organisationID).
		Find(&projects).Select("ID", "Name", "Description").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Projects fetched successfully", "projects": projects})
}

func DeleteProject(c *gin.Context) {

	if err := database.DB.First(&models.Project{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	if err := database.DB.Delete(&models.Project{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
