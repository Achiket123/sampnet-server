package people

import (
	"context"
)

type UseCase interface {
	// Contacts
	GetContacts(ctx context.Context, filter ContactsFilter, orgID uint) (*ContactsResponse, error)
	GetContactByID(ctx context.Context, id uint, orgID uint) (*ContactDetailResponse, error)
	CreateContact(ctx context.Context, params *CreateContactParams, orgID uint) (*PeopleContact, error)
	UpdateContact(ctx context.Context, id uint, params *UpdateContactParams, orgID uint, currentUserID uint) (*PeopleContact, error)
	DeleteContact(ctx context.Context, id uint, orgID uint) error
	BulkUpdateStage(ctx context.Context, ids []uint, stage string, orgID uint, changedByID uint) error
	BulkDelete(ctx context.Context, ids []uint, orgID uint) error
	AssignContact(ctx context.Context, contactID uint, employeeID uint, orgID uint) error

	// Interactions
	GetInteractions(ctx context.Context, contactID uint, orgID uint) ([]PeopleInteraction, error)
	AddInteraction(ctx context.Context, contactID uint, params *AddInteractionParams, orgID uint, createdByID uint) (*PeopleInteraction, error)
	DeleteInteraction(ctx context.Context, interactionID uint, reqUserID uint, reqUserIsAdmin bool) error

	// Pipeline Stages
	GetPipelineStages(ctx context.Context, orgID uint) ([]PeoplePipelineStage, error)
	CreatePipelineStage(ctx context.Context, params *CreateStageParams, orgID uint) (*PeoplePipelineStage, error)
	UpdatePipelineStage(ctx context.Context, id uint, params *UpdateStageParams, orgID uint) (*PeoplePipelineStage, error)
	DeletePipelineStage(ctx context.Context, id uint, orgID uint, reassignToStage *string) error
	ReorderPipelineStages(ctx context.Context, orderedIDs []uint, orgID uint) error

	// Lists
	GetLists(ctx context.Context, orgID uint) ([]PeopleList, error)
	CreateList(ctx context.Context, params *CreateListParams, orgID uint) (*PeopleList, error)
	UpdateList(ctx context.Context, id uint, params *UpdateListParams, orgID uint) (*PeopleList, error)
	DeleteList(ctx context.Context, id uint, orgID uint) error
	RemoveContactFromList(ctx context.Context, listID uint, contactID uint, orgID uint) error

	// Analytics
	GetAnalytics(ctx context.Context, orgID uint) (*PeopleAnalytics, error)
}
