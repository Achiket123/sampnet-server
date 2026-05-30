package auth

import "context"

// UseCase defines the application-facing ports for authentication operations.
type UseCase interface {
	SignUp(ctx context.Context, user *User, password string) (string, error)
	SignIn(ctx context.Context, email, password string) (string, error)
	CompleteSignIn(ctx context.Context, email, phone, password, city, country, profilePic string) (string, error)
	ValidateEmployee(ctx context.Context, userID uint) (string, error)
}
