package organisations

import (
	"net/http"
	"server/database"
	"server/database/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RegisterOrganisation handles the registration of a new organization
func RegisterOrganisation(c *gin.Context) {
	var organisation models.Organisation

	// Bind the JSON request body to the organisation struct
	if err := c.BindJSON(&organisation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Attempt to create the organisation in the database
	if err := database.DB.Create(&organisation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organisation"})
		return
	}

	// Respond with success message and the created organisation
	c.JSON(http.StatusCreated, gin.H{
		"message":      "Organisation created successfully",
		"organisation": organisation,
	})
}

func GetOrganisation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var organisation models.Organisation
	err := database.DB.First(&organisation, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organisation not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"organisation": organisation})
}

func UpdateOrganisation(c *gin.Context) {
	
	id, _ := strconv.Atoi(c.Param("id"))
	var organisation models.Organisation
	err := database.DB.First(&organisation, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organisation not found"})
		return
	}
	if err := c.BindJSON(&organisation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := database.DB.Save(&organisation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organisation"})
		return
	}	
	c.JSON(http.StatusOK, gin.H{"organisation": organisation})
}

