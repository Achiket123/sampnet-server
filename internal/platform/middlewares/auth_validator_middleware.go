package middlewares

import (
	"net/http"
	"server/internal/platform/miscallenous" 
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
			if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok {
				if sub, ok := claims["sub"].(float64); ok {
					c.Set("userID", uint(sub))
				}
			}
			c.Next()
			return
		}

		// If we reach here, the token was invalid
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
}
