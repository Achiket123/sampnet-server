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

func (h *Handler) CreateProject(c *gin.Context) {
	orgID, userID, ok := getOrgAndUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User or Organisation ID missing from token"})
		return
	}

	var project domain.Project
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	project.OrganisationID = orgID
	project.CreatedBy = userID

	if err := h.uc.CreateProject(c.Request.Context(), &project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project created successfully", "project": project})
}

func (h *Handler) GetProject(c *gin.Context) {
	orgIDVal, existsOrg := c.Get("organisationID")
	if !existsOrg {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organisation ID missing from token"})
		return
	}
	orgID := orgIDVal.(uint)

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

	if project.OrganisationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Project does not belong to your organisation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project fetched successfully", "project": project})
}

func (h *Handler) GetProjectsByOrganisation(c *gin.Context) {
	orgIDVal, existsOrg := c.Get("organisationID")
	if !existsOrg {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organisation ID missing from token"})
		return
	}
	orgID := orgIDVal.(uint)

	projects, err := h.uc.GetProjectsByOrganisation(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Projects fetched successfully", "projects": projects})
}

func (h *Handler) UpdateProject(c *gin.Context) {
	orgIDVal, existsOrg := c.Get("organisationID")
	if !existsOrg {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organisation ID missing from token"})
		return
	}
	orgID := orgIDVal.(uint)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	existingProject, err := h.uc.GetProject(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if existingProject.OrganisationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Project does not belong to your organisation"})
		return
	}

	var project domain.Project
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}
	project.ID = uint(id)
	project.OrganisationID = orgID // Enforce multi-tenancy integrity

	if err := h.uc.UpdateProject(c.Request.Context(), &project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully", "project": project})
}

func (h *Handler) GetProjectsByTeam(c *gin.Context) {
	orgIDVal, existsOrg := c.Get("organisationID")
	if !existsOrg {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organisation ID missing from token"})
		return
	}
	orgID := orgIDVal.(uint)

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

	// Filter elements or verify that they belong to the correct organisation
	var filteredProjects []domain.Project
	for _, p := range projects {
		if p.OrganisationID == orgID {
			filteredProjects = append(filteredProjects, p)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Projects fetched successfully", "projects": filteredProjects})
}

func (h *Handler) GetProjectsWithLessData(c *gin.Context) {
	orgIDVal, existsOrg := c.Get("organisationID")
	if !existsOrg {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organisation ID missing from token"})
		return
	}
	orgID := orgIDVal.(uint)

	projects, err := h.uc.GetProjectsWithLessData(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Projects fetched successfully", "projects": projects})
}

func (h *Handler) DeleteProject(c *gin.Context) {
	orgIDVal, existsOrg := c.Get("organisationID")
	if !existsOrg {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organisation ID missing from token"})
		return
	}
	orgID := orgIDVal.(uint)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	existingProject, err := h.uc.GetProject(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if existingProject.OrganisationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Project does not belong to your organisation"})
		return
	}

	if err := h.uc.DeleteProject(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

func (h *Handler) CreateMilestone(c *gin.Context) {
	orgIDVal, existsOrg := c.Get("organisationID")
	if !existsOrg {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organisation ID missing from token"})
		return
	}
	orgID := orgIDVal.(uint)

	var milestone domain.Milestone
	if err := c.BindJSON(&milestone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := h.uc.GetProject(c.Request.Context(), uint(projectID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if project.OrganisationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Project does not belong to your organisation"})
		return
	}

	milestone.ProjectID = uint(projectID)
	milestone.OrganisationID = orgID

	if err := h.uc.CreateMilestone(c.Request.Context(), &milestone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create milestone"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Milestone created successfully", "milestone": milestone})
}

func (h *Handler) UpdateMilestone(c *gin.Context) {
	orgIDVal, existsOrg := c.Get("organisationID")
	if !existsOrg {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organisation ID missing from token"})
		return
	}
	orgID := orgIDVal.(uint)

	milestoneID, err := strconv.Atoi(c.Param("milestone_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone ID"})
		return
	}

	existingMilestone, err := h.uc.GetMilestoneByID(c.Request.Context(), uint(milestoneID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found"})
		return
	}

	if existingMilestone.OrganisationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Milestone does not belong to your organisation"})
		return
	}

	var milestone domain.Milestone
	if err := c.BindJSON(&milestone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	milestone.ID = uint(milestoneID)
	milestone.ProjectID = existingMilestone.ProjectID
	milestone.OrganisationID = orgID

	if err := h.uc.UpdateMilestone(c.Request.Context(), &milestone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update milestone"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Milestone updated successfully", "milestone": milestone})
}

func (h *Handler) DeleteMilestone(c *gin.Context) {
	orgIDVal, existsOrg := c.Get("organisationID")
	if !existsOrg {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organisation ID missing from token"})
		return
	}
	orgID := orgIDVal.(uint)

	milestoneID, err := strconv.Atoi(c.Param("milestone_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone ID"})
		return
	}

	existingMilestone, err := h.uc.GetMilestoneByID(c.Request.Context(), uint(milestoneID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found"})
		return
	}

	if existingMilestone.OrganisationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Milestone does not belong to your organisation"})
		return
	}

	if err := h.uc.DeleteMilestone(c.Request.Context(), uint(milestoneID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete milestone"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Milestone deleted successfully"})
}
