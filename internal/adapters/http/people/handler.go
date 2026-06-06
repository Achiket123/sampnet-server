package people

import (
	"net/http"
	"strconv"

	"server/internal/domain/people"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc people.UseCase
}

func NewHandler(uc people.UseCase) *Handler {
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

func parseFilter(c *gin.Context) people.ContactsFilter {
	filter := people.ContactsFilter{}
	if stage := c.Query("stage"); stage != "" {
		filter.Stage = &stage
	}
	if t := c.Query("type"); t != "" {
		filter.Type = &t
	}
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}
	if listIDStr := c.Query("list_id"); listIDStr != "" {
		if listID, err := strconv.Atoi(listIDStr); err == nil {
			filter.ListID = &listID
		}
	}
	filter.Search = c.Query("search")
	filter.SortBy = c.DefaultQuery("sort_by", "created_at")
	filter.SortOrder = c.DefaultQuery("sort_order", "desc")
	
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 {
		limit = 20
	}
	filter.Limit = limit

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	filter.Page = page

	return filter
}

func (h *Handler) GetContacts(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	filter := parseFilter(c)
	res, err := h.uc.GetContacts(c.Request.Context(), filter, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetContactByID(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	res, err := h.uc.GetContactByID(c.Request.Context(), uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateContact(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req people.CreateContactParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	if req.FirstName == "" && req.LastName == "" && (req.Company == nil || *req.Company == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Must provide either name or company"})
		return
	}

	res, err := h.uc.CreateContact(c.Request.Context(), &req, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *Handler) UpdateContact(c *gin.Context) {
	orgID, userID, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	var req people.UpdateContactParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	res, err := h.uc.UpdateContact(c.Request.Context(), uint(id), &req, orgID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) DeleteContact(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	if err := h.uc.DeleteContact(c.Request.Context(), uint(id), orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}

func (h *Handler) BulkUpdateStage(c *gin.Context) {
	orgID, userID, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		IDs   []uint `json:"ids"`
		Stage string `json:"stage"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	if err := h.uc.BulkUpdateStage(c.Request.Context(), req.IDs, req.Stage, orgID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contacts updated successfully"})
}

func (h *Handler) AddInteraction(c *gin.Context) {
	orgID, userID, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	var req people.AddInteractionParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	res, err := h.uc.AddInteraction(c.Request.Context(), uint(id), &req, orgID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *Handler) GetAnalytics(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	res, err := h.uc.GetAnalytics(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Lists handlers
func (h *Handler) GetLists(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	res, err := h.uc.GetLists(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lists": res})
}

func (h *Handler) CreateList(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req people.CreateListParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	res, err := h.uc.CreateList(c.Request.Context(), &req, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Pipeline handlers
func (h *Handler) GetPipelineStages(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	res, err := h.uc.GetPipelineStages(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stages": res})
}

func (h *Handler) CreatePipelineStage(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req people.CreateStageParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	res, err := h.uc.CreatePipelineStage(c.Request.Context(), &req, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *Handler) ReorderStages(c *gin.Context) {
	orgID, _, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		OrderedIDs []uint `json:"ordered_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.uc.ReorderPipelineStages(c.Request.Context(), req.OrderedIDs, orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stages reordered successfully"})
}