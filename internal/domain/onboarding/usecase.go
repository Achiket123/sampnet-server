package onboarding

import (
	"context"
)

type UseCase interface {
	GetOnboardingProgress(ctx context.Context, userID uint) (*OnboardingProgress, error)
	UpdateOnboardingProgress(ctx context.Context, progress *OnboardingProgress) error
}
