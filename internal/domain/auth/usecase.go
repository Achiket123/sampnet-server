package auth

import "context"

// UseCase defines the application-facing ports for authentication operations.
type UseCase interface {
	SignUp(ctx context.Context, user *User, password string) (TokenPair, error)
	SignIn(ctx context.Context, email, password string) (TokenPair, error)
	CompleteSignIn(ctx context.Context, email, phone, password, city, country, profilePic string) (TokenPair, error)
	ValidateEmployee(ctx context.Context, userID uint) (string, error)
	RefreshToken(ctx context.Context, refreshToken string) (TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
