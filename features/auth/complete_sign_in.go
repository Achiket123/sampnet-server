package auth

import (
	"net/http"
	"server/database"
	"server/database/models"
	"server/miscallenous"
	"time"
	"github.com/gin-gonic/gin"
)

func CompleteSignIn(c *gin.Context) {
	var userModel models.UserModel
	var user struct {
		Email       string `json:"email" binding:"required,email"`
		PhoneNumber int64  `json:"phone_number" binding:"required"`
		Password    string `json:"password" binding:"required"`
		City        string `json:"city" binding:"required"`
		Country     string `json:"country" binding:"required"`
		ProfilePic  string `json:"profile_pic" binding:"required"`
		
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Where("email = ? OR phone_number = ?", user.Email, user.PhoneNumber).Find(&userModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userModel.LastLoginAt = time.Now()
	userModel.City = user.City
	userModel.Country = user.Country
	userModel.ProfilePic = user.ProfilePic
	userModel.IsVerified = true

	hashedPassword, err := miscallenous.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	userModel.HashedPassword = hashedPassword

	if err := database.DB.Save(&userModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	token, err := miscallenous.GenerateJWTToken(userModel,"user",userModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.Header("Authorization", token)
	c.JSON(http.StatusOK, gin.H{"message": "User Updated Successfully"})
}
