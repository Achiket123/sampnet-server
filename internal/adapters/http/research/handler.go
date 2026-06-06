package research

import (
	"fmt"
	"log"
	"net/http"
	domain "server/internal/domain/research"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc domain.UseCase
}

func NewHandler(uc domain.UseCase) *Handler {
	return &Handler{uc: uc}
}

func getAuthContext(c *gin.Context) (uint, uint, string, bool) {
	userIDVal, existsUser := c.Get("userID")
	orgIDVal, existsOrg := c.Get("organisationID")
	roleVal, existsRole := c.Get("role")
	if !existsUser || !existsOrg || !existsRole {
		return 0, 0, "", false
	}
	userID, okUser := userIDVal.(uint)
	orgID, okOrg := orgIDVal.(uint)
	role, okRole := roleVal.(string)
	if !okUser || !okOrg || !okRole {
		return 0, 0, "", false
	}
	return orgID, userID, role, true
}

func (h *Handler) Create(c *gin.Context) {
	orgID, userID, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}

	var req domain.CreateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	req.AuthorID = userID
	req.OrganisationID = orgID

	entry, err := h.uc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Research entry created successfully", "data": entry})
}

func (h *Handler) GetByID(c *gin.Context) {
	orgID, _, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid research entry ID"})
		return
	}

	entry, err := h.uc.GetByID(c.Request.Context(), uint(id), orgID)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Research entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": entry})
}

func (h *Handler) List(c *gin.Context) {
	orgID, _, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	projectIDStr := c.Query("project_id")
	teamIDStr := c.Query("team_id")

	req := domain.ListRequest{
		OrganisationID: orgID,
		Status:         c.Query("status"),
		Search:         c.Query("q"),
		Limit:          limit,
		Offset:         offset,
	}

	if projectIDStr != "" {
		pID, _ := strconv.Atoi(projectIDStr)
		uPID := uint(pID)
		req.ProjectID = &uPID
	}
	if teamIDStr != "" {
		tID, _ := strconv.Atoi(teamIDStr)
		uTID := uint(tID)
		req.TeamID = &uTID
	}

	items, total, err := h.uc.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": total,
	})
}

func (h *Handler) Update(c *gin.Context) {
	orgID, userID, role, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid research entry ID"})
		return
	}

	var req domain.UpdateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req.ID = uint(id)

	err = h.uc.Update(c.Request.Context(), req, userID, role, orgID)
	if err != nil {
		if err == domain.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Research entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Research entry updated successfully"})
}

func (h *Handler) Delete(c *gin.Context) {
	orgID, userID, role, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid research entry ID"})
		return
	}

	err = h.uc.Delete(c.Request.Context(), uint(id), userID, role, orgID)
	if err != nil {
		if err == domain.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Research entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Research entry deleted successfully"})
}

// --- Folder Handlers ---

func (h *Handler) CreateFolder(c *gin.Context) {
	orgID, userID, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}

	var req domain.CreateFolderRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req.CreatedBy = userID

	folder, err := h.uc.CreateFolder(c.Request.Context(), req, userID, orgID)
	if err != nil {
		if err == domain.ErrDuplicateFolderName || err == domain.ErrInvalidFolderParent {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": folder})
}

func (h *Handler) GetFolder(c *gin.Context) {
	orgID, _, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	folder, err := h.uc.GetFolderByID(c.Request.Context(), uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": folder})
}

func (h *Handler) ListFolders(c *gin.Context) {
	orgID, _, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}
	researchID, _ := strconv.Atoi(c.Query("research_id"))
	if researchID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "research_id is required"})
		return
	}

	var parentID *uint
	parentIDStr := c.Query("parent_id")
	if parentIDStr != "" {
		pid, _ := strconv.Atoi(parentIDStr)
		upid := uint(pid)
		parentID = &upid
	}

	folders, err := h.uc.GetFoldersByResearchID(c.Request.Context(), uint(researchID), parentID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": folders})
}

func (h *Handler) UpdateFolder(c *gin.Context) {
	orgID, userID, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var req domain.UpdateFolderRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req.ID = uint(id)

	err := h.uc.UpdateFolder(c.Request.Context(), req, userID, orgID)
	if err != nil {
		if err == domain.ErrDuplicateFolderName || err == domain.ErrInvalidFolderParent {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Folder updated successfully"})
}

func (h *Handler) DeleteFolder(c *gin.Context) {
	orgID, userID, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.uc.DeleteFolder(c.Request.Context(), uint(id), userID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Folder deleted successfully"})
}

// --- Document Handlers ---

func (h *Handler) CreateDocument(c *gin.Context) {
	orgID, userID, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}

	var req domain.CreateDocumentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req.CreatedBy = userID

	doc, err := h.uc.CreateDocument(c.Request.Context(), req, userID, orgID)
	if err != nil {
		if err == domain.ErrDuplicateDocTitle || err == domain.ErrInvalidFolderParent {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": doc})
}

func (h *Handler) GetDocument(c *gin.Context) {
	orgID, _, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	doc, err := h.uc.GetDocumentByID(c.Request.Context(), uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": doc})
}

func (h *Handler) ListDocuments(c *gin.Context) {
	orgID, _, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}
	researchID, _ := strconv.Atoi(c.Query("research_id"))
	if researchID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "research_id is required"})
		return
	}

	var folderID *uint
	folderIDStr := c.Query("folder_id")
	if folderIDStr != "" {
		fid, _ := strconv.Atoi(folderIDStr)
		ufid := uint(fid)
		folderID = &ufid
	}

	docs, err := h.uc.GetDocumentsByFolderID(c.Request.Context(), uint(researchID), folderID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": docs})
}

func (h *Handler) UpdateDocument(c *gin.Context) {
	orgID, userID, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var req domain.UpdateDocumentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req.ID = uint(id)

	err := h.uc.UpdateDocument(c.Request.Context(), req, userID, orgID)
	if err != nil {
		if err == domain.ErrDuplicateDocTitle || err == domain.ErrInvalidFolderParent {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

func (h *Handler) DeleteDocument(c *gin.Context) {
	orgID, userID, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.uc.DeleteDocument(c.Request.Context(), uint(id), userID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// --- Artifact Handlers ---

func (h *Handler) UploadFile(c *gin.Context) {
	orgID, userID, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}

	researchID, _ := strconv.Atoi(c.PostForm("research_id"))
	documentIDStr := c.PostForm("document_id")
	folderIDStr := c.PostForm("folder_id")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	req := domain.CreateFileRequest{
		ResearchID: uint(researchID),
		FileName:   file.Filename,
		Size:       file.Size,
		MimeType:   file.Header.Get("Content-Type"),
		CreatedBy:  userID,
	}

	if documentIDStr != "" {
		did, _ := strconv.Atoi(documentIDStr)
		udid := uint(did)
		req.DocumentID = &udid
	}
	if folderIDStr != "" {
		fid, _ := strconv.Atoi(folderIDStr)
		ufid := uint(fid)
		req.FolderID = &ufid
	}

	// Save the file to disk
	uploadDir := "uploads/research"
	filePath := fmt.Sprintf("%s/%d_%s", uploadDir, time.Now().Unix(), file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	req.StoragePath = filePath

	resFile, err := h.uc.UploadFile(c.Request.Context(), req, userID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": resFile})
}

func (h *Handler) DeleteFile(c *gin.Context) {
	orgID, userID, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.uc.DeleteFile(c.Request.Context(), uint(id), userID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func (h *Handler) ListFiles(c *gin.Context) {
	orgID, _, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}
	documentID, _ := strconv.Atoi(c.Query("document_id"))
	if documentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document_id is required"})
		return
	}

	files, err := h.uc.GetFilesByDocumentID(c.Request.Context(), uint(documentID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": files})
}

func (h *Handler) AddReference(c *gin.Context) {
	orgID, userID, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}

	var req domain.CreateReferenceRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req.CreatedBy = userID

	ref, err := h.uc.AddReference(c.Request.Context(), req, userID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": ref})
}

func (h *Handler) DeleteReference(c *gin.Context) {
	orgID, userID, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.uc.DeleteReference(c.Request.Context(), uint(id), userID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reference deleted successfully"})
}

func (h *Handler) ListReferences(c *gin.Context) {
	orgID, _, _, ok := getAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}
	documentID, _ := strconv.Atoi(c.Query("document_id"))
	if documentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document_id is required"})
		return
	}

	refs, err := h.uc.GetReferencesByDocumentID(c.Request.Context(), uint(documentID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": refs})
}

func (h *Handler) DownloadFile(c *gin.Context) {
	orgID, _, _, ok := getAuthContext(c)
	log.Default().Println("Create", orgID, ok)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	file, err := h.uc.GetFileByID(c.Request.Context(), uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(file.StoragePath)
}
