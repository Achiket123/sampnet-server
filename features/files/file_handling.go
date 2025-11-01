package files

import (
	
	"io"
	"net/http"
	"server/database"
	"server/database/models"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	var filemodel models.File

	// Attempt to bind JSON and check for errors
	file, _, err := c.Request.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// Read the image data
	imageData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	filemodel.FileName = c.PostForm("file_name")
	filemodel.FileType = c.PostForm("file_type")
	filemodel.FileSize = int64(len(imageData))
	filemodel.Data = imageData
	// Attempt to save file to database and handle any errors
	if err := database.DB.Create(&filemodel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "message": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_id": filemodel.ID})
}

// GetFile retrieves a file from the database and sends it as a response
func GetFile(c *gin.Context) {
	// Get the file ID from the URL parameter
	fileID := c.Param("id")

	// Initialize a File model
	var file models.File

	// Attempt to retrieve the file from the database
	if err := database.DB.First(&file, fileID).Error; err != nil {
		// If the file is not found, return a 404 error
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	// Set the Content-Type header to match the file type
	c.Header("Content-Type", file.FileType)

	// Send the file data as the response body
	c.Data(http.StatusOK, file.FileType, file.Data)
}
