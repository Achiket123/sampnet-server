package research

import (
	"context"
	"errors"
	domain "server/internal/domain/research"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) domain.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, entry *domain.ResearchEntry) error {
	m := toModel(entry)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	entry.ID = m.ID
	entry.CreatedAt = m.CreatedAt
	entry.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint, orgID uint) (*domain.ResearchEntry, error) {
	var m models.ResearchEntry
	if err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Project").
		Preload("Team").
		Where("id = ? AND organisation_id = ?", id, orgID).
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *gormRepository) List(ctx context.Context, orgID uint, filters domain.ListFilters, limit int, offset int) ([]domain.ResearchEntry, int, error) {
	var list []models.ResearchEntry
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ResearchEntry{}).Where("organisation_id = ?", orgID)

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.ProjectID != nil {
		query = query.Where("project_id = ?", *filters.ProjectID)
	}
	if filters.TeamID != nil {
		query = query.Where("team_id = ?", *filters.TeamID)
	}
	if filters.Search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Author").
		Preload("Project").
		Preload("Team").
		Order("updated_at desc").
		Limit(limit).
		Offset(offset).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return toDomainList(list), int(total), nil
}

func (r *gormRepository) Update(ctx context.Context, entry *domain.ResearchEntry) error {
	m := toModel(entry)
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return err
	}
	entry.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint, orgID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND organisation_id = ?", id, orgID).Delete(&models.ResearchEntry{}).Error
}

func (r *gormRepository) FolderNameExists(ctx context.Context, researchID uint, parentID *uint, name string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.ResearchFolder{}).Where("research_id = ? AND name = ?", researchID, name)
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *gormRepository) DocumentTitleExists(ctx context.Context, researchID uint, folderID *uint, title string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.ResearchDocument{}).Where("research_id = ? AND title = ?", researchID, title)
	if folderID != nil {
		query = query.Where("folder_id = ?", *folderID)
	} else {
		query = query.Where("folder_id IS NULL")
	}
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func toModel(d *domain.ResearchEntry) *models.ResearchEntry {
	m := &models.ResearchEntry{
		Title:          d.Title,
		Description:    d.Description,
		Thumbnail:      d.Thumbnail,
		Status:         d.Status,
		Visibility:     d.Visibility,
		CreatedBy:      d.AuthorID, // Mapping domain AuthorID to DB CreatedBy
		ProjectID:      d.ProjectID,
		TeamID:         d.TeamID,
		OrganisationID: d.OrganisationID,
		Tags:           models.TagsList(d.Tags),
	}
	m.ID = d.ID
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt
	return m
}

func toDomain(m *models.ResearchEntry) *domain.ResearchEntry {
	authorName := ""
	if m.Author.ID != 0 {
		authorName = m.Author.FirstName + " " + m.Author.LastName
	}

	projectName := ""
	if m.Project.ID != 0 {
		projectName = m.Project.Name
	}

	teamName := ""
	if m.Team.ID != 0 {
		teamName = m.Team.Name
	}

	return &domain.ResearchEntry{
		ID:             m.ID,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		Title:          m.Title,
		Description:    m.Description,
		Thumbnail:      m.Thumbnail,
		Status:         m.Status,
		Visibility:     m.Visibility,
		AuthorID:       m.CreatedBy, // Mapping DB CreatedBy to domain AuthorID
		AuthorName:     authorName,
		ProjectID:      m.ProjectID,
		ProjectName:    projectName,
		TeamID:         m.TeamID,
		TeamName:       teamName,
		OrganisationID: m.OrganisationID,
		Tags:           []string(m.Tags),
	}
}

func toDomainList(list []models.ResearchEntry) []domain.ResearchEntry {
	res := make([]domain.ResearchEntry, len(list))
	for i, m := range list {
		res[i] = *toDomain(&m)
	}
	return res
}

// --- Folder Operations ---

func (r *gormRepository) CreateFolder(ctx context.Context, folder *domain.ResearchFolder) error {
	m := &models.ResearchFolder{
		ResearchID:     folder.ResearchID,
		OrganisationID: folder.OrganisationID,
		ParentID:       folder.ParentID,
		Name:           folder.Name,
		CreatedBy:      folder.CreatedBy,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	folder.ID = m.ID
	folder.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *gormRepository) GetFolderByID(ctx context.Context, id uint) (*domain.ResearchFolder, error) {
	var m models.ResearchFolder
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &domain.ResearchFolder{
		ID:             m.ID,
		ResearchID:     m.ResearchID,
		OrganisationID: m.OrganisationID,
		ParentID:       m.ParentID,
		Name:           m.Name,
		CreatedBy:      m.CreatedBy,
		UpdatedAt:      m.UpdatedAt,
	}, nil
}

func (r *gormRepository) GetFoldersByResearchID(ctx context.Context, researchID uint, parentID *uint) ([]domain.ResearchFolder, error) {
	var list []models.ResearchFolder
	query := r.db.WithContext(ctx).Where("research_id = ?", researchID)
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}

	res := make([]domain.ResearchFolder, len(list))
	for i, m := range list {
		res[i] = domain.ResearchFolder{
			ID:             m.ID,
			ResearchID:     m.ResearchID,
			OrganisationID: m.OrganisationID,
			ParentID:       m.ParentID,
			Name:           m.Name,
			CreatedBy:      m.CreatedBy,
			UpdatedAt:      m.UpdatedAt,
		}
	}
	return res, nil
}

func (r *gormRepository) UpdateFolder(ctx context.Context, folder *domain.ResearchFolder) error {
	return r.db.WithContext(ctx).Model(&models.ResearchFolder{}).Where("id = ?", folder.ID).Updates(map[string]interface{}{
		"name":      folder.Name,
		"parent_id": folder.ParentID,
	}).Error
}

func (r *gormRepository) DeleteFolder(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ResearchFolder{}, id).Error
}

// --- Document Operations ---

func (r *gormRepository) CreateDocument(ctx context.Context, doc *domain.ResearchDocument) error {
	m := &models.ResearchDocument{
		ResearchID:     doc.ResearchID,
		OrganisationID: doc.OrganisationID,
		FolderID:       doc.FolderID,
		Title:          doc.Title,
		Content:        doc.Content,
		Status:         doc.Status,
		CreatedBy:      doc.CreatedBy,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	doc.ID = m.ID
	doc.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *gormRepository) GetDocumentByID(ctx context.Context, id uint) (*domain.ResearchDocument, error) {
	var m models.ResearchDocument
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &domain.ResearchDocument{
		ID:             m.ID,
		ResearchID:     m.ResearchID,
		OrganisationID: m.OrganisationID,
		FolderID:       m.FolderID,
		Title:          m.Title,
		Content:        m.Content,
		IsPinned:       m.IsPinned,
		Status:         m.Status,
		CreatedBy:      m.CreatedBy,
		UpdatedAt:      m.UpdatedAt,
	}, nil
}

func (r *gormRepository) GetDocumentsByFolderID(ctx context.Context, researchID uint, folderID *uint) ([]domain.ResearchDocument, error) {
	var list []models.ResearchDocument
	query := r.db.WithContext(ctx).Where("research_id = ?", researchID)
	if folderID != nil {
		query = query.Where("folder_id = ?", *folderID)
	} else {
		query = query.Where("folder_id IS NULL")
	}
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}

	res := make([]domain.ResearchDocument, len(list))
	for i, m := range list {
		res[i] = domain.ResearchDocument{
			ID:             m.ID,
			ResearchID:     m.ResearchID,
			OrganisationID: m.OrganisationID,
			FolderID:       m.FolderID,
			Title:          m.Title,
			Content:        m.Content,
			IsPinned:       m.IsPinned,
			Status:         m.Status,
			CreatedBy:      m.CreatedBy,
			UpdatedAt:      m.UpdatedAt,
		}
	}
	return res, nil
}

func (r *gormRepository) UpdateDocument(ctx context.Context, doc *domain.ResearchDocument) error {
	return r.db.WithContext(ctx).Model(&models.ResearchDocument{}).Where("id = ?", doc.ID).Updates(map[string]interface{}{
		"title":     doc.Title,
		"content":   doc.Content,
		"folder_id": doc.FolderID,
	}).Error
}

func (r *gormRepository) DeleteDocument(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ResearchDocument{}, id).Error
}

// --- File Operations ---

func (r *gormRepository) CreateFile(ctx context.Context, file *domain.ResearchFile) error {
	m := &models.ResearchFile{
		ResearchID:     file.ResearchID,
		DocumentID:     file.DocumentID,
		FolderID:       file.FolderID,
		OrganisationID: file.OrganisationID,
		FileName:       file.FileName,
		OriginalName:   file.FileName,
		MimeType:       file.MimeType,
		Size:           file.Size,
		StoragePath:    file.StoragePath,
		CreatedBy:      file.CreatedBy,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	file.ID = m.ID
	file.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *gormRepository) GetFileByID(ctx context.Context, id uint) (*domain.ResearchFile, error) {
	var m models.ResearchFile
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &domain.ResearchFile{
		ID:             m.ID,
		ResearchID:     m.ResearchID,
		DocumentID:     m.DocumentID,
		OrganisationID: m.OrganisationID,
		FolderID:       m.FolderID,
		FileName:       m.FileName,
		MimeType:       m.MimeType,
		Size:           m.Size,
		StoragePath:    m.StoragePath,
		CreatedBy:      m.CreatedBy,
		UpdatedAt:      m.UpdatedAt,
	}, nil
}

func (r *gormRepository) GetFilesByDocumentID(ctx context.Context, docID uint) ([]domain.ResearchFile, error) {
	var list []models.ResearchFile
	if err := r.db.WithContext(ctx).Where("document_id = ?", docID).Find(&list).Error; err != nil {
		return nil, err
	}
	res := make([]domain.ResearchFile, len(list))
	for i, m := range list {
		res[i] = domain.ResearchFile{
			ID:             m.ID,
			ResearchID:     m.ResearchID,
			DocumentID:     m.DocumentID,
			OrganisationID: m.OrganisationID,
			FileName:       m.FileName,
			MimeType:       m.MimeType,
			Size:           m.Size,
			StoragePath:    m.StoragePath,
			CreatedBy:      m.CreatedBy,
			UpdatedAt:      m.UpdatedAt,
		}
	}
	return res, nil
}

func (r *gormRepository) GetFilesByResearchID(ctx context.Context, researchID uint) ([]domain.ResearchFile, error) {
	var list []models.ResearchFile
	if err := r.db.WithContext(ctx).Where("research_id = ? AND document_id IS NULL", researchID).Find(&list).Error; err != nil {
		return nil, err
	}
	res := make([]domain.ResearchFile, len(list))
	for i, m := range list {
		res[i] = domain.ResearchFile{
			ID:             m.ID,
			ResearchID:     m.ResearchID,
			DocumentID:     m.DocumentID,
			OrganisationID: m.OrganisationID,
			FileName:       m.FileName,
			MimeType:       m.MimeType,
			Size:           m.Size,
			StoragePath:    m.StoragePath,
			CreatedBy:      m.CreatedBy,
			UpdatedAt:      m.UpdatedAt,
		}
	}
	return res, nil
}

func (r *gormRepository) UpdateFile(ctx context.Context, file *domain.ResearchFile) error {
	return r.db.WithContext(ctx).Model(&models.ResearchFile{}).Where("id = ?", file.ID).Updates(map[string]interface{}{
		"file_name": file.FileName,
		"folder_id": file.FolderID,
		"is_pinned": file.IsPinned,
	}).Error
}

func (r *gormRepository) DeleteFile(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ResearchFile{}, id).Error
}

// --- Reference Operations ---

func (r *gormRepository) CreateReference(ctx context.Context, ref *domain.ResearchReference) error {
	m := &models.ResearchReference{
		ResearchID:     ref.ResearchID,
		DocumentID:     ref.DocumentID,
		OrganisationID: ref.OrganisationID,
		Title:          ref.Title,
		URL:            ref.URL,
		Authors:        ref.Authors,
		CreatedBy:      ref.CreatedBy,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	ref.ID = m.ID
	return nil
}

func (r *gormRepository) GetReferenceByID(ctx context.Context, id uint) (*domain.ResearchReference, error) {
	var m models.ResearchReference
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &domain.ResearchReference{
		ID:             m.ID,
		ResearchID:     m.ResearchID,
		DocumentID:     m.DocumentID,
		OrganisationID: m.OrganisationID,
		Title:          m.Title,
		URL:            m.URL,
		Authors:        m.Authors,
		CreatedBy:      m.CreatedBy,
	}, nil
}

func (r *gormRepository) GetReferencesByDocumentID(ctx context.Context, docID uint) ([]domain.ResearchReference, error) {
	var list []models.ResearchReference
	if err := r.db.WithContext(ctx).Where("document_id = ?", docID).Find(&list).Error; err != nil {
		return nil, err
	}
	res := make([]domain.ResearchReference, len(list))
	for i, m := range list {
		res[i] = domain.ResearchReference{
			ID:             m.ID,
			ResearchID:     m.ResearchID,
			DocumentID:     m.DocumentID,
			OrganisationID: m.OrganisationID,
			Title:          m.Title,
			URL:            m.URL,
			Authors:        m.Authors,
			CreatedBy:      m.CreatedBy,
		}
	}
	return res, nil
}

func (r *gormRepository) GetReferencesByResearchID(ctx context.Context, researchID uint) ([]domain.ResearchReference, error) {
	var list []models.ResearchReference
	if err := r.db.WithContext(ctx).Where("research_id = ? AND document_id IS NULL", researchID).Find(&list).Error; err != nil {
		return nil, err
	}
	res := make([]domain.ResearchReference, len(list))
	for i, m := range list {
		res[i] = domain.ResearchReference{
			ID:             m.ID,
			ResearchID:     m.ResearchID,
			DocumentID:     m.DocumentID,
			OrganisationID: m.OrganisationID,
			Title:          m.Title,
			URL:            m.URL,
			Authors:        m.Authors,
			CreatedBy:      m.CreatedBy,
		}
	}
	return res, nil
}

func (r *gormRepository) UpdateReference(ctx context.Context, ref *domain.ResearchReference) error {
	return r.db.WithContext(ctx).Model(&models.ResearchReference{}).Where("id = ?", ref.ID).Updates(map[string]interface{}{
		"title":   ref.Title,
		"url":     ref.URL,
		"authors": ref.Authors,
	}).Error
}

func (r *gormRepository) DeleteReference(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ResearchReference{}, id).Error
}
