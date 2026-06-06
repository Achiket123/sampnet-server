package onboarding

import (
	"context"
	"server/internal/domain/onboarding"
)

type service struct {
	repo onboarding.Repository
}

func NewService(repo onboarding.Repository) onboarding.UseCase {
	return &service{repo: repo}
}

func (s *service) GetOnboardingProgress(ctx context.Context, userID uint) (*onboarding.OnboardingProgress, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *service) UpdateOnboardingProgress(ctx context.Context, progress *onboarding.OnboardingProgress) error {
	return s.repo.Update(ctx, progress)
}
