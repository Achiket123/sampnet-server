package organisation

import "context"

// UseCase defines application-facing ports for organisation operations.
type UseCase interface {
	Register(ctx context.Context, org *Entity, ownerUserID uint) (*OwnerEmployeeRow, error)
	Get(ctx context.Context, id uint) (*Entity, error)
	Update(ctx context.Context, org *Entity) error
}
