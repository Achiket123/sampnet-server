package callstate

import "context"

type Repository interface {
	CreateOrUpdate(ctx context.Context, state *State) error
	GetByID(ctx context.Context, id uint) (*State, error)
	UpdateOffer(ctx context.Context, id uint, callingID, firstName, lastName, offer string) error
	EndCall(ctx context.Context, id uint) error
}
