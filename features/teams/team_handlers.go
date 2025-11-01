package teams

import (
	"log"
	"net/http"
	"server/database"
	"server/database/models"

	"github.com/gin-gonic/gin"
)

func CreateTeam(c *gin.Context) {

	type CreateTeamRequest struct {
		Name           string ` json:"name"`
		Description    string ` json:"description"`
		OrganisationID uint   ` json:"organisation_id"`
		CreatedBy      uint   ` json:"created_by"`
		TeamLead       uint   ` json:"team_lead"`
		Members        []uint `json:"members"`
		ProjectId      int    `json:"project_id"`
	}

	var createTeamReq CreateTeamRequest
	var team models.Team

	// Parse JSON body into the `team` struct
	if err := c.ShouldBindJSON(&createTeamReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	// Create the team
	team.Name = createTeamReq.Name
	team.Description = createTeamReq.Description
	team.OrganisationID = createTeamReq.OrganisationID
	team.CreatedBy = createTeamReq.CreatedBy
	team.TeamLead = createTeamReq.TeamLead
	team.IsActive = true
	if err := database.DB.Create(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team", "message": err.Error()})
		return
	}
	// create Team Members
	for _, memberID := range createTeamReq.Members {
		log.Println(memberID)
		log.Println(team.ID)
		var role string
		if team.TeamLead == memberID {
			role = "Team Lead"
		} else {
			role = "Member"
		}
		member := models.TeamMember{
			UserID:     memberID,
			TeamID:     team.ID,
			Role:       role,
			IsActive:   true,
			IsLeader:   team.TeamLead == memberID,
			IsAdmin:    false,
			IsManager:  false,
			IsEmployee: false,
			IsBoss:     false,
		}
		if err := database.DB.Create(&member).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team member", "message": err.Error()})
			return
		}
	}

	// Create the team and its relationships

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Team created successfully", "team": team})
}

func GetTeam(c *gin.Context) {
	var team models.Team
	if err := database.DB.Preload("Organisation").Preload("CreatedByUser").Preload("TeamLeadUser").First(&team, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get team", "message": err.Error()})
		return
	}

	var teamMembers []models.TeamMember
	if err := database.DB.Preload("User").Where("team_id = ?", team.ID).Find(&teamMembers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get team members", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team fetched successfully", "team": team, "team_members": teamMembers})
}
func GetTeamsWithLessData(c *gin.Context) {
	var teams []models.Team
	if err := database.DB.Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get teams"})
		return
	}
	c.JSON(http.StatusOK, teams)
}

func GetTeamsByOrganisation(c *gin.Context) {
	var teams []models.Team
	if err := database.DB.Where("organisation_id = ?", c.Param("organisation_id")).Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get teams", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Teams fetched successfully", "teams": teams})
}

func UpdateTeam(c *gin.Context) {
	var team models.Team
	if err := c.BindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := database.DB.Save(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update team"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team updated successfully"})
}

func DelteTeam(c *gin.Context) {
	if err := database.DB.First(&models.Team{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	if err := database.DB.Delete(&models.Team{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete team"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}
