package auth

import "context"

// Repository defines the persistence-facing ports for user-related operations.
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhoneNumber(ctx context.Context, phone string) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	Update(ctx context.Context, user *User) error
}
