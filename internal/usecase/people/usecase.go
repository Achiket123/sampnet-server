package people

import (
	"context"
	"fmt"
	"server/internal/domain/people"
)

type usecase struct {
	repo people.Repository
}

func NewUseCase(repo people.Repository) people.UseCase {
	return &usecase{repo: repo}
}

func (u *usecase) GetContacts(ctx context.Context, filter people.ContactsFilter, orgID uint) (*people.ContactsResponse, error) {
	return u.repo.GetContacts(ctx, filter, orgID)
}

func (u *usecase) GetContactByID(ctx context.Context, id uint, orgID uint) (*people.ContactDetailResponse, error) {
	contact, err := u.repo.GetContactByID(ctx, id, orgID)
	if err != nil {
		return nil, err
	}
	interactions, err := u.repo.GetInteractions(ctx, id, orgID)
	if err != nil {
		return nil, err
	}
	return &people.ContactDetailResponse{
		Contact:      *contact,
		Interactions: interactions,
	}, nil
}

func (u *usecase) CreateContact(ctx context.Context, params *people.CreateContactParams, orgID uint) (*people.PeopleContact, error) {
	if params.Stage == "" {
		params.Stage = "new" // Default
	}
	if params.Currency == "" {
		params.Currency = "USD"
	}
	if params.Status == "" {
		params.Status = "active"
	}
	
	contact := &people.PeopleContact{
		OrganisationID: orgID,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		Email:          params.Email,
		Phone:          params.Phone,
		Company:        params.Company,
		JobTitle:       params.JobTitle,
		AvatarUrl:      params.AvatarUrl,
		Type:           params.Type,
		Status:         params.Status,
		Stage:          params.Stage,
		Tags:           params.Tags,
		ListIDs:        params.ListIDs,
		AssignedToID:   params.AssignedToID,
		Source:         params.Source,
		DealValue:      params.DealValue,
		Currency:       params.Currency,
		Notes:          params.Notes,
		Address:        params.Address,
		City:           params.City,
		Country:        params.Country,
		LinkedinUrl:    params.LinkedinUrl,
		WebsiteUrl:     params.WebsiteUrl,
		CustomFields:   params.CustomFields,
		LastContactedAt: params.LastContactedAt,
	}

	err := u.repo.CreateContact(ctx, contact)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (u *usecase) UpdateContact(ctx context.Context, id uint, params *people.UpdateContactParams, orgID uint, currentUserID uint) (*people.PeopleContact, error) {
	contact, err := u.repo.GetContactByID(ctx, id, orgID)
	if err != nil {
		return nil, err
	}

	stageChanged := false
	var oldStage string
	var newStage string

	if params.FirstName != nil {
		contact.FirstName = *params.FirstName
	}
	if params.LastName != nil {
		contact.LastName = *params.LastName
	}
	if params.Email != nil {
		contact.Email = params.Email
	}
	if params.Phone != nil {
		contact.Phone = params.Phone
	}
	if params.Company != nil {
		contact.Company = params.Company
	}
	if params.JobTitle != nil {
		contact.JobTitle = params.JobTitle
	}
	if params.AvatarUrl != nil {
		contact.AvatarUrl = params.AvatarUrl
	}
	if params.Type != nil {
		contact.Type = *params.Type
	}
	if params.Status != nil {
		contact.Status = *params.Status
	}
	if params.Stage != nil && *params.Stage != contact.Stage {
		stageChanged = true
		oldStage = contact.Stage
		newStage = *params.Stage
		contact.Stage = *params.Stage
	}
	if params.Tags != nil {
		contact.Tags = params.Tags
	}
	if params.ListIDs != nil {
		contact.ListIDs = params.ListIDs
	}
	if params.AssignedToID != nil {
		contact.AssignedToID = params.AssignedToID
	}
	if params.Source != nil {
		contact.Source = params.Source
	}
	if params.DealValue != nil {
		contact.DealValue = params.DealValue
	}
	if params.Currency != nil {
		contact.Currency = *params.Currency
	}
	if params.Notes != nil {
		contact.Notes = params.Notes
	}
	if params.Address != nil {
		contact.Address = params.Address
	}
	if params.City != nil {
		contact.City = params.City
	}
	if params.Country != nil {
		contact.Country = params.Country
	}
	if params.LinkedinUrl != nil {
		contact.LinkedinUrl = params.LinkedinUrl
	}
	if params.WebsiteUrl != nil {
		contact.WebsiteUrl = params.WebsiteUrl
	}
	if params.CustomFields != nil {
		contact.CustomFields = params.CustomFields
	}

	err = u.repo.UpdateContact(ctx, contact)
	if err != nil {
		return nil, err
	}

	if stageChanged {
		_ = u.repo.AddInteraction(ctx, &people.PeopleInteraction{
			ContactID:   contact.ID,
			CreatedByID: currentUserID,
			Type:        "stage_change",
			Content:     fmt.Sprintf("Stage changed from %s to %s", oldStage, newStage),
		})
	}

	return contact, nil
}

func (u *usecase) DeleteContact(ctx context.Context, id uint, orgID uint) error {
	return u.repo.DeleteContact(ctx, id, orgID)
}

func (u *usecase) BulkUpdateStage(ctx context.Context, ids []uint, stage string, orgID uint, changedByID uint) error {
	return u.repo.BulkUpdateStage(ctx, ids, stage, orgID, changedByID)
}

func (u *usecase) BulkDelete(ctx context.Context, ids []uint, orgID uint) error {
	return u.repo.BulkDelete(ctx, ids, orgID)
}

func (u *usecase) AssignContact(ctx context.Context, contactID uint, employeeID uint, orgID uint) error {
	return u.repo.AssignContact(ctx, contactID, employeeID, orgID)
}

func (u *usecase) GetInteractions(ctx context.Context, contactID uint, orgID uint) ([]people.PeopleInteraction, error) {
	return u.repo.GetInteractions(ctx, contactID, orgID)
}

func (u *usecase) AddInteraction(ctx context.Context, contactID uint, params *people.AddInteractionParams, orgID uint, createdByID uint) (*people.PeopleInteraction, error) {
	// Verify contact belongs to org
	_, err := u.repo.GetContactByID(ctx, contactID, orgID)
	if err != nil {
		return nil, err
	}

	in := &people.PeopleInteraction{
		ContactID:    contactID,
		CreatedByID:  createdByID,
		Type:         params.Type,
		Content:      params.Content,
		Outcome:      params.Outcome,
		LinkedTaskID: params.LinkedTaskID,
		OccurredAt:   params.OccurredAt,
	}

	err = u.repo.AddInteraction(ctx, in)
	if err != nil {
		return nil, err
	}
	return in, nil
}

func (u *usecase) DeleteInteraction(ctx context.Context, interactionID uint, reqUserID uint, reqUserIsAdmin bool) error {
	return u.repo.DeleteInteraction(ctx, interactionID, reqUserID, reqUserIsAdmin)
}

func (u *usecase) GetPipelineStages(ctx context.Context, orgID uint) ([]people.PeoplePipelineStage, error) {
	return u.repo.GetPipelineStages(ctx, orgID)
}

func (u *usecase) CreatePipelineStage(ctx context.Context, params *people.CreateStageParams, orgID uint) (*people.PeoplePipelineStage, error) {
	stage := &people.PeoplePipelineStage{
		OrganisationID: orgID,
		Key:            params.Key,
		Label:          params.Label,
		Color:          params.Color,
		Order:          params.Order,
		IsWon:          params.IsWon,
		IsLost:         params.IsLost,
	}
	err := u.repo.CreatePipelineStage(ctx, stage)
	if err != nil {
		return nil, err
	}
	return stage, nil
}

func (u *usecase) UpdatePipelineStage(ctx context.Context, id uint, params *people.UpdateStageParams, orgID uint) (*people.PeoplePipelineStage, error) {
	stage := &people.PeoplePipelineStage{
		ID:             id,
		OrganisationID: orgID,
	}
	
	if params.Label != nil {
		stage.Label = *params.Label
	}
	if params.Color != nil {
		stage.Color = *params.Color
	}
	if params.Order != nil {
		stage.Order = *params.Order
	}
	if params.IsWon != nil {
		stage.IsWon = *params.IsWon
	}
	if params.IsLost != nil {
		stage.IsLost = *params.IsLost
	}

	err := u.repo.UpdatePipelineStage(ctx, stage)
	if err != nil {
		return nil, err
	}
	return stage, nil
}

func (u *usecase) DeletePipelineStage(ctx context.Context, id uint, orgID uint, reassignToStage *string) error {
	return u.repo.DeletePipelineStage(ctx, id, orgID, reassignToStage)
}

func (u *usecase) ReorderPipelineStages(ctx context.Context, orderedIDs []uint, orgID uint) error {
	return u.repo.ReorderPipelineStages(ctx, orderedIDs, orgID)
}

func (u *usecase) GetLists(ctx context.Context, orgID uint) ([]people.PeopleList, error) {
	return u.repo.GetLists(ctx, orgID)
}

func (u *usecase) CreateList(ctx context.Context, params *people.CreateListParams, orgID uint) (*people.PeopleList, error) {
	list := &people.PeopleList{
		OrganisationID: orgID,
		Name:           params.Name,
		Description:    params.Description,
		Color:          params.Color,
	}
	err := u.repo.CreateList(ctx, list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (u *usecase) UpdateList(ctx context.Context, id uint, params *people.UpdateListParams, orgID uint) (*people.PeopleList, error) {
	list := &people.PeopleList{
		ID:             id,
		OrganisationID: orgID,
	}
	if params.Name != nil {
		list.Name = *params.Name
	}
	if params.Description != nil {
		list.Description = *params.Description
	}
	if params.Color != nil {
		list.Color = *params.Color
	}
	err := u.repo.UpdateList(ctx, list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (u *usecase) DeleteList(ctx context.Context, id uint, orgID uint) error {
	return u.repo.DeleteList(ctx, id, orgID)
}

func (u *usecase) RemoveContactFromList(ctx context.Context, listID uint, contactID uint, orgID uint) error {
	return u.repo.RemoveContactFromList(ctx, listID, contactID, orgID)
}

func (u *usecase) GetAnalytics(ctx context.Context, orgID uint) (*people.PeopleAnalytics, error) {
	return u.repo.GetAnalytics(ctx, orgID)
}