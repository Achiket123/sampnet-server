package employees

import (
	"context"
	domain "server/internal/domain/employees"
	authDomain "server/internal/domain/auth"
	orgDomain "server/internal/domain/organisation"
	"server/internal/platform/miscallenous"
	"errors"
)

type service struct {
	repo     domain.Repository
	userRepo authDomain.Repository
	orgRepo  orgDomain.Repository
}

func NewService(repo domain.Repository, userRepo authDomain.Repository, orgRepo orgDomain.Repository) domain.UseCase {
	return &service{
		repo:     repo,
		userRepo: userRepo,
		orgRepo:  orgRepo,
	}
}

func (s *service) AddEmployee(ctx context.Context, emp *domain.Employee) error {
	// If UserID is 0, we might need to create the user first (based on original logic)
	// But usually, we'd expect the user to exist or provide user details.
	// Looking at AddEmployees, it takes user details.
	
	if emp.UserID == 0 {
		user := &authDomain.User{
			FirstName:   emp.User.FirstName,
			LastName:    emp.User.LastName,
			Email:       emp.User.Email,
			PhoneNumber: emp.User.PhoneNumber,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return err
		}
		emp.UserID = user.ID
	}

	return s.repo.Create(ctx, emp)
}

func (s *service) GetEmployees(ctx context.Context, orgID uint) ([]domain.Employee, error) {
	return s.repo.GetEmployeesByOrg(ctx, orgID)
}

func (s *service) GetEmployee(ctx context.Context, id uint) (*domain.Employee, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) UpdateEmployee(ctx context.Context, emp *domain.Employee) error {
	return s.repo.Update(ctx, emp)
}

func (s *service) DeleteEmployee(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) SearchEmployees(ctx context.Context, query string) ([]domain.Employee, error) {
	return s.repo.Search(ctx, query)
}

func (s *service) MakeManager(ctx context.Context, manager *domain.Manager) error {
	// Validate user
	if _, err := s.userRepo.GetByID(ctx, manager.UserID); err != nil {
		return errors.New("user not found")
	}
	// Validate organisation
	if _, err := s.orgRepo.GetByID(ctx, manager.OrganisationID); err != nil {
		return errors.New("organisation not found")
	}
	return s.repo.CreateManager(ctx, manager)
}

func (s *service) IsEmployeeOrManager(ctx context.Context, userID uint) (string, interface{}, string, error) {
	manager, err := s.repo.GetManagerByUserID(ctx, userID)
	if err == nil {
		token, err := miscallenous.GenerateJWTToken(manager, "manager", manager.UserID)
		return "Manager found", manager, token, err
	}

	employee, err := s.repo.GetByUserID(ctx, userID)
	if err == nil {
		token, err := miscallenous.GenerateJWTToken(employee, "employee", employee.UserID)
		return "Employee found", employee, token, err
	}

	return "", nil, "", errors.New("user not found as employee or manager")
}

func (s *service) CreateBoss(ctx context.Context, boss *domain.Boss) error {
	return s.repo.CreateBoss(ctx, boss)
}

func (s *service) GetBoss(ctx context.Context, userID uint) (*domain.Boss, error) {
	return s.repo.GetBossByUserID(ctx, userID)
}
