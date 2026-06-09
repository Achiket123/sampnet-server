package resources

import (
	"net/http"
	"strconv"

	"server/internal/domain/resources"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc resources.UseCase
}

func NewHandler(uc resources.UseCase) *Handler {
	return &Handler{uc: uc}
}

func getOrgAndUser(c *gin.Context) (uint, uint, bool) {
	userIDVal, existsUser := c.Get("userID")
	orgIDVal, existsOrg := c.Get("organisationID")
	if !existsUser || !existsOrg {
		return 0, 0, false
	}
	userID, okUser := userIDVal.(uint)
	orgID, okOrg := orgIDVal.(uint)
	if !okUser || !okOrg {
		return 0, 0, false
	}
	return orgID, userID, true
}

func handleValidationError(c *gin.Context, err error) bool {
	if valErr, ok := err.(*resources.ValidationError); ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "validation_failed",
			"details": valErr.Errors,
		})
		return true
	}
	return false
}

// Collections Handlers

func (h *Handler) CreateCollection(c *gin.Context) {
	orgID, userID, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Name        string                      `json:"name"`
		Description string                      `json:"description"`
		Icon        *string                     `json:"icon"`
		Colour      *string                     `json:"colour"`
		Fields      []resources.FieldDefinition `json:"fields"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	coll, err := h.uc.CreateCollection(c.Request.Context(), orgID, userID, req.Name, req.Description, req.Icon, req.Colour, req.Fields)
	if err != nil {
		if handleValidationError(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, coll)
}

func (h *Handler) GetCollection(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	coll, err := h.uc.GetCollection(c.Request.Context(), uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	c.JSON(http.StatusOK, coll)
}

func (h *Handler) ListCollections(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	colls, total, err := h.uc.ListCollections(c.Request.Context(), orgID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"collections": colls,
		"total":       total,
	})
}

func (h *Handler) UpdateCollection(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Icon        *string `json:"icon"`
		Colour      *string `json:"colour"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	coll, err := h.uc.UpdateCollection(c.Request.Context(), uint(id), orgID, req.Name, req.Description, req.Icon, req.Colour)
	if err != nil {
		if handleValidationError(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coll)
}

func (h *Handler) DeleteCollection(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	forceStr := c.DefaultQuery("force", "false")
	force := forceStr == "true"

	if err := h.uc.DeleteCollection(c.Request.Context(), uint(id), orgID, force); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collection and associated records deleted successfully"})
}

func (h *Handler) AddField(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	var req resources.FieldDefinition
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	coll, err := h.uc.AddFieldToCollection(c.Request.Context(), uint(id), orgID, req)
	if err != nil {
		if handleValidationError(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coll)
}

func (h *Handler) UpdateField(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	fieldKey := c.Param("key")

	var req resources.FieldDefinition
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	coll, err := h.uc.UpdateFieldInCollection(c.Request.Context(), uint(id), orgID, fieldKey, req)
	if err != nil {
		if handleValidationError(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coll)
}

func (h *Handler) RemoveField(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	fieldKey := c.Param("key")

	coll, warning, err := h.uc.RemoveFieldFromCollection(c.Request.Context(), uint(id), orgID, fieldKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"collection": coll,
		"warning":    warning,
	})
}

// Records Handlers

func (h *Handler) CreateRecord(c *gin.Context) {
	orgID, userID, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	collIDStr := c.Param("id")
	collID, err := strconv.ParseUint(collIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	var req struct {
		Data      map[string]interface{} `json:"data"`
		ProjectID *uint                  `json:"project_id"`
		TeamID    *uint                  `json:"team_id"`
		TaskID    *uint                  `json:"task_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	record, err := h.uc.CreateRecord(c.Request.Context(), orgID, uint(collID), userID, req.Data, req.ProjectID, req.TeamID, req.TaskID)
	if err != nil {
		if handleValidationError(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

func (h *Handler) GetRecord(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	collIDStr := c.Param("id")
	collID, err := strconv.ParseUint(collIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	recordIDStr := c.Param("record_id")
	recordID, err := strconv.ParseUint(recordIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	record, err := h.uc.GetRecord(c.Request.Context(), uint(recordID), uint(collID), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func (h *Handler) ListRecords(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	collIDStr := c.Param("id")
	collID, err := strconv.ParseUint(collIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	search := c.Query("search")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// All other query params are treated as exact match JSONB containment filters
	filters := make(map[string]interface{})
	for k, v := range c.Request.URL.Query() {
		if k == "search" || k == "sort_by" || k == "sort_order" || k == "limit" || k == "offset" {
			continue
		}
		if len(v) > 0 {
			filters[k] = v[0]
		}
	}

	recordFilters := resources.RecordFilters{
		Search:    search,
		Filters:   filters,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Offset:    offset,
		Limit:     limit,
	}

	recs, total, err := h.uc.ListRecords(c.Request.Context(), uint(collID), orgID, recordFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"records": recs,
		"total":   total,
	})
}

func (h *Handler) UpdateRecord(c *gin.Context) {
	orgID, userID, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	collIDStr := c.Param("id")
	collID, err := strconv.ParseUint(collIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	recordIDStr := c.Param("record_id")
	recordID, err := strconv.ParseUint(recordIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	var req struct {
		Data map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	record, err := h.uc.UpdateRecord(c.Request.Context(), uint(recordID), uint(collID), orgID, userID, req.Data)
	if err != nil {
		if handleValidationError(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

func (h *Handler) DeleteRecord(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	collIDStr := c.Param("id")
	collID, err := strconv.ParseUint(collIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	recordIDStr := c.Param("record_id")
	recordID, err := strconv.ParseUint(recordIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	if err := h.uc.DeleteRecord(c.Request.Context(), uint(recordID), uint(collID), orgID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

func (h *Handler) BulkCreate(c *gin.Context) {
	orgID, userID, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	collIDStr := c.Param("id")
	collID, err := strconv.ParseUint(collIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	res, err := h.uc.BulkCreateRecords(c.Request.Context(), orgID, uint(collID), userID, req["records"].([]map[string]interface{}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Export(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	collIDStr := c.Param("id")
	collID, err := strconv.ParseUint(collIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	exported, headers, err := h.uc.ExportRecords(c.Request.Context(), uint(collID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"headers": headers,
		"records": exported,
	})
}
