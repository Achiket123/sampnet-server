package auth

import (
	"net/http"
	"server/database"
	"server/database/models"
	"server/miscallenous"
	"time"

	"github.com/gin-gonic/gin"
)

// SignUp handles the user registration process.
// It creates a new user in the database and generates a JWT token for the user.
func SignUp(c *gin.Context) {
	var user models.UserModel
	currentTime := time.Now()

	// Hash the password
	password := c.Request.FormValue("password")
	hashedPassword, err := miscallenous.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Populate user model
	user = models.UserModel{
		FirstName:      c.Request.FormValue("first_name"),
		LastName:       c.Request.FormValue("last_name"),
		Email:          c.Request.FormValue("email"),
		City:           c.Request.FormValue("city"),
		Country:        c.Request.FormValue("country"),
		LastLoginAt:    currentTime,
		PhoneNumber:    c.Request.FormValue("phone_number"),
		HashedPassword: hashedPassword,
		ProfilePic:     c.Request.FormValue("profile_pic"),
	}

	// Create the user in the database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "message": err.Error()})
		return
	}

	// Generate a JWT token for the new user
	token, err := miscallenous.GenerateJWTToken(user, "user", user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.Header("Authorization", token)
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}
