package organisation

import "context"

// Repository defines persistence-facing ports for organisation operations.
type Repository interface {
	Create(ctx context.Context, org *Entity) error
	CreateWithOwner(ctx context.Context, org *Entity, ownerUserID uint) (*OwnerEmployeeRow, error)
	GetByID(ctx context.Context, id uint) (*Entity, error)
	Update(ctx context.Context, org *Entity) error
}
