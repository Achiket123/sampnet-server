package search

import (
	"net/http"
	"strconv"
	"strings"

	domain "server/internal/domain/search"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the global search feature.
type Handler struct {
	uc domain.UseCase
}

// NewHandler creates a new search Handler.
func NewHandler(uc domain.UseCase) *Handler {
	return &Handler{uc: uc}
}

// Search handles GET /search?q=...&types=...&limit=...&offset=...
func (h *Handler) Search(c *gin.Context) {
	// Authenticate from JWT middleware context values.
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orgIDRaw, exists := c.Get("organisationID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing organisation"})
		return
	}
	orgID, ok := orgIDRaw.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid organisation id"})
		return
	}

	q := c.Query("q")

	// Parse optional comma-separated type filter.
	var types []string
	if raw := c.Query("types"); raw != "" {
		for _, t := range strings.Split(raw, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				types = append(types, t)
			}
		}
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filters := &domain.SearchFilters{
		Query:          q,
		OrganisationID: orgID,
		Types:          types,
		Limit:          limit,
		Offset:         offset,
	}

	results, err := h.uc.Search(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
