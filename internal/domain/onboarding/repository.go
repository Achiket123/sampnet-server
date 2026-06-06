package onboarding

import (
	"context"
)

type Repository interface {
	GetByUserID(ctx context.Context, userID uint) (*OnboardingProgress, error)
	Update(ctx context.Context, progress *OnboardingProgress) error
}
