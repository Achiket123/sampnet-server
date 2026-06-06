package people

import (
	"context"
	"fmt"
	"server/internal/domain/people"
	"server/internal/platform/database/models"
	"time"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) people.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) GetContacts(ctx context.Context, filter people.ContactsFilter, orgID uint) (*people.ContactsResponse, error) {
	var contacts []models.PeopleContact
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PeopleContact{}).Where("organisation_id = ?", orgID)

	if filter.Stage != nil && *filter.Stage != "" {
		query = query.Where("stage = ?", *filter.Stage)
	}
	if filter.Type != nil && *filter.Type != "" {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.ListID != nil {
		query = query.Where("id IN (SELECT contact_id FROM people_list_contacts WHERE list_id = ?)", *filter.ListID)
	}
	if filter.Search != "" {
		searchStr := "%" + filter.Search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR company ILIKE ?", searchStr, searchStr, searchStr, searchStr)
	}

	query.Count(&total)

	if filter.SortBy != "" {
		order := "DESC"
		if filter.SortOrder == "asc" {
			order = "ASC"
		}
		query = query.Order(fmt.Sprintf("%s %s", filter.SortBy, order))
	} else {
		query = query.Order("created_at DESC")
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.Preload("AssignedTo").Preload("AssignedTo.User").Offset(offset).Limit(filter.Limit).Find(&contacts).Error
	if err != nil {
		return nil, err
	}

	res := make([]people.PeopleContact, len(contacts))
	for i, c := range contacts {
		res[i] = toDomainContact(&c)
	}

	return &people.ContactsResponse{
		Contacts: res,
		Total:    int(total),
		Page:     filter.Page,
		Limit:    filter.Limit,
	}, nil
}

func (r *gormRepository) GetContactByID(ctx context.Context, id uint, orgID uint) (*people.PeopleContact, error) {
	var contact models.PeopleContact
	err := r.db.WithContext(ctx).Preload("AssignedTo").Preload("AssignedTo.User").Where("id = ? AND organisation_id = ?", id, orgID).First(&contact).Error
	if err != nil {
		return nil, err
	}
	c := toDomainContact(&contact)
	return &c, nil
}

func (r *gormRepository) CreateContact(ctx context.Context, contact *people.PeopleContact) error {
	m := fromDomainContact(contact)
	err := r.db.WithContext(ctx).Create(m).Error
	if err != nil {
		return err
	}
	contact.ID = m.ID
	contact.CreatedAt = m.CreatedAt
	contact.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *gormRepository) UpdateContact(ctx context.Context, contact *people.PeopleContact) error {
	m := fromDomainContact(contact)
	return r.db.WithContext(ctx).Model(m).Updates(m).Error
}

func (r *gormRepository) DeleteContact(ctx context.Context, id uint, orgID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND organisation_id = ?", id, orgID).Delete(&models.PeopleContact{}).Error
}

func (r *gormRepository) BulkUpdateStage(ctx context.Context, ids []uint, stage string, orgID uint, changedByID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.PeopleContact{}).Where("id IN ? AND organisation_id = ?", ids, orgID).Update("stage", stage).Error
		if err != nil {
			return err
		}
		
		for _, id := range ids {
			interaction := models.PeopleInteraction{
				ContactID:   id,
				CreatedByID: changedByID,
				Type:        "stage_change",
				Content:     fmt.Sprintf("Stage changed to %s", stage),
				OccurredAt:  time.Now(),
			}
			if err := tx.Create(&interaction).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *gormRepository) BulkDelete(ctx context.Context, ids []uint, orgID uint) error {
	return r.db.WithContext(ctx).Where("id IN ? AND organisation_id = ?", ids, orgID).Delete(&models.PeopleContact{}).Error
}

func (r *gormRepository) AssignContact(ctx context.Context, contactID uint, employeeID uint, orgID uint) error {
	return r.db.WithContext(ctx).Model(&models.PeopleContact{}).Where("id = ? AND organisation_id = ?", contactID, orgID).Update("assigned_to_id", employeeID).Error
}

func (r *gormRepository) GetInteractions(ctx context.Context, contactID uint, orgID uint) ([]people.PeopleInteraction, error) {
	var interactions []models.PeopleInteraction
	err := r.db.WithContext(ctx).
		Joins("JOIN people_contacts on people_contacts.id = people_interactions.contact_id").
		Preload("CreatedBy").Preload("CreatedBy.User").
		Where("people_interactions.contact_id = ? AND people_contacts.organisation_id = ?", contactID, orgID).
		Order("occurred_at DESC").Find(&interactions).Error
	if err != nil {
		return nil, err
	}

	res := make([]people.PeopleInteraction, len(interactions))
	for i, in := range interactions {
		res[i] = toDomainInteraction(&in)
	}
	return res, nil
}

func (r *gormRepository) AddInteraction(ctx context.Context, interaction *people.PeopleInteraction) error {
	m := fromDomainInteraction(interaction)
	err := r.db.WithContext(ctx).Create(m).Error
	if err != nil {
		return err
	}
	interaction.ID = m.ID
	interaction.CreatedAt = m.CreatedAt

	if interaction.Type == "call" || interaction.Type == "email" || interaction.Type == "meeting" {
		r.db.WithContext(ctx).Model(&models.PeopleContact{}).Where("id = ?", interaction.ContactID).Update("last_contacted_at", interaction.OccurredAt)
	}

	return nil
}

func (r *gormRepository) DeleteInteraction(ctx context.Context, interactionID uint, reqUserID uint, reqUserIsAdmin bool) error {
	query := r.db.WithContext(ctx).Where("id = ?", interactionID)
	if !reqUserIsAdmin {
		query = query.Where("created_by_id = ?", reqUserID)
	}
	return query.Delete(&models.PeopleInteraction{}).Error
}

func (r *gormRepository) GetPipelineStages(ctx context.Context, orgID uint) ([]people.PeoplePipelineStage, error) {
	var stages []models.PeoplePipelineStage
	err := r.db.WithContext(ctx).Where("organisation_id = ?", orgID).Order("\"order\" ASC").Find(&stages).Error
	if err != nil {
		return nil, err
	}
	res := make([]people.PeoplePipelineStage, len(stages))
	for i, s := range stages {
		res[i] = toDomainStage(&s)
	}
	return res, nil
}

func (r *gormRepository) CreatePipelineStage(ctx context.Context, stage *people.PeoplePipelineStage) error {
	m := fromDomainStage(stage)
	err := r.db.WithContext(ctx).Create(m).Error
	if err != nil {
		return err
	}
	stage.ID = m.ID
	return nil
}

func (r *gormRepository) UpdatePipelineStage(ctx context.Context, stage *people.PeoplePipelineStage) error {
	m := fromDomainStage(stage)
	return r.db.WithContext(ctx).Model(m).Updates(m).Error
}

func (r *gormRepository) DeletePipelineStage(ctx context.Context, id uint, orgID uint, reassignToStage *string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if reassignToStage != nil {
			var stage models.PeoplePipelineStage
			if err := tx.Where("id = ? AND organisation_id = ?", id, orgID).First(&stage).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.PeopleContact{}).Where("stage = ? AND organisation_id = ?", stage.Key, orgID).Update("stage", *reassignToStage).Error; err != nil {
				return err
			}
		} else {
			var stage models.PeoplePipelineStage
			if err := tx.Where("id = ? AND organisation_id = ?", id, orgID).First(&stage).Error; err != nil {
				return err
			}
			var count int64
			tx.Model(&models.PeopleContact{}).Where("stage = ? AND organisation_id = ?", stage.Key, orgID).Count(&count)
			if count > 0 {
				return fmt.Errorf("Cannot delete stage with active contacts")
			}
		}
		return tx.Where("id = ? AND organisation_id = ?", id, orgID).Delete(&models.PeoplePipelineStage{}).Error
	})
}

func (r *gormRepository) ReorderPipelineStages(ctx context.Context, orderedIDs []uint, orgID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, id := range orderedIDs {
			if err := tx.Model(&models.PeoplePipelineStage{}).Where("id = ? AND organisation_id = ?", id, orgID).Update("order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *gormRepository) GetLists(ctx context.Context, orgID uint) ([]people.PeopleList, error) {
	var lists []models.PeopleList
	err := r.db.WithContext(ctx).Where("organisation_id = ?", orgID).Find(&lists).Error
	if err != nil {
		return nil, err
	}
	res := make([]people.PeopleList, len(lists))
	for i, l := range lists {
		var count int64
		r.db.WithContext(ctx).Model(&models.PeopleListContact{}).Where("list_id = ?", l.ID).Count(&count)
		res[i] = toDomainList(&l)
		res[i].ContactCount = int(count)
	}
	return res, nil
}

func (r *gormRepository) CreateList(ctx context.Context, list *people.PeopleList) error {
	m := fromDomainList(list)
	err := r.db.WithContext(ctx).Create(m).Error
	if err != nil {
		return err
	}
	list.ID = m.ID
	list.CreatedAt = m.CreatedAt
	return nil
}

func (r *gormRepository) UpdateList(ctx context.Context, list *people.PeopleList) error {
	m := fromDomainList(list)
	return r.db.WithContext(ctx).Model(m).Updates(m).Error
}

func (r *gormRepository) DeleteList(ctx context.Context, id uint, orgID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND organisation_id = ?", id, orgID).Delete(&models.PeopleList{}).Error
}

func (r *gormRepository) RemoveContactFromList(ctx context.Context, listID uint, contactID uint, orgID uint) error {
	return r.db.WithContext(ctx).Where("list_id = ? AND contact_id = ?", listID, contactID).Delete(&models.PeopleListContact{}).Error
}

func (r *gormRepository) GetAnalytics(ctx context.Context, orgID uint) (*people.PeopleAnalytics, error) {
	var res people.PeopleAnalytics
	
	// Default values
	res.ByStage = make(map[string]int)
	res.ByType = make(map[string]int)
	res.BySource = make(map[string]int)

	var totalContacts int64
	r.db.WithContext(ctx).Model(&models.PeopleContact{}).Where("organisation_id = ?", orgID).Count(&totalContacts)
	res.TotalContacts = int(totalContacts)

	type GroupResult struct {
		Key   string
		Count int
	}

	var stages []GroupResult
	r.db.WithContext(ctx).Model(&models.PeopleContact{}).Select("stage as key, count(*) as count").Where("organisation_id = ?", orgID).Group("stage").Scan(&stages)
	for _, s := range stages {
		res.ByStage[s.Key] = s.Count
	}

	var types []GroupResult
	r.db.WithContext(ctx).Model(&models.PeopleContact{}).Select("type as key, count(*) as count").Where("organisation_id = ?", orgID).Group("type").Scan(&types)
	for _, t := range types {
		res.ByType[t.Key] = t.Count
	}

	var sources []GroupResult
	r.db.WithContext(ctx).Model(&models.PeopleContact{}).Select("source as key, count(*) as count").Where("organisation_id = ? AND source IS NOT NULL", orgID).Group("source").Scan(&sources)
	for _, s := range sources {
		res.BySource[s.Key] = s.Count
	}

	var totals struct {
		TotalDealValue float64
	}
	r.db.WithContext(ctx).Model(&models.PeopleContact{}).Select("SUM(deal_value) as total_deal_value").Where("organisation_id = ? AND deal_value IS NOT NULL", orgID).Scan(&totals)
	res.TotalDealValue = totals.TotalDealValue

	var wonTotals struct {
		TotalDealValue float64
	}
	r.db.WithContext(ctx).Model(&models.PeopleContact{}).Select("SUM(deal_value) as total_deal_value").Joins("JOIN people_pipeline_stages ON people_pipeline_stages.key = people_contacts.stage").Where("people_contacts.organisation_id = ? AND people_contacts.deal_value IS NOT NULL AND people_pipeline_stages.is_won = ?", orgID, true).Scan(&wonTotals)
	res.WonDealValue = wonTotals.TotalDealValue

	if res.TotalDealValue > 0 {
		res.ConversionRate = (res.WonDealValue / res.TotalDealValue) * 100
	}

	if totalContacts > 0 {
		r.db.WithContext(ctx).Model(&models.PeopleContact{}).Select("AVG(deal_value) as avg_deal_value").Where("organisation_id = ? AND deal_value IS NOT NULL", orgID).Scan(&res.AvgDealValue)
	}

	var recent int64
	r.db.WithContext(ctx).Model(&models.PeopleInteraction{}).Joins("JOIN people_contacts on people_contacts.id = people_interactions.contact_id").Where("people_contacts.organisation_id = ? AND people_interactions.occurred_at > ?", orgID, time.Now().AddDate(0, 0, -7)).Count(&recent)
	res.RecentActivity = int(recent)

	return &res, nil
}

func (r *gormRepository) HasContactsInStage(ctx context.Context, stageID uint, orgID uint) (bool, error) {
	var stage models.PeoplePipelineStage
	if err := r.db.WithContext(ctx).Where("id = ? AND organisation_id = ?", stageID, orgID).First(&stage).Error; err != nil {
		return false, err
	}
	var count int64
	r.db.WithContext(ctx).Model(&models.PeopleContact{}).Where("stage = ? AND organisation_id = ?", stage.Key, orgID).Count(&count)
	return count > 0, nil
}

// Mappers
func toDomainContact(m *models.PeopleContact) people.PeopleContact {
	c := people.PeopleContact{
		ID:             m.ID,
		OrganisationID: m.OrganisationID,
		FirstName:      m.FirstName,
		LastName:       m.LastName,
		Email:          m.Email,
		Phone:          m.Phone,
		Company:        m.Company,
		JobTitle:       m.JobTitle,
		AvatarUrl:      m.AvatarUrl,
		Type:           m.Type,
		Status:         m.Status,
		Stage:          m.Stage,
		Tags:           []string(m.Tags),
		ListIDs:        []int(m.ListIDs),
		AssignedToID:   m.AssignedToID,
		Source:         m.Source,
		DealValue:      m.DealValue,
		Currency:       m.Currency,
		Notes:          m.Notes,
		Address:        m.Address,
		City:           m.City,
		Country:        m.Country,
		LinkedinUrl:    m.LinkedinUrl,
		WebsiteUrl:     m.WebsiteUrl,
		CustomFields:   map[string]interface{}(m.CustomFields),
		LastContactedAt: m.LastContactedAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
	if m.AssignedTo != nil && m.AssignedTo.User.ID != 0 {
		name := m.AssignedTo.User.FirstName + " " + m.AssignedTo.User.LastName
		c.AssignedToName = &name
		if m.AssignedTo.User.ProfilePic != "" {
			c.AssignedToAvatar = &m.AssignedTo.User.ProfilePic
		}
	}
	return c
}

func fromDomainContact(c *people.PeopleContact) *models.PeopleContact {
	return &models.PeopleContact{
		Model:          gorm.Model{ID: c.ID, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt},
		OrganisationID: c.OrganisationID,
		FirstName:      c.FirstName,
		LastName:       c.LastName,
		Email:          c.Email,
		Phone:          c.Phone,
		Company:        c.Company,
		JobTitle:       c.JobTitle,
		AvatarUrl:      c.AvatarUrl,
		Type:           c.Type,
		Status:         c.Status,
		Stage:          c.Stage,
		Tags:           models.StringArray(c.Tags),
		ListIDs:        models.IntArray(c.ListIDs),
		AssignedToID:   c.AssignedToID,
		Source:         c.Source,
		DealValue:      c.DealValue,
		Currency:       c.Currency,
		Notes:          c.Notes,
		Address:        c.Address,
		City:           c.City,
		Country:        c.Country,
		LinkedinUrl:    c.LinkedinUrl,
		WebsiteUrl:     c.WebsiteUrl,
		CustomFields:   models.JSONMap(c.CustomFields),
		LastContactedAt: c.LastContactedAt,
	}
}

func toDomainInteraction(m *models.PeopleInteraction) people.PeopleInteraction {
	in := people.PeopleInteraction{
		ID:           m.ID,
		ContactID:    m.ContactID,
		CreatedByID:  m.CreatedByID,
		Type:         m.Type,
		Content:      m.Content,
		Outcome:      m.Outcome,
		LinkedTaskID: m.LinkedTaskID,
		OccurredAt:   m.OccurredAt,
		CreatedAt:    m.CreatedAt,
	}
	// Note: You would usually populate linked task title here if joined
	if m.CreatedBy.User.ID != 0 {
		in.CreatedByName = m.CreatedBy.User.FirstName + " " + m.CreatedBy.User.LastName
		if m.CreatedBy.User.ProfilePic != "" {
			in.CreatedByAvatarUrl = &m.CreatedBy.User.ProfilePic
		}
	}
	return in
}

func fromDomainInteraction(in *people.PeopleInteraction) *models.PeopleInteraction {
	return &models.PeopleInteraction{
		Model:        gorm.Model{ID: in.ID, CreatedAt: in.CreatedAt},
		ContactID:    in.ContactID,
		CreatedByID:  in.CreatedByID,
		Type:         in.Type,
		Content:      in.Content,
		Outcome:      in.Outcome,
		LinkedTaskID: in.LinkedTaskID,
		OccurredAt:   in.OccurredAt,
	}
}

func toDomainStage(m *models.PeoplePipelineStage) people.PeoplePipelineStage {
	return people.PeoplePipelineStage{
		ID:             m.ID,
		OrganisationID: m.OrganisationID,
		Key:            m.Key,
		Label:          m.Label,
		Color:          m.Color,
		Order:          m.Order,
		IsDefault:      m.IsDefault,
		IsWon:          m.IsWon,
		IsLost:         m.IsLost,
	}
}

func fromDomainStage(s *people.PeoplePipelineStage) *models.PeoplePipelineStage {
	return &models.PeoplePipelineStage{
		Model:          gorm.Model{ID: s.ID},
		OrganisationID: s.OrganisationID,
		Key:            s.Key,
		Label:          s.Label,
		Color:          s.Color,
		Order:          s.Order,
		IsDefault:      s.IsDefault,
		IsWon:          s.IsWon,
		IsLost:         s.IsLost,
	}
}

func toDomainList(m *models.PeopleList) people.PeopleList {
	return people.PeopleList{
		ID:             m.ID,
		OrganisationID: m.OrganisationID,
		Name:           m.Name,
		Description:    m.Description,
		Color:          m.Color,
		CreatedAt:      m.CreatedAt,
	}
}

func fromDomainList(l *people.PeopleList) *models.PeopleList {
	return &models.PeopleList{
		Model:          gorm.Model{ID: l.ID, CreatedAt: l.CreatedAt},
		OrganisationID: l.OrganisationID,
		Name:           l.Name,
		Description:    l.Description,
		Color:          l.Color,
	}
}
