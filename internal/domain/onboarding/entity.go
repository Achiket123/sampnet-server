package onboarding

type OnboardingProgress struct {
	UserID           uint `json:"userId"`
	OrganisationID   uint `json:"organisationId"`
	ProfileCompleted bool `json:"profileCompleted"`
	TeamJoined       bool `json:"teamJoined"`
	TaskCreated      bool `json:"taskCreated"`
	InviteSent       bool `json:"inviteSent"`
}
