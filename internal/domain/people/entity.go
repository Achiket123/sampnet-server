package people

import (
	"time"
)

type PeopleContact struct {
	ID             uint                   `json:"id"`
	OrganisationID uint                   `json:"organisation_id"`
	FirstName      string                 `json:"first_name"`
	LastName       string                 `json:"last_name"`
	Email          *string                `json:"email"`
	Phone          *string                `json:"phone"`
	Company        *string                `json:"company"`
	JobTitle       *string                `json:"job_title"`
	AvatarUrl      *string                `json:"avatar_url"`
	Type           string                 `json:"type"`
	Status         string                 `json:"status"`
	Stage          string                 `json:"stage"`
	Tags           []string               `json:"tags"`
	ListIDs        []int                  `json:"list_ids"`
	AssignedToID   *uint                  `json:"assigned_to_id"`
	AssignedToName *string                `json:"assigned_to_name"`
	AssignedToAvatar *string              `json:"assigned_to_avatar_url"`
	Source         *string                `json:"source"`
	DealValue      *float64               `json:"deal_value"`
	Currency       string                 `json:"currency"`
	Notes          *string                `json:"notes"`
	Address        *string                `json:"address"`
	City           *string                `json:"city"`
	Country        *string                `json:"country"`
	LinkedinUrl    *string                `json:"linkedin_url"`
	WebsiteUrl     *string                `json:"website_url"`
	CustomFields   map[string]interface{} `json:"custom_fields"`
	LastContactedAt *time.Time             `json:"last_contacted_at"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type PeopleInteraction struct {
	ID                 uint       `json:"id"`
	ContactID          uint       `json:"contact_id"`
	CreatedByID        uint       `json:"created_by_id"`
	CreatedByName      string     `json:"created_by_name"`
	CreatedByAvatarUrl *string    `json:"created_by_avatar_url"`
	Type               string     `json:"type"`
	Content            string     `json:"content"`
	Outcome            *string    `json:"outcome"`
	LinkedTaskID       *uint      `json:"linked_task_id"`
	LinkedTaskTitle    *string    `json:"linked_task_title"`
	OccurredAt         time.Time  `json:"occurred_at"`
	CreatedAt          time.Time  `json:"created_at"`
}

type PeoplePipelineStage struct {
	ID             uint   `json:"id"`
	OrganisationID uint   `json:"organisation_id"`
	Key            string `json:"key"`
	Label          string `json:"label"`
	Color          string `json:"color"`
	Order          int    `json:"order"`
	IsDefault      bool   `json:"is_default"`
	IsWon          bool   `json:"is_won"`
	IsLost         bool   `json:"is_lost"`
}

type PeopleList struct {
	ID             uint      `json:"id"`
	OrganisationID uint      `json:"organisation_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Color          string    `json:"color"`
	ContactCount   int       `json:"contact_count"`
	CreatedAt      time.Time `json:"created_at"`
}

type PeopleAnalytics struct {
	TotalContacts   int            `json:"total_contacts"`
	ByStage         map[string]int `json:"by_stage"`
	ByType          map[string]int `json:"by_type"`
	BySource        map[string]int `json:"by_source"`
	TotalDealValue  float64        `json:"total_deal_value"`
	WonDealValue    float64        `json:"won_deal_value"`
	ConversionRate  float64        `json:"conversion_rate"`
	AvgDealValue    float64        `json:"avg_deal_value"`
	RecentActivity  int            `json:"recent_activity"`
}

type CreateContactParams struct {
	FirstName    string                 `json:"first_name" binding:"required"`
	LastName     string                 `json:"last_name" binding:"required"`
	Email        *string                `json:"email"`
	Phone        *string                `json:"phone"`
	Company      *string                `json:"company"`
	JobTitle     *string                `json:"job_title"`
	AvatarUrl    *string                `json:"avatar_url"`
	Type         string                 `json:"type" binding:"required"`
	Status       string                 `json:"status"`
	Stage        string                 `json:"stage"`
	Tags         []string               `json:"tags"`
	ListIDs      []int                  `json:"list_ids"`
	AssignedToID *uint                  `json:"assigned_to_id"`
	Source       *string                `json:"source"`
	DealValue    *float64               `json:"deal_value"`
	Currency     string                 `json:"currency"`
	Notes        *string                `json:"notes"`
	Address      *string                `json:"address"`
	City         *string                `json:"city"`
	Country      *string                `json:"country"`
	LinkedinUrl  *string                `json:"linkedin_url"`
	WebsiteUrl   *string                `json:"website_url"`
	CustomFields map[string]interface{} `json:"custom_fields"`
	LastContactedAt *time.Time             `json:"last_contacted_at"`
}

type UpdateContactParams struct {
	FirstName    *string                 `json:"first_name"`
	LastName     *string                 `json:"last_name"`
	Email        *string                 `json:"email"`
	Phone        *string                 `json:"phone"`
	Company      *string                 `json:"company"`
	JobTitle     *string                 `json:"job_title"`
	AvatarUrl    *string                 `json:"avatar_url"`
	Type         *string                 `json:"type"`
	Status       *string                 `json:"status"`
	Stage        *string                 `json:"stage"`
	Tags         []string                `json:"tags"`
	ListIDs      []int                   `json:"list_ids"`
	AssignedToID *uint                   `json:"assigned_to_id"`
	Source       *string                 `json:"source"`
	DealValue    *float64                `json:"deal_value"`
	Currency     *string                 `json:"currency"`
	Notes        *string                 `json:"notes"`
	Address      *string                 `json:"address"`
	City         *string                 `json:"city"`
	Country      *string                 `json:"country"`
	LinkedinUrl  *string                 `json:"linkedin_url"`
	WebsiteUrl   *string                 `json:"website_url"`
	CustomFields map[string]interface{}  `json:"custom_fields"`
}

type AddInteractionParams struct {
	Type         string     `json:"type" binding:"required"`
	Content      string     `json:"content" binding:"required"`
	Outcome      *string    `json:"outcome"`
	LinkedTaskID *uint      `json:"linked_task_id"`
	OccurredAt   time.Time  `json:"occurred_at"`
}

type CreateStageParams struct {
	Key       string `json:"key" binding:"required"`
	Label     string `json:"label" binding:"required"`
	Color     string `json:"color" binding:"required"`
	Order     int    `json:"order"`
	IsWon     bool   `json:"is_won"`
	IsLost    bool   `json:"is_lost"`
}

type UpdateStageParams struct {
	Label  *string `json:"label"`
	Color  *string `json:"color"`
	Order  *int    `json:"order"`
	IsWon  *bool   `json:"is_won"`
	IsLost *bool   `json:"is_lost"`
}

type CreateListParams struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type UpdateListParams struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
}

type ContactsFilter struct {
	Stage    *string
	Type     *string
	Status   *string
	ListID   *int
	Search   string
	SortBy   string
	SortOrder string
	Page     int
	Limit    int
}

type ContactsResponse struct {
	Contacts []PeopleContact `json:"contacts"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	Limit    int             `json:"limit"`
}

type ContactDetailResponse struct {
	Contact      PeopleContact       `json:"contact"`
	Interactions []PeopleInteraction `json:"interactions"`
}
