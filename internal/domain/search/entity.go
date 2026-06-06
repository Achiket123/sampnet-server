package search

// SearchResultItem is a single result from a global search query.
type SearchResultItem struct {
	ID             uint                   `json:"id"`
	Type           string                 `json:"type"`
	Title          string                 `json:"title"`
	Subtitle       string                 `json:"subtitle"`
	OrganisationID uint                   `json:"organisation_id"`
	ExtraData      map[string]interface{} `json:"extra_data,omitempty"`
}

// SearchResults is the top-level response from a global search query.
type SearchResults struct {
	Query         string             `json:"query"`
	TotalCount    int                `json:"total_count"`
	Items         []SearchResultItem `json:"items"`
	TaskCount     int                `json:"task_count"`
	ProjectCount  int                `json:"project_count"`
	TeamCount     int                `json:"team_count"`
	EmployeeCount int                `json:"employee_count"`
	ChatCount     int                `json:"chat_count"`
	ResourceCount int                `json:"resource_count"`
	ResearchCount int                `json:"research_count"`
}

// SearchFilters carries all parameters for a search request.
type SearchFilters struct {
	Query          string
	OrganisationID uint
	// Types is an optional whitelist of result types (e.g. ["task","project"]).
	// An empty slice means all types are searched.
	Types  []string
	Limit  int
	Offset int
}