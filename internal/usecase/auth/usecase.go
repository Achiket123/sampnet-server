package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	domain "server/internal/domain/auth"
	empDomain "server/internal/domain/employees"
	"server/internal/platform/database/models"
	mailerPlatform "server/internal/platform/mailer"
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
	ErrInvitePending      = errors.New("please accept your invitation to set your password")
)

type service struct {
	repo    domain.Repository
	empRepo empDomain.Repository
	mailer  mailerPlatform.Mailer
}

func NewService(repo domain.Repository, empRepo empDomain.Repository, mailer mailerPlatform.Mailer) domain.UseCase {
	return &service{repo: repo, empRepo: empRepo, mailer: mailer}
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

	if user.HashedPassword == "" {
		return domain.TokenPair{}, ErrInvitePending
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

	token, err := miscallenous.GenerateJWTToken(employee, "employee", employee.UserID)
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
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		IsVerified:  user.IsVerified,
		ProfilePic:  user.ProfilePic,
		City:        user.City,
		Country:     user.Country,
		DateOfBirth: user.DateOfBirth,
		LastLoginAt: user.LastLoginAt,
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

func (s *service) SendVerificationEmail(ctx context.Context, userID uint) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if user.IsVerified {
		return errors.New("user is already verified")
	}

	var token string
	active, err := s.repo.GetActiveEmailVerificationByUserID(ctx, userID)
	if err == nil && active != nil {
		token = active.Token
	} else {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			return err
		}
		token = hex.EncodeToString(b)
		ev := &domain.EmailVerification{
			UserID:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		if err := s.repo.CreateEmailVerification(ctx, ev); err != nil {
			return err
		}
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "https://sampnet.achiket.site"
	}
	verificationURL := fmt.Sprintf("%s/#/verify-email-landing?token=%s", frontendURL, token)

	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333333; margin: 0; padding: 20px;">
  <div style="max-width: 600px; margin: 0 auto; border: 1px solid #e0e0e0; border-radius: 8px; padding: 20px; box-shadow: 0 4px 6px rgba(0,0,0,0.05);">
    <h2 style="color: #2c3e50; border-bottom: 2px solid #e74c3c; padding-bottom: 10px; margin-top: 0;">Verify Your Email Address</h2>
    <p>Hello %s,</p>
    <p>Thank you for signing up! Please verify your email address by clicking the button below:</p>
    <div style="text-align: center; margin: 30px 0;">
      <a href="%s" style="background-color: #e74c3c; color: #ffffff; text-decoration: none; padding: 12px 24px; border-radius: 4px; font-weight: bold; display: inline-block;">Verify Email</a>
    </div>
    <p style="font-size: 0.9em; color: #7f8c8d;">If the button doesn't work, copy and paste this link into your web browser:</p>
    <p style="font-size: 0.9em; word-break: break-all; color: #e74c3c;"><a href="%s" style="color: #e74c3c;">%s</a></p>
    <hr style="border: 0; border-top: 1px solid #e0e0e0; margin: 20px 0;" />
    <p style="font-size: 0.8em; color: #95a5a6; text-align: center;">This link will expire in <strong>24 hours</strong>.</p>
  </div>
</body>
</html>`, user.FirstName, verificationURL, verificationURL, verificationURL)

	return s.mailer.SendMail(user.Email, "Verify your email address", htmlBody)
}

func (s *service) VerifyEmail(ctx context.Context, token string) error {
	ev, err := s.repo.GetEmailVerificationByToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired verification token")
	}

	if ev.UsedAt != nil {
		return errors.New("this verification token has already been used")
	}

	if time.Now().After(ev.ExpiresAt) {
		return errors.New("this verification token has expired")
	}

	user, err := s.repo.GetByID(ctx, ev.UserID)
	if err != nil {
		return ErrUserNotFound
	}

	user.IsVerified = true
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	now := time.Now()
	ev.UsedAt = &now
	return s.repo.UpdateEmailVerification(ctx, ev)
}

func (s *service) GetMe(ctx context.Context, userID uint) (domain.TokenPair, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return domain.TokenPair{}, ErrUserNotFound
	}
	return s.issueTokenPair(ctx, user)
}
