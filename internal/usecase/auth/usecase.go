package auth

import (
	"context"
	"errors"
	"log"
	domain "server/internal/domain/auth"
	empDomain "server/internal/domain/employees"
	"server/internal/platform/database/models"
	"server/internal/platform/miscallenous"
	"time"

	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrNotAnEmployee      = errors.New("user is not an employee")
	ErrInvalidToken       = errors.New("invalid or expired refresh token")
)

type service struct {
	repo    domain.Repository
	empRepo empDomain.Repository
}

func NewService(repo domain.Repository, empRepo empDomain.Repository) domain.UseCase {
	return &service{repo: repo, empRepo: empRepo}
}

func (s *service) SignUp(ctx context.Context, user *domain.User, password string) (domain.TokenPair, error) {
	hashedPassword, err := miscallenous.HashPassword(password)
	if err != nil {
		return domain.TokenPair{}, err
	}
	user.HashedPassword = hashedPassword
	user.LastLoginAt = time.Now()

	existing, err := s.repo.GetByEmail(ctx, user.Email)
	if existing != nil {
		if err != nil && err != gorm.ErrRecordNotFound {
			return domain.TokenPair{}, err
		}
		if existing.HashedPassword != "" {
			return domain.TokenPair{}, gorm.ErrRegistered
		}
		user.ID = existing.ID
		if err := s.repo.Update(ctx, user); err != nil {
			return domain.TokenPair{}, err
		}
	} else {
		if err := s.repo.Create(ctx, user); err != nil {
			return domain.TokenPair{}, err
		}
	}

	return s.issueTokenPair(ctx, user)
}

func (s *service) SignIn(ctx context.Context, email, password string) (domain.TokenPair, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return domain.TokenPair{}, ErrInvalidCredentials
	}

	if !miscallenous.VerifyPassword(user.HashedPassword, password) {
		return domain.TokenPair{}, ErrInvalidCredentials
	}

	return s.issueTokenPair(ctx, user)
}

func (s *service) CompleteSignIn(ctx context.Context, email, phone, password, city, country, profilePic string) (domain.TokenPair, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		user, err = s.repo.GetByPhoneNumber(ctx, phone)
		if err != nil {
			return domain.TokenPair{}, ErrUserNotFound
		}
	}

	hashedPassword, err := miscallenous.HashPassword(password)
	if err != nil {
		return domain.TokenPair{}, err
	}

	user.HashedPassword = hashedPassword
	user.City = city
	user.Country = country
	user.ProfilePic = profilePic
	user.IsVerified = true
	user.LastLoginAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return domain.TokenPair{}, err
	}

	return s.issueTokenPair(ctx, user)
}

func (s *service) ValidateEmployee(ctx context.Context, userID uint) (string, error) {
	employee, err := s.empRepo.GetByUserID(ctx, userID)
	if err != nil {
		return "", ErrNotAnEmployee
	}

	empModel := models.Employee{
		UserID:         employee.UserID,
		OrganisationID: employee.OrganisationID,
		Type:           employee.Type,
	}

	token, err := miscallenous.GenerateJWTToken(empModel, "employee", employee.UserID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) RefreshToken(ctx context.Context, rawRefreshToken string) (domain.TokenPair, error) {
	tokenHash := miscallenous.HashRefreshToken(rawRefreshToken)

	stored, err := s.repo.GetRefreshToken(ctx, tokenHash)
	if err != nil || stored == nil {
		return domain.TokenPair{}, ErrInvalidToken
	}

	if stored.Revoked {
		// Token reuse detected — revoke all tokens for this user (rotation breach)
		_ = s.repo.RevokeAllUserRefreshTokens(ctx, stored.UserID)
		return domain.TokenPair{}, ErrInvalidToken
	}

	if stored.ExpiresAt < time.Now().Unix() {
		return domain.TokenPair{}, ErrInvalidToken
	}

	// Revoke old token before issuing new pair
	if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return domain.TokenPair{}, err
	}

	user, err := s.repo.GetByID(ctx, stored.UserID)
	if err != nil {
		return domain.TokenPair{}, ErrUserNotFound
	}

	return s.issueTokenPair(ctx, user)
}

func (s *service) Logout(ctx context.Context, rawRefreshToken string) error {
	tokenHash := miscallenous.HashRefreshToken(rawRefreshToken)
	return s.repo.RevokeRefreshToken(ctx, tokenHash)
}

// issueTokenPair generates a new access + refresh token pair and persists the refresh token.
func (s *service) issueTokenPair(ctx context.Context, user *domain.User) (domain.TokenPair, error) {
	userModel := models.UserModel{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	userModel.ID = user.ID

	log.Default().Printf("Issuing token pair for user %d", user.ID)

	accessToken, err := miscallenous.GenerateJWTToken(userModel, "user", user.ID)
	if err != nil {
		return domain.TokenPair{}, err
	}

	rawRefresh, err := miscallenous.GenerateRefreshToken()
	if err != nil {
		return domain.TokenPair{}, err
	}

	refreshRecord := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: miscallenous.HashRefreshToken(rawRefresh),
		ExpiresAt: miscallenous.RefreshTokenExpiresAt(),
		Revoked:   false,
	}

	if err := s.repo.SaveRefreshToken(ctx, refreshRecord); err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh, // send raw; only hash is stored in DB
	}, nil
}
