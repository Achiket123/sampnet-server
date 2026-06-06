package research

import (
	"context"
	"errors"
	"fmt"
	"log"
	"server/internal/domain/notifications"
	"server/internal/domain/projects"
	domain "server/internal/domain/research"
	"server/internal/domain/teams"
	"strings"
)

type useCase struct {
	repo           domain.Repository
	projectRepo    projects.Repository
	teamRepo       teams.Repository
	notificationUc notifications.UseCase
}

func NewUseCase(repo domain.Repository, projectRepo projects.Repository, teamRepo teams.Repository, notificationUc notifications.UseCase) domain.UseCase {
	return &useCase{
		repo:           repo,
		projectRepo:    projectRepo,
		teamRepo:       teamRepo,
		notificationUc: notificationUc,
	}
}

func (u *useCase) Create(ctx context.Context, req domain.CreateRequest) (*domain.ResearchEntry, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	if !isValidStatus(req.Status) {
		return nil, domain.ErrInvalidStatus
	}

	if req.ProjectID != nil {
		p, err := u.projectRepo.GetByID(ctx, *req.ProjectID)
		if err != nil || p.OrganisationID != req.OrganisationID {
			return nil, errors.New("invalid project id")
		}
	}

	if req.TeamID != nil {
		t, err := u.teamRepo.GetByID(ctx, *req.TeamID)
		if err != nil || t.OrganisationID != req.OrganisationID {
			return nil, errors.New("invalid team id")
		}
	}

	entry := &domain.ResearchEntry{
		Title:          req.Title,
		Description:    req.Description,
		Status:         req.Status,
		AuthorID:       req.AuthorID,
		ProjectID:      req.ProjectID,
		TeamID:         req.TeamID,
		OrganisationID: req.OrganisationID,
		Tags:           req.Tags,
	}

	if err := u.repo.Create(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}

func (u *useCase) GetByID(ctx context.Context, id uint, orgID uint) (*domain.ResearchEntry, error) {
	return u.repo.GetByID(ctx, id, orgID)
}

func (u *useCase) List(ctx context.Context, req domain.ListRequest) ([]domain.ResearchEntry, int, error) {
	filters := domain.ListFilters{
		Status:    req.Status,
		ProjectID: req.ProjectID,
		TeamID:    req.TeamID,
		Search:    req.Search,
	}
	return u.repo.List(ctx, req.OrganisationID, filters, req.Limit, req.Offset)
}

func (u *useCase) Update(ctx context.Context, req domain.UpdateRequest, actorID uint, actorRole string, orgID uint) error {
	entry, err := u.repo.GetByID(ctx, req.ID, orgID)
	if err != nil {
		return err
	}

	if !u.isAuthorized(entry, actorID, actorRole) {
		return domain.ErrUnauthorized
	}

	oldStatus := entry.Status

	if req.Title != "" {
		entry.Title = req.Title
	}

	if req.Status != "" {
		if !isValidStatus(req.Status) {
			return domain.ErrInvalidStatus
		}
		entry.Status = req.Status
	}

	entry.Description = req.Description
	entry.Tags = req.Tags

	if req.ProjectID != nil {
		p, err := u.projectRepo.GetByID(ctx, *req.ProjectID)
		if err != nil || p.OrganisationID != orgID {
			return errors.New("invalid project id")
		}
		entry.ProjectID = req.ProjectID
	} else {
		entry.ProjectID = nil
	}

	if req.TeamID != nil {
		t, err := u.teamRepo.GetByID(ctx, *req.TeamID)
		if err != nil || t.OrganisationID != orgID {
			return errors.New("invalid team id")
		}
		entry.TeamID = req.TeamID
	} else {
		entry.TeamID = nil
	}

	if err := u.repo.Update(ctx, entry); err != nil {
		return err
	}

	if (entry.Status == "review" || entry.Status == "published") && entry.Status != oldStatus {
		u.notifyStakeholders(ctx, entry)
	}

	return nil
}

func (u *useCase) Delete(ctx context.Context, id uint, actorID uint, actorRole string, orgID uint) error {
	entry, err := u.repo.GetByID(ctx, id, orgID)
	if err != nil {
		return err
	}

	if !u.isAuthorized(entry, actorID, actorRole) {
		return domain.ErrUnauthorized
	}

	return u.repo.Delete(ctx, id, orgID)
}

// --- Folder Operations ---

func (u *useCase) CreateFolder(ctx context.Context, req domain.CreateFolderRequest, actorID uint, orgID uint) (*domain.ResearchFolder, error) {
	entry, err := u.GetByID(ctx, req.ResearchID, orgID)
	if err != nil {
		return nil, err
	}
	if !u.isAuthorized(entry, actorID, "") { // simplistic check for now
		return nil, domain.ErrUnauthorized
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, errors.New("folder name required")
	}
	if err := u.validateFolderParent(ctx, req.ResearchID, req.ParentID); err != nil {
		return nil, err
	}
	exists, err := u.repo.FolderNameExists(ctx, req.ResearchID, req.ParentID, req.Name, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrDuplicateFolderName
	}

	folder := &domain.ResearchFolder{
		ResearchID: req.ResearchID,
		ParentID:   req.ParentID,
		Name:       req.Name,
		CreatedBy:  req.CreatedBy,
	}
	if err := u.repo.CreateFolder(ctx, folder); err != nil {
		return nil, err
	}
	return folder, nil
}

func (u *useCase) GetFolderByID(ctx context.Context, id uint, orgID uint) (*domain.ResearchFolder, error) {
	// First get the folder
	folder, err := u.repo.GetFolderByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Verify org access via the parent research entry
	_, err = u.GetByID(ctx, folder.ResearchID, orgID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	return folder, nil
}

func (u *useCase) GetFoldersByResearchID(ctx context.Context, researchID uint, parentID *uint, orgID uint) ([]domain.ResearchFolder, error) {
	_, err := u.GetByID(ctx, researchID, orgID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	return u.repo.GetFoldersByResearchID(ctx, researchID, parentID)
}

func (u *useCase) UpdateFolder(ctx context.Context, req domain.UpdateFolderRequest, actorID uint, orgID uint) error {
	folder, err := u.GetFolderByID(ctx, req.ID, orgID)
	if err != nil {
		return err
	}
	entry, err := u.GetByID(ctx, folder.ResearchID, orgID)
	if err != nil {
		return err
	}
	if !u.isAuthorized(entry, actorID, "") {
		return domain.ErrUnauthorized
	}

	if req.Name != "" {
		folder.Name = strings.TrimSpace(req.Name)
		if folder.Name == "" {
			return errors.New("folder name required")
		}
	}
	if err := u.validateFolderMove(ctx, folder, req.ParentID); err != nil {
		return err
	}
	folder.ParentID = req.ParentID
	excludeID := folder.ID
	exists, err := u.repo.FolderNameExists(ctx, folder.ResearchID, folder.ParentID, folder.Name, &excludeID)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrDuplicateFolderName
	}

	return u.repo.UpdateFolder(ctx, folder)
}

func (u *useCase) DeleteFolder(ctx context.Context, id uint, actorID uint, orgID uint) error {
	folder, err := u.GetFolderByID(ctx, id, orgID)
	if err != nil {
		return err
	}
	entry, err := u.GetByID(ctx, folder.ResearchID, orgID)
	if err != nil {
		return err
	}
	if !u.isAuthorized(entry, actorID, "") {
		return domain.ErrUnauthorized
	}
	return u.repo.DeleteFolder(ctx, id)
}

// --- Document Operations ---

func (u *useCase) CreateDocument(ctx context.Context, req domain.CreateDocumentRequest, actorID uint, orgID uint) (*domain.ResearchDocument, error) {
	entry, err := u.GetByID(ctx, req.ResearchID, orgID)
	if err != nil {
		return nil, err
	}
	if !u.isAuthorized(entry, actorID, "") {
		return nil, domain.ErrUnauthorized
	}
	if req.Title == "" {
		return nil, errors.New("document title required")
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return nil, errors.New("document title required")
	}
	if err := u.validateFolderParent(ctx, req.ResearchID, req.FolderID); err != nil {
		return nil, err
	}
	exists, err := u.repo.DocumentTitleExists(ctx, req.ResearchID, req.FolderID, req.Title, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrDuplicateDocTitle
	}

	doc := &domain.ResearchDocument{
		ResearchID: req.ResearchID,
		FolderID:   req.FolderID,
		Title:      req.Title,
		Content:    req.Content,
		Status:     "active",
		CreatedBy:  req.CreatedBy,
	}
	if err := u.repo.CreateDocument(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (u *useCase) GetDocumentByID(ctx context.Context, id uint, orgID uint) (*domain.ResearchDocument, error) {
	doc, err := u.repo.GetDocumentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_, err = u.GetByID(ctx, doc.ResearchID, orgID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	return doc, nil
}

func (u *useCase) GetDocumentsByFolderID(ctx context.Context, researchID uint, folderID *uint, orgID uint) ([]domain.ResearchDocument, error) {
	_, err := u.GetByID(ctx, researchID, orgID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	return u.repo.GetDocumentsByFolderID(ctx, researchID, folderID)
}

func (u *useCase) UpdateDocument(ctx context.Context, req domain.UpdateDocumentRequest, actorID uint, orgID uint) error {
	doc, err := u.GetDocumentByID(ctx, req.ID, orgID)
	if err != nil {
		return err
	}
	entry, err := u.GetByID(ctx, doc.ResearchID, orgID)
	if err != nil {
		return err
	}
	if !u.isAuthorized(entry, actorID, "") {
		return domain.ErrUnauthorized
	}

	if req.Title != "" {
		doc.Title = strings.TrimSpace(req.Title)
		if doc.Title == "" {
			return errors.New("document title required")
		}
	}
	doc.Content = req.Content
	if err := u.validateFolderParent(ctx, doc.ResearchID, req.FolderID); err != nil {
		return err
	}
	doc.FolderID = req.FolderID
	excludeID := doc.ID
	exists, err := u.repo.DocumentTitleExists(ctx, doc.ResearchID, doc.FolderID, doc.Title, &excludeID)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrDuplicateDocTitle
	}

	return u.repo.UpdateDocument(ctx, doc)
}

func (u *useCase) DeleteDocument(ctx context.Context, id uint, actorID uint, orgID uint) error {
	doc, err := u.GetDocumentByID(ctx, id, orgID)
	if err != nil {
		return err
	}
	entry, err := u.GetByID(ctx, doc.ResearchID, orgID)
	if err != nil {
		return err
	}
	if !u.isAuthorized(entry, actorID, "") {
		return domain.ErrUnauthorized
	}
	return u.repo.DeleteDocument(ctx, id)
}

// --- Artifact Operations ---

func (u *useCase) UploadFile(ctx context.Context, req domain.CreateFileRequest, actorID uint, orgID uint) (*domain.ResearchFile, error) {
	_, err := u.GetByID(ctx, req.ResearchID, orgID)
	if err != nil {
		return nil, err
	}

	if req.DocumentID != nil {
		_, err = u.GetDocumentByID(ctx, *req.DocumentID, orgID)
		if err != nil {
			return nil, err
		}
	}

	file := &domain.ResearchFile{
		ResearchID:     req.ResearchID,
		DocumentID:     req.DocumentID,
		FolderID:       req.FolderID,
		OrganisationID: orgID,
		FileName:       req.FileName,
		MimeType:       req.MimeType,
		Size:           req.Size,
		StoragePath:    req.StoragePath,
		CreatedBy:      actorID,
	}

	if err := u.repo.CreateFile(ctx, file); err != nil {
		return nil, err
	}
	return file, nil
}

func (u *useCase) GetFileByID(ctx context.Context, id uint, orgID uint) (*domain.ResearchFile, error) {
	file, err := u.repo.GetFileByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if file.OrganisationID != orgID {
		return nil, domain.ErrUnauthorized
	}
	return file, nil
}

func (u *useCase) DeleteFile(ctx context.Context, id uint, actorID uint, orgID uint) error {
	file, err := u.repo.GetFileByID(ctx, id)
	if err != nil {
		return err
	}
	if file.OrganisationID != orgID {
		return domain.ErrUnauthorized
	}
	return u.repo.DeleteFile(ctx, id)
}

func (u *useCase) GetFilesByDocumentID(ctx context.Context, docID uint, orgID uint) ([]domain.ResearchFile, error) {
	_, err := u.GetDocumentByID(ctx, docID, orgID)
	if err != nil {
		return nil, err
	}
	return u.repo.GetFilesByDocumentID(ctx, docID)
}

func (u *useCase) AddReference(ctx context.Context, req domain.CreateReferenceRequest, actorID uint, orgID uint) (*domain.ResearchReference, error) {
	_, err := u.GetByID(ctx, req.ResearchID, orgID)
	if err != nil {
		return nil, err
	}

	if req.DocumentID != nil {
		_, err = u.GetDocumentByID(ctx, *req.DocumentID, orgID)
		if err != nil {
			return nil, err
		}
	}

	ref := &domain.ResearchReference{
		ResearchID:     req.ResearchID,
		DocumentID:     req.DocumentID,
		OrganisationID: orgID,
		Title:          req.Title,
		URL:            req.URL,
		Authors:        req.Authors,
	}

	if err := u.repo.CreateReference(ctx, ref); err != nil {
		return nil, err
	}
	return ref, nil
}

func (u *useCase) DeleteReference(ctx context.Context, id uint, actorID uint, orgID uint) error {
	ref, err := u.repo.GetReferenceByID(ctx, id)
	if err != nil {
		return err
	}
	if ref.OrganisationID != orgID {
		return domain.ErrUnauthorized
	}
	return u.repo.DeleteReference(ctx, id)
}

func (u *useCase) GetReferencesByDocumentID(ctx context.Context, docID uint, orgID uint) ([]domain.ResearchReference, error) {
	_, err := u.GetDocumentByID(ctx, docID, orgID)
	if err != nil {
		return nil, err
	}
	return u.repo.GetReferencesByDocumentID(ctx, docID)
}

func (u *useCase) validateFolderParent(ctx context.Context, researchID uint, parentID *uint) error {
	if parentID == nil {
		return nil
	}

	parent, err := u.repo.GetFolderByID(ctx, *parentID)
	if err != nil {
		return err
	}
	if parent.ResearchID != researchID {
		return domain.ErrInvalidFolderParent
	}
	return nil
}

func (u *useCase) validateFolderMove(ctx context.Context, folder *domain.ResearchFolder, nextParentID *uint) error {
	if err := u.validateFolderParent(ctx, folder.ResearchID, nextParentID); err != nil {
		return err
	}
	if nextParentID == nil {
		return nil
	}
	if *nextParentID == folder.ID {
		return domain.ErrInvalidFolderParent
	}

	currentID := nextParentID
	for currentID != nil {
		if *currentID == folder.ID {
			return domain.ErrInvalidFolderParent
		}
		parent, err := u.repo.GetFolderByID(ctx, *currentID)
		if err != nil {
			return err
		}
		currentID = parent.ParentID
	}
	return nil
}

func (u *useCase) isAuthorized(entry *domain.ResearchEntry, actorID uint, actorRole string) bool {
	if entry.AuthorID == actorID {
		return true
	}
	if actorRole == "manager" || actorRole == "boss" {
		return true
	}
	return false
}

func (u *useCase) notifyStakeholders(ctx context.Context, entry *domain.ResearchEntry) {
	var userIDs []uint
	if entry.TeamID != nil {
		members, err := u.teamRepo.GetMembersByTeam(ctx, *entry.TeamID)
		if err == nil {
			for _, m := range members {
				userIDs = append(userIDs, m.UserID)
			}
		}
	}

	title := "Research Status Updated"
	message := fmt.Sprintf("Research '%s' is now %s", entry.Title, entry.Status)
	link := fmt.Sprintf("/research/%d", entry.ID)

	for _, userID := range userIDs {
		if userID == entry.AuthorID {
			continue
		}
		go func(uid uint) {
			err := u.notificationUc.CreateNotification(context.Background(), uid, entry.OrganisationID, title, message, "research_status_updated", link)
			if err != nil {
				log.Printf("failed to create notification for research: %v", uid)
			}
		}(userID)
	}
}

func isValidStatus(status string) bool {
	switch status {
	case "draft", "in_progress", "review", "published":
		return true
	}
	return false
}
