package auth

import (
	"context"
	"errors"
	domain "server/internal/domain/auth"
	empDomain "server/internal/domain/employees"
	"server/internal/platform/miscallenous"
	"server/internal/platform/database/models"
	"time"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound      = errors.New("user not found")
	ErrNotAnEmployee     = errors.New("user is not an employee")
)

type service struct {
	repo    domain.Repository
	empRepo empDomain.Repository
}

func NewService(repo domain.Repository, empRepo empDomain.Repository) domain.UseCase {
	return &service{
		repo:    repo,
		empRepo: empRepo,
	}
}

func (s *service) SignUp(ctx context.Context, user *domain.User, password string) (string, error) {
	hashedPassword, err := miscallenous.HashPassword(password)
	if err != nil {
		return "", err
	}
	user.HashedPassword = hashedPassword
	user.LastLoginAt = time.Now()

	if err := s.repo.Create(ctx, user); err != nil {
		return "", err
	}

	// Map domain user to model for JWT generation (as current JWT helper expects model)
	// Ideally, the JWT helper should take a domain entity or a generic interface.
	userModel := models.UserModel{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	userModel.ID = user.ID

	token, err := miscallenous.GenerateJWTToken(userModel, "user", user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !miscallenous.VerifyPassword(user.HashedPassword, password) {
		return "", ErrInvalidCredentials
	}

	userModel := models.UserModel{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	userModel.ID = user.ID

	token, err := miscallenous.GenerateJWTToken(userModel, "user", user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) CompleteSignIn(ctx context.Context, email, phone, password, city, country, profilePic string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		user, err = s.repo.GetByPhoneNumber(ctx, phone)
		if err != nil {
			return "", ErrUserNotFound
		}
	}

	hashedPassword, err := miscallenous.HashPassword(password)
	if err != nil {
		return "", err
	}

	user.HashedPassword = hashedPassword
	user.City = city
	user.Country = country
	user.ProfilePic = profilePic
	user.IsVerified = true
	user.LastLoginAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return "", err
	}

	userModel := models.UserModel{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	userModel.ID = user.ID

	token, err := miscallenous.GenerateJWTToken(userModel, "user", user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) ValidateEmployee(ctx context.Context, userID uint) (string, error) {
	employee, err := s.empRepo.GetByUserID(ctx, userID)
	if err != nil {
		return "", ErrNotAnEmployee
	}

	// Map to model for JWT generation
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
