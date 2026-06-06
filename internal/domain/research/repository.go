package research

import "context"

type ListFilters struct {
	Status    string
	ProjectID *uint
	TeamID    *uint
	Search    string
}

type Repository interface {
	Create(ctx context.Context, entry *ResearchEntry) error
	GetByID(ctx context.Context, id uint, orgID uint) (*ResearchEntry, error)
	List(ctx context.Context, orgID uint, filters ListFilters, limit int, offset int) ([]ResearchEntry, int, error)
	Update(ctx context.Context, entry *ResearchEntry) error
	Delete(ctx context.Context, id uint, orgID uint) error
	FolderNameExists(ctx context.Context, researchID uint, parentID *uint, name string, excludeID *uint) (bool, error)
	DocumentTitleExists(ctx context.Context, researchID uint, folderID *uint, title string, excludeID *uint) (bool, error)

	// Folder Operations
	CreateFolder(ctx context.Context, folder *ResearchFolder) error
	GetFolderByID(ctx context.Context, id uint) (*ResearchFolder, error)
	GetFoldersByResearchID(ctx context.Context, researchID uint, parentID *uint) ([]ResearchFolder, error)
	UpdateFolder(ctx context.Context, folder *ResearchFolder) error
	DeleteFolder(ctx context.Context, id uint) error

	// Document Operations
	CreateDocument(ctx context.Context, doc *ResearchDocument) error
	GetDocumentByID(ctx context.Context, id uint) (*ResearchDocument, error)
	GetDocumentsByFolderID(ctx context.Context, researchID uint, folderID *uint) ([]ResearchDocument, error)
	UpdateDocument(ctx context.Context, doc *ResearchDocument) error
	DeleteDocument(ctx context.Context, id uint) error

	// File Operations (Artifacts)
	CreateFile(ctx context.Context, file *ResearchFile) error
	GetFileByID(ctx context.Context, id uint) (*ResearchFile, error)
	GetFilesByDocumentID(ctx context.Context, docID uint) ([]ResearchFile, error)
	GetFilesByResearchID(ctx context.Context, researchID uint) ([]ResearchFile, error)
	UpdateFile(ctx context.Context, file *ResearchFile) error
	DeleteFile(ctx context.Context, id uint) error

	// Reference/Link Operations
	CreateReference(ctx context.Context, ref *ResearchReference) error
	GetReferenceByID(ctx context.Context, id uint) (*ResearchReference, error)
	GetReferencesByDocumentID(ctx context.Context, docID uint) ([]ResearchReference, error)
	GetReferencesByResearchID(ctx context.Context, researchID uint) ([]ResearchReference, error)
	UpdateReference(ctx context.Context, ref *ResearchReference) error
	DeleteReference(ctx context.Context, id uint) error
}
