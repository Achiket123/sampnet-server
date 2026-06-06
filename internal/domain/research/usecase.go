package research

import (
	"context"
	"errors"
)

type CreateRequest struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Status         string   `json:"status"`
	AuthorID       uint     `json:"author_id"`
	ProjectID      *uint    `json:"project_id"`
	TeamID         *uint    `json:"team_id"`
	OrganisationID uint     `json:"organisation_id"`
	Tags           []string `json:"tags"`
}

type UpdateRequest struct {
	ID          uint     `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	ProjectID   *uint    `json:"project_id"`
	TeamID      *uint    `json:"team_id"`
	Tags        []string `json:"tags"`
}

type ListRequest struct {
	OrganisationID uint
	Status         string
	ProjectID      *uint
	TeamID         *uint
	Search         string
	Limit          int
	Offset         int
}

type CreateFolderRequest struct {
	ResearchID uint   `json:"research_id"`
	ParentID   *uint  `json:"parent_id"`
	Name       string `json:"name"`
	CreatedBy  uint   `json:"created_by"`
}

type UpdateFolderRequest struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"` // For moving folders
}

type CreateDocumentRequest struct {
	ResearchID uint   `json:"research_id"`
	FolderID   *uint  `json:"folder_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreatedBy  uint   `json:"created_by"`
}

type UpdateDocumentRequest struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	FolderID *uint  `json:"folder_id"` // For moving documents
}

type CreateFileRequest struct {
	ResearchID  uint   `json:"research_id"`
	DocumentID  *uint  `json:"document_id"`
	FolderID    *uint  `json:"folder_id"`
	FileName    string `json:"file_name"`
	MimeType    string `json:"mime_type"`
	Size        int64  `json:"size"`
	Checksum    string `json:"checksum"`
	StoragePath string `json:"storage_path"`
	CreatedBy   uint   `json:"created_by"`
}

type CreateReferenceRequest struct {
	ResearchID uint   `json:"research_id"`
	DocumentID *uint  `json:"document_id"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	Authors    string `json:"authors"`
	CreatedBy  uint   `json:"created_by"`
}

type UseCase interface {
	Create(ctx context.Context, req CreateRequest) (*ResearchEntry, error)
	GetByID(ctx context.Context, id uint, orgID uint) (*ResearchEntry, error)
	List(ctx context.Context, req ListRequest) ([]ResearchEntry, int, error)
	Update(ctx context.Context, req UpdateRequest, actorID uint, actorRole string, orgID uint) error
	Delete(ctx context.Context, id uint, actorID uint, actorRole string, orgID uint) error

	// Folder Operations
	CreateFolder(ctx context.Context, req CreateFolderRequest, actorID uint, orgID uint) (*ResearchFolder, error)
	GetFolderByID(ctx context.Context, id uint, orgID uint) (*ResearchFolder, error)
	GetFoldersByResearchID(ctx context.Context, researchID uint, parentID *uint, orgID uint) ([]ResearchFolder, error)
	UpdateFolder(ctx context.Context, req UpdateFolderRequest, actorID uint, orgID uint) error
	DeleteFolder(ctx context.Context, id uint, actorID uint, orgID uint) error

	// Document Operations
	CreateDocument(ctx context.Context, req CreateDocumentRequest, actorID uint, orgID uint) (*ResearchDocument, error)
	GetDocumentByID(ctx context.Context, id uint, orgID uint) (*ResearchDocument, error)
	GetDocumentsByFolderID(ctx context.Context, researchID uint, folderID *uint, orgID uint) ([]ResearchDocument, error)
	UpdateDocument(ctx context.Context, req UpdateDocumentRequest, actorID uint, orgID uint) error
	DeleteDocument(ctx context.Context, id uint, actorID uint, orgID uint) error

	// Artifact Operations
	UploadFile(ctx context.Context, req CreateFileRequest, actorID uint, orgID uint) (*ResearchFile, error)
	GetFileByID(ctx context.Context, id uint, orgID uint) (*ResearchFile, error)
	DeleteFile(ctx context.Context, id uint, actorID uint, orgID uint) error
	GetFilesByDocumentID(ctx context.Context, docID uint, orgID uint) ([]ResearchFile, error)

	AddReference(ctx context.Context, req CreateReferenceRequest, actorID uint, orgID uint) (*ResearchReference, error)
	DeleteReference(ctx context.Context, id uint, actorID uint, orgID uint) error
	GetReferencesByDocumentID(ctx context.Context, docID uint, orgID uint) ([]ResearchReference, error)
}

var (
	ErrInvalidStatus       = errors.New("invalid status")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrNotFound            = errors.New("research entry not found")
	ErrDuplicateFolderName = errors.New("folder with the same name already exists in this location")
	ErrDuplicateDocTitle   = errors.New("document with the same title already exists in this location")
	ErrInvalidFolderParent = errors.New("invalid folder parent")
)
