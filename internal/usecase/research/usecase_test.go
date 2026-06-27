package research

import (
	"context"
	"errors"
	"testing"

	"server/internal/domain/notifications"
	"server/internal/domain/projects"
	domain "server/internal/domain/research"
	"server/internal/domain/teams"
)

type mockRepo struct {
	createFunc                 func(ctx context.Context, entry *domain.ResearchEntry) error
	getByIDFunc                func(ctx context.Context, id uint, orgID uint) (*domain.ResearchEntry, error)
	listFunc                   func(ctx context.Context, orgID uint, filters domain.ListFilters, limit int, offset int) ([]domain.ResearchEntry, int, error)
	updateFunc                 func(ctx context.Context, entry *domain.ResearchEntry) error
	deleteFunc                 func(ctx context.Context, id uint, orgID uint) error
	folderNameExistsFunc       func(ctx context.Context, researchID uint, parentID *uint, name string, excludeID *uint) (bool, error)
	docTitleExistsFunc         func(ctx context.Context, researchID uint, folderID *uint, title string, excludeID *uint) (bool, error)
	createFolderFunc           func(ctx context.Context, folder *domain.ResearchFolder) error
	getFolderByIDFunc          func(ctx context.Context, id uint) (*domain.ResearchFolder, error)
	getFoldersByResearchIDFunc func(ctx context.Context, researchID uint, parentID *uint) ([]domain.ResearchFolder, error)
	updateFolderFunc           func(ctx context.Context, folder *domain.ResearchFolder) error
	deleteFolderFunc           func(ctx context.Context, id uint) error
	createDocFunc              func(ctx context.Context, doc *domain.ResearchDocument) error
	getDocByIDFunc             func(ctx context.Context, id uint) (*domain.ResearchDocument, error)
	getDocsByFolderIDFunc      func(ctx context.Context, researchID uint, folderID *uint) ([]domain.ResearchDocument, error)
	updateDocFunc              func(ctx context.Context, doc *domain.ResearchDocument) error
	deleteDocFunc              func(ctx context.Context, id uint) error
	createFileFunc             func(ctx context.Context, file *domain.ResearchFile) error
	getFileByIDFunc            func(ctx context.Context, id uint) (*domain.ResearchFile, error)
	getFilesByDocIDFunc        func(ctx context.Context, docID uint) ([]domain.ResearchFile, error)
	getFilesByResearchIDFunc   func(ctx context.Context, researchID uint) ([]domain.ResearchFile, error)
	updateFileFunc             func(ctx context.Context, file *domain.ResearchFile) error
	deleteFileFunc             func(ctx context.Context, id uint) error
	createRefFunc              func(ctx context.Context, ref *domain.ResearchReference) error
	getRefByIDFunc             func(ctx context.Context, id uint) (*domain.ResearchReference, error)
	getRefsByDocIDFunc         func(ctx context.Context, docID uint) ([]domain.ResearchReference, error)
	getRefsByResearchIDFunc    func(ctx context.Context, researchID uint) ([]domain.ResearchReference, error)
	updateRefFunc              func(ctx context.Context, ref *domain.ResearchReference) error
	deleteRefFunc              func(ctx context.Context, id uint) error
}

func (m *mockRepo) Create(ctx context.Context, entry *domain.ResearchEntry) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, entry)
	}
	entry.ID = 1
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id uint, orgID uint) (*domain.ResearchEntry, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id, orgID)
	}
	return &domain.ResearchEntry{ID: id, OrganisationID: orgID, AuthorID: 1, Status: "draft"}, nil
}

func (m *mockRepo) List(ctx context.Context, orgID uint, filters domain.ListFilters, limit int, offset int) ([]domain.ResearchEntry, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, orgID, filters, limit, offset)
	}
	return []domain.ResearchEntry{{ID: 1, OrganisationID: orgID}}, 1, nil
}

func (m *mockRepo) Update(ctx context.Context, entry *domain.ResearchEntry) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, entry)
	}
	return nil
}

func (m *mockRepo) Delete(ctx context.Context, id uint, orgID uint) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id, orgID)
	}
	return nil
}

func (m *mockRepo) FolderNameExists(ctx context.Context, researchID uint, parentID *uint, name string, excludeID *uint) (bool, error) {
	if m.folderNameExistsFunc != nil {
		return m.folderNameExistsFunc(ctx, researchID, parentID, name, excludeID)
	}
	return false, nil
}

func (m *mockRepo) DocumentTitleExists(ctx context.Context, researchID uint, folderID *uint, title string, excludeID *uint) (bool, error) {
	if m.docTitleExistsFunc != nil {
		return m.docTitleExistsFunc(ctx, researchID, folderID, title, excludeID)
	}
	return false, nil
}

func (m *mockRepo) CreateFolder(ctx context.Context, folder *domain.ResearchFolder) error {
	if m.createFolderFunc != nil {
		return m.createFolderFunc(ctx, folder)
	}
	folder.ID = 10
	return nil
}

func (m *mockRepo) GetFolderByID(ctx context.Context, id uint) (*domain.ResearchFolder, error) {
	if m.getFolderByIDFunc != nil {
		return m.getFolderByIDFunc(ctx, id)
	}
	return &domain.ResearchFolder{ID: id, ResearchID: 1}, nil
}

func (m *mockRepo) GetFoldersByResearchID(ctx context.Context, researchID uint, parentID *uint) ([]domain.ResearchFolder, error) {
	if m.getFoldersByResearchIDFunc != nil {
		return m.getFoldersByResearchIDFunc(ctx, researchID, parentID)
	}
	return []domain.ResearchFolder{{ID: 10, ResearchID: researchID}}, nil
}

func (m *mockRepo) UpdateFolder(ctx context.Context, folder *domain.ResearchFolder) error {
	if m.updateFolderFunc != nil {
		return m.updateFolderFunc(ctx, folder)
	}
	return nil
}

func (m *mockRepo) DeleteFolder(ctx context.Context, id uint) error {
	if m.deleteFolderFunc != nil {
		return m.deleteFolderFunc(ctx, id)
	}
	return nil
}

func (m *mockRepo) CreateDocument(ctx context.Context, doc *domain.ResearchDocument) error {
	if m.createDocFunc != nil {
		return m.createDocFunc(ctx, doc)
	}
	doc.ID = 100
	return nil
}

func (m *mockRepo) GetDocumentByID(ctx context.Context, id uint) (*domain.ResearchDocument, error) {
	if m.getDocByIDFunc != nil {
		return m.getDocByIDFunc(ctx, id)
	}
	return &domain.ResearchDocument{ID: id, ResearchID: 1}, nil
}

func (m *mockRepo) GetDocumentsByFolderID(ctx context.Context, researchID uint, folderID *uint) ([]domain.ResearchDocument, error) {
	if m.getDocsByFolderIDFunc != nil {
		return m.getDocsByFolderIDFunc(ctx, researchID, folderID)
	}
	return []domain.ResearchDocument{{ID: 100, ResearchID: researchID}}, nil
}

func (m *mockRepo) UpdateDocument(ctx context.Context, doc *domain.ResearchDocument) error {
	if m.updateDocFunc != nil {
		return m.updateDocFunc(ctx, doc)
	}
	return nil
}

func (m *mockRepo) DeleteDocument(ctx context.Context, id uint) error {
	if m.deleteDocFunc != nil {
		return m.deleteDocFunc(ctx, id)
	}
	return nil
}

func (m *mockRepo) CreateFile(ctx context.Context, file *domain.ResearchFile) error {
	if m.createFileFunc != nil {
		return m.createFileFunc(ctx, file)
	}
	file.ID = 1000
	return nil
}

func (m *mockRepo) GetFileByID(ctx context.Context, id uint) (*domain.ResearchFile, error) {
	if m.getFileByIDFunc != nil {
		return m.getFileByIDFunc(ctx, id)
	}
	return &domain.ResearchFile{ID: id, ResearchID: 1, OrganisationID: 10}, nil
}

func (m *mockRepo) GetFilesByDocumentID(ctx context.Context, docID uint) ([]domain.ResearchFile, error) {
	if m.getFilesByDocIDFunc != nil {
		return m.getFilesByDocIDFunc(ctx, docID)
	}
	return []domain.ResearchFile{{ID: 1000}}, nil
}

func (m *mockRepo) GetFilesByResearchID(ctx context.Context, researchID uint) ([]domain.ResearchFile, error) {
	if m.getFilesByResearchIDFunc != nil {
		return m.getFilesByResearchIDFunc(ctx, researchID)
	}
	return []domain.ResearchFile{}, nil
}

func (m *mockRepo) UpdateFile(ctx context.Context, file *domain.ResearchFile) error {
	if m.updateFileFunc != nil {
		return m.updateFileFunc(ctx, file)
	}
	return nil
}

func (m *mockRepo) DeleteFile(ctx context.Context, id uint) error {
	if m.deleteFileFunc != nil {
		return m.deleteFileFunc(ctx, id)
	}
	return nil
}

func (m *mockRepo) CreateReference(ctx context.Context, ref *domain.ResearchReference) error {
	if m.createRefFunc != nil {
		return m.createRefFunc(ctx, ref)
	}
	ref.ID = 500
	return nil
}

func (m *mockRepo) GetReferenceByID(ctx context.Context, id uint) (*domain.ResearchReference, error) {
	if m.getRefByIDFunc != nil {
		return m.getRefByIDFunc(ctx, id)
	}
	return &domain.ResearchReference{ID: id, ResearchID: 1, OrganisationID: 10}, nil
}

func (m *mockRepo) GetReferencesByDocumentID(ctx context.Context, docID uint) ([]domain.ResearchReference, error) {
	if m.getRefsByDocIDFunc != nil {
		return m.getRefsByDocIDFunc(ctx, docID)
	}
	return []domain.ResearchReference{{ID: 500}}, nil
}

func (m *mockRepo) GetReferencesByResearchID(ctx context.Context, researchID uint) ([]domain.ResearchReference, error) {
	if m.getRefsByResearchIDFunc != nil {
		return m.getRefsByResearchIDFunc(ctx, researchID)
	}
	return []domain.ResearchReference{}, nil
}

func (m *mockRepo) UpdateReference(ctx context.Context, ref *domain.ResearchReference) error {
	if m.updateRefFunc != nil {
		return m.updateRefFunc(ctx, ref)
	}
	return nil
}

func (m *mockRepo) DeleteReference(ctx context.Context, id uint) error {
	if m.deleteRefFunc != nil {
		return m.deleteRefFunc(ctx, id)
	}
	return nil
}

type mockProjectRepo struct {
	projects.Repository
	getByIDFunc func(ctx context.Context, id uint) (*projects.Project, error)
}

func (m *mockProjectRepo) GetByID(ctx context.Context, id uint) (*projects.Project, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return &projects.Project{ID: id, OrganisationID: 10}, nil
}

type mockTeamRepo struct {
	teams.Repository
	getByIDFunc          func(ctx context.Context, id uint) (*teams.Team, error)
	getMembersByTeamFunc func(ctx context.Context, teamID uint) ([]teams.TeamMember, error)
}

func (m *mockTeamRepo) GetByID(ctx context.Context, id uint) (*teams.Team, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return &teams.Team{ID: id, OrganisationID: 10}, nil
}

func (m *mockTeamRepo) GetMembersByTeam(ctx context.Context, teamID uint) ([]teams.TeamMember, error) {
	if m.getMembersByTeamFunc != nil {
		return m.getMembersByTeamFunc(ctx, teamID)
	}
	return []teams.TeamMember{{UserID: 2}}, nil
}

type mockNotificationUc struct {
	notifications.UseCase
	createNotificationFunc func(ctx context.Context, userID uint, orgID uint, title string, message string, notifType string, link string) error
}

func (m *mockNotificationUc) CreateNotification(ctx context.Context, userID uint, orgID uint, title string, message string, notifType string, link string) error {
	if m.createNotificationFunc != nil {
		return m.createNotificationFunc(ctx, userID, orgID, title, message, notifType, link)
	}
	return nil
}

func TestCreate(t *testing.T) {
	repo := &mockRepo{}
	pRepo := &mockProjectRepo{}
	tRepo := &mockTeamRepo{}
	nUc := &mockNotificationUc{}
	uc := NewUseCase(repo, pRepo, tRepo, nUc)

	ctx := context.Background()

	t.Run("empty title", func(t *testing.T) {
		_, err := uc.Create(ctx, domain.CreateRequest{Title: "", Status: "draft"})
		if err == nil || err.Error() != "title is required" {
			t.Errorf("expected title is required error, got %v", err)
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		_, err := uc.Create(ctx, domain.CreateRequest{Title: "Test", Status: "invalid"})
		if !errors.Is(err, domain.ErrInvalidStatus) {
			t.Errorf("expected ErrInvalidStatus, got %v", err)
		}
	})

	t.Run("invalid project id", func(t *testing.T) {
		pID := uint(5)
		pRepoFail := &mockProjectRepo{
			getByIDFunc: func(ctx context.Context, id uint) (*projects.Project, error) {
				return nil, errors.New("not found")
			},
		}
		ucFail := NewUseCase(repo, pRepoFail, tRepo, nUc)
		_, err := ucFail.Create(ctx, domain.CreateRequest{Title: "Test", Status: "draft", ProjectID: &pID, OrganisationID: 10})
		if err == nil || err.Error() != "invalid project id" {
			t.Errorf("expected invalid project id error, got %v", err)
		}
	})

	t.Run("invalid team id", func(t *testing.T) {
		tID := uint(5)
		tRepoFail := &mockTeamRepo{
			getByIDFunc: func(ctx context.Context, id uint) (*teams.Team, error) {
				return nil, errors.New("not found")
			},
		}
		ucFail := NewUseCase(repo, pRepo, tRepoFail, nUc)
		_, err := ucFail.Create(ctx, domain.CreateRequest{Title: "Test", Status: "draft", TeamID: &tID, OrganisationID: 10})
		if err == nil || err.Error() != "invalid team id" {
			t.Errorf("expected invalid team id error, got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		pID := uint(1)
		tID := uint(2)
		entry, err := uc.Create(ctx, domain.CreateRequest{
			Title:          "Test Research",
			Status:         "draft",
			ProjectID:      &pID,
			TeamID:         &tID,
			OrganisationID: 10,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if entry.Title != "Test Research" {
			t.Errorf("expected title 'Test Research', got %s", entry.Title)
		}
	})
}

func TestGetByIDAndList(t *testing.T) {
	repo := &mockRepo{}
	uc := NewUseCase(repo, &mockProjectRepo{}, &mockTeamRepo{}, &mockNotificationUc{})
	ctx := context.Background()

	t.Run("GetByID", func(t *testing.T) {
		res, err := uc.GetByID(ctx, 1, 10)
		if err != nil || res.ID != 1 {
			t.Errorf("unexpected get by id result: %v, %v", res, err)
		}
	})

	t.Run("List", func(t *testing.T) {
		items, total, err := uc.List(ctx, domain.ListRequest{OrganisationID: 10, Limit: 10, Offset: 0})
		if err != nil || total != 1 || len(items) != 1 {
			t.Errorf("unexpected list result: items=%v total=%d err=%v", items, total, err)
		}
	})
}

func TestUpdateAndDelete(t *testing.T) {
	repo := &mockRepo{}
	pRepo := &mockProjectRepo{}
	tRepo := &mockTeamRepo{}
	nUc := &mockNotificationUc{}
	uc := NewUseCase(repo, pRepo, tRepo, nUc)
	ctx := context.Background()

	t.Run("Update not found", func(t *testing.T) {
		repoFail := &mockRepo{
			getByIDFunc: func(ctx context.Context, id uint, orgID uint) (*domain.ResearchEntry, error) {
				return nil, domain.ErrNotFound
			},
		}
		ucFail := NewUseCase(repoFail, pRepo, tRepo, nUc)
		err := ucFail.Update(ctx, domain.UpdateRequest{ID: 1}, 1, "user", 10)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("Update unauthorized", func(t *testing.T) {
		repoUnauth := &mockRepo{
			getByIDFunc: func(ctx context.Context, id uint, orgID uint) (*domain.ResearchEntry, error) {
				return &domain.ResearchEntry{ID: id, AuthorID: 99}, nil
			},
		}
		ucUnauth := NewUseCase(repoUnauth, pRepo, tRepo, nUc)
		err := ucUnauth.Update(ctx, domain.UpdateRequest{ID: 1}, 1, "employee", 10)
		if !errors.Is(err, domain.ErrUnauthorized) {
			t.Errorf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("Update invalid status", func(t *testing.T) {
		err := uc.Update(ctx, domain.UpdateRequest{ID: 1, Status: "invalid"}, 1, "user", 10)
		if !errors.Is(err, domain.ErrInvalidStatus) {
			t.Errorf("expected ErrInvalidStatus, got %v", err)
		}
	})

	t.Run("Update success and notify", func(t *testing.T) {
		tID := uint(2)
		repoNotify := &mockRepo{
			getByIDFunc: func(ctx context.Context, id uint, orgID uint) (*domain.ResearchEntry, error) {
				return &domain.ResearchEntry{ID: id, AuthorID: 1, Status: "draft", TeamID: &tID, OrganisationID: 10, Title: "Res"}, nil
			},
		}
		ucNotify := NewUseCase(repoNotify, pRepo, tRepo, nUc)
		err := ucNotify.Update(ctx, domain.UpdateRequest{ID: 1, Status: "review"}, 1, "manager", 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Delete unauthorized", func(t *testing.T) {
		repoUnauth := &mockRepo{
			getByIDFunc: func(ctx context.Context, id uint, orgID uint) (*domain.ResearchEntry, error) {
				return &domain.ResearchEntry{ID: id, AuthorID: 99}, nil
			},
		}
		ucUnauth := NewUseCase(repoUnauth, pRepo, tRepo, nUc)
		err := ucUnauth.Delete(ctx, 1, 1, "employee", 10)
		if !errors.Is(err, domain.ErrUnauthorized) {
			t.Errorf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("Delete success", func(t *testing.T) {
		err := uc.Delete(ctx, 1, 1, "user", 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestFolderOperations(t *testing.T) {
	repo := &mockRepo{}
	uc := NewUseCase(repo, &mockProjectRepo{}, &mockTeamRepo{}, &mockNotificationUc{})
	ctx := context.Background()

	t.Run("CreateFolder empty name", func(t *testing.T) {
		_, err := uc.CreateFolder(ctx, domain.CreateFolderRequest{ResearchID: 1, Name: "   "}, 1, 10)
		if err == nil || err.Error() != "folder name required" {
			t.Errorf("expected folder name required error, got %v", err)
		}
	})

	t.Run("CreateFolder duplicate name", func(t *testing.T) {
		repoDup := &mockRepo{
			folderNameExistsFunc: func(ctx context.Context, researchID uint, parentID *uint, name string, excludeID *uint) (bool, error) {
				return true, nil
			},
		}
		ucDup := NewUseCase(repoDup, &mockProjectRepo{}, &mockTeamRepo{}, &mockNotificationUc{})
		_, err := ucDup.CreateFolder(ctx, domain.CreateFolderRequest{ResearchID: 1, Name: "Docs"}, 1, 10)
		if !errors.Is(err, domain.ErrDuplicateFolderName) {
			t.Errorf("expected ErrDuplicateFolderName, got %v", err)
		}
	})

	t.Run("CreateFolder success", func(t *testing.T) {
		folder, err := uc.CreateFolder(ctx, domain.CreateFolderRequest{ResearchID: 1, Name: "Docs"}, 1, 10)
		if err != nil || folder.Name != "Docs" {
			t.Errorf("unexpected result: %v, %v", folder, err)
		}
	})

	t.Run("GetFolderByID and GetFoldersByResearchID", func(t *testing.T) {
		f, err := uc.GetFolderByID(ctx, 10, 10)
		if err != nil || f.ID != 10 {
			t.Errorf("unexpected GetFolderByID result: %v, %v", f, err)
		}

		fs, err := uc.GetFoldersByResearchID(ctx, 1, nil, 10)
		if err != nil || len(fs) != 1 {
			t.Errorf("unexpected GetFoldersByResearchID result: %v, %v", fs, err)
		}
	})

	t.Run("UpdateFolder and DeleteFolder", func(t *testing.T) {
		err := uc.UpdateFolder(ctx, domain.UpdateFolderRequest{ID: 10, Name: "New Docs"}, 1, 10)
		if err != nil {
			t.Errorf("unexpected UpdateFolder error: %v", err)
		}

		err = uc.DeleteFolder(ctx, 10, 1, 10)
		if err != nil {
			t.Errorf("unexpected DeleteFolder error: %v", err)
		}
	})
}

func TestDocumentOperations(t *testing.T) {
	repo := &mockRepo{}
	uc := NewUseCase(repo, &mockProjectRepo{}, &mockTeamRepo{}, &mockNotificationUc{})
	ctx := context.Background()

	t.Run("CreateDocument empty title", func(t *testing.T) {
		_, err := uc.CreateDocument(ctx, domain.CreateDocumentRequest{ResearchID: 1, Title: ""}, 1, 10)
		if err == nil || err.Error() != "document title required" {
			t.Errorf("expected document title required error, got %v", err)
		}
	})

	t.Run("CreateDocument duplicate title", func(t *testing.T) {
		repoDup := &mockRepo{
			docTitleExistsFunc: func(ctx context.Context, researchID uint, folderID *uint, title string, excludeID *uint) (bool, error) {
				return true, nil
			},
		}
		ucDup := NewUseCase(repoDup, &mockProjectRepo{}, &mockTeamRepo{}, &mockNotificationUc{})
		_, err := ucDup.CreateDocument(ctx, domain.CreateDocumentRequest{ResearchID: 1, Title: "Doc1"}, 1, 10)
		if !errors.Is(err, domain.ErrDuplicateDocTitle) {
			t.Errorf("expected ErrDuplicateDocTitle, got %v", err)
		}
	})

	t.Run("CreateDocument success", func(t *testing.T) {
		doc, err := uc.CreateDocument(ctx, domain.CreateDocumentRequest{ResearchID: 1, Title: "Doc1"}, 1, 10)
		if err != nil || doc.Title != "Doc1" {
			t.Errorf("unexpected result: %v, %v", doc, err)
		}
	})

	t.Run("GetDocumentByID and GetDocumentsByFolderID", func(t *testing.T) {
		doc, err := uc.GetDocumentByID(ctx, 100, 10)
		if err != nil || doc.ID != 100 {
			t.Errorf("unexpected GetDocumentByID result: %v, %v", doc, err)
		}

		docs, err := uc.GetDocumentsByFolderID(ctx, 1, nil, 10)
		if err != nil || len(docs) != 1 {
			t.Errorf("unexpected GetDocumentsByFolderID result: %v, %v", docs, err)
		}
	})

	t.Run("UpdateDocument and DeleteDocument", func(t *testing.T) {
		err := uc.UpdateDocument(ctx, domain.UpdateDocumentRequest{ID: 100, Title: "Updated Doc"}, 1, 10)
		if err != nil {
			t.Errorf("unexpected UpdateDocument error: %v", err)
		}

		err = uc.DeleteDocument(ctx, 100, 1, 10)
		if err != nil {
			t.Errorf("unexpected DeleteDocument error: %v", err)
		}
	})
}

func TestFileAndReferenceOperations(t *testing.T) {
	repo := &mockRepo{}
	uc := NewUseCase(repo, &mockProjectRepo{}, &mockTeamRepo{}, &mockNotificationUc{})
	ctx := context.Background()

	t.Run("UploadFile success", func(t *testing.T) {
		file, err := uc.UploadFile(ctx, domain.CreateFileRequest{ResearchID: 1, FileName: "test.pdf"}, 1, 10)
		if err != nil || file.FileName != "test.pdf" {
			t.Errorf("unexpected UploadFile result: %v, %v", file, err)
		}
	})

	t.Run("GetFileByID unauthorized", func(t *testing.T) {
		repoUnauth := &mockRepo{
			getFileByIDFunc: func(ctx context.Context, id uint) (*domain.ResearchFile, error) {
				return &domain.ResearchFile{ID: id, OrganisationID: 999}, nil
			},
		}
		ucUnauth := NewUseCase(repoUnauth, &mockProjectRepo{}, &mockTeamRepo{}, &mockNotificationUc{})
		_, err := ucUnauth.GetFileByID(ctx, 1000, 10)
		if !errors.Is(err, domain.ErrUnauthorized) {
			t.Errorf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("DeleteFile success", func(t *testing.T) {
		err := uc.DeleteFile(ctx, 1000, 1, 10)
		if err != nil {
			t.Errorf("unexpected DeleteFile error: %v", err)
		}
	})

	t.Run("GetFilesByDocumentID success", func(t *testing.T) {
		files, err := uc.GetFilesByDocumentID(ctx, 100, 10)
		if err != nil || len(files) != 1 {
			t.Errorf("unexpected GetFilesByDocumentID result: %v, %v", files, err)
		}
	})

	t.Run("AddReference and DeleteReference", func(t *testing.T) {
		ref, err := uc.AddReference(ctx, domain.CreateReferenceRequest{ResearchID: 1, Title: "Ref1", URL: "https://example.com"}, 1, 10)
		if err != nil || ref.Title != "Ref1" {
			t.Errorf("unexpected AddReference result: %v, %v", ref, err)
		}

		err = uc.DeleteReference(ctx, 500, 1, 10)
		if err != nil {
			t.Errorf("unexpected DeleteReference error: %v", err)
		}

		refs, err := uc.GetReferencesByDocumentID(ctx, 100, 10)
		if err != nil || len(refs) != 1 {
			t.Errorf("unexpected GetReferencesByDocumentID result: %v, %v", refs, err)
		}
	})
}
