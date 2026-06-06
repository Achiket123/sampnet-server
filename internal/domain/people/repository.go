package people

import (
	"context"
)

type Repository interface {
	// Contacts
	GetContacts(ctx context.Context, filter ContactsFilter, orgID uint) (*ContactsResponse, error)
	GetContactByID(ctx context.Context, id uint, orgID uint) (*PeopleContact, error)
	CreateContact(ctx context.Context, contact *PeopleContact) error
	UpdateContact(ctx context.Context, contact *PeopleContact) error
	DeleteContact(ctx context.Context, id uint, orgID uint) error
	BulkUpdateStage(ctx context.Context, ids []uint, stage string, orgID uint, changedByID uint) error
	BulkDelete(ctx context.Context, ids []uint, orgID uint) error
	AssignContact(ctx context.Context, contactID uint, employeeID uint, orgID uint) error

	// Interactions
	GetInteractions(ctx context.Context, contactID uint, orgID uint) ([]PeopleInteraction, error)
	AddInteraction(ctx context.Context, interaction *PeopleInteraction) error
	DeleteInteraction(ctx context.Context, interactionID uint, reqUserID uint, reqUserIsAdmin bool) error

	// Pipeline Stages
	GetPipelineStages(ctx context.Context, orgID uint) ([]PeoplePipelineStage, error)
	CreatePipelineStage(ctx context.Context, stage *PeoplePipelineStage) error
	UpdatePipelineStage(ctx context.Context, stage *PeoplePipelineStage) error
	DeletePipelineStage(ctx context.Context, id uint, orgID uint, reassignToStage *string) error
	ReorderPipelineStages(ctx context.Context, orderedIDs []uint, orgID uint) error

	// Lists
	GetLists(ctx context.Context, orgID uint) ([]PeopleList, error)
	CreateList(ctx context.Context, list *PeopleList) error
	UpdateList(ctx context.Context, list *PeopleList) error
	DeleteList(ctx context.Context, id uint, orgID uint) error
	RemoveContactFromList(ctx context.Context, listID uint, contactID uint, orgID uint) error

	// Analytics
	GetAnalytics(ctx context.Context, orgID uint) (*PeopleAnalytics, error)

	// Utils
	HasContactsInStage(ctx context.Context, stageID uint, orgID uint) (bool, error)
}
