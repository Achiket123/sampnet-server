package people

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	domain "server/internal/domain/people"

	"github.com/redis/go-redis/v9"
)

type cachedRepository struct {
	client *redis.Client
	inner  domain.Repository
}

// NewCachedRepository wraps a people Repository with a Redis caching layer.
func NewCachedRepository(client *redis.Client, inner domain.Repository) domain.Repository {
	return &cachedRepository{
		client: client,
		inner:  inner,
	}
}

// Contacts

func (r *cachedRepository) GetContacts(ctx context.Context, filter domain.ContactsFilter, orgID uint) (*domain.ContactsResponse, error) {
	version := "0"
	vKey := "people:contacts_version:" + strconv.FormatUint(uint64(orgID), 10)
	vVal, err := r.client.Get(ctx, vKey).Result()
	if err == nil {
		version = vVal
	} else if err != redis.Nil {
		log.Printf("Redis error getting contacts version: %v", err)
	}

	stageStr := "none"
	if filter.Stage != nil {
		stageStr = *filter.Stage
	}

	typeStr := "none"
	if filter.Type != nil {
		typeStr = *filter.Type
	}

	statusStr := "none"
	if filter.Status != nil {
		statusStr = *filter.Status
	}

	listStr := "none"
	if filter.ListID != nil {
		listStr = strconv.Itoa(*filter.ListID)
	}

	searchStr := "none"
	trimmedSearch := strings.TrimSpace(strings.ToLower(filter.Search))
	if trimmedSearch != "" {
		searchStr = trimmedSearch
	}

	sortByStr := "none"
	if filter.SortBy != "" {
		sortByStr = filter.SortBy
	}

	sortOrderStr := "none"
	if filter.SortOrder != "" {
		sortOrderStr = filter.SortOrder
	}

	pageStr := strconv.Itoa(filter.Page)
	limitStr := strconv.Itoa(filter.Limit)

	sig := strings.Join([]string{
		stageStr,
		typeStr,
		statusStr,
		listStr,
		searchStr,
		sortByStr,
		sortOrderStr,
		pageStr,
		limitStr,
	}, ":")

	cacheKey := fmt.Sprintf("people:contacts:%d:v%s:%s", orgID, version, sig)

	cachedVal, err := r.client.Get(ctx, cacheKey).Result()
	if err == nil {
		var resp domain.ContactsResponse
		if err := json.Unmarshal([]byte(cachedVal), &resp); err == nil {
			return &resp, nil
		}
		log.Printf("Error unmarshaling cached contacts: %v", err)
	} else if err != redis.Nil {
		log.Printf("Redis error getting cached contacts: %v", err)
	}

	resp, err := r.inner.GetContacts(ctx, filter, orgID)
	if err != nil {
		return nil, err
	}

	respBytes, err := json.Marshal(resp)
	if err == nil {
		if err := r.client.Set(ctx, cacheKey, respBytes, 30*time.Second).Err(); err != nil {
			log.Printf("Redis error setting cached contacts: %v", err)
		}
	} else {
		log.Printf("Error marshaling contacts response for cache: %v", err)
	}

	return resp, nil
}

func (r *cachedRepository) GetContactByID(ctx context.Context, id uint, orgID uint) (*domain.PeopleContact, error) {
	key := "people:contact:" + strconv.FormatUint(uint64(id), 10)
	cachedVal, err := r.client.Get(ctx, key).Result()
	if err == nil {
		var contact domain.PeopleContact
		if err := json.Unmarshal([]byte(cachedVal), &contact); err == nil {
			if contact.OrganisationID == orgID {
				return &contact, nil
			}
		} else {
			log.Printf("Error unmarshaling cached contact: %v", err)
		}
	} else if err != redis.Nil {
		log.Printf("Redis error getting cached contact: %v", err)
	}

	contact, err := r.inner.GetContactByID(ctx, id, orgID)
	if err != nil {
		return nil, err
	}

	contactBytes, err := json.Marshal(contact)
	if err == nil {
		if err := r.client.Set(ctx, key, contactBytes, 60*time.Second).Err(); err != nil {
			log.Printf("Redis error setting cached contact: %v", err)
		}
	} else {
		log.Printf("Error marshaling contact for cache: %v", err)
	}

	return contact, nil
}

func (r *cachedRepository) CreateContact(ctx context.Context, contact *domain.PeopleContact) error {
	err := r.inner.CreateContact(ctx, contact)
	if err != nil {
		return err
	}
	vKey := "people:contacts_version:" + strconv.FormatUint(uint64(contact.OrganisationID), 10)
	if err := r.client.Incr(ctx, vKey).Err(); err != nil {
		log.Printf("Redis error incrementing contacts version: %v", err)
	}
	return nil
}

func (r *cachedRepository) UpdateContact(ctx context.Context, contact *domain.PeopleContact) error {
	err := r.inner.UpdateContact(ctx, contact)
	if err != nil {
		return err
	}
	key := "people:contact:" + strconv.FormatUint(uint64(contact.ID), 10)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		log.Printf("Redis error deleting contact key: %v", err)
	}
	vKey := "people:contacts_version:" + strconv.FormatUint(uint64(contact.OrganisationID), 10)
	if err := r.client.Incr(ctx, vKey).Err(); err != nil {
		log.Printf("Redis error incrementing contacts version: %v", err)
	}
	return nil
}

func (r *cachedRepository) DeleteContact(ctx context.Context, id uint, orgID uint) error {
	err := r.inner.DeleteContact(ctx, id, orgID)
	if err != nil {
		return err
	}
	key := "people:contact:" + strconv.FormatUint(uint64(id), 10)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		log.Printf("Redis error deleting contact key: %v", err)
	}
	vKey := "people:contacts_version:" + strconv.FormatUint(uint64(orgID), 10)
	if err := r.client.Incr(ctx, vKey).Err(); err != nil {
		log.Printf("Redis error incrementing contacts version: %v", err)
	}
	return nil
}

func (r *cachedRepository) BulkUpdateStage(ctx context.Context, ids []uint, stage string, orgID uint, changedByID uint) error {
	err := r.inner.BulkUpdateStage(ctx, ids, stage, orgID, changedByID)
	if err != nil {
		return err
	}
	for _, id := range ids {
		key := "people:contact:" + strconv.FormatUint(uint64(id), 10)
		if err := r.client.Del(ctx, key).Err(); err != nil {
			log.Printf("Redis error deleting contact key: %v", err)
		}
	}
	vKey := "people:contacts_version:" + strconv.FormatUint(uint64(orgID), 10)
	if err := r.client.Incr(ctx, vKey).Err(); err != nil {
		log.Printf("Redis error incrementing contacts version: %v", err)
	}
	return nil
}

func (r *cachedRepository) BulkDelete(ctx context.Context, ids []uint, orgID uint) error {
	err := r.inner.BulkDelete(ctx, ids, orgID)
	if err != nil {
		return err
	}
	for _, id := range ids {
		key := "people:contact:" + strconv.FormatUint(uint64(id), 10)
		if err := r.client.Del(ctx, key).Err(); err != nil {
			log.Printf("Redis error deleting contact key: %v", err)
		}
	}
	vKey := "people:contacts_version:" + strconv.FormatUint(uint64(orgID), 10)
	if err := r.client.Incr(ctx, vKey).Err(); err != nil {
		log.Printf("Redis error incrementing contacts version: %v", err)
	}
	return nil
}

func (r *cachedRepository) AssignContact(ctx context.Context, contactID uint, employeeID uint, orgID uint) error {
	err := r.inner.AssignContact(ctx, contactID, employeeID, orgID)
	if err != nil {
		return err
	}
	key := "people:contact:" + strconv.FormatUint(uint64(contactID), 10)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		log.Printf("Redis error deleting contact key: %v", err)
	}
	vKey := "people:contacts_version:" + strconv.FormatUint(uint64(orgID), 10)
	if err := r.client.Incr(ctx, vKey).Err(); err != nil {
		log.Printf("Redis error incrementing contacts version: %v", err)
	}
	return nil
}

// Interactions

func (r *cachedRepository) GetInteractions(ctx context.Context, contactID uint, orgID uint) ([]domain.PeopleInteraction, error) {
	return r.inner.GetInteractions(ctx, contactID, orgID)
}

func (r *cachedRepository) AddInteraction(ctx context.Context, interaction *domain.PeopleInteraction) error {
	err := r.inner.AddInteraction(ctx, interaction)
	if err != nil {
		return err
	}
	// This method only receives a PeopleInteraction value which does not carry an OrganisationID field, so the version counter for the owning organisation cannot be located without an extra database lookup, and the contacts list view is allowed to show a slightly stale last contacted timestamp for up to the thirty second list cache ttl, which is an acceptable and deliberate tradeoff, not a bug.
	key := "people:contact:" + strconv.FormatUint(uint64(interaction.ContactID), 10)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		log.Printf("Redis error deleting contact key: %v", err)
	}
	return nil
}

func (r *cachedRepository) DeleteInteraction(ctx context.Context, interactionID uint, reqUserID uint, reqUserIsAdmin bool) error {
	return r.inner.DeleteInteraction(ctx, interactionID, reqUserID, reqUserIsAdmin)
}

// Pipeline Stages

func (r *cachedRepository) GetPipelineStages(ctx context.Context, orgID uint) ([]domain.PeoplePipelineStage, error) {
	key := "people:pipeline_stages:" + strconv.FormatUint(uint64(orgID), 10)
	cachedVal, err := r.client.Get(ctx, key).Result()
	if err == nil {
		var stages []domain.PeoplePipelineStage
		if err := json.Unmarshal([]byte(cachedVal), &stages); err == nil {
			return stages, nil
		}
		log.Printf("Error unmarshaling cached pipeline stages: %v", err)
	} else if err != redis.Nil {
		log.Printf("Redis error getting cached pipeline stages: %v", err)
	}

	stages, err := r.inner.GetPipelineStages(ctx, orgID)
	if err != nil {
		return nil, err
	}

	stagesBytes, err := json.Marshal(stages)
	if err == nil {
		if err := r.client.Set(ctx, key, stagesBytes, 300*time.Second).Err(); err != nil {
			log.Printf("Redis error setting cached pipeline stages: %v", err)
		}
	} else {
		log.Printf("Error marshaling pipeline stages for cache: %v", err)
	}

	return stages, nil
}

func (r *cachedRepository) CreatePipelineStage(ctx context.Context, stage *domain.PeoplePipelineStage) error {
	err := r.inner.CreatePipelineStage(ctx, stage)
	if err != nil {
		return err
	}
	key := "people:pipeline_stages:" + strconv.FormatUint(uint64(stage.OrganisationID), 10)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		log.Printf("Redis error deleting pipeline stages: %v", err)
	}
	return nil
}

func (r *cachedRepository) UpdatePipelineStage(ctx context.Context, stage *domain.PeoplePipelineStage) error {
	err := r.inner.UpdatePipelineStage(ctx, stage)
	if err != nil {
		return err
	}
	key := "people:pipeline_stages:" + strconv.FormatUint(uint64(stage.OrganisationID), 10)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		log.Printf("Redis error deleting pipeline stages: %v", err)
	}
	return nil
}

func (r *cachedRepository) DeletePipelineStage(ctx context.Context, id uint, orgID uint, reassignToStage *string) error {
	err := r.inner.DeletePipelineStage(ctx, id, orgID, reassignToStage)
	if err != nil {
		return err
	}
	key := "people:pipeline_stages:" + strconv.FormatUint(uint64(orgID), 10)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		log.Printf("Redis error deleting pipeline stages: %v", err)
	}
	return nil
}

func (r *cachedRepository) ReorderPipelineStages(ctx context.Context, orderedIDs []uint, orgID uint) error {
	err := r.inner.ReorderPipelineStages(ctx, orderedIDs, orgID)
	if err != nil {
		return err
	}
	key := "people:pipeline_stages:" + strconv.FormatUint(uint64(orgID), 10)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		log.Printf("Redis error deleting pipeline stages: %v", err)
	}
	return nil
}

// Lists

func (r *cachedRepository) GetLists(ctx context.Context, orgID uint) ([]domain.PeopleList, error) {
	return r.inner.GetLists(ctx, orgID)
}

func (r *cachedRepository) CreateList(ctx context.Context, list *domain.PeopleList) error {
	return r.inner.CreateList(ctx, list)
}

func (r *cachedRepository) UpdateList(ctx context.Context, list *domain.PeopleList) error {
	return r.inner.UpdateList(ctx, list)
}

func (r *cachedRepository) DeleteList(ctx context.Context, id uint, orgID uint) error {
	return r.inner.DeleteList(ctx, id, orgID)
}

func (r *cachedRepository) RemoveContactFromList(ctx context.Context, listID uint, contactID uint, orgID uint) error {
	return r.inner.RemoveContactFromList(ctx, listID, contactID, orgID)
}

// Analytics

func (r *cachedRepository) GetAnalytics(ctx context.Context, orgID uint) (*domain.PeopleAnalytics, error) {
	return r.inner.GetAnalytics(ctx, orgID)
}

// Utils

func (r *cachedRepository) HasContactsInStage(ctx context.Context, stageID uint, orgID uint) (bool, error) {
	return r.inner.HasContactsInStage(ctx, stageID, orgID)
}
