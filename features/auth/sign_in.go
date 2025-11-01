package auth

import (
	"log"
	"net/http"
	"server/database"
	"server/database/models"
	"server/miscallenous"

	"github.com/gin-gonic/gin"
)

// SignIn handles the user login process.
func SignIn(c *gin.Context) {
	var user models.UserModel
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	database.DB.Where("email = ?", email).Find(&user)
	log.Println(password, email)
	if !miscallenous.VerifyPassword(user.HashedPassword, password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	token, err := miscallenous.GenerateJWTToken(user, "user", user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.Header("Authorization", token)
	c.JSON(http.StatusOK, gin.H{"message": "Sign in successful"})
}
