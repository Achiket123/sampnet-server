package research

import "time"

type ResearchEntry struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Thumbnail      string    `json:"thumbnail,omitempty"`
	Status         string    `json:"status"`     // draft, active, archived, completed
	Visibility     string    `json:"visibility"` // private, team, organisation
	AuthorID       uint      `json:"author_id"`
	AuthorName     string    `json:"author_name"`
	ProjectID      *uint     `json:"project_id,omitempty"`
	ProjectName    string    `json:"project_name,omitempty"`
	TeamID         *uint     `json:"team_id,omitempty"`
	TeamName       string    `json:"team_name,omitempty"`
	OrganisationID uint      `json:"organisation_id"`
	Tags           []string  `json:"tags"`

	// Sub-entities (populated on deep fetch)
	Folders       []ResearchFolder       `json:"folders,omitempty"`
	Documents     []ResearchDocument     `json:"documents,omitempty"`
	Files         []ResearchFile         `json:"files,omitempty"`
	Collaborators []ResearchCollaborator `json:"collaborators,omitempty"`
}

type ResearchFolder struct {
	ID             uint               `json:"id"`
	ResearchID     uint               `json:"research_id"`
	OrganisationID uint               `json:"organisation_id"`
	ParentID       *uint              `json:"parent_id,omitempty"`
	Name           string             `json:"name"`
	CreatedBy      uint               `json:"created_by"`
	UpdatedAt      time.Time          `json:"updated_at"`
	SubFolders     []ResearchFolder   `json:"sub_folders,omitempty"`
	Documents      []ResearchDocument `json:"documents,omitempty"`
	Files          []ResearchFile     `json:"files,omitempty"`
}

type ResearchDocument struct {
	ID             uint      `json:"id"`
	ResearchID     uint      `json:"research_id"`
	OrganisationID uint      `json:"organisation_id"`
	FolderID       *uint     `json:"folder_id,omitempty"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	IsPinned       bool      `json:"is_pinned"`
	Status         string    `json:"status"`
	CreatedBy      uint      `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ResearchFile struct {
	ID             uint      `json:"id"`
	ResearchID     uint      `json:"research_id"`
	DocumentID     *uint     `json:"document_id,omitempty"`
	OrganisationID uint      `json:"organisation_id"`
	FolderID       *uint     `json:"folder_id,omitempty"`
	FileName       string    `json:"file_name"`
	MimeType       string    `json:"mime_type"`
	Size           int64     `json:"size"`
	StoragePath    string    `json:"storage_path"`
	PreviewPath    string    `json:"preview_path,omitempty"`
	ThumbnailPath  string    `json:"thumbnail_path,omitempty"`
	IsPinned       bool      `json:"is_pinned"`
	CreatedBy      uint      `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ResearchReference struct {
	ID             uint   `json:"id"`
	ResearchID     uint   `json:"research_id"`
	DocumentID     *uint  `json:"document_id,omitempty"`
	OrganisationID uint   `json:"organisation_id"`
	Title          string `json:"title"`
	Authors        string `json:"authors"`
	PublicationDate string `json:"publication_date,omitempty"`
	URL            string `json:"url,omitempty"`
	CreatedBy      uint   `json:"created_by"`
}

type ResearchCollaborator struct {
	ID         uint   `json:"id"`
	ResearchID uint   `json:"research_id"`
	UserID     uint   `json:"user_id"`
	UserName   string `json:"user_name"`
	Role       string `json:"role"` // owner, editor, commenter, viewer
}

type ResearchActivity struct {
	ID         uint      `json:"id"`
	ResearchID uint      `json:"research_id"`
	ActorName  string    `json:"actor_name"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	CreatedAt  time.Time `json:"created_at"`
}
