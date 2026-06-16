package invites

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	authDomain "server/internal/domain/auth"
	employeesDomain "server/internal/domain/employees"
	domain "server/internal/domain/invites"
	"server/internal/platform/database/models"
	"server/internal/platform/mailer"
	"server/internal/platform/miscallenous"
)

type service struct {
	repo         domain.Repository
	userRepo     authDomain.Repository
	employeeRepo employeesDomain.Repository
	mailer       mailer.Mailer
}

func NewService(repo domain.Repository, userRepo authDomain.Repository, employeeRepo employeesDomain.Repository, m mailer.Mailer) domain.UseCase {
	return &service{
		repo:         repo,
		userRepo:     userRepo,
		employeeRepo: employeeRepo,
		mailer:       m,
	}
}

func (s *service) InviteEmployee(ctx context.Context, invite *domain.EmployeeInvite) error {
	// 1. Check duplicate
	existing, err := s.repo.GetByEmail(ctx, invite.Email, invite.OrganisationID)
	if err == nil && existing != nil {
		return errors.New("an active invitation already exists for this email in the organisation")
	}

	// 2. Generate secure token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	invite.Token = hex.EncodeToString(b)

	// 3. Set properties
	invite.Status = "pending"
	invite.ExpiresAt = time.Now().Add(72 * time.Hour)

	// 4. Persist to DB
	if err := s.repo.Create(ctx, invite); err != nil {
		return err
	}

	// 5. Build accept URL
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "https://sampnet.achiket.site" // fallback default
	}
	inviteURL := fmt.Sprintf("%s/accept-invite?token=%s", frontendURL, invite.Token)

	// 6. Build HTML body
	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333333; margin: 0; padding: 20px;">
  <div style="max-width: 600px; margin: 0 auto; border: 1px solid #e0e0e0; border-radius: 8px; padding: 20px; box-shadow: 0 4px 6px rgba(0,0,0,0.05);">
    <h2 style="color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 10px; margin-top: 0;">You're Invited!</h2>
    <p>Hello %s,</p>
    <p>You have been invited to join the organisation as an employee. Click the button below to set up your account and get started:</p>
    <div style="text-align: center; margin: 30px 0;">
      <a href="%s" style="background-color: #3498db; color: #ffffff; text-decoration: none; padding: 12px 24px; border-radius: 4px; font-weight: bold; display: inline-block;">Accept Invitation</a>
    </div>
    <p style="font-size: 0.9em; color: #7f8c8d;">If the button doesn't work, copy and paste this link into your web browser:</p>
    <p style="font-size: 0.9em; word-break: break-all; color: #3498db;"><a href="%s" style="color: #3498db;">%s</a></p>
    <hr style="border: 0; border-top: 1px solid #e0e0e0; margin: 20px 0;" />
    <p style="font-size: 0.8em; color: #95a5a6; text-align: center;">This invitation expires in <strong>72 hours</strong>.</p>
  </div>
</body>
</html>`, invite.FirstName, inviteURL, inviteURL, inviteURL)

	// 7. Send the email
	err = s.mailer.SendMail(invite.Email, "You're invited", htmlBody)
	log.Default().Println(err)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AcceptInvite(ctx context.Context, token string, password string) (authDomain.TokenPair, error) {
	// 1. Look up invite
	invite, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		return authDomain.TokenPair{}, errors.New("invalid or expired invitation link")
	}

	// 2. Validate state
	if invite.Status != "pending" {
		return authDomain.TokenPair{}, errors.New("this invitation has already been used")
	}
	if time.Now().After(invite.ExpiresAt) {
		return authDomain.TokenPair{}, errors.New("this invitation has expired")
	}

	// 3. Check if user already exists
	var user *authDomain.User
	hashed, errPassword := miscallenous.HashPassword(password)
	if errPassword != nil {
		return authDomain.TokenPair{}, errPassword
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, invite.Email)
	if err != nil {
		// User does not exist — create them
		user = &authDomain.User{
			FirstName:      invite.FirstName,
			LastName:       invite.LastName,
			Email:          invite.Email,
			PhoneNumber:    invite.PhoneNumber,
			HashedPassword: hashed,
			IsVerified:     true,
		}
	} else {
		// User exists — update verification status and password
		user = existingUser
		user.HashedPassword = hashed
		user.IsVerified = true
	}

	// 4. Build employee record
	employee := &employeesDomain.Employee{
		UserID:         user.ID,
		OrganisationID: invite.OrganisationID,
		EmploymentID:   invite.EmploymentID,
		Email:          invite.Email,
	}

	// 5. Delegate atomically
	if err := s.repo.AcceptInvite(ctx, invite.ID, user, employee); err != nil {
		return authDomain.TokenPair{}, err
	}

	// 6. Generate token pair
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

	accessToken, err := miscallenous.GenerateJWTToken(userModel, "user", user.ID)
	if err != nil {
		return authDomain.TokenPair{}, err
	}

	rawRefresh, err := miscallenous.GenerateRefreshToken()
	if err != nil {
		return authDomain.TokenPair{}, err
	}

	refreshRecord := &authDomain.RefreshToken{
		UserID:    user.ID,
		TokenHash: miscallenous.HashRefreshToken(rawRefresh),
		ExpiresAt: miscallenous.RefreshTokenExpiresAt(),
		Revoked:   false,
	}

	if err := s.userRepo.SaveRefreshToken(ctx, refreshRecord); err != nil {
		return authDomain.TokenPair{}, err
	}

	return authDomain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}

func (s *service) GetInvitesByOrg(ctx context.Context, orgID uint) ([]domain.EmployeeInvite, error) {
	return s.repo.GetByOrg(ctx, orgID)
}

func (s *service) ResendInvite(ctx context.Context, email string, orgID uint) error {
	// 1. Fetch invite
	invite, err := s.repo.GetByEmail(ctx, email, orgID)
	if err != nil {
		return err
	}

	// 2. Validate organisation match
	if invite.OrganisationID != orgID {
		return errors.New("unauthorized: invite belongs to a different organisation")
	}

	// 3. Check if invite is already accepted
	if invite.Status == "accepted" {
		return errors.New("cannot resend: invitation has already been accepted")
	}

	// 4. Check if the user is registered and verified
	existingUser, err := s.userRepo.GetByEmail(ctx, invite.Email)
	if err == nil && existingUser != nil && existingUser.IsVerified {
		return errors.New("cannot resend: user is already registered and verified")
	}

	// 5. Generate secure token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	invite.Token = hex.EncodeToString(b)

	// 6. Reset properties
	invite.Status = "pending"
	invite.ExpiresAt = time.Now().Add(72 * time.Hour)

	// 7. Update in DB
	if err := s.repo.Update(ctx, invite); err != nil {
		return err
	}

	// 8. Build accept URL
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "https://sampnet.achiket.site"
	}
	inviteURL := fmt.Sprintf("%s/#/accept-invite?token=%s", frontendURL, invite.Token)

	// 9. Build HTML body
	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333333; margin: 0; padding: 20px;">
  <div style="max-width: 600px; margin: 0 auto; border: 1px solid #e0e0e0; border-radius: 8px; padding: 20px; box-shadow: 0 4px 6px rgba(0,0,0,0.05);">
    <h2 style="color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 10px; margin-top: 0;">You're Invited!</h2>
    <p>Hello %s,</p>
    <p>You have been invited to join the organisation as an employee. Click the button below to set up your account and get started:</p>
    <div style="text-align: center; margin: 30px 0;">
      <a href="%s" style="background-color: #3498db; color: #ffffff; text-decoration: none; padding: 12px 24px; border-radius: 4px; font-weight: bold; display: inline-block;">Accept Invitation</a>
    </div>
    <p style="font-size: 0.9em; color: #7f8c8d;">If the button doesn't work, copy and paste this link into your web browser:</p>
    <p style="font-size: 0.9em; word-break: break-all; color: #3498db;"><a href="%s" style="color: #3498db;">%s</a></p>
    <hr style="border: 0; border-top: 1px solid #e0e0e0; margin: 20px 0;" />
    <p style="font-size: 0.8em; color: #95a5a6; text-align: center;">This invitation expires in <strong>72 hours</strong>.</p>
  </div>
</body>
</html>`, invite.FirstName, inviteURL, inviteURL, inviteURL)

	// 10. Send the email
	err = s.mailer.SendMail(invite.Email, "You're invited", htmlBody)
	if err != nil {
		return err
	}
	return nil
}
