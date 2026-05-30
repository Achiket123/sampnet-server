package projects

import (
	"net/http"
	domain "server/internal/domain/projects"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc domain.UseCase
}

func NewHandler(uc domain.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) CreateProject(c *gin.Context) {
	var project domain.Project
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.uc.CreateProject(c.Request.Context(), &project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project created successfully", "project": project})
}

func (h *Handler) GetProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := h.uc.GetProject(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project fetched successfully", "project": project})
}

func (h *Handler) GetProjectsByOrganisation(c *gin.Context) {
	orgID, err := strconv.Atoi(c.Param("organisation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID"})
		return
	}

	projects, err := h.uc.GetProjectsByOrganisation(c.Request.Context(), uint(orgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Projects fetched successfully", "projects": projects})
}

func (h *Handler) UpdateProject(c *gin.Context) {
	var project domain.Project
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	project.ID = uint(id)

	if err := h.uc.UpdateProject(c.Request.Context(), &project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully", "project": project})
}

func (h *Handler) GetProjectsByTeam(c *gin.Context) {
	teamID, err := strconv.Atoi(c.Param("team_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	projects, err := h.uc.GetProjectsByTeam(c.Request.Context(), uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Projects fetched successfully", "projects": projects})
}

func (h *Handler) GetProjectsWithLessData(c *gin.Context) {
	orgID, err := strconv.Atoi(c.Param("organisation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID"})
		return
	}

	projects, err := h.uc.GetProjectsWithLessData(c.Request.Context(), uint(orgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Projects fetched successfully", "projects": projects})
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := h.uc.DeleteProject(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
