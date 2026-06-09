package middlewares

import (
	"log"
	"net/http"
	"server/internal/platform/miscallenous"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ── Parse Authorization header ────────────────────────────────────────
		// Expected formats:
		//   "<accessToken>"
		//   "<accessToken> <employeeToken>"
		//   "Bearer <accessToken>"
		//   "Bearer <accessToken> <employeeToken>"
		rawHeader := c.GetHeader("Authorization")
		rawHeader = strings.TrimPrefix(rawHeader, "Bearer ")
		rawHeader = strings.TrimSpace(rawHeader)

		parts := strings.Fields(rawHeader) // splits on any whitespace, skips empties
		if len(parts) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		accessToken := parts[0]
		empToken := ""
		if len(parts) >= 2 && parts[1] != "null" {
			empToken = parts[1]
		}

		// ── Validate access token (required) ──────────────────────────────────
		decodedToken, err := miscallenous.DecodeJWTToken(accessToken)
		if err != nil || !decodedToken.Valid {
			log.Default().Println("Error: Token: ", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Set userID from access token sub claim
		if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(float64); ok {
				c.Set("userID", uint(sub))
			}

			// Access token may itself embed employee/manager/boss data
			// (this happens when IsEmployeeOrManager generates the token)
			setRoleAndOrgFromClaims(c, claims)
		}

		// ── Validate employee token (optional) ────────────────────────────────
		// A failure here is non-fatal — the access token already authenticated
		// the user. We simply won't have org/role context from the emp token.
		if empToken != "" {
			decodedEmpToken, err := miscallenous.DecodeJWTToken(empToken)
			if err != nil || !decodedEmpToken.Valid {
				// Log it but DO NOT abort — access token is still valid
				log.Default().Println("Warning: employee token invalid, ignoring: ", err)
			} else {
				if claims, ok := decodedEmpToken.Claims.(jwt.MapClaims); ok {
					// Employee token claims take precedence over access token claims
					// for org/role since they are more specific
					setRoleAndOrgFromClaims(c, claims)
				}
			}
		}

		c.Next()
	}
}

// setRoleAndOrgFromClaims reads employee/manager/boss claims from a decoded
// JWT payload and sets "organisationID" and "role" on the gin context.
// All owner/boss variants are normalised to "boss" for consistent handler checks.
func setRoleAndOrgFromClaims(c *gin.Context, claims jwt.MapClaims) {
	if emp, ok := claims["employee"].(map[string]interface{}); ok {
		if orgID, ok := emp["organisation_id"].(float64); ok && orgID != 0 {
			c.Set("organisationID", uint(orgID))
		}
		role := "employee"
		if t, ok := emp["type"].(string); ok {
			role = normaliseRole(t)
		}
		c.Set("role", role)
		return
	}

	if mgr, ok := claims["manager"].(map[string]interface{}); ok {
		if orgID, ok := mgr["organisation_id"].(float64); ok && orgID != 0 {
			c.Set("organisationID", uint(orgID))
		}
		c.Set("role", "manager")
		return
	}

	if boss, ok := claims["boss"].(map[string]interface{}); ok {
		if orgID, ok := boss["organisation_id"].(float64); ok && orgID != 0 {
			c.Set("organisationID", uint(orgID))
		}
		c.Set("role", "boss")
		return
	}
}

// normaliseRole maps all owner/boss type strings to the single canonical
// value "boss" that handlers and RoleMiddleware check against.
func normaliseRole(t string) string {
	switch t {
	case "owner", "boss":
		return "boss"
	case "manager":
		return "manager"
	default:
		return "employee"
	}
}

func RoleMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: No role found"})
			c.Abort()
			return
		}

		userRole, ok := userRoleVal.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Invalid role type"})
			c.Abort()
			return
		}

		for _, r := range roles {
			if r == userRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
		c.Abort()
	}
}
