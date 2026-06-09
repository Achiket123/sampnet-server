package middlewares

import (
	"log"
	"net/http"
	"server/internal/platform/miscallenous"
	"strings"

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
		if strings.Contains(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		// Attempt to decode and validate the token
		decodedToken, err := miscallenous.DecodeJWTToken(token)
		// If there's an error in decoding, respond with Unauthorized
		if err != nil {
			log.Default().Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// If the token is valid, allow the request to proceed
		if decodedToken.Valid {
			log.Default().Println("token is valid")
			if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok {
				if sub, ok := claims["sub"].(float64); ok {
					c.Set("userID", uint(sub))
				}
				if emp, ok := claims["employee"].(map[string]interface{}); ok {
					if orgID, ok := emp["organisation_id"].(float64); ok {
						c.Set("organisationID", uint(orgID))
					}
					role := "employee"
					if t, ok := emp["type"].(string); ok {
						switch t {
						case "owner", "boss":
							role = "boss"
						case "manager":
							role = "manager"
						}
					}
					c.Set("role", role)
				} else if mgr, ok := claims["manager"].(map[string]interface{}); ok {
					if orgID, ok := mgr["organisation_id"].(float64); ok {
						c.Set("organisationID", uint(orgID))
					}
					c.Set("role", "manager")
				} else if boss, ok := claims["boss"].(map[string]interface{}); ok {
					if orgID, ok := boss["organisation_id"].(float64); ok {
						c.Set("organisationID", uint(orgID))
					}
					c.Set("role", "boss")
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

func RoleMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleStr, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: No role found"})
			c.Abort()
			return
		}

		userRole, ok := userRoleStr.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Invalid role type"})
			c.Abort()
			return
		}

		allowed := false
		for _, r := range roles {
			if r == userRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
