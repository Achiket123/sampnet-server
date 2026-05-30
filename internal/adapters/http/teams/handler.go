package teams

import (
	"net/http"
	domain "server/internal/domain/teams"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc domain.UseCase
}

func NewHandler(uc domain.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) CreateTeam(c *gin.Context) {
	var req struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
		OrganisationID uint   `json:"organisation_id"`
		CreatedBy      uint   `json:"created_by"`
		TeamLead       uint   `json:"team_lead"`
		Members        []uint `json:"members"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	team := &domain.Team{
		Name:           req.Name,
		Description:    req.Description,
		OrganisationID: req.OrganisationID,
		CreatedBy:      req.CreatedBy,
		TeamLead:       req.TeamLead,
	}

	if err := h.uc.CreateTeam(c.Request.Context(), team, req.Members); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team created successfully", "team": team})
}

func (h *Handler) GetTeam(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	team, members, err := h.uc.GetTeam(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get team", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team fetched successfully", "team": team, "team_members": members})
}

func (h *Handler) GetTeamsByOrganisation(c *gin.Context) {
	orgID, err := strconv.Atoi(c.Param("organisation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID"})
		return
	}

	teams, err := h.uc.GetTeamsByOrganisation(c.Request.Context(), uint(orgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get teams", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Teams fetched successfully", "teams": teams})
}

func (h *Handler) UpdateTeam(c *gin.Context) {
	var team domain.Team
	if err := c.BindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}
	team.ID = uint(id)

	if err := h.uc.UpdateTeam(c.Request.Context(), &team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update team"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team updated successfully"})
}

func (h *Handler) DeleteTeam(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	if err := h.uc.DeleteTeam(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete team"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}
