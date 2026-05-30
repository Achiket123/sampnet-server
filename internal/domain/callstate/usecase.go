package callstate

import "context"

type UseCase interface {
	CreateOrUpdate(ctx context.Context, state *State) error
	Get(ctx context.Context, id uint) (*State, error)
	CreateOffer(ctx context.Context, id uint, callingID, firstName, lastName, offer string) error
	EndCall(ctx context.Context, id uint) error
}
