package middlewares

import (
	"net/http"
	"server/miscallenous" 
	"github.com/gin-gonic/gin"
)

// ValidateToken is a middleware function that checks for a valid JWT token in the request header.
// It decodes the token and verifies its validity before allowing the request to proceed.
// If the token is invalid or missing, it responds with an "Unauthorized" error.
func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		token := c.GetHeader("Authorization")

		// Attempt to decode and validate the token
		decodedToken, err := miscallenous.DecodeJWTToken(token)

		// If there's an error in decoding, respond with Unauthorized
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// If the token is valid, allow the request to proceed
		if decodedToken.Valid {
			c.Next()
			return
		}

		// If we reach here, the token was invalid
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
}
