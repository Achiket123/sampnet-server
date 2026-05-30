package callstate

import (
	"context"
	domain "server/internal/domain/callstate"
	"server/internal/platform/database/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }

func (r *GormRepository) CreateOrUpdate(ctx context.Context, state *domain.State) error {
	m := models.CallState{ID: state.ID, FirstName: state.FirstName, LastName: state.LastName, Email: state.Email, OrganisationID: state.OrganisationID, InCall: state.InCall, LastCall: state.LastCall, CallingID: state.CallingID, CallingFirstName: state.CallingFirstName, CallingLastName: state.CallingLastName, Offer: state.Offer}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&m).Error
}

func (r *GormRepository) GetByID(ctx context.Context, id uint) (*domain.State, error) {
	var row models.CallState
	if err := r.db.WithContext(ctx).First(&row, id).Error; err != nil {
		return nil, err
	}
	return &domain.State{ID: row.ID, FirstName: row.FirstName, LastName: row.LastName, Email: row.Email, OrganisationID: row.OrganisationID, InCall: row.InCall, LastCall: row.LastCall, CallingID: row.CallingID, CallingFirstName: row.CallingFirstName, CallingLastName: row.CallingLastName, Offer: row.Offer}, nil
}

func (r *GormRepository) UpdateOffer(ctx context.Context, id uint, callingID, firstName, lastName, offer string) error {
	now := time.Now().UTC()
	updates := map[string]any{"calling_id": callingID, "calling_first_name": firstName, "calling_last_name": lastName, "offer": offer, "in_call": true, "last_call": &now}
	return r.db.WithContext(ctx).Model(&models.CallState{}).Where("id = ?", id).Updates(updates).Error
}

func (r *GormRepository) EndCall(ctx context.Context, id uint) error {
	now := time.Now().UTC()
	updates := map[string]any{"calling_id": nil, "calling_first_name": nil, "calling_last_name": nil, "offer": nil, "in_call": false, "last_call": &now}
	return r.db.WithContext(ctx).Model(&models.CallState{}).Where("id = ?", id).Updates(updates).Error
}
